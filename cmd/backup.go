package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/saurabhdhingra/backyard-backup/internal/archiver"
	"github.com/saurabhdhingra/backyard-backup/internal/db"
	"github.com/saurabhdhingra/backyard-backup/internal/notify"
	"github.com/saurabhdhingra/backyard-backup/internal/storage"
	"github.com/spf13/cobra"
)

var backupCmd = &cobra.Command{
	Use:   "backup",
	Short: "Perform a database backup",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Starting backup...")
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
		tmpDir, err := os.MkdirTemp("", "backyard-backup")
		if err != nil {
			fmt.Printf("Error creating temp dir: %v\n", err)
			os.Exit(1)
		}
		defer os.RemoveAll(tmpDir)

		// 4. Dump Database
		fmt.Println("Dumping database...")
		dumpPath, err := database.Dump(tmpDir)
		if err != nil {
			fmt.Printf("Error dumping database: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Database dumped to: %s\n", dumpPath)

		finalPath := dumpPath

		// 5. Compress if enabled
		if AppConfig.Backup.Compression {
			fmt.Println("Compressing backup...")
			compressedPath := dumpPath + ".gz"
			if err := archiver.Compress(dumpPath, compressedPath); err != nil {
				fmt.Printf("Error compressing backup: %v\n", err)
				os.Exit(1)
			}
			finalPath = compressedPath
			fmt.Printf("Backup compressed to: %s\n", finalPath)
		}

		// 6. Upload to Storage
		fmt.Println("Uploading to storage...")
		remotePath := filepath.Base(finalPath)
		if err := store.Upload(finalPath, remotePath); err != nil {
			fmt.Printf("Error uploading to storage: %v\n", err)
			os.Exit(1)
		}

		duration := time.Since(startTime)
		successMsg := fmt.Sprintf("Backup completed successfully in %s", duration)
		fmt.Println(successMsg)

		if AppConfig.Notify.Enabled && AppConfig.Notify.SlackWebhook != "" {
			fmt.Println("Sending Slack notification...")
			if err := notify.SendSlackNotification(AppConfig.Notify.SlackWebhook, successMsg); err != nil {
				fmt.Printf("Error sending notification: %v\n", err)
				// Don't fail the backup just because notification failed
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(backupCmd)
}
