package steps

import tea "github.com/charmbracelet/bubbletea"

// Navigation messages used by the wizard state machine.

// NextStepMsg is sent when a step wants to advance.
type NextStepMsg struct{}

// PrevStepMsg is sent when a step wants to go back.
type PrevStepMsg struct{}

// NextStep sends a message to advance to the next wizard step.
func NextStep() tea.Msg { return NextStepMsg{} }

// PrevStep sends a message to go back to the previous wizard step.
func PrevStep() tea.Msg { return PrevStepMsg{} }
