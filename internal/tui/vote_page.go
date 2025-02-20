package tui

import (
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/govote-sh/govote/internal/api"
	"github.com/govote-sh/govote/internal/listManager"
)

func (m model) UpdateVote(msg tea.Msg) (tea.Model, tea.Cmd) {
	if m.lm != nil {
		var cmd tea.Cmd
		m.lm, cmd = m.lm.UpdateActiveList(msg)
		if cmd != nil {
			return m, cmd
		}
	}

	if keyMsg, ok := msg.(tea.KeyMsg); ok {
		switch keyMsg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		case "enter":
			_, ok := m.lm.SelectedItem().(api.PollingPlace)
			if ok {
				m.currPage = pollingPlacePage
			}
			return m, nil
		case "tab":
			if m.lm != nil {
				m.lm.CycleNext()
			}
		case "shift+tab":
			if m.lm != nil {
				m.lm.CyclePrev()
			}
		}
	}

	return m, nil
}

func (m model) viewVote() string {
	if m.lm == nil {
		return "building list..."
	}
	return m.render.NewStyle().Margin(1, 1).MaxWidth(m.width).MaxHeight(m.height).Render(lipgloss.JoinVertical(
		lipgloss.Top,
		m.HeaderView(),
		m.render.NewStyle().Foreground(lipgloss.Color("63")).MarginLeft(3).Render("Use tab to cycle through the lists of voting options"),
		m.lm.ActiveList().View(),
	))
}

func (m model) InitVotePageListManager() *listManager.ListManager {
	// Type conversions
	var pollingLocationItems, earlyVoteItems, dropOffItems []list.Item

	for _, pollingPlace := range m.electionData.PollingLocations {
		pollingLocationItems = append(pollingLocationItems, pollingPlace)
	}

	for _, earlyVoteSite := range m.electionData.EarlyVoteSites {
		earlyVoteItems = append(earlyVoteItems, earlyVoteSite)
	}

	for _, dropOffLocation := range m.electionData.DropOffLocations {
		dropOffItems = append(dropOffItems, dropOffLocation)
	}

	return listManager.InitListManager(
		[][]list.Item{
			pollingLocationItems,
			earlyVoteItems,
			dropOffItems,
		},
		[]string{
			"Polling Locations",
			"Early Voting Sites",
			"Drop Off Locations",
		},
		m.width, m.height-4,
	)
}
