package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/saurabhdhingra/backyard-backup/internal/archiver"
	"github.com/saurabhdhingra/backyard-backup/internal/db"
	"github.com/saurabhdhingra/backyard-backup/internal/storage"
	"github.com/spf13/cobra"
)

var restoreFile string

var restoreCmd = &cobra.Command{
	Use:   "restore",
	Short: "Restore a database from a backup",
	Run: func(cmd *cobra.Command, args []string) {
		if restoreFile == "" {
			fmt.Println("Error: --file flag is required")
			os.Exit(1)
		}

		fmt.Println("Starting restore...")
		startTime := time.Now()

		// 1. Initialize Database
		dbConfig := db.Config{
			Type:     AppConfig.Database.Type,
			Host:     AppConfig.Database.Host,
			Port:     AppConfig.Database.Port,
			User:     AppConfig.Database.User,
			Password: AppConfig.Database.Password,
			DBName:   AppConfig.Database.DBName,
			DSN:      AppConfig.Database.DSN,
		}
		database, err := db.NewDatabase(dbConfig)
		if err != nil {
			fmt.Printf("Error initializing database: %v\n", err)
			os.Exit(1)
		}
		// Note: For restore, we might need a connection, but pg_dump/psql usually handles it via cli args.
		// However, establishing checking connectivity is good practice.
		if err := database.Connect(); err != nil {
			fmt.Printf("Error connecting to database: %v\n", err)
			os.Exit(1)
		}
		defer database.Close()

		// 2. Initialize Storage
		storeConfig := storage.Config{
			Type:      AppConfig.Storage.Type,
			BasePath:  AppConfig.Storage.Path,
			Bucket:    AppConfig.Storage.Bucket,
			Region:    AppConfig.Storage.Region,
			AccessKey: AppConfig.Storage.AccessKey,
			SecretKey: AppConfig.Storage.SecretKey,
		}
		store, err := storage.NewStorage(storeConfig)
		if err != nil {
			fmt.Printf("Error initializing storage: %v\n", err)
			os.Exit(1)
		}

		// 3. Create Temp Directory
		tmpDir, err := os.MkdirTemp("", "backyard-restore")
		if err != nil {
			fmt.Printf("Error creating temp dir: %v\n", err)
			os.Exit(1)
		}
		defer os.RemoveAll(tmpDir)

		// 4. Download from Storage
		localDownloadPath := filepath.Join(tmpDir, filepath.Base(restoreFile))
		fmt.Printf("Downloading backup from storage: %s\n", restoreFile)
		if err := store.Download(restoreFile, localDownloadPath); err != nil {
			fmt.Printf("Error downloading file: %v\n", err)
			os.Exit(1)
		}

		finalRestorePath := localDownloadPath

		// 5. Decompress if needed
		// Simple check: if ends with .gz
		if strings.HasSuffix(restoreFile, ".gz") {
			fmt.Println("Decompressing backup...")
			decompressedPath := strings.TrimSuffix(localDownloadPath, ".gz")
			if err := archiver.Decompress(localDownloadPath, decompressedPath); err != nil {
				fmt.Printf("Error decompressing file: %v\n", err)
				os.Exit(1)
			}
			finalRestorePath = decompressedPath
		}

		// 6. Restore to Database
		fmt.Println("Restoring to database...")
		if err := database.Restore(finalRestorePath); err != nil {
			fmt.Printf("Error restoring database: %v\n", err)
			os.Exit(1)
		}

		duration := time.Since(startTime)
		fmt.Printf("Restore completed successfully in %s\n", duration)
	},
}

func init() {
	rootCmd.AddCommand(restoreCmd)
	restoreCmd.Flags().StringVarP(&restoreFile, "file", "f", "", "Path to the backup file in storage to restore")
}
