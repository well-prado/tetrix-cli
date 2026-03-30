package steps

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/well-prado/tetrix-cli/internal/config"
	"github.com/well-prado/tetrix-cli/internal/tui/styles"
)

// OAuthModel collects GitHub (and optionally Google/GitLab) OAuth credentials.
type OAuthModel struct {
	inputs  []textinput.Model
	focused int
	cfg     *config.Config
	err     string
}

const (
	oauthGHClientID = iota
	oauthGHClientSecret
)

// NewOAuthModel creates the OAuth configuration step.
func NewOAuthModel(cfg *config.Config) OAuthModel {
	inputs := make([]textinput.Model, 2)

	inputs[oauthGHClientID] = textinput.New()
	inputs[oauthGHClientID].Placeholder = "Ov23li..."
	inputs[oauthGHClientID].Focus()
	inputs[oauthGHClientID].CharLimit = 40
	inputs[oauthGHClientID].Width = 40
	inputs[oauthGHClientID].SetValue(cfg.Auth.GitHubClientID)

	inputs[oauthGHClientSecret] = textinput.New()
	inputs[oauthGHClientSecret].Placeholder = "github_secret_..."
	inputs[oauthGHClientSecret].EchoMode = textinput.EchoPassword
	inputs[oauthGHClientSecret].CharLimit = 80
	inputs[oauthGHClientSecret].Width = 40
	inputs[oauthGHClientSecret].SetValue(cfg.Auth.GitHubClientSecret)

	return OAuthModel{
		inputs:  inputs,
		focused: 0,
		cfg:     cfg,
	}
}

func (m OAuthModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m OAuthModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "tab", "down":
			m.focused = (m.focused + 1) % len(m.inputs)
			return m, m.updateFocus()
		case "shift+tab", "up":
			m.focused = (m.focused - 1 + len(m.inputs)) % len(m.inputs)
			return m, m.updateFocus()
		case "enter":
			if m.focused == len(m.inputs)-1 {
				if err := m.validate(); err != "" {
					m.err = err
					return m, nil
				}
				m.save()
				return m, NextStep
			}
			m.focused++
			return m, m.updateFocus()
		case "ctrl+c":
			return m, tea.Quit
		case "esc":
			return m, PrevStep
		}
	}

	cmd := m.updateInputs(msg)
	return m, cmd
}

func (m OAuthModel) View() string {
	s := styles.Title.Render("GitHub OAuth Setup") + "\n"
	s += styles.StepIndicator.Render("Step 2/7") + "\n\n"
	s += styles.Dimmed.Render("Create a GitHub OAuth App for user sign-in:") + "\n"
	s += styles.Dimmed.Render("  Settings → Developer Settings → OAuth Apps → New") + "\n"
	s += styles.Dimmed.Render(fmt.Sprintf("  Callback URL: %s/api/auth/callback/github", m.cfg.AppURL)) + "\n\n"

	labels := []string{"GitHub Client ID:", "GitHub Client Secret:"}
	for i, input := range m.inputs {
		label := styles.Label.Render(labels[i])
		cursor := " "
		if i == m.focused {
			cursor = styles.ActiveInput.Render(styles.Arrow)
		}
		s += fmt.Sprintf(" %s %s %s\n", cursor, label, input.View())
	}

	if m.err != "" {
		s += "\n" + styles.Error.Render("  "+m.err)
	}

	s += "\n" + styles.StatusBar.Render("[Tab] Next field  [Enter] Continue  [Esc] Back  [Ctrl+C] Quit")
	return s
}

func (m *OAuthModel) updateFocus() tea.Cmd {
	cmds := make([]tea.Cmd, len(m.inputs))
	for i := range m.inputs {
		if i == m.focused {
			cmds[i] = m.inputs[i].Focus()
		} else {
			m.inputs[i].Blur()
		}
	}
	return tea.Batch(cmds...)
}

func (m *OAuthModel) updateInputs(msg tea.Msg) tea.Cmd {
	cmds := make([]tea.Cmd, len(m.inputs))
	for i := range m.inputs {
		m.inputs[i], cmds[i] = m.inputs[i].Update(msg)
	}
	return tea.Batch(cmds...)
}

func (m *OAuthModel) validate() string {
	clientID := strings.TrimSpace(m.inputs[oauthGHClientID].Value())
	clientSecret := strings.TrimSpace(m.inputs[oauthGHClientSecret].Value())
	if clientID == "" || clientSecret == "" {
		return "GitHub Client ID and Secret are required."
	}
	return ""
}

func (m *OAuthModel) save() {
	m.cfg.Auth.GitHubClientID = strings.TrimSpace(m.inputs[oauthGHClientID].Value())
	m.cfg.Auth.GitHubClientSecret = strings.TrimSpace(m.inputs[oauthGHClientSecret].Value())
}

// GetConfig returns the updated config.
func (m OAuthModel) GetConfig() *config.Config { return m.cfg }
