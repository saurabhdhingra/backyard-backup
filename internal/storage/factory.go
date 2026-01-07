package storage

import (
	"fmt"
)

func NewStorage(cfg Config) (Storage, error) {
	switch cfg.Type {
	case "local":
		return NewLocal(cfg), nil
	case "s3", "aws":
		return NewS3(cfg)
	default:
		return nil, fmt.Errorf("unsupported storage type: %s", cfg.Type)
	}
}
