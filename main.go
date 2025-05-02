package main

import (
	"os"

	"github.com/spf13/cobra"
)

var (
	temporalAddress string
	jsonOutput      bool
)

// rootCmd is the CLI root command.
var rootCmd = &cobra.Command{
	Use:     "cryptellation",
	Version: "1.0.0",
	Short:   "cryptellation - a CLI to manage Cryptellation system",
}

func main() {
	var errCode int

	// Set flags
	rootCmd.PersistentFlags().StringVarP(&temporalAddress,
		"temporal-address", "t", "localhost:7233", "Set output to JSON format")
	rootCmd.PersistentFlags().BoolVarP(&jsonOutput,
		"json", "j", false, "Set output to JSON format")

	// Set commands
	setExchangesCommands(rootCmd)
	setServicesCommands(rootCmd)

	// Execute command
	if err := rootCmd.Execute(); err != nil {
		os.Exit(errCode)
	}
}
