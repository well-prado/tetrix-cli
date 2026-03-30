package steps

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/well-prado/tetrix-cli/internal/docker"
	"github.com/well-prado/tetrix-cli/internal/tui/styles"
)

// PrerequisitesModel checks system requirements.
type PrerequisitesModel struct {
	results []docker.PrereqResult
	checked bool
	allPass bool
	ports   map[string]int
}

// NewPrerequisitesModel creates the prerequisites check step.
func NewPrerequisitesModel(ports map[string]int) PrerequisitesModel {
	return PrerequisitesModel{ports: ports}
}

type checkDoneMsg struct {
	results []docker.PrereqResult
}

func (m PrerequisitesModel) Init() tea.Cmd {
	return func() tea.Msg {
		results := docker.CheckAll(m.ports)
		return checkDoneMsg{results: results}
	}
}

func (m PrerequisitesModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case checkDoneMsg:
		m.results = msg.results
		m.checked = true
		m.allPass = docker.AllPassed(msg.results)
		return m, nil
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			if m.allPass {
				return m, NextStep
			}
		case "r":
			if m.checked && !m.allPass {
				m.checked = false
				return m, m.Init()
			}
		case "q", "ctrl+c":
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m PrerequisitesModel) View() string {
	s := styles.Title.Render("Prerequisites Check") + "\n"
	s += styles.StepIndicator.Render("Step 1/7") + "\n\n"

	if !m.checked {
		s += styles.Dimmed.Render("  Checking system requirements...") + "\n"
		return s
	}

	for _, r := range m.results {
		var icon, detail string
		if r.Passed {
			icon = styles.Success.Render(styles.CheckMark)
			detail = styles.Dimmed.Render(r.Detail)
		} else {
			icon = styles.Error.Render(styles.CrossMark)
			detail = styles.Error.Render(r.Detail)
		}
		s += fmt.Sprintf("  %s %s  %s\n", icon, r.Name, detail)
		if r.Warning != "" {
			s += styles.Warning.Render(fmt.Sprintf("    ⚠ %s", r.Warning)) + "\n"
		}
	}

	s += "\n"
	if m.allPass {
		s += styles.Success.Render("All checks passed!") + "\n\n"
		s += styles.StatusBar.Render("[Enter] Continue  [q] Quit")
	} else {
		s += styles.Error.Render("Some checks failed. Please fix the issues above.") + "\n\n"
		s += styles.StatusBar.Render("[r] Re-check  [q] Quit")
	}
	return s
}
