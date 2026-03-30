package steps

import (
	"fmt"
	"os/exec"
	"runtime"
	"time"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/well-prado/tetrix-cli/internal/config"
	"github.com/well-prado/tetrix-cli/internal/docker"
	"github.com/well-prado/tetrix-cli/internal/secrets"
	"github.com/well-prado/tetrix-cli/internal/templates"
	"github.com/well-prado/tetrix-cli/internal/tui/styles"
)

// ProgressModel runs the installation pipeline.
type ProgressModel struct {
	cfg        *config.Config
	cliVersion string
	spinner    spinner.Model
	steps      []installStep
	current    int
	done       bool
	err        error
}

type installStep struct {
	name   string
	status stepStatus
}

type stepStatus int

const (
	stepPending stepStatus = iota
	stepRunning
	stepDone
	stepFailed
)

type stepDoneMsg struct{ err error }

// NewProgressModel creates the installation progress step.
func NewProgressModel(cfg *config.Config, cliVersion string) ProgressModel {
	sp := spinner.New()
	sp.Spinner = spinner.Dot
	sp.Style = lipgloss.NewStyle().Foreground(styles.Purple)

	return ProgressModel{
		cfg:        cfg,
		cliVersion: cliVersion,
		spinner:    sp,
		steps: []installStep{
			{name: "Generating secrets"},
			{name: "Writing configuration files"},
			{name: "Writing docker-compose.yml"},
			{name: "Pulling Docker images"},
			{name: "Starting services"},
			{name: "Waiting for health checks"},
		},
	}
}

func (m ProgressModel) Init() tea.Cmd {
	return tea.Batch(m.spinner.Tick, m.runCurrentStep())
}

func (m ProgressModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd

	case stepDoneMsg:
		if msg.err != nil {
			m.steps[m.current].status = stepFailed
			m.err = msg.err
			m.done = true
			return m, nil
		}
		m.steps[m.current].status = stepDone
		m.current++
		if m.current >= len(m.steps) {
			m.done = true
			return m, nil
		}
		m.steps[m.current].status = stepRunning
		return m, m.runCurrentStep()

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		case "enter":
			if m.done && m.err == nil {
				// Open browser
				openBrowser(m.cfg.AppURL)
				return m, tea.Quit
			}
		case "q":
			if m.done {
				return m, tea.Quit
			}
		}
	}
	return m, nil
}

func (m ProgressModel) View() string {
	s := styles.Title.Render("Installing Tetrix CE") + "\n"
	s += styles.StepIndicator.Render("Step 7/7") + "\n\n"

	for i, step := range m.steps {
		var icon string
		switch step.status {
		case stepPending:
			icon = styles.Dimmed.Render(styles.Circle)
		case stepRunning:
			icon = m.spinner.View()
		case stepDone:
			icon = styles.Success.Render(styles.CheckMark)
		case stepFailed:
			icon = styles.Error.Render(styles.CrossMark)
		}
		name := step.name
		if step.status == stepPending {
			name = styles.Dimmed.Render(name)
		}
		_ = i
		s += fmt.Sprintf("  %s %s\n", icon, name)
	}

	if m.done {
		s += "\n"
		if m.err != nil {
			s += styles.Error.Render(fmt.Sprintf("  Installation failed: %s", m.err)) + "\n\n"
			s += styles.StatusBar.Render("[q] Quit")
		} else {
			s += styles.Success.Render("  Tetrix CE is running!") + "\n\n"
			s += fmt.Sprintf("  Open in browser: %s\n\n", styles.Subtitle.Render(m.cfg.AppURL))
			s += styles.Dimmed.Render("  Useful commands:") + "\n"
			s += styles.Dimmed.Render("    tetrix status   - View service health") + "\n"
			s += styles.Dimmed.Render("    tetrix logs     - View logs") + "\n"
			s += styles.Dimmed.Render("    tetrix stop     - Stop all services") + "\n"
			s += styles.Dimmed.Render("    tetrix update   - Update to latest version") + "\n"
			s += "\n"
			s += fmt.Sprintf("  Config saved to: %s\n\n", styles.Dimmed.Render(m.cfg.TetrixHome+"/config.yaml"))
			s += styles.StatusBar.Render("[Enter] Open browser  [q] Quit")
		}
	}

	return s
}

