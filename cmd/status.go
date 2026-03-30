package cmd

import (
	"fmt"

	"github.com/well-prado/tetrix-cli/internal/docker"
	"github.com/spf13/cobra"
)

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show the health status of all Tetrix CE services",
	RunE: func(cmd *cobra.Command, args []string) error {
		cm := docker.NewComposeManager(tetrixHome)
		if !cm.IsRunning() {
			fmt.Println("Tetrix CE is not running. Use 'tetrix start' to start it.")
			return nil
		}
		output, err := cm.PS()
		if err != nil {
			return fmt.Errorf("failed to get status: %w", err)
		}
		fmt.Println("Tetrix CE Service Status:")
		fmt.Println()
		fmt.Print(output)
		return nil
	},
}
