package tui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/well-prado/tetrix-cli/internal/config"
	"github.com/well-prado/tetrix-cli/internal/tui/steps"
)

// WizardModel orchestrates the installation wizard as a step-by-step state machine.
type WizardModel struct {
	cfg        *config.Config
	cliVersion string
	step       int
	models     []tea.Model
}

// NewWizard creates the full installation wizard.
func NewWizard(cfg *config.Config, cliVersion string) WizardModel {
	ports := map[string]int{
		"web":           cfg.Ports.Web,
		"api":           cfg.Ports.API,
		"multi-agent":   cfg.Ports.MultiAgent,
		"openai-proxy":  cfg.Ports.OpenAIProxy,
		"minio-api":     cfg.Ports.MinioAPI,
		"minio-console": cfg.Ports.MinioConsole,
		"vault":         cfg.Ports.Vault,
	}

	return WizardModel{
		cfg:        cfg,
		cliVersion: cliVersion,
		step:       0,
		models: []tea.Model{
			steps.NewWelcomeModel(cliVersion),        // 0
			steps.NewPrerequisitesModel(ports),        // 1
			steps.NewOAuthModel(cfg),                  // 2
			steps.NewAIKeysModel(cfg),                 // 3
			steps.NewPortsModel(cfg),                  // 4
			steps.NewReviewModel(cfg),                 // 5
			steps.NewProgressModel(cfg, cliVersion),   // 6
		},
	}
}

func (m WizardModel) Init() tea.Cmd {
	return m.models[m.step].Init()
}

func (m WizardModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// Handle navigation messages from steps
	switch msg.(type) {
	case steps.NextStepMsg:
		if m.step < len(m.models)-1 {
			m.step++
			// Rebuild review model to reflect latest config
			if m.step == 5 {
				m.models[5] = steps.NewReviewModel(m.cfg)
			}
			return m, m.models[m.step].Init()
		}
		return m, nil
	case steps.PrevStepMsg:
		if m.step > 0 {
			m.step--
			return m, m.models[m.step].Init()
		}
		return m, nil
	}

	// Delegate to current step
	var cmd tea.Cmd
	m.models[m.step], cmd = m.models[m.step].Update(msg)
	return m, cmd
}

func (m WizardModel) View() string {
	return m.models[m.step].View()
}
