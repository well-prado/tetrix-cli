package steps

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/well-prado/tetrix-cli/internal/config"
	"github.com/well-prado/tetrix-cli/internal/tui/styles"
)

// ReviewModel shows a summary of all settings before installation.
type ReviewModel struct {
	cfg *config.Config
}

// NewReviewModel creates the review step.
func NewReviewModel(cfg *config.Config) ReviewModel {
	return ReviewModel{cfg: cfg}
}

func (m ReviewModel) Init() tea.Cmd { return nil }

func (m ReviewModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if keyMsg, ok := msg.(tea.KeyMsg); ok {
		switch keyMsg.String() {
		case "enter", "y":
			return m, NextStep
		case "esc", "b":
			return m, PrevStep
		case "q", "ctrl+c":
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m ReviewModel) View() string {
	s := styles.Title.Render("Review Configuration") + "\n"
	s += styles.StepIndicator.Render("Step 5/7") + "\n\n"

	// Authentication
	s += styles.Subtitle.Render("  Authentication") + "\n"
	s += configLine("GitHub OAuth", checkOrDash(m.cfg.Auth.GitHubClientID))
	s += configLine("Google OAuth", checkOrDash(m.cfg.Auth.GoogleClientID))
	s += configLine("GitLab OAuth", checkOrDash(m.cfg.Auth.GitLabClientID))
	s += configLine("AUTH_SECRET", styles.Success.Render("Auto-generated"))
	s += "\n"

	// AI Providers
	s += styles.Subtitle.Render("  AI Providers") + "\n"
	s += configLine("OpenAI", maskKey(m.cfg.AI.OpenAIAPIKey))
	s += configLine("Context7", maskKey(m.cfg.AI.Context7APIKey))
	s += "\n"

	// Ports
	s += styles.Subtitle.Render("  Ports") + "\n"
	s += configLine("Web UI", fmt.Sprintf(":%d", m.cfg.Ports.Web))
	s += configLine("API", fmt.Sprintf(":%d", m.cfg.Ports.API))
	s += configLine("Multi-Agent", fmt.Sprintf(":%d", m.cfg.Ports.MultiAgent))
	s += configLine("OpenAI Proxy", fmt.Sprintf(":%d", m.cfg.Ports.OpenAIProxy))
	s += "\n"

	// Infrastructure
	s += styles.Subtitle.Render("  Infrastructure") + "\n"
	s += configLine("Passwords", styles.Success.Render("All auto-generated"))
	s += configLine("Install path", m.cfg.TetrixHome)
	s += "\n"

	s += styles.StatusBar.Render("[Enter] Install  [Esc] Back  [q] Quit")
	return s
}

func configLine(label, value string) string {
	return fmt.Sprintf("    %s %s\n", styles.Label.Render(label+":"), value)
}

func checkOrDash(val string) string {
	if val != "" {
		return styles.Success.Render(styles.CheckMark + " Configured")
	}
	return styles.Dimmed.Render("- Not configured")
}

func maskKey(key string) string {
	if key == "" {
		return styles.Error.Render("Missing!")
	}
	if len(key) > 8 {
		return styles.Success.Render(styles.CheckMark + " " + key[:4] + "..." + key[len(key)-4:])
	}
	return styles.Success.Render(styles.CheckMark + " Set")
}

func maskKeyOptional(key string) string {
	if key == "" {
		return styles.Dimmed.Render("- Not configured")
	}
	return maskKey(key)
}
