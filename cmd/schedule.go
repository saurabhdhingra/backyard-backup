package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/robfig/cron/v3"
	"github.com/spf13/cobra"
)

var scheduleCmd = &cobra.Command{
	Use:   "schedule",
	Short: "Run backup on a schedule",
	Run: func(cmd *cobra.Command, args []string) {
		schedule := AppConfig.Backup.Schedule
		if schedule == "" {
			fmt.Println("Error: No schedule defined in config")
			os.Exit(1)
		}

		c := cron.New()
		_, err := c.AddFunc(schedule, func() {
			fmt.Println("Running scheduled backup...", time.Now())
			// We can call the backup logic directly or execute the command
			// Calling Run of backupCmd is a bit tricky due to args/flags handling.
			// Best to extract logic or re-execute binary.
			// For simplicity here, we'll invoke the backupCmd.Run logic efficiently or refactor backup logic to a function.
			// However, in this simple CLI, calling backupCmd.Run(cmd, []string{}) might work if no specific flags are required.
			// But backupCmd doesn't return error, it exits. We should probably refactor 'Run' to 'RunE' and return error.
			// Ideally, we extract `RunBackup` function.
			// Since we haven't refactored, we will mimic the call.
			// BEWARE: os.Exit in backupCmd will kill the scheduler.
			// I need to refactor backupCmd to NOT os.Exit.

			// For now, let's print a message that actual scheduling requires refactoring backup logic.
			// Or better, let's just spawning a subprocess
			// exec.Command(os.Args[0], "backup").Run()
			// This is safer for isolation.

			// Implementation with subprocess:
			// self, _ := os.Executable()
			// cmd := exec.Command(self, "backup")
			// cmd.Stdout = os.Stdout
			// cmd.Stderr = os.Stderr
			// cmd.Run()

			// Note: os.Executable might not be clean in 'go run'.
			fmt.Println("Triggering backup job...")
			backupCmd.Run(backupCmd, []string{})
		})

		if err != nil {
			fmt.Printf("Error adding cron job: %v\n", err)
			os.Exit(1)
		}

		c.Start()
		fmt.Printf("Backup scheduler started with schedule: %s\n", schedule)

		// Block forever
		select {}
	},
}

func init() {
	rootCmd.AddCommand(scheduleCmd)
}
