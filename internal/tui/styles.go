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

func sectionTitleStyle(r *lipgloss.Renderer, text string) string {
	return r.NewStyle().
		Foreground(lipgloss.Color("205")).
		Bold(true).
		Render(text)
}

func fieldLabelStyle(r *lipgloss.Renderer, text string) string {
	return r.NewStyle().
		Foreground(lipgloss.Color("255")).
		Render(text)
}

func fieldValueStyle(r *lipgloss.Renderer, text string) string {
	return r.NewStyle().
		Foreground(lipgloss.Color("63")).
		Render(text)
}

// joinNonEmptyVertical does a lipgloss.JoinVertical, but skips empty arguments (avoiding empty lines)
func joinNonEmptyVertical(pos lipgloss.Position, items ...string) string {
	nonEmptyItems := []string{}
	for _, item := range items {
		if item != "" {
			nonEmptyItems = append(nonEmptyItems, item)
		}
	}
	return lipgloss.JoinVertical(pos, nonEmptyItems...)
}
