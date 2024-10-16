package tui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/govote-sh/govote/internal/http"
)

func (m model) UpdateVote(msg tea.Msg) (tea.Model, tea.Cmd) {
	// Handle list updates
	if m.pollingLocationListCreated {
		var cmd tea.Cmd
		m.pollingLocationList, cmd = m.pollingLocationList.Update(msg)
		if cmd != nil {
			return m, cmd
		}
	}

	// Allow the user to exit by pressing "q" or "ctrl+c"
	if keyMsg, ok := msg.(tea.KeyMsg); ok {
		switch keyMsg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		case "enter":
			selectedItem, ok := m.pollingLocationList.SelectedItem().(http.PollingPlace)
			if ok {
				// Move to the polling place detail page with the selected item
				m.selectedPollingPlace = &selectedItem
				m.currPage = pollingPlacePage
			}
			return m, nil
		}
	}
	return m, nil
}

func (m model) viewVote() string {
	// headerText := fmt.Sprintf("Upcoming %s on %s", m.electionData.Election.Name, m.electionData.Election.ElectionDay)
	// header := m.headerStyle.Render(headerText)
	// subtitleText := fmt.Sprintf("Results for: %s", m.electionData.NormalizedInput.String())
	// subtitle := m.headerStyle.MarginBottom(1).Render(subtitleText)
	if !m.pollingLocationListCreated {
		return "building list..."
	}
	return m.render.NewStyle().Margin(1, 1).MaxWidth(m.width).MaxHeight(m.height).Render(lipgloss.JoinVertical(
		lipgloss.Top,
		m.HeaderView(),
		m.pollingLocationList.View(),
	))
}
