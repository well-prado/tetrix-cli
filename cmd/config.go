package cmd

import (
	"fmt"
	"os"

	"github.com/well-prado/tetrix-cli/internal/config"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "View or edit Tetrix CE configuration",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.Load(tetrixHome)
		if err != nil {
			return fmt.Errorf("no installation found at %s — run 'tetrix install' first", tetrixHome)
		}

		// Redact secrets for display
		display := *cfg
		display.Auth.AuthSecret = redact(display.Auth.AuthSecret)
		display.Auth.GitHubClientSecret = redact(display.Auth.GitHubClientSecret)
		display.Auth.GoogleClientSecret = redact(display.Auth.GoogleClientSecret)
		display.Auth.GitLabClientSecret = redact(display.Auth.GitLabClientSecret)
		display.AI.OpenAIAPIKey = redact(display.AI.OpenAIAPIKey)
		display.AI.Context7APIKey = redact(display.AI.Context7APIKey)
		display.Infrastructure.PostgresPassword = redact(display.Infrastructure.PostgresPassword)
		display.Infrastructure.MeilisearchKey = redact(display.Infrastructure.MeilisearchKey)
		display.Infrastructure.MinioSecretKey = redact(display.Infrastructure.MinioSecretKey)
		display.Infrastructure.VaultToken = redact(display.Infrastructure.VaultToken)

		out, err := yaml.Marshal(display)
		if err != nil {
			return err
		}

		fmt.Printf("Tetrix CE Configuration (%s/config.yaml)\n\n", tetrixHome)
		os.Stdout.Write(out)
		return nil
	},
}

func redact(s string) string {
	if s == "" {
		return ""
	}
	if len(s) <= 8 {
		return "****"
	}
	return s[:4] + "..." + s[len(s)-4:]
}