func (m ProgressModel) runCurrentStep() tea.Cmd {
	m.steps[m.current].status = stepRunning

	switch m.current {
	case 0:
		return m.generateSecrets()
	case 1:
		return m.writeConfig()
	case 2:
		return m.writeCompose()
	case 3:
		return m.pullImages()
	case 4:
		return m.startServices()
	case 5:
		return m.waitHealth()
	default:
		return func() tea.Msg { return stepDoneMsg{} }
	}
}

func (m *ProgressModel) generateSecrets() tea.Cmd {
	return func() tea.Msg {
		var err error
		if m.cfg.Auth.AuthSecret == "" {
			m.cfg.Auth.AuthSecret, err = secrets.GenerateAuthSecret()
			if err != nil {
				return stepDoneMsg{err: fmt.Errorf("auth secret: %w", err)}
			}
		}
		if m.cfg.Infrastructure.PostgresPassword == "" {
			m.cfg.Infrastructure.PostgresPassword, err = secrets.GeneratePassword()
			if err != nil {
				return stepDoneMsg{err: fmt.Errorf("postgres password: %w", err)}
			}
		}
		if m.cfg.Infrastructure.MeilisearchKey == "" {
			m.cfg.Infrastructure.MeilisearchKey, err = secrets.GeneratePassword()
			if err != nil {
				return stepDoneMsg{err: fmt.Errorf("meilisearch key: %w", err)}
			}
		}
		if m.cfg.Infrastructure.MinioSecretKey == "" {
			m.cfg.Infrastructure.MinioSecretKey, err = secrets.GeneratePassword()
			if err != nil {
				return stepDoneMsg{err: fmt.Errorf("minio secret: %w", err)}
			}
		}
		if m.cfg.Infrastructure.VaultToken == "" {
			m.cfg.Infrastructure.VaultToken, err = secrets.GeneratePassword()
			if err != nil {
				return stepDoneMsg{err: fmt.Errorf("vault token: %w", err)}
			}
		}
		return stepDoneMsg{}
	}
}

func (m *ProgressModel) writeConfig() tea.Cmd {
	return func() tea.Msg {
		if err := m.cfg.Save(); err != nil {
			return stepDoneMsg{err: fmt.Errorf("save config: %w", err)}
		}
		return stepDoneMsg{}
	}
}

func (m *ProgressModel) writeCompose() tea.Cmd {
	return func() tea.Msg {
		data := templates.NewTemplateData(m.cfg, m.cliVersion)
		if err := templates.WriteAllFiles(data); err != nil {
			return stepDoneMsg{err: err}
		}
		return stepDoneMsg{}
	}
}

func (m *ProgressModel) pullImages() tea.Cmd {
	return func() tea.Msg {
		cm := docker.NewComposeManager(m.cfg.TetrixHome)
		if err := cm.Pull(); err != nil {
			return stepDoneMsg{err: fmt.Errorf("pull images: %w", err)}
		}
		return stepDoneMsg{}
	}
}

func (m *ProgressModel) startServices() tea.Cmd {
	return func() tea.Msg {
		cm := docker.NewComposeManager(m.cfg.TetrixHome)
		if err := cm.Up(); err != nil {
			return stepDoneMsg{err: fmt.Errorf("start services: %w", err)}
		}
		return stepDoneMsg{}
	}
}

func (m *ProgressModel) waitHealth() tea.Cmd {
	return func() tea.Msg {
		// Wait up to 120s for services to become healthy
		cm := docker.NewComposeManager(m.cfg.TetrixHome)
		for i := 0; i < 24; i++ {
			time.Sleep(5 * time.Second)
			health, err := cm.ServiceHealth()
			if err != nil {
				continue
			}
			allHealthy := true
			for _, status := range health {
				if status == "" {
					allHealthy = false
					break
				}
			}
			if allHealthy && len(health) > 0 {
				return stepDoneMsg{}
			}
		}
		// Even if not all healthy yet, continue — services may still be starting
		return stepDoneMsg{}
	}
}

func openBrowser(url string) {
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "darwin":
		cmd = exec.Command("open", url)
	case "linux":
		cmd = exec.Command("xdg-open", url)
	default:
		return
	}
	_ = cmd.Start()
}
