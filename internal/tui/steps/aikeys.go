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
	aiOpenAIKey = iota
	aiContext7Key
)

// NewAIKeysModel creates the AI keys step.
func NewAIKeysModel(cfg *config.Config) AIKeysModel {
	inputs := make([]textinput.Model, 2)

	inputs[aiOpenAIKey] = textinput.New()
	inputs[aiOpenAIKey].Placeholder = "sk-..."
	inputs[aiOpenAIKey].EchoMode = textinput.EchoPassword
	inputs[aiOpenAIKey].CharLimit = 120
	inputs[aiOpenAIKey].Width = 44
	inputs[aiOpenAIKey].Focus()
	inputs[aiOpenAIKey].SetValue(cfg.AI.OpenAIAPIKey)

	inputs[aiContext7Key] = textinput.New()
	inputs[aiContext7Key].Placeholder = "ctx7-..."
	inputs[aiContext7Key].EchoMode = textinput.EchoPassword
	inputs[aiContext7Key].CharLimit = 80
	inputs[aiContext7Key].Width = 44
	inputs[aiContext7Key].SetValue(cfg.AI.Context7APIKey)

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

func (m AIKeysModel) View() string {
	s := styles.Title.Render("AI API Keys") + "\n"
	s += styles.StepIndicator.Render("Step 3/7") + "\n\n"
	s += styles.Dimmed.Render("Tetrix uses OpenAI for AI features and Context7 for documentation.") + "\n\n"

	labels := []string{
		"OpenAI API Key:       ",
		"Context7 API Key:     ",
	}
	required := []bool{true, true}

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
	openai := strings.TrimSpace(m.inputs[aiOpenAIKey].Value())
	context7 := strings.TrimSpace(m.inputs[aiContext7Key].Value())
	if openai == "" {
		return "OpenAI API Key is required."
	}
	if context7 == "" {
		return "Context7 API Key is required."
	}
	return ""
}

func (m *AIKeysModel) save() {
	m.cfg.AI.OpenAIAPIKey = strings.TrimSpace(m.inputs[aiOpenAIKey].Value())
	m.cfg.AI.Context7APIKey = strings.TrimSpace(m.inputs[aiContext7Key].Value())
}

// GetConfig returns the updated config.
func (m AIKeysModel) GetConfig() *config.Config { return m.cfg }
