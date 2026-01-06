package db

// Database interface defines the methods that any database provider must implement
type Database interface {
	// Connect establishes a connection to the database
	Connect() error

	// Dump performs a backup of the database to the specified path
	// returns the path to the backup file and an error if any
	Dump(destinationPath string) (string, error)

	// Restore attempts to restore the database from the specified file
	Restore(sourcePath string) error

	// Close closes the database connection
	Close() error
}

// Config holds common database configuration parameters
type Config struct {
	Type     string
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
	DSN      string
}
