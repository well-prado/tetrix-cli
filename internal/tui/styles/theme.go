package styles

import "github.com/charmbracelet/lipgloss"

// Brand colors
var (
	Purple    = lipgloss.Color("#7C3AED")
	Blue      = lipgloss.Color("#3B82F6")
	Green     = lipgloss.Color("#10B981")
	Red       = lipgloss.Color("#EF4444")
	Yellow    = lipgloss.Color("#F59E0B")
	Gray      = lipgloss.Color("#6B7280")
	DarkGray  = lipgloss.Color("#374151")
	LightGray = lipgloss.Color("#9CA3AF")
	White     = lipgloss.Color("#F9FAFB")
	Dim       = lipgloss.Color("#4B5563")
)

// Reusable styles
var (
	// Title is for section headers
	Title = lipgloss.NewStyle().
		Bold(true).
		Foreground(Purple).
		MarginBottom(1)

	// Subtitle is for secondary headers
	Subtitle = lipgloss.NewStyle().
			Bold(true).
			Foreground(Blue)

	// Success marks passed checks
	Success = lipgloss.NewStyle().
		Foreground(Green).
		Bold(true)

	// Error marks failed checks
	Error = lipgloss.NewStyle().
		Foreground(Red).
		Bold(true)

	// Warning marks non-critical issues
	Warning = lipgloss.NewStyle().
		Foreground(Yellow)

	// Dimmed is for secondary text
	Dimmed = lipgloss.NewStyle().
		Foreground(Dim)

	// Label is for form field labels
	Label = lipgloss.NewStyle().
		Foreground(LightGray).
		Width(26)

	// Value is for form field values
	Value = lipgloss.NewStyle().
		Foreground(White)

	// StatusBar is the bottom navigation bar
	StatusBar = lipgloss.NewStyle().
			Foreground(Gray).
			MarginTop(1)

	// Box wraps content in a bordered container
	Box = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(DarkGray).
		Padding(1, 2)

	// ActiveInput highlights the current input
	ActiveInput = lipgloss.NewStyle().
			Foreground(Purple)

	// Cursor style for text inputs
	Cursor = lipgloss.NewStyle().
		Foreground(Purple)

	// StepIndicator shows current step progress (e.g., "Step 2/7")
	StepIndicator = lipgloss.NewStyle().
			Foreground(Gray).
			Italic(true)
)

// Icons
const (
	CheckMark = "✓"
	CrossMark = "✗"
	Bullet    = "●"
	Circle    = "○"
	Arrow     = "→"
	Spinner   = "◐"
)

// Logo returns the Tetrix ASCII art logo.
func Logo() string {
	logo := `
 ████████╗███████╗████████╗██████╗ ██╗██╗  ██╗
 ╚══██╔══╝██╔════╝╚══██╔══╝██╔══██╗██║╚██╗██╔╝
    ██║   █████╗     ██║   ██████╔╝██║ ╚███╔╝
    ██║   ██╔══╝     ██║   ██╔══██╗██║ ██╔██╗
    ██║   ███████╗   ██║   ██║  ██║██║██╔╝ ██╗
    ╚═╝   ╚══════╝   ╚═╝   ╚═╝  ╚═╝╚═╝╚═╝  ╚═╝`
	return lipgloss.NewStyle().Foreground(Purple).Bold(true).Render(logo)
}
