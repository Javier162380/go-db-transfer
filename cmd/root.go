package cmd

import (
	"github.com/spf13/cobra"
)

// config file
var CONFIGFILE string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "go-db-transfer",
	Short: "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:
Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
}

// Init App Cli 
func Init() {

	rootCmd.Flags().StringVar(&CONFIGFILE, "config", "", "Config file for the app")
	rootCmd.Execute()
}