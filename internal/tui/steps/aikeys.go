package steps

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/well-prado/tetrix-cli/internal/config"
	"github.com/well-prado/tetrix-cli/internal/tui/styles"
)

// AIKeysModel collects AI provider API keys.
type AIKeysModel struct {
	inputs  []textinput.Model
	focused int
	cfg     *config.Config
	err     string
}

const (
	aiAnthropicKey = iota
	aiVoyageKey
	aiOpenAIKey
	aiGoogleKey
)

// NewAIKeysModel creates the AI keys step.
func NewAIKeysModel(cfg *config.Config) AIKeysModel {
	inputs := make([]textinput.Model, 4)

	inputs[aiAnthropicKey] = textinput.New()
	inputs[aiAnthropicKey].Placeholder = "sk-ant-api03-..."
	inputs[aiAnthropicKey].EchoMode = textinput.EchoPassword
	inputs[aiAnthropicKey].CharLimit = 120
	inputs[aiAnthropicKey].Width = 44
	inputs[aiAnthropicKey].Focus()
	inputs[aiAnthropicKey].SetValue(cfg.AI.AnthropicAPIKey)

	inputs[aiVoyageKey] = textinput.New()
	inputs[aiVoyageKey].Placeholder = "pa-..."
	inputs[aiVoyageKey].EchoMode = textinput.EchoPassword
	inputs[aiVoyageKey].CharLimit = 80
	inputs[aiVoyageKey].Width = 44
	inputs[aiVoyageKey].SetValue(cfg.AI.VoyageAPIKey)

	inputs[aiOpenAIKey] = textinput.New()
	inputs[aiOpenAIKey].Placeholder = "sk-... (optional)"
	inputs[aiOpenAIKey].EchoMode = textinput.EchoPassword
	inputs[aiOpenAIKey].CharLimit = 80
	inputs[aiOpenAIKey].Width = 44
	inputs[aiOpenAIKey].SetValue(cfg.AI.OpenAIAPIKey)

	inputs[aiGoogleKey] = textinput.New()
	inputs[aiGoogleKey].Placeholder = "AIza... (optional)"
	inputs[aiGoogleKey].EchoMode = textinput.EchoPassword
	inputs[aiGoogleKey].CharLimit = 80
	inputs[aiGoogleKey].Width = 44
	inputs[aiGoogleKey].SetValue(cfg.AI.GoogleAPIKey)

	return AIKeysModel{inputs: inputs, focused: 0, cfg: cfg}
}

func (m AIKeysModel) Init() tea.Cmd { return textinput.Blink }

func (m AIKeysModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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
			if m.focused >= aiOpenAIKey-1 {
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

func (m AIKeysModel) View() string {
	s := styles.Title.Render("AI API Keys") + "\n"
	s += styles.StepIndicator.Render("Step 3/7") + "\n\n"
	s += styles.Dimmed.Render("Tetrix uses AI for code understanding and semantic search.") + "\n\n"

	labels := []string{
		"Anthropic API Key:    ",
		"Voyage AI API Key:    ",
		"OpenAI API Key:       ",
		"Google API Key:       ",
	}
	required := []bool{true, true, false, false}

	for i, input := range m.inputs {
		label := styles.Label.Render(labels[i])
		cursor := " "
		if i == m.focused {
			cursor = styles.ActiveInput.Render(styles.Arrow)
		}
		tag := ""
		if required[i] {
			tag = styles.Error.Render(" *")
		} else {
			tag = styles.Dimmed.Render("  ")
		}
		s += fmt.Sprintf(" %s%s %s %s\n", cursor, tag, label, input.View())
	}

	s += "\n" + styles.Dimmed.Render("  * Required fields") + "\n"

	if m.err != "" {
		s += "\n" + styles.Error.Render("  "+m.err)
	}

	s += "\n" + styles.StatusBar.Render("[Tab] Next field  [Enter] Continue  [Esc] Back  [Ctrl+C] Quit")
	return s
}

func (m *AIKeysModel) updateFocus() tea.Cmd {
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

func (m *AIKeysModel) updateInputs(msg tea.Msg) tea.Cmd {
	cmds := make([]tea.Cmd, len(m.inputs))
	for i := range m.inputs {
		m.inputs[i], cmds[i] = m.inputs[i].Update(msg)
	}
	return tea.Batch(cmds...)
}

func (m *AIKeysModel) validate() string {
	anthropic := strings.TrimSpace(m.inputs[aiAnthropicKey].Value())
	voyage := strings.TrimSpace(m.inputs[aiVoyageKey].Value())
	if anthropic == "" {
		return "Anthropic API Key is required."
	}
	if voyage == "" {
		return "Voyage AI API Key is required for code embeddings."
	}
	return ""
}

func (m *AIKeysModel) save() {
	m.cfg.AI.AnthropicAPIKey = strings.TrimSpace(m.inputs[aiAnthropicKey].Value())
	m.cfg.AI.VoyageAPIKey = strings.TrimSpace(m.inputs[aiVoyageKey].Value())
	m.cfg.AI.OpenAIAPIKey = strings.TrimSpace(m.inputs[aiOpenAIKey].Value())
	m.cfg.AI.GoogleAPIKey = strings.TrimSpace(m.inputs[aiGoogleKey].Value())
}

// GetConfig returns the updated config.
func (m AIKeysModel) GetConfig() *config.Config { return m.cfg }
