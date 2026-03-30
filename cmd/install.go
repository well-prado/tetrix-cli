package cmd

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/well-prado/tetrix-cli/internal/config"
	"github.com/well-prado/tetrix-cli/internal/tui"
	"github.com/spf13/cobra"
)

var installCmd = &cobra.Command{
	Use:   "install",
	Short: "Run the interactive setup wizard to install Tetrix CE",
	Long:  "Walks you through configuring and installing all Tetrix CE services via Docker.",
	RunE: func(cmd *cobra.Command, args []string) error {
		if config.Exists(tetrixHome) {
			fmt.Println("Tetrix CE is already installed at", tetrixHome)
			fmt.Println("Run 'tetrix config' to modify settings or 'tetrix uninstall' first.")
			os.Exit(1)
		}

		cfg := config.NewDefault(tetrixHome)
		cfg.CLIVersion = cliVersion

		wizard := tui.NewWizard(cfg, cliVersion)
		p := tea.NewProgram(wizard, tea.WithAltScreen())
		if _, err := p.Run(); err != nil {
			return fmt.Errorf("wizard error: %w", err)
		}
		return nil
	},
}
