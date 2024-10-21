package tui

import (
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/govote-sh/govote/internal/http"
	"github.com/govote-sh/govote/internal/listManager"
)

func (m model) UpdateVote(msg tea.Msg) (tea.Model, tea.Cmd) {
	// Handle list updates

	// Allow the user to exit by pressing "q" or "ctrl+c"
	if keyMsg, ok := msg.(tea.KeyMsg); ok {
		switch keyMsg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		case "enter":
			if m.lm != nil && m.lm.SettingFilter() {
				break
			}
			selectedItem, ok := m.lm.SelectedItem().(http.PollingPlace)
			if ok {
				// Move to the polling place detail page with the selected item
				m.selectedPollingPlace = &selectedItem
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

	if m.lm != nil {
		var cmd tea.Cmd
		m.lm, cmd = m.lm.UpdateActiveList(msg)
		if cmd != nil {
			return m, cmd
		}
	}
	return m, nil
}

func (m model) viewVote() string {
	// headerText := fmt.Sprintf("Upcoming %s on %s", m.electionData.Election.Name, m.electionData.Election.ElectionDay)
	// header := m.headerStyle.Render(headerText)
	// subtitleText := fmt.Sprintf("Results for: %s", m.electionData.NormalizedInput.String())
	// subtitle := m.headerStyle.MarginBottom(1).Render(subtitleText)
	if m.lm == nil {
		return "building list..."
	}
	return m.render.NewStyle().Margin(1, 1).MaxWidth(m.width).MaxHeight(m.height).Render(lipgloss.JoinVertical(
		lipgloss.Top,
		m.HeaderView(),
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
		m.width, m.height,
	)
}

// refactor: take in a List struct (list and createdBool) and title
/*
func (m *model) initList(width, height int) {
	if m.electionData == nil {
		fmt.Println("electionData is nil")
		return
	}

	// Check if PollingLocations is nil or empty
	if m.electionData.PollingLocations == nil || len(m.electionData.PollingLocations) == 0 {
		fmt.Println("PollingLocations is nil or empty")
		return
	}

	// Convert []PollingPlace to []list.Item
	var items []list.Item
	for _, pollingPlace := range m.electionData.EarlyVoteSites {
		items = append(items, pollingPlace)
	}

	m.pollingLocationList = list.New(items, list.NewDefaultDelegate(), width, height)
	m.pollingLocationList.Title = "Polling Locations"
	m.pollingLocationListCreated = true
}
*/
