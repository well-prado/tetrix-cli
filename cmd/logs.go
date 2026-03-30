package cmd

import (
	"fmt"

	"github.com/well-prado/tetrix-cli/internal/docker"
	"github.com/spf13/cobra"
)

var logsTail int
var logsFollow bool

var logsCmd = &cobra.Command{
	Use:   "logs [service]",
	Short: "View logs from Tetrix CE services",
	Long:  "Stream logs from all or a specific service. Valid services: api, web, worker, multi-agent, openai-proxy, credential-service, postgres, neo4j, redis, meilisearch, minio, vault, mongodb",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		service := ""
		if len(args) > 0 {
			service = args[0]
		}
		cm := docker.NewComposeManager(tetrixHome)
		logCmd := cm.Logs(service, logsTail, logsFollow)
		if err := logCmd.Run(); err != nil {
			return fmt.Errorf("failed to stream logs: %w", err)
		}
		return nil
	},
}

func init() {
	logsCmd.Flags().IntVar(&logsTail, "tail", 100, "Number of lines to show from the end")
	logsCmd.Flags().BoolVarP(&logsFollow, "follow", "f", true, "Follow log output")
}
