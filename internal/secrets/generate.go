package secrets

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
)

// GenerateSecret creates a cryptographically secure random string.
// The result is base64url-encoded (no padding) for safe use in env vars and YAML.
func GenerateSecret(byteLength int) (string, error) {
	b := make([]byte, byteLength)
	if _, err := rand.Read(b); err != nil {
		return "", fmt.Errorf("failed to generate secret: %w", err)
	}
	return base64.RawURLEncoding.EncodeToString(b), nil
}

// MustGenerateSecret generates a secret or panics. Use only during setup.
func MustGenerateSecret(byteLength int) string {
	s, err := GenerateSecret(byteLength)
	if err != nil {
		panic(err)
	}
	return s
}

// GenerateAuthSecret creates a 32-byte AUTH_SECRET suitable for Auth.js.
func GenerateAuthSecret() (string, error) {
	return GenerateSecret(32)
}

// GeneratePassword creates a 24-byte password for database/infrastructure use.
func GeneratePassword() (string, error) {
	return GenerateSecret(24)
}
