package main

import "github.com/charmbracelet/lipgloss"

// Styles for the TUI
var (
	TitleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#FAFAFA")).
			Background(lipgloss.Color("#7D56F4")).
			Padding(0, 1)

	SelectedStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FAFAFA")).
			Background(lipgloss.Color("#7D56F4")).
			Bold(true)

	NormalStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#939093"))

	BorderStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#7D56F4")).
			Padding(0, 1)

	ErrorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FF5F56")).
			Bold(true)

	WarningStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFBD2E")).
			Bold(true)

	InfoStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#28CA42")).
			Bold(true)

	SuccessStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#28CA42")).
			Bold(true)

	// Pod status colors
	RunningStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#28CA42")).
			Bold(true)

	PendingStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFBD2E")).
			Bold(true)

	FailedStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FF5F56")).
			Bold(true)

	TerminatingStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#FF7F50")).
				Bold(true)

	UnknownStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#939093")).
			Bold(true)

	// Pod box styles
	PodBoxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#666666")).
			Padding(1, 2).
			Width(45).
			Height(6)

	SelectedPodBoxStyle = lipgloss.NewStyle().
				Border(lipgloss.RoundedBorder()).
				BorderForeground(lipgloss.Color("#7D56F4")).
				Background(lipgloss.Color("#1a1a2e")).
				Padding(1, 2).
				Width(45).
				Height(6)
)
