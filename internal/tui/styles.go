package tui

import "github.com/charmbracelet/lipgloss"

var (
	HeaderStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("205")).Align(lipgloss.Center).Bold(true).Padding(0, 1)
	SubtitleStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("240")).Align(lipgloss.Center).Padding(0, 1)
)
