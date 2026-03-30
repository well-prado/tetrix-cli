package templates

import (
	"bytes"
	"embed"
	"fmt"
	"os"
	"path/filepath"
	"text/template"

	"github.com/well-prado/tetrix-cli/internal/config"
)

//go:embed all:files
var embeddedFS embed.FS

// TemplateData holds all values needed to render templates.
type TemplateData struct {
	// Version info
	CLIVersion    string
	TetrixVersion string

	// Auth
	AuthSecret         string
	GitHubClientID     string
	GitHubClientSecret string
	GoogleClientID     string
	GoogleClientSecret string
	GoogleOAuthDomain  string
	GitLabClientID     string
	GitLabClientSecret string
	GitLabBaseURL      string

	// AI
	AnthropicAPIKey string
	VoyageAPIKey    string
	OpenAIAPIKey    string
	GoogleAPIKey    string

	// Ports
	WebPort          int
	APIPort          int
	MultiAgentPort   int
	OpenAIProxyPort  int
	MinioAPIPort     int
	MinioConsolePort int
	VaultPort        int

	// Infrastructure
	PostgresPassword string
	MeilisearchKey   string
	MinioSecretKey   string
	VaultToken       string

	// Worker
	WorkerConcurrency int
	LogLevel          string

	// Paths
	AppURL     string
	TetrixHome string

	// Docker
	DockerHubOrg string
}

// NewTemplateData converts a Config into template rendering data.
func NewTemplateData(cfg *config.Config, cliVersion string) *TemplateData {
	return &TemplateData{
		CLIVersion:         cliVersion,
		TetrixVersion:      cfg.Version,
		AuthSecret:         cfg.Auth.AuthSecret,
		GitHubClientID:     cfg.Auth.GitHubClientID,
		GitHubClientSecret: cfg.Auth.GitHubClientSecret,
		GoogleClientID:     cfg.Auth.GoogleClientID,
		GoogleClientSecret: cfg.Auth.GoogleClientSecret,
		GoogleOAuthDomain:  cfg.Auth.GoogleOAuthDomain,
		GitLabClientID:     cfg.Auth.GitLabClientID,
		GitLabClientSecret: cfg.Auth.GitLabClientSecret,
		GitLabBaseURL:      cfg.Auth.GitLabBaseURL,
		AnthropicAPIKey:    cfg.AI.AnthropicAPIKey,
		VoyageAPIKey:       cfg.AI.VoyageAPIKey,
		OpenAIAPIKey:       cfg.AI.OpenAIAPIKey,
		GoogleAPIKey:       cfg.AI.GoogleAPIKey,
		WebPort:            cfg.Ports.Web,
		APIPort:            cfg.Ports.API,
		MultiAgentPort:     cfg.Ports.MultiAgent,
		OpenAIProxyPort:    cfg.Ports.OpenAIProxy,
		MinioAPIPort:       cfg.Ports.MinioAPI,
		MinioConsolePort:   cfg.Ports.MinioConsole,
		VaultPort:          cfg.Ports.Vault,
		PostgresPassword:   cfg.Infrastructure.PostgresPassword,
		MeilisearchKey:     cfg.Infrastructure.MeilisearchKey,
		MinioSecretKey:     cfg.Infrastructure.MinioSecretKey,
		VaultToken:         cfg.Infrastructure.VaultToken,
		WorkerConcurrency:  cfg.Worker.Concurrency,
		LogLevel:           cfg.Worker.LogLevel,
		AppURL:             cfg.AppURL,
		TetrixHome:         cfg.TetrixHome,
		DockerHubOrg:       config.DockerHubOrg,
	}
}

// RenderDockerCompose renders the docker-compose.yml template with the given data.
func RenderDockerCompose(data *TemplateData) (string, error) {
	return renderTemplate("files/docker-compose.yml.tmpl", data)
}

// RenderFile writes a rendered template to the target path.
func RenderFile(data *TemplateData, templateName, targetPath string) error {
	content, err := renderTemplate("files/"+templateName, data)
	if err != nil {
		return err
	}
	return os.WriteFile(targetPath, []byte(content), 0644)
}

// CopyStaticFile copies a static (non-template) file from embedded FS to target.
func CopyStaticFile(name, targetPath string) error {
	data, err := embeddedFS.ReadFile("files/" + name)
	if err != nil {
		return fmt.Errorf("embedded file %s not found: %w", name, err)
	}
	return os.WriteFile(targetPath, data, 0644)
}

// WriteAllFiles renders and writes all config files to the Tetrix home directory.
func WriteAllFiles(data *TemplateData) error {
	home := data.TetrixHome
	if err := os.MkdirAll(home, 0700); err != nil {
		return fmt.Errorf("failed to create directory %s: %w", home, err)
	}

	// Render docker-compose.yml
	compose, err := RenderDockerCompose(data)
	if err != nil {
		return fmt.Errorf("failed to render docker-compose.yml: %w", err)
	}
	if err := os.WriteFile(filepath.Join(home, "docker-compose.yml"), []byte(compose), 0644); err != nil {
		return fmt.Errorf("failed to write docker-compose.yml: %w", err)
	}

	// Copy init-pgvector.sql (static, no templating needed)
	if err := CopyStaticFile("init-pgvector.sql", filepath.Join(home, "init-pgvector.sql")); err != nil {
		return fmt.Errorf("failed to write init-pgvector.sql: %w", err)
	}

	return nil
}

func renderTemplate(name string, data *TemplateData) (string, error) {
	content, err := embeddedFS.ReadFile(name)
	if err != nil {
		return "", fmt.Errorf("template %s not found: %w", name, err)
	}
	tmpl, err := template.New(name).Parse(string(content))
	if err != nil {
		return "", fmt.Errorf("failed to parse template %s: %w", name, err)
	}
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", fmt.Errorf("failed to render template %s: %w", name, err)
	}
	return buf.String(), nil
}
