package cmd

import (
	"fmt"

	"github.com/well-prado/tetrix-cli/internal/docker"
	"github.com/spf13/cobra"
)

var stopCmd = &cobra.Command{
	Use:   "stop",
	Short: "Stop all Tetrix CE services",
	RunE: func(cmd *cobra.Command, args []string) error {
		cm := docker.NewComposeManager(tetrixHome)
		fmt.Println("Stopping Tetrix CE...")
		if err := cm.Down(); err != nil {
			return fmt.Errorf("failed to stop: %w", err)
		}
		fmt.Println("All services stopped.")
		return nil
	},
}
