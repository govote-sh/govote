package tui

import "github.com/charmbracelet/lipgloss"

func (m model) RenderErrorBox(text string) string {
	const HEADER_HEIGHT = 3
	return lipgloss.Place(
		m.width, m.height-HEADER_HEIGHT, lipgloss.Center, lipgloss.Center,
		lipgloss.NewStyle().
			Foreground(lipgloss.Color("205")).
			Background(lipgloss.Color("52")).
			Bold(true).
			Padding(1, 2).
			Border(lipgloss.RoundedBorder()).
			Width(m.width/3).
			Render(text),
	)
}

// JoinNonEmptyVertical does a lipgloss.JoinVertical, but skips empty arguments (avoiding empty lines)
func JoinNonEmptyVertical(pos lipgloss.Position, items ...string) string {
	nonEmptyItems := []string{}
	for _, item := range items {
		if item != "" {
			nonEmptyItems = append(nonEmptyItems, item)
		}
	}
	return lipgloss.JoinVertical(pos, nonEmptyItems...)
}
