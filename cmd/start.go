package cmd

import (
	"fmt"

	"github.com/well-prado/tetrix-cli/internal/docker"
	"github.com/spf13/cobra"
)

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start all Tetrix CE services",
	RunE: func(cmd *cobra.Command, args []string) error {
		cm := docker.NewComposeManager(tetrixHome)
		fmt.Println("Starting Tetrix CE...")
		if err := cm.Up(); err != nil {
			return fmt.Errorf("failed to start: %w", err)
		}
		fmt.Println("Tetrix CE is running! Open http://localhost:3000")
		return nil
	},
}
