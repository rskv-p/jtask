package cmd

import (
	"fmt"
	"os"

	"github.com/rskv-p/jtask/pkg/x_config"
	"github.com/rskv-p/jtask/pkg/x_log"
	"github.com/spf13/cobra"
)

// ---------- Global Flag ----------
var pathFlag string // Global flag for task file path
var cfg x_config.Config

// ---------- Root Command Definition ----------
var rootCmd = &cobra.Command{
	Use:   "jt",                                                           // CLI command name
	Short: "JsonTask - lightning fast and small task runner",              // Short description
	Long:  "JT is a task runner for managing async and sequential tasks.", // Detailed description
	Run: func(cmd *cobra.Command, args []string) {
		// If no args are passed, show help
		if len(args) == 0 {
			cmd.Help() // Show the help message for the root command
		}
	},
}

// ---------- Command Execution ----------
func Execute() {
	// Log the command execution start
	x_log.Info().
		Str("command", rootCmd.Use).
		Msg("executing root command")

	// Execute the root command, handling errors if any
	if err := rootCmd.Execute(); err != nil {
		// Log error if command execution fails
		x_log.Error().
			Err(err).
			Msg("command execution failed")

		os.Exit(1)
	}

	// Log success if the command executes correctly
	x_log.Info().Msg("command executed successfully")
}

// ---------- Command Initialization ----------
func init() {
	cfg, err := x_config.LoadConfig()
	if err != nil {
		fmt.Println("Failed to load config:", err)
		os.Exit(1)
	}

	// Apply logger configuration from the config
	x_log.InitWithConfig(&cfg.Logger, "main")

	// Log the initialization of the root command
	x_log.Debug().Msg("initializing root command with --config flag")

	// Define the --config flag for the root command
	rootCmd.PersistentFlags().
		StringVarP(&pathFlag, "config", "p", "./tasks.json", "Path to the tasks file")

	// Log the default path of tasks file for debugging
	x_log.Debug().
		Str("default_path", "./tasks.json").
		Msg("initialized config file flag")
}
