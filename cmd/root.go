package cmd

import (
	"fmt"
	"os"

	"github.com/saurabhdhingra/backyard-backup/internal/config"
	"github.com/spf13/cobra"
)

var cfgFile string
var AppConfig *config.Config

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "dbbackup",
	Short: "A CLI utility for backing up databases",
	Long: `Backyard Backup is a CLI tool to backup and restore various databases
supported (PostgreSQL, MySQL, MongoDB, SQLite) to local or cloud storage.`,
}

func Execute(){
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init(){
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "confif file (default id $HOME/.backyard-backup/config.yaml)")
}

func initConfig() {
	var err error
	AppConfig, err = config.LoadConfig(cfgFile)
	if err != nil {
		if cfgFile != "" {
			fmt.Printf("Error loading config file: %v\n", err)
			os.Exit(1)
		}
	} 
}