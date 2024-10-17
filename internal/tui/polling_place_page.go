package tui

import (
	"strings"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/govote-sh/govote/internal/http"
)

func (m model) viewPollingPlace() string {
	if m.selectedPollingPlace == nil {
		return "No polling place selected."
	}
	hoursTable := newPollingPlaceHoursTable(*m.selectedPollingPlace)
	titleStyle := m.render.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("205")).
		Background(lipgloss.Color("63")).
		Padding(0, 1).
		Render
	boldStyle := m.render.NewStyle().Bold(true).Render
	title := titleStyle("Polling Place Details")

	return m.render.NewStyle().Margin(1, 1).MaxWidth(m.width).MaxHeight(m.height).Render(lipgloss.JoinVertical(
		lipgloss.Top,
		m.HeaderView(),
		title,
		boldStyle(m.selectedPollingPlace.Address.String()),
		hoursTable.View(),
	))
}

func (m model) updatePollingPlace(msg tea.Msg) (tea.Model, tea.Cmd) {
	// Allow the user to exit by pressing "q" or "ctrl+c"
	if keyMsg, ok := msg.(tea.KeyMsg); ok {
		switch keyMsg.String() {
		case "esc":
			if m.pollingLocationListCreated && m.pollingLocationList.IsFiltered() || m.pollingLocationList.FilterState() == list.Unfiltered {
				m.currPage = votePage
			}
			return m, nil
		}
	}
	return m, nil
}

func newPollingPlaceHoursTable(p http.PollingPlace) table.Model {
	pollingHours := parsePollingHours(p.PollingHours)

	// Define columns for the table
	columns := []table.Column{
		{Title: "Day", Width: 20},
		{Title: "Hours", Width: 15},
	}

	// Create rows based on polling hours
	var rows []table.Row
	for _, entry := range pollingHours {
		day := entry[0]
		hours := entry[1]
		rows = append(rows, table.Row{day, hours})
	}

	// Create the table model with the rows and columns
	t := table.New(table.WithColumns(columns), table.WithRows(rows), table.WithHeight(15))

	// Style the table (optional)
	t.SetStyles(table.Styles{
		Header: lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("205")),
		Cell:   lipgloss.NewStyle().Foreground(lipgloss.Color("240")),
	})

	return t
}

func parsePollingHours(pollingHours string) [][2]string {
	var result [][2]string

	// Split the input string by newline
	lines := strings.Split(pollingHours, "\n")

	// Split each line by the first colon to get day and hours
	for _, line := range lines {
		parts := strings.SplitN(line, ": ", 2)
		if len(parts) == 2 {
			day := parts[0]
			hours := parts[1]
			result = append(result, [2]string{day, hours})
		}
	}
	return result
}
