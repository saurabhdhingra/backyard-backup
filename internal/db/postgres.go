package db

import (
	"database/sql"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	_ "github.com/lib/pq"
)

type Postgres struct {
	Config Config
	conn   *sql.DB
}

func NewPostgres(cfg Config) *Postgres {
	return &Postgres{Config: cfg}
}

func (p *Postgres) Connect() error {
	var connStr string
	if p.Config.DSN != "" {
		connStr = p.Config.DSN
	} else {
		connStr = fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
			p.Config.Host, p.Config.Port, p.Config.User, p.Config.Password, p.Config.DBName)
	}

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return fmt.Errorf("failed to open connection: %w", err)
	}

	if err := db.Ping(); err != nil {
		return fmt.Errorf("failed to ping database: %w", err)
	}

	p.conn = db
	return nil
}

func (p *Postgres) Close() error {
	if p.conn != nil {
		return p.conn.Close()
	}
	return nil
}

func (p *Postgres) Dump(destinationPath string) (string, error) {
	// Generate filename if destination is a directory
	// Use DBName from config if available, otherwise "db"
	dbName := p.Config.DBName
	if dbName == "" {
		dbName = "db"
	}
	fileName := fmt.Sprintf("%s_%s.sql", dbName, time.Now().Format("20060102_150405"))
	fullPath := filepath.Join(destinationPath, fileName)

	// We'll use pg_dump command
	var cmd *exec.Cmd

	if p.Config.DSN != "" {
		// If DSN is provided, use it directly as the dbname argument
		cmd = exec.Command("pg_dump", p.Config.DSN, "-f", fullPath)
	} else {
		// PGPASSWORD environment variable is used to pass password to pg_dump to avoid prompt
		cmd = exec.Command("pg_dump",
			"-h", p.Config.Host,
			"-p", fmt.Sprintf("%d", p.Config.Port),
			"-U", p.Config.User,
			"-d", p.Config.DBName,
			"-f", fullPath,
		)
		cmd.Env = append(os.Environ(), fmt.Sprintf("PGPASSWORD=%s", p.Config.Password))
	}

	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("pg_dump failed: %s, output: %s", err, string(output))
	}

	return fullPath, nil
}

func (p *Postgres) Restore(sourcePath string) error {
	// PGPASSWORD environment variable is used
	var cmd *exec.Cmd

	if p.Config.DSN != "" {
		cmd = exec.Command("psql", p.Config.DSN, "-f", sourcePath)
	} else {
		cmd = exec.Command("psql",
			"-h", p.Config.Host,
			"-p", fmt.Sprintf("%d", p.Config.Port),
			"-U", p.Config.User,
			"-d", p.Config.DBName,
			"-f", sourcePath,
		)
		cmd.Env = append(os.Environ(), fmt.Sprintf("PGPASSWORD=%s", p.Config.Password))
	}

	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("restore failed: %s, output: %s", err, string(output))
	}

	return nil
}
