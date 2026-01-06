package db

import (
	"database/sql"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type MySQL struct {
	Config Config
	conn   *sql.DB
}

func NewMySQL(cfg Config) *MySQL {
	return &MySQL{Config: cfg}
}

func (m *MySQL) Connect() error {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s",
		m.Config.User, m.Config.Password, m.Config.Host, m.Config.Port, m.Config.DBName)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return fmt.Errorf("failed to open mysql connection: %w", err)
	}

	if err := db.Ping(); err != nil {
		return fmt.Errorf("failed to ping mysql database: %w", err)
	}

	m.conn = db
	return nil
}

func (m *MySQL) Close() error {
	if m.conn != nil {
		return m.conn.Close()
	}
	return nil
}

func (m *MySQL) Dump(destinationPath string) (string, error) {
	// Generate filename
	fileName := fmt.Sprintf("%s_%s.sql", m.Config.DBName, time.Now().Format("20060102_150405"))
	fullPath := filepath.Join(destinationPath, fileName)

	// mysqldump command
	// mysqldump -h host -P port -u user -p[password] dbname > outfile
	// Note: putting password in command args is insecure, better to use cnf file or ENV.
	// MYSQL_PWD env var is supported by mysqldump.

	cmd := exec.Command("mysqldump",
		"-h", m.Config.Host,
		"-P", fmt.Sprintf("%d", m.Config.Port),
		"-u", m.Config.User,
		m.Config.DBName,
	)

	cmd.Env = append(os.Environ(), fmt.Sprintf("MYSQL_PWD=%s", m.Config.Password))

	outFile, err := os.Create(fullPath)
	if err != nil {
		return "", fmt.Errorf("failed to create dump file: %w", err)
	}
	defer outFile.Close()

	cmd.Stdout = outFile

	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("mysqldump failed: %w", err)
	}

	return fullPath, nil
}

func (m *MySQL) Restore(sourcePath string) error {
	// mysql -u user -p dbname < infile
	cmd := exec.Command("mysql",
		"-h", m.Config.Host,
		"-P", fmt.Sprintf("%d", m.Config.Port),
		"-u", m.Config.User,
		m.Config.DBName,
	)

	cmd.Env = append(os.Environ(), fmt.Sprintf("MYSQL_PWD=%s", m.Config.Password))

	inFile, err := os.Open(sourcePath)
	if err != nil {
		return fmt.Errorf("failed to open restore file: %w", err)
	}
	defer inFile.Close()

	cmd.Stdin = inFile

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("mysql restore failed: %w", err)
	}

	return nil
}
