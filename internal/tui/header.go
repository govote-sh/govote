package tui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
)

func (m model) HeaderUpdate(msg tea.Msg) (model, tea.Cmd) {
	if !m.hasMenu {
		return m, nil
	}
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		// Or directly navigate to a specific tab
		case "p":
			m.page = pollingLocationPage
		// case "e":
		// 	m.page = earlyVotePage
		// case "d":
		// 	m.page = ballotDropOffPage
		case "c":
			m.page = contestsPage
		case "r":
			m.page = registerPage
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

	// Define the tabs
	title := activeTabStyle("govote.sh")
	electionDay := inactiveTabStyle("Vote")
	// earlyVote := inactiveTabStyle("Early Vote")
	// ballotDropOff := inactiveTabStyle("Drop Off")
	contests := inactiveTabStyle("Contests")
	register := inactiveTabStyle("Register")

	// Bold the active tab based on the current page
	switch m.page {
	case pollingLocationPage:
		electionDay = activeTabStyle("Vote")
	// case earlyVotePage:
	// 	earlyVote = activeTabStyle("Early Vote")
	// case ballotDropOffPage:
	// 	if len(m.electionData.DropOffLocations) > 0 {
	// 		ballotDropOff = activeTabStyle("Drop Off")
	// 	}
	case contestsPage:
		contests = activeTabStyle("Contests")
	case registerPage:
		register = activeTabStyle("Register")
	}

	var tabs []string
	tabs = []string{title, electionDay, contests, register}
	return table.New().
		Border(lipgloss.NormalBorder()).
		Row(tabs...).
		Width(m.width).
		StyleFunc(func(row, col int) lipgloss.Style {
			return m.render.NewStyle().
				Padding(0, 1).
				AlignHorizontal(lipgloss.Center)
		}).
		Render()
}
