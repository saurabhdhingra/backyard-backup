package db

import (
	"database/sql"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type SQLite struct {
	Config Config
	conn   *sql.DB
}

func NewSQLite(cfg Config) *SQLite {
	return &SQLite{Config: cfg}
}

func (s *SQLite) Connect() error {
	// For SQLite, Host/Port/User/Pass are irrelevant, usually just DBName is the path
	if s.Config.DBName == "" {
		return fmt.Errorf("sqlite database path (dbname) is required")
	}

	db, err := sql.Open("sqlite3", s.Config.DBName)
	if err != nil {
		return fmt.Errorf("failed to open sqlite database: %w", err)
	}

	if err := db.Ping(); err != nil {
		return fmt.Errorf("failed to ping sqlite database: %w", err)
	}

	s.conn = db
	return nil
}

func (s *SQLite) Close() error {
	if s.conn != nil {
		return s.conn.Close()
	}
	return nil
}

func (s *SQLite) Dump(destinationPath string) (string, error) {
	// Destination file
	baseName := filepath.Base(s.Config.DBName)
	fileName := fmt.Sprintf("%s_%s.sql", baseName, time.Now().Format("20060102_150405"))
	fullPath := filepath.Join(destinationPath, fileName)

	// Use sqlite3 command line tool to dump
	// syntax: sqlite3 <dbfile> .dump > <outfile>
	cmd := exec.Command("sqlite3", s.Config.DBName, ".dump")

	outFile, err := os.Create(fullPath)
	if err != nil {
		return "", fmt.Errorf("failed to create dump file: %w", err)
	}
	defer outFile.Close()

	cmd.Stdout = outFile

	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("sqlite3 dump failed: %w", err)
	}

	return fullPath, nil
}

func (s *SQLite) Restore(sourcePath string) error {
	// syntax: sqlite3 <dbfile> < <infile>
	cmd := exec.Command("sqlite3", s.Config.DBName)

	inFile, err := os.Open(sourcePath)
	if err != nil {
		return fmt.Errorf("failed to open source dump file: %w", err)
	}
	defer inFile.Close()

	cmd.Stdin = inFile

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("sqlite3 restore failed: %w", err)
	}

	return nil
}
