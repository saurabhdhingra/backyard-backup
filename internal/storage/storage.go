package storage

import "io"

// Storage interface defines the methods for storage backends
type Storage interface {
	Upload(localPath string, remotePath string) error
	Download(remotePath string, localPath string) error
	StreamUpload(reader io.Reader, remotePath string) error
}

// Config holds common storage configuration parameters
type Config struct {
	Type      string 
	Region    string 
	Bucket    string 
	BasePath  string 
	AccessKey string
	SecretKey string
}
