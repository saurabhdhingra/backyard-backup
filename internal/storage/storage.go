package storage

import "io"

// Storage interface defines the methods for storage backends
type Storage interface {
	// Upload pushes a file to the storage
	Upload(localPath string, remotePath string) error

	// Download retrieves a file from the storage
	Download(remotePath string, localPath string) error

	// StreamUpload allows uploading from a reader (useful for piping compressed data)
	StreamUpload(reader io.Reader, remotePath string) error
}

// Config holds common storage configuration parameters
type Config struct {
	Type      string // "local", "s3", "gcs", "azure"
	Region    string // for cloud providers
	Bucket    string // for cloud providers
	BasePath  string // for local storage or prefix in cloud
	AccessKey string
	SecretKey string
}
