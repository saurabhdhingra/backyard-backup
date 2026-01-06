package archiver

import (
	"compress/gzip"
	"fmt"
	"io"
	"os"
)

// Compress compresses the source file to the destination file using Gzip
func Compress(sourcePath, destPath string) error {
	srcFile, err := os.Open(sourcePath)
	if err != nil {
		return fmt.Errorf("failed to open source file: %w", err)
	}
	defer srcFile.Close()

	destFile, err := os.Create(destPath)
	if err != nil {
		return fmt.Errorf("failed to create destination file: %w", err)
	}
	defer destFile.Close()

	gw := gzip.NewWriter(destFile)
	defer gw.Close()

	if _, err := io.Copy(gw, srcFile); err != nil {
		return fmt.Errorf("failed to compress file: %w", err)
	}

	return nil
}

// Decompress decompresses the source file to the destination file using Gzip
func Decompress(sourcePath, destPath string) error {
	srcFile, err := os.Open(sourcePath)
	if err != nil {
		return fmt.Errorf("failed to open source file: %w", err)
	}
	defer srcFile.Close()

	gr, err := gzip.NewReader(srcFile)
	if err != nil {
		return fmt.Errorf("failed to create gzip reader: %w", err)
	}
	defer gr.Close()

	destFile, err := os.Create(destPath)
	if err != nil {
		return fmt.Errorf("failed to create destination file: %w", err)
	}
	defer destFile.Close()

	if _, err := io.Copy(destFile, gr); err != nil {
		return fmt.Errorf("failed to decompress file: %w", err)
	}

	return nil
}
