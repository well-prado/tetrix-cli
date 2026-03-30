package docker

import (
	"fmt"
	"net"
	"os/exec"
	"regexp"
	"runtime"
	"strconv"
	"strings"
)

// PrereqResult holds the result of a single prerequisite check.
type PrereqResult struct {
	Name    string
	Passed  bool
	Detail  string
	Warning string // Non-blocking issue
}

// CheckAll runs all prerequisite checks and returns results.
func CheckAll(ports map[string]int) []PrereqResult {
	var results []PrereqResult
	results = append(results, checkDockerEngine())
	results = append(results, checkDockerCompose())
	results = append(results, checkRAM())
	for name, port := range ports {
		results = append(results, checkPort(name, port))
	}
	return results
}

// AllPassed returns true if every check passed.
func AllPassed(results []PrereqResult) bool {
	for _, r := range results {
		if !r.Passed {
			return false
		}
	}
	return true
}

func checkDockerEngine() PrereqResult {
	out, err := exec.Command("docker", "version", "--format", "{{.Server.Version}}").CombinedOutput()
	if err != nil {
		return PrereqResult{
			Name:   "Docker Engine",
			Passed: false,
			Detail: "Docker is not installed or not running. Please install Docker Desktop: https://docker.com/get-started",
		}
	}
	version := strings.TrimSpace(string(out))
	major := parseMajorVersion(version)
	if major < 20 {
		return PrereqResult{
			Name:   "Docker Engine",
			Passed: false,
			Detail: fmt.Sprintf("Docker %s found, minimum 20.0 required", version),
		}
	}
	return PrereqResult{
		Name:   "Docker Engine",
		Passed: true,
		Detail: fmt.Sprintf("v%s", version),
	}
}

func checkDockerCompose() PrereqResult {
	out, err := exec.Command("docker", "compose", "version", "--short").CombinedOutput()
	if err != nil {
		return PrereqResult{
			Name:   "Docker Compose",
			Passed: false,
			Detail: "Docker Compose v2 not found. Upgrade Docker Desktop or install the compose plugin.",
		}
	}
	version := strings.TrimSpace(string(out))
	// Remove leading 'v' if present
	version = strings.TrimPrefix(version, "v")
	major := parseMajorVersion(version)
	if major < 2 {
		return PrereqResult{
			Name:   "Docker Compose",
			Passed: false,
			Detail: fmt.Sprintf("Compose %s found, minimum v2.0 required", version),
		}
	}
	return PrereqResult{
		Name:   "Docker Compose",
		Passed: true,
		Detail: fmt.Sprintf("v%s", version),
	}
}

func checkRAM() PrereqResult {
	totalGB := getTotalRAMGB()
	if totalGB < 4 {
		return PrereqResult{
			Name:   "Available RAM",
			Passed: false,
			Detail: fmt.Sprintf("%.1f GB detected, minimum 8 GB recommended", totalGB),
		}
	}
	result := PrereqResult{
		Name:   "Available RAM",
		Passed: true,
		Detail: fmt.Sprintf("%.1f GB", totalGB),
	}
	if totalGB < 8 {
		result.Warning = "8 GB recommended for optimal performance"
	}
	return result
}

func checkPort(name string, port int) PrereqResult {
	ln, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return PrereqResult{
			Name:   fmt.Sprintf("Port %d (%s)", port, name),
			Passed: false,
			Detail: fmt.Sprintf("Port %d is already in use", port),
		}
	}
	ln.Close()
	return PrereqResult{
		Name:   fmt.Sprintf("Port %d (%s)", port, name),
		Passed: true,
		Detail: "available",
	}
}

func parseMajorVersion(version string) int {
	re := regexp.MustCompile(`^(\d+)`)
	matches := re.FindStringSubmatch(version)
	if len(matches) < 2 {
		return 0
	}
	v, _ := strconv.Atoi(matches[1])
	return v
}

func getTotalRAMGB() float64 {
	switch runtime.GOOS {
	case "darwin":
		out, err := exec.Command("sysctl", "-n", "hw.memsize").CombinedOutput()
		if err != nil {
			return 0
		}
		bytes, err := strconv.ParseFloat(strings.TrimSpace(string(out)), 64)
		if err != nil {
			return 0
		}
		return bytes / (1024 * 1024 * 1024)
	case "linux":
		out, err := exec.Command("grep", "MemTotal", "/proc/meminfo").CombinedOutput()
		if err != nil {
			return 0
		}
		parts := strings.Fields(string(out))
		if len(parts) < 2 {
			return 0
		}
		kb, err := strconv.ParseFloat(parts[1], 64)
		if err != nil {
			return 0
		}
		return kb / (1024 * 1024)
	default:
		return 16 // Assume sufficient on unknown OS
	}
}
