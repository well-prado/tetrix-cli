package docker

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// ComposeManager handles Docker Compose operations for the Tetrix stack.
type ComposeManager struct {
	Home        string
	ProjectName string
}

// NewComposeManager creates a compose manager pointing to the Tetrix home directory.
func NewComposeManager(home string) *ComposeManager {
	return &ComposeManager{
		Home:        home,
		ProjectName: "tetrix",
	}
}

func (m *ComposeManager) composeFile() string {
	return filepath.Join(m.Home, "docker-compose.yml")
}

func (m *ComposeManager) baseArgs() []string {
	return []string{
		"compose",
		"-f", m.composeFile(),
		"-p", m.ProjectName,
	}
}

// Pull downloads all Docker images defined in the compose file.
func (m *ComposeManager) Pull() error {
	args := append(m.baseArgs(), "pull")
	return m.run(args...)
}

// Up starts all services in detached mode.
func (m *ComposeManager) Up() error {
	args := append(m.baseArgs(), "up", "-d", "--remove-orphans")
	return m.run(args...)
}

// Down stops and removes all containers.
func (m *ComposeManager) Down() error {
	args := append(m.baseArgs(), "down")
	return m.run(args...)
}

// DownWithVolumes stops containers and removes volumes.
func (m *ComposeManager) DownWithVolumes() error {
	args := append(m.baseArgs(), "down", "-v")
	return m.run(args...)
}

// PS returns the output of docker compose ps.
func (m *ComposeManager) PS() (string, error) {
	args := append(m.baseArgs(), "ps", "--format", "table {{.Name}}\t{{.Status}}\t{{.Ports}}")
	return m.output(args...)
}

// Logs streams logs for a service (or all services if service is empty).
func (m *ComposeManager) Logs(service string, tail int, follow bool) *exec.Cmd {
	args := m.baseArgs()
	args = append(args, "logs")
	if follow {
		args = append(args, "-f")
	}
	if tail > 0 {
		args = append(args, "--tail", fmt.Sprintf("%d", tail))
	}
	if service != "" {
		args = append(args, service)
	}
	cmd := exec.Command("docker", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd
}

// PullProgress runs docker compose pull and pipes output to stdout.
func (m *ComposeManager) PullProgress() *exec.Cmd {
	args := append(m.baseArgs(), "pull")
	cmd := exec.Command("docker", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd
}

// IsRunning checks if at least one Tetrix container is running.
func (m *ComposeManager) IsRunning() bool {
	args := append(m.baseArgs(), "ps", "-q")
	out, err := m.output(args...)
	if err != nil {
		return false
	}
	return strings.TrimSpace(out) != ""
}

// ServiceHealth returns a map of service name → health status string.
func (m *ComposeManager) ServiceHealth() (map[string]string, error) {
	args := append(m.baseArgs(), "ps", "--format", "{{.Name}}|{{.Status}}")
	out, err := m.output(args...)
	if err != nil {
		return nil, err
	}
	health := make(map[string]string)
	for _, line := range strings.Split(strings.TrimSpace(out), "\n") {
		if line == "" {
			continue
		}
		parts := strings.SplitN(line, "|", 2)
		if len(parts) == 2 {
			health[parts[0]] = parts[1]
		}
	}
	return health, nil
}

func (m *ComposeManager) run(args ...string) error {
	cmd := exec.Command("docker", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func (m *ComposeManager) output(args ...string) (string, error) {
	out, err := exec.Command("docker", args...).CombinedOutput()
	return string(out), err
}
