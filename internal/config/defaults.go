package config

// TetrixVersion is the platform version this CLI release targets.
// Updated at build time or manually for each CLI release.
const TetrixVersion = "latest"

// DockerHubOrg is the DockerHub organization/user for published images.
const DockerHubOrg = "mcastillo507"

// ImageNames maps service names to their DockerHub image names.
var ImageNames = map[string]string{
	"api":                "tetrix-api",
	"web":                "tetrix-web",
	"worker":             "tetrix-worker",
	"multi-agent":        "tetrix-multiagent",
	"openai-proxy":       "tetrix-openaiproxy",
	"credential-service": "tetrix-credential-service",
}

// InfraImages maps infrastructure service names to their pinned image references.
var InfraImages = map[string]string{
	"postgres":    "ankane/pgvector:latest",
	"neo4j":       "neo4j:5.15.0",
	"redis":       "redis:7-alpine",
	"meilisearch": "getmeili/meilisearch:v1.5",
	"minio":       "minio/minio:latest",
	"minio-init":  "minio/mc:latest",
	"vault":       "hashicorp/vault:1.15",
	"mongodb":     "mongo:7",
}

// DefaultPorts returns a map of service name → default port.
func DefaultPorts() map[string]int {
	return map[string]int{
		"web":           3000,
		"api":           4000,
		"multi-agent":   7777,
		"openai-proxy":  8085,
		"minio-api":     9000,
		"minio-console": 9001,
		"vault":         8200,
	}
}
