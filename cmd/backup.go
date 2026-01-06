package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/saurabhdhingra/backyard-backup/internal/archiver"
)

var backupCmd = &cobra.Command{
	Use:	"backup",
	Short:	"Perform a database backup",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Starting backup...")
		startTime := time.Now()

		dbConfig := db.Config{
			Type:	AppConfig.Database.Type,
			Host:	AppConfig.Database.Host,
			Port:	AppConfig.Database.Port,
			User: 	AppConfig.Database.Password,
			DBName: AppConfig.Database.DBName,
			DSN: 	AppConfig.Database.DSN,
		}

		database, err := db.NewDatabase(dbConfig)
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

		
	}
}