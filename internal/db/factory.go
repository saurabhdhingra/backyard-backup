package db

import (
	"fmt"
)

func NewDatabase(cfg Config) (Database, error) {
	switch cfg.Type {
	case "postgres", "postgresql":
		return NewPostgres(cfg), nil
	case "sqlite", "sqlite3":
		return NewSQLite(cfg), nil
	case "mysql":
		return NewMySQL(cfg), nil
	case "mongodb", "mongo":
		return NewMongoDB(cfg), nil
	default:
		return nil, fmt.Errorf("unsupported database type: %s", cfg.Type)
	}
}
