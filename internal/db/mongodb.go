package db

import (
	"context"
	"fmt"
	"os/exec"
	"path/filepath"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type MongoDB struct {
	Config Config
	client *mongo.Client
}

func NewMongoDB(cfg Config) *MongoDB {
	return &MongoDB{Config: cfg}
}

func (m *MongoDB) Connect() error {
	var uri string
	if m.Config.DSN != "" {
		uri = m.Config.DSN
	} else {
		// Construct URI: mongodb://[user:pass@]host:port[/dbname]
		creds := ""
		if m.Config.User != "" && m.Config.Password != "" {
			creds = fmt.Sprintf("%s:%s@", m.Config.User, m.Config.Password)
		}
		uri = fmt.Sprintf("mongodb://%s%s:%d", creds, m.Config.Host, m.Config.Port)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		return fmt.Errorf("failed to create mongo client: %w", err)
	}

	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		return fmt.Errorf("failed to ping mongodb: %w", err)
	}

	m.client = client
	return nil
}

func (m *MongoDB) Close() error {
	if m.client != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		return m.client.Disconnect(ctx)
	}
	return nil
}

func (m *MongoDB) Dump(destinationPath string) (string, error) {
	// Generate filename
	dbName := m.Config.DBName
	if dbName == "" {
		dbName = "all_dbs"
	}
	fileName := fmt.Sprintf("%s_%s.archive", dbName, time.Now().Format("20060102_150405"))
	fullPath := filepath.Join(destinationPath, fileName)

	// Build mongodump command.
	// We use --archive to output a single file.
	args := []string{"--archive=" + fullPath}

	if m.Config.DSN != "" {
		args = append(args, "--uri="+m.Config.DSN)
	} else {
		args = append(args, "--host", m.Config.Host)
		args = append(args, "--port", fmt.Sprintf("%d", m.Config.Port))
		if m.Config.User != "" {
			args = append(args, "--username", m.Config.User)
		}
		if m.Config.Password != "" {
			args = append(args, "--password", m.Config.Password)
		}
		if m.Config.DBName != "" {
			args = append(args, "--db", m.Config.DBName)
		}
	}

	cmd := exec.Command("mongodump", args...)

	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("mongodump failed: %s, output: %s", err, string(output))
	}

	return fullPath, nil
}

func (m *MongoDB) Restore(sourcePath string) error {
	// Build mongorestore command
	args := []string{"--archive=" + sourcePath}

	if m.Config.DSN != "" {
		args = append(args, "--uri="+m.Config.DSN)
	} else {
		args = append(args, "--host", m.Config.Host)
		args = append(args, "--port", fmt.Sprintf("%d", m.Config.Port))
		if m.Config.User != "" {
			args = append(args, "--username", m.Config.User)
		}
		if m.Config.Password != "" {
			args = append(args, "--password", m.Config.Password)
		}
		// Dropping before restore is often safer for clean restore, but risky.
		// We'll leave it to user to clean or use flags if needed,
		// but standard restore adds data.
	}

	cmd := exec.Command("mongorestore", args...)

	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("mongorestore failed: %s, output: %s", err, string(output))
	}

	return nil
}
