package storage

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
)

type Local struct {
	Config Config
}

func NewLocal(cfg Config) *Local {
	return &Local{Config: cfg}
}

func (l *Local) Upload(localPath string, remotePath string) error {
	// In local storage, remotePath is relative to the BasePath in config
	destPath := filepath.Join(l.Config.BasePath, remotePath)

	// Create directory if not exists
	if err := os.MkdirAll(filepath.Dir(destPath), 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	// Copy file
	srcFile, err := os.Open(localPath)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	destFile, err := os.Create(destPath)
	if err != nil {
		return err
	}
	defer destFile.Close()

	if _, err := io.Copy(destFile, srcFile); err != nil {
		return err
	}

	return nil
}

func (l *Local) Download(remotePath string, localPath string) error {
	srcPath := filepath.Join(l.Config.BasePath, remotePath)

	srcFile, err := os.Open(srcPath)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	destFile, err := os.Create(localPath)
	if err != nil {
		return err
	}
	defer destFile.Close()

	if _, err := io.Copy(destFile, srcFile); err != nil {
		return err
	}

	return nil
}

func (l *Local) StreamUpload(reader io.Reader, remotePath string) error {
	destPath := filepath.Join(l.Config.BasePath, remotePath)

	if err := os.MkdirAll(filepath.Dir(destPath), 0755); err != nil {
		return err
	}

	destFile, err := os.Create(destPath)
	if err != nil {
		return err
	}
	defer destFile.Close()

	if _, err := io.Copy(destFile, reader); err != nil {
		return err
	}

	return nil
}
