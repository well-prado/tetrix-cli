package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

var (
	cliVersion = "dev"
	cliCommit  = "none"
	tetrixHome string
	verbose    bool
	noColor    bool
)

// SetVersionInfo is called from main to inject build-time values.
func SetVersionInfo(v, c string) {
	cliVersion = v
	cliCommit = c
}

// DefaultHome returns the default Tetrix installation directory.
func DefaultHome() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return filepath.Join(".", ".tetrix")
	}
	return filepath.Join(home, ".tetrix")
}

var rootCmd = &cobra.Command{
	Use:   "tetrix",
	Short: "Tetrix CE — AI-powered code understanding platform",
	Long: `Tetrix CLI installs and manages the Tetrix CE platform.

Run 'tetrix install' to get started with the interactive setup wizard.`,
}

// Execute runs the root command.
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.PersistentFlags().StringVar(&tetrixHome, "home", DefaultHome(), "Tetrix installation directory")
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "Enable verbose output")
	rootCmd.PersistentFlags().BoolVar(&noColor, "no-color", false, "Disable colored output")

	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(installCmd)
	rootCmd.AddCommand(startCmd)
	rootCmd.AddCommand(stopCmd)
	rootCmd.AddCommand(statusCmd)
	rootCmd.AddCommand(logsCmd)
	rootCmd.AddCommand(updateCmd)
	rootCmd.AddCommand(configCmd)
	rootCmd.AddCommand(uninstallCmd)
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print CLI version information",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("tetrix CLI %s (commit: %s)\n", cliVersion, cliCommit)
	},
}
