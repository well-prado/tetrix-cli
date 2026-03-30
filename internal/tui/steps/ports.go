package steps

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/well-prado/tetrix-cli/internal/config"
	"github.com/well-prado/tetrix-cli/internal/tui/styles"
)

// PortsModel lets users configure service ports.
type PortsModel struct {
	inputs  []textinput.Model
	labels  []string
	focused int
	cfg     *config.Config
	err     string
}

// NewPortsModel creates the port configuration step.
func NewPortsModel(cfg *config.Config) PortsModel {
	labels := []string{
		"Web UI:",
		"API Server:",
		"Multi-Agent:",
		"OpenAI Proxy:",
		"MinIO API:",
		"MinIO Console:",
		"Vault:",
	}
	defaults := []int{
		cfg.Ports.Web,
		cfg.Ports.API,
		cfg.Ports.MultiAgent,
		cfg.Ports.OpenAIProxy,
		cfg.Ports.MinioAPI,
		cfg.Ports.MinioConsole,
		cfg.Ports.Vault,
	}

	inputs := make([]textinput.Model, len(labels))
	for i := range inputs {
		inputs[i] = textinput.New()
		inputs[i].SetValue(strconv.Itoa(defaults[i]))
		inputs[i].CharLimit = 5
		inputs[i].Width = 8
	}
	inputs[0].Focus()

	return PortsModel{inputs: inputs, labels: labels, focused: 0, cfg: cfg}
}

func (m PortsModel) Init() tea.Cmd { return textinput.Blink }

func (m PortsModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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
			if err := m.validate(); err != "" {
				m.err = err
				return m, nil
			}
			m.save()
			return m, NextStep
		case "ctrl+c":
			return m, tea.Quit
		case "esc":
			return m, PrevStep
		}
	}
	cmd := m.updateInputs(msg)
	return m, cmd
}

func (m PortsModel) View() string {
	s := styles.Title.Render("Port Configuration") + "\n"
	s += styles.StepIndicator.Render("Step 4/7") + "\n\n"
	s += styles.Dimmed.Render("Press Enter to accept defaults, or type new values.") + "\n\n"

	for i, input := range m.inputs {
		label := styles.Label.Render(m.labels[i])
		cursor := " "
		if i == m.focused {
			cursor = styles.ActiveInput.Render(styles.Arrow)
		}
		s += fmt.Sprintf("  %s %s %s\n", cursor, label, input.View())
	}

	if m.err != "" {
		s += "\n" + styles.Error.Render("  "+m.err)
	}

	s += "\n" + styles.StatusBar.Render("[Tab] Next field  [Enter] Continue  [Esc] Back  [Ctrl+C] Quit")
	return s
}

func (m *PortsModel) updateFocus() tea.Cmd {
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

func (m *PortsModel) updateInputs(msg tea.Msg) tea.Cmd {
	cmds := make([]tea.Cmd, len(m.inputs))
	for i := range m.inputs {
		m.inputs[i], cmds[i] = m.inputs[i].Update(msg)
	}
	return tea.Batch(cmds...)
}

func (m *PortsModel) validate() string {
	seen := make(map[int]string)
	for i, input := range m.inputs {
		val := strings.TrimSpace(input.Value())
		port, err := strconv.Atoi(val)
		if err != nil || port < 1 || port > 65535 {
			return fmt.Sprintf("%s must be a valid port (1-65535)", m.labels[i])
		}
		if prev, ok := seen[port]; ok {
			return fmt.Sprintf("Port %d is used by both %s and %s", port, prev, m.labels[i])
		}
		seen[port] = m.labels[i]
	}
	return ""
}

func (m *PortsModel) save() {
	ports := make([]int, len(m.inputs))
	for i, input := range m.inputs {
		ports[i], _ = strconv.Atoi(strings.TrimSpace(input.Value()))
	}
	m.cfg.Ports.Web = ports[0]
	m.cfg.Ports.API = ports[1]
	m.cfg.Ports.MultiAgent = ports[2]
	m.cfg.Ports.OpenAIProxy = ports[3]
	m.cfg.Ports.MinioAPI = ports[4]
	m.cfg.Ports.MinioConsole = ports[5]
	m.cfg.Ports.Vault = ports[6]
	m.cfg.AppURL = fmt.Sprintf("http://localhost:%d", ports[0])
}

// GetConfig returns the updated config.
func (m PortsModel) GetConfig() *config.Config { return m.cfg }
