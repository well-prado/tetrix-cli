package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/well-prado/tetrix-cli/internal/docker"
	"github.com/spf13/cobra"
)

var keepData bool

var uninstallCmd = &cobra.Command{
	Use:   "uninstall",
	Short: "Remove Tetrix CE installation",
	Long:  "Stops all services, removes containers, and optionally removes all data volumes.",
	RunE: func(cmd *cobra.Command, args []string) error {
		cm := docker.NewComposeManager(tetrixHome)

		if keepData {
			fmt.Println("This will stop and remove all Tetrix CE containers but KEEP your data volumes.")
		} else {
			fmt.Println("WARNING: This will stop all services, remove containers, AND DELETE ALL DATA.")
			fmt.Println("Your databases, indexes, and uploaded files will be permanently lost.")
		}
		fmt.Print("\nAre you sure? Type 'yes' to confirm: ")

		reader := bufio.NewReader(os.Stdin)
		answer, _ := reader.ReadString('\n')
		if strings.TrimSpace(answer) != "yes" {
			fmt.Println("Cancelled.")
			return nil
		}

		fmt.Println("Stopping services...")
		if keepData {
			if err := cm.Down(); err != nil {
				fmt.Printf("Warning: failed to stop services: %v\n", err)
			}
		} else {
			if err := cm.DownWithVolumes(); err != nil {
				fmt.Printf("Warning: failed to stop services: %v\n", err)
			}
		}

		// Remove configuration files
		fmt.Printf("Removing %s...\n", tetrixHome)
		if err := os.RemoveAll(tetrixHome); err != nil {
			return fmt.Errorf("failed to remove %s: %w", tetrixHome, err)
		}

		fmt.Println("Tetrix CE has been uninstalled.")
		if keepData {
			fmt.Println("Data volumes were preserved. Remove them manually with: docker volume ls | grep tetrix")
		}
		return nil
	},
}

func init() {
	uninstallCmd.Flags().BoolVar(&keepData, "keep-data", false, "Remove containers but preserve data volumes")
}
