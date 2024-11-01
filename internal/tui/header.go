package tui

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
)

func (m model) HeaderUpdate(msg tea.Msg) (model, tea.Cmd) {
	if !m.hasMenu || (m.lm != nil && m.lm.SettingFilter()) {
		return m, nil
	}
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		// Or directly navigate to a specific tab
		case "v":
			m.currPage = votePage
		case "c":
			m.currPage = contestsPage
		case "r":
			m.currPage = registerPage
		case "q": // Quit
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m model) HeaderView() string {
	// Define the styles for active and inactive tabs
	activeTabStyle := m.render.NewStyle().Bold(true).Foreground(lipgloss.Color("205")).Render
	inactiveTabStyle := m.render.NewStyle().Foreground(lipgloss.Color("240")).Render
	letterStyle := m.render.NewStyle().Foreground(lipgloss.Color("205")).Render // Always active style for letter indicators

	// Define the tabs with letter indicators
	title := activeTabStyle("govote.sh")
	esc := fmt.Sprintf("%s %s", letterStyle("[ESC]"), inactiveTabStyle("Back"))
	electionDay := fmt.Sprintf("%s %s", letterStyle("[V]"), inactiveTabStyle("Vote"))
	contests := fmt.Sprintf("%s %s", letterStyle("[C]"), inactiveTabStyle("Contests"))
	register := fmt.Sprintf("%s %s", letterStyle("[R]"), inactiveTabStyle("Register"))

	// Bold the active tab based on the current page
	switch m.currPage {
	case votePage:
		electionDay = fmt.Sprintf("%s %s", letterStyle("[V]"), activeTabStyle("Vote"))
	case contestsPage:
		contests = fmt.Sprintf("%s %s", letterStyle("[C]"), activeTabStyle("Contests"))
	case registerPage:
		register = fmt.Sprintf("%s %s", letterStyle("[R]"), activeTabStyle("Register"))
	}

	// Combine the tabs and ensure proper padding to avoid the bar cutting off
	var tabs []string
	if m.currPage != pollingPlacePage && m.currPage != contestContentPage {
		tabs = []string{title, electionDay, contests, register}
	} else {
		tabs = []string{title, esc}
	}
	return table.New().
		Border(lipgloss.NormalBorder()).
		Row(tabs...).
		Width(m.width - 2). // Add extra space to account for borders
		StyleFunc(func(row, col int) lipgloss.Style {
			return m.render.NewStyle().
				Padding(0, 2). // Padding on both sides, including right
				AlignHorizontal(lipgloss.Center)
		}).
		Render()
}
