package config

import (
	"os"
	"path/filepath"
	"time"

	"gopkg.in/yaml.v3"
)

// Config holds all Tetrix CE installation configuration.
type Config struct {
	Version    string    `yaml:"version"`
	CLIVersion string    `yaml:"cli_version"`
	InstallAt  time.Time `yaml:"installed_at"`

	Auth           AuthConfig           `yaml:"auth"`
	AI             AIConfig             `yaml:"ai"`
	Ports          PortsConfig          `yaml:"ports"`
	Infrastructure InfrastructureConfig `yaml:"infrastructure"`
	Worker         WorkerConfig         `yaml:"worker"`

	AppURL     string `yaml:"app_url"`
	TetrixHome string `yaml:"tetrix_home"`
}

// AuthConfig holds authentication provider settings.
type AuthConfig struct {
	AuthSecret         string `yaml:"auth_secret"`
	GitHubClientID     string `yaml:"github_client_id"`
	GitHubClientSecret string `yaml:"github_client_secret"`
	GoogleClientID     string `yaml:"google_client_id"`
	GoogleClientSecret string `yaml:"google_client_secret"`
	GoogleOAuthDomain  string `yaml:"google_oauth_domain"`
	GitLabClientID     string `yaml:"gitlab_client_id"`
	GitLabClientSecret string `yaml:"gitlab_client_secret"`
	GitLabBaseURL      string `yaml:"gitlab_base_url"`
}

// AIConfig holds AI provider API keys.
type AIConfig struct {
	OpenAIAPIKey  string `yaml:"openai_api_key"`
	Context7APIKey string `yaml:"context7_api_key"`
}

// PortsConfig holds port assignments for each service.
type PortsConfig struct {
	Web          int `yaml:"web"`
	API          int `yaml:"api"`
	MultiAgent   int `yaml:"multi_agent"`
	OpenAIProxy  int `yaml:"openai_proxy"`
	MinioAPI     int `yaml:"minio_api"`
	MinioConsole int `yaml:"minio_console"`
	Vault        int `yaml:"vault"`
}

// InfrastructureConfig holds auto-generated infrastructure passwords.
type InfrastructureConfig struct {
	PostgresPassword string `yaml:"postgres_password"`
	MeilisearchKey   string `yaml:"meilisearch_key"`
	MinioSecretKey   string `yaml:"minio_secret_key"`
	VaultToken       string `yaml:"vault_token"`
}

// WorkerConfig holds worker process settings.
type WorkerConfig struct {
	Concurrency int    `yaml:"concurrency"`
	LogLevel    string `yaml:"log_level"`
}

// NewDefault creates a Config with sensible defaults.
func NewDefault(home string) *Config {
	return &Config{
		Version:    TetrixVersion,
		CLIVersion: "",
		InstallAt:  time.Now(),
		Auth: AuthConfig{
			GitLabBaseURL: "https://gitlab.com",
		},
		Ports: PortsConfig{
			Web:          3000,
			API:          4000,
			MultiAgent:   7777,
			OpenAIProxy:  8085,
			MinioAPI:     9000,
			MinioConsole: 9001,
			Vault:        8200,
		},
		Worker: WorkerConfig{
			Concurrency: 2,
			LogLevel:    "info",
		},
		AppURL:     "http://localhost:3000",
		TetrixHome: home,
	}
}

// Save writes the config to ~/.tetrix/config.yaml.
func (c *Config) Save() error {
	if err := os.MkdirAll(c.TetrixHome, 0700); err != nil {
		return err
	}
	data, err := yaml.Marshal(c)
	if err != nil {
		return err
	}
	return os.WriteFile(filepath.Join(c.TetrixHome, "config.yaml"), data, 0600)
}

// Load reads config from ~/.tetrix/config.yaml.
func Load(home string) (*Config, error) {
	data, err := os.ReadFile(filepath.Join(home, "config.yaml"))
	if err != nil {
		return nil, err
	}
	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}

// Exists checks if a config file exists at the given home directory.
func Exists(home string) bool {
	_, err := os.Stat(filepath.Join(home, "config.yaml"))
	return err == nil
}
