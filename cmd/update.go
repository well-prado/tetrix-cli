package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/well-prado/tetrix-cli/internal/config"
	"github.com/well-prado/tetrix-cli/internal/docker"
	"github.com/well-prado/tetrix-cli/internal/templates"
	"github.com/spf13/cobra"
)

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Pull latest Docker images and restart services",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.Load(tetrixHome)
		if err != nil {
			return fmt.Errorf("no installation found at %s — run 'tetrix install' first", tetrixHome)
		}

		fmt.Println("This will pull the latest Tetrix CE images and restart all services.")
		fmt.Print("Continue? [y/N] ")
		reader := bufio.NewReader(os.Stdin)
		answer, _ := reader.ReadString('\n')
		if !strings.HasPrefix(strings.ToLower(strings.TrimSpace(answer)), "y") {
			fmt.Println("Cancelled.")
			return nil
		}

		cm := docker.NewComposeManager(tetrixHome)

		// Re-render compose file (picks up any template changes from CLI update)
		data := templates.NewTemplateData(cfg, cliVersion)
		if err := templates.WriteAllFiles(data); err != nil {
			return fmt.Errorf("failed to update compose files: %w", err)
		}

		fmt.Println("Pulling latest images...")
		if err := cm.Pull(); err != nil {
			return fmt.Errorf("failed to pull images: %w", err)
		}

		fmt.Println("Restarting services...")
		if err := cm.Down(); err != nil {
			return fmt.Errorf("failed to stop services: %w", err)
		}
		if err := cm.Up(); err != nil {
			return fmt.Errorf("failed to start services: %w", err)
		}

		fmt.Println("Update complete! Tetrix CE is running.")
		return nil
	},
}
