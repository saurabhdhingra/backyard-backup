package cmd

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/robfig/cron/v3"
	"github.com/saurabhdhingra/backyard-backup/internal/notify"
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
			fmt.Printf("[%s] Running scheduled backup...\n", time.Now().Format(time.RFC3339))

			if err := RunBackup(); err != nil {
				fmt.Printf("[%s] Scheduled backup failed: %v\n", time.Now().Format(time.RFC3339), err)
				// Send failure notification if configured
				if AppConfig.Notify.Enabled && AppConfig.Notify.SlackWebhook != "" {
					notifyErr := notifyBackupFailure(err)
					if notifyErr != nil {
						fmt.Printf("Warning: failed to send failure notification: %v\n", notifyErr)
					}
				}
				return
			}

			fmt.Printf("[%s] Scheduled backup completed successfully\n", time.Now().Format(time.RFC3339))
		})

		if err != nil {
			fmt.Printf("Error adding cron job: %v\n", err)
			os.Exit(1)
		}

		c.Start()
		fmt.Printf("Backup scheduler started with schedule: %s\n", schedule)
		fmt.Println("Press Ctrl+C to stop the scheduler")

		// Handle graceful shutdown
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
		<-sigChan

		fmt.Println("\nShutting down scheduler...")
		c.Stop()
		fmt.Println("Scheduler stopped")
	},
}

// notifyBackupFailure sends a Slack notification about backup failure
func notifyBackupFailure(err error) error {
	msg := fmt.Sprintf("🚨 Backup failed: %v", err)
	return notify.SendSlackNotification(AppConfig.Notify.SlackWebhook, msg)
}

func init() {
	rootCmd.AddCommand(scheduleCmd)
}
