package steps

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/well-prado/tetrix-cli/internal/config"
	"github.com/well-prado/tetrix-cli/internal/tui/styles"
)

// WelcomeModel is the first wizard step.
type WelcomeModel struct {
	version string
}

// NewWelcomeModel creates the welcome step.
func NewWelcomeModel(version string) WelcomeModel {
	return WelcomeModel{version: version}
}

func (m WelcomeModel) Init() tea.Cmd { return nil }

func (m WelcomeModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if keyMsg, ok := msg.(tea.KeyMsg); ok {
		switch keyMsg.String() {
		case "enter":
			return m, NextStep
		case "q", "ctrl+c":
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m WelcomeModel) View() string {
	s := styles.Logo() + "\n\n"
	s += styles.Subtitle.Render("AI-Powered Code Understanding Platform") + "\n"
	s += styles.Dimmed.Render(fmt.Sprintf("Community Edition %s", config.TetrixVersion)) + "\n\n"
	s += "This wizard will configure and install Tetrix CE on your machine.\n"
	s += "You'll need:\n\n"
	s += fmt.Sprintf("  %s GitHub OAuth App (Client ID + Secret)\n", styles.Bullet)
	s += fmt.Sprintf("  %s Anthropic API Key (for AI features)\n", styles.Bullet)
	s += fmt.Sprintf("  %s Voyage AI API Key (for code embeddings)\n", styles.Bullet)
	s += "\n" + styles.StatusBar.Render("[Enter] Continue  [q] Quit")
	return s
}
