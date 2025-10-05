package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/govote-sh/govote/internal/api"
)

func (m model) viewPollingPlace() string {
	// Check if list manager is nil
	if m.lm == nil {
		return m.renderPageError("No polling place selected")
	}

	// Check if selected item exists
	selectedItem := m.lm.SelectedItem()
	if selectedItem == nil {
		return m.renderPageError("No polling place selected")
	}

	// Type assert with safety check
	selectedPollingPlace, ok := selectedItem.(api.PollingPlace)
	if !ok {
		return m.renderPageError("Invalid polling place data")
	}

	// Title and bold styles
	titleStyle := m.render.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("205")).
		Background(lipgloss.Color("63")).
		Padding(0, 1).
		Render
	boldStyle := m.render.NewStyle().Bold(true).Render

	title := titleStyle("Polling Place Details")

	address := boldStyle(selectedPollingPlace.Address.String())

	hoursTable := newPollingPlaceHoursTable(selectedPollingPlace.PollingHours)

	// Notes (if any)
	var notes string
	if selectedPollingPlace.Notes != "" {
		notes = boldStyle("Notes: ") + fieldValueStyle(m.render, selectedPollingPlace.Notes)
	}

	// Voter services (if any)
	var voterServices string
	if selectedPollingPlace.VoterServices != "" {
		voterServices = boldStyle("Voter Services: ") + fieldValueStyle(m.render, selectedPollingPlace.VoterServices)
	}

	// Start and end dates (if any)
	var dates string
	if selectedPollingPlace.StartDate != "" && selectedPollingPlace.EndDate != "" {
		if selectedPollingPlace.StartDate == selectedPollingPlace.EndDate {
			dates = fmt.Sprintf("%s: %s", boldStyle("Date"), fieldValueStyle(m.render, selectedPollingPlace.StartDate))
		} else {
			dates = fmt.Sprintf("%s: %s → %s", boldStyle("Available Dates"), selectedPollingPlace.StartDate, selectedPollingPlace.EndDate)
		}
	} else {
		dates = ""
	}

	// Latitude and Longitude (if any)
	var coordinates string
	if url, err := selectedPollingPlace.GetMapsUrl(); err == nil {
		coordinates = boldStyle("Map link: ") + fieldValueStyle(m.render, url)
	}

	return m.render.NewStyle().Margin(1, 1).MaxWidth(m.width).MaxHeight(m.height).Render(
		joinNonEmptyVertical(
			lipgloss.Top,
			m.HeaderView(),
			title,
			address,
			"\t",
			hoursTable.View(),
			"\t",
			notes,
			voterServices,
			dates,
			coordinates,
		),
	)
}

func (m model) updatePollingPlace(msg tea.Msg) (tea.Model, tea.Cmd) {
	// Allow the user to exit by pressing "q" or "ctrl+c"
	if keyMsg, ok := msg.(tea.KeyMsg); ok {
		switch keyMsg.String() {
		case "esc":
			if m.lm != nil && !m.lm.SettingFilter() {
				m.currPage = votePage
			}
			return m, nil
		}
	}
	return m, nil
}

func newPollingPlaceHoursTable(hours string) table.Model {
	pollingHours := parsePollingHours(hours)

	// Define columns for the table
	columns := []table.Column{
		{Title: "Day", Width: 20},
		{Title: "Hours", Width: 24},
	}

	// Create rows based on polling hours
	var rows []table.Row
	for _, entry := range pollingHours {
		day := entry[0]
		hours := entry[1]
		rows = append(rows, table.Row{day, hours})
	}

	tableHeight := min(15+1, len(rows)+1)

	// Create the table model with the rows and columns
	t := table.New(table.WithColumns(columns), table.WithRows(rows), table.WithHeight(tableHeight))

	t.SetStyles(table.Styles{
		Header: lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("205")),
		Cell:   lipgloss.NewStyle().Foreground(lipgloss.Color("255")),
	})

	return t
}

func parsePollingHours(pollingHours string) [][2]string {
	var result [][2]string

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
