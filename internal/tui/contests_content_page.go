package tui

import (
	"fmt"

	"charm.land/bubbles/v2/table"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/govote-sh/govote/internal/api"
	"github.com/govote-sh/govote/internal/utils"
)

func (m model) updateContestContent(msg tea.Msg) (tea.Model, tea.Cmd) {
	if keyMsg, ok := msg.(tea.KeyPressMsg); ok {
		switch keyMsg.String() {
		case "esc":
			m.currPage = contestsPage
			return m, nil
		}
	}
	return m, nil
}

func (m model) viewContestContent() string {
	// Check if contestsList is nil
	if m.contestsList == nil {
		return m.renderPageError("No contest selected")
	}

	// Check if selected item exists
	selectedItem := m.contestsList.SelectedItem()
	if selectedItem == nil {
		return m.renderPageError("No contest selected")
	}

	// Type assert with safety check
	selectedContest, ok := selectedItem.(api.Contest)
	if !ok {
		return m.renderPageError("Invalid contest data")
	}

	// Title styling
	title := sectionTitleStyle("Contest Details")

	// Contest basic info
	var basicInfo []string
	if selectedContest.BallotTitle != "" {
		basicInfo = append(basicInfo, fmt.Sprintf("%s: %s", fieldLabelStyle("Ballot Title"), fieldValueStyle(selectedContest.BallotTitle)))
	}
	if selectedContest.Office != "" {
		basicInfo = append(basicInfo, fmt.Sprintf("%s: %s", fieldLabelStyle("Office"), fieldValueStyle(selectedContest.Office)))
	}
	if selectedContest.NumberElected != "" {
		basicInfo = append(basicInfo, fmt.Sprintf("%s: %s", fieldLabelStyle("Number Elected"), fieldValueStyle(selectedContest.NumberElected)))
	}
	if selectedContest.BallotPlacement != "" {
		basicInfo = append(basicInfo, fmt.Sprintf("%s: %s", fieldLabelStyle("Ballot Placement"), fieldValueStyle(selectedContest.BallotPlacement)))
	}

	// Electorate Specifications
	var electorateSpecs string
	if selectedContest.ElectorateSpecifications != "" {
		electorateSpecs = fmt.Sprintf("%s: %s", fieldLabelStyle("Electorate Specifications"), fieldValueStyle(selectedContest.ElectorateSpecifications))
	}

	// Referendum Information for ballot-measure contests
	var referendumInfo []string
	if selectedContest.ReferendumTitle != "" {
		referendumInfo = append(referendumInfo, fmt.Sprintf("%s: %s", fieldLabelStyle("Referendum Title"), fieldValueStyle(selectedContest.ReferendumTitle)))
	}
	if selectedContest.ReferendumText != "" {
		referendumInfo = append(referendumInfo, fmt.Sprintf("%s:\n%s", fieldLabelStyle("Referendum Text"), fieldValueStyle(utils.Wrap(selectedContest.ReferendumText, m.width-4))))
	}
	if selectedContest.ReferendumSubtitle != "" {
		referendumInfo = append(referendumInfo, fmt.Sprintf("%s: %s", fieldLabelStyle("Subtitle"), fieldValueStyle(selectedContest.ReferendumSubtitle)))
	}
	if selectedContest.ReferendumBrief != "" {
		referendumInfo = append(referendumInfo, fmt.Sprintf("%s: %s", fieldLabelStyle("Description"), fieldValueStyle(selectedContest.ReferendumBrief)))
	}
	if selectedContest.ReferendumProStatement != "" {
		referendumInfo = append(referendumInfo, fmt.Sprintf("%s: %s", fieldLabelStyle("Pro Statement"), fieldValueStyle(selectedContest.ReferendumProStatement)))
	}
	if selectedContest.ReferendumConStatement != "" {
		referendumInfo = append(referendumInfo, fmt.Sprintf("%s: %s", fieldLabelStyle("Con Statement"), fieldValueStyle(selectedContest.ReferendumConStatement)))
	}
	if selectedContest.ReferendumUrl != "" {
		referendumInfo = append(referendumInfo, fmt.Sprintf("%s: %s", fieldLabelStyle("URL"), fieldValueStyle(selectedContest.ReferendumUrl)))
	}

	// Candidate Information for office contests
	var candidateTable string
	if len(selectedContest.Candidates) > 0 {
		candidateTable = sectionTitleStyle("Candidates") + "\n" + newCandidateTable(selectedContest.Candidates).View()
	}

	return lipgloss.NewStyle().Margin(1, 1).MaxWidth(m.width).MaxHeight(m.height).Render(
		joinNonEmptyVertical(
			lipgloss.Top,
			m.HeaderView(),
			title,
			joinNonEmptyVertical(lipgloss.Top, basicInfo...),
			electorateSpecs,
			joinNonEmptyVertical(lipgloss.Top, referendumInfo...),
			candidateTable,
		),
	)
}

// Helper function to create a candidate table
func newCandidateTable(candidates []api.Candidate) table.Model {
	columns := []table.Column{
		{Title: "Name", Width: 45},
		{Title: "Party", Width: 20},
	}
	var rows []table.Row
	for _, candidate := range candidates {
		rows = append(rows, table.Row{
			candidate.Name,
			candidate.Party,
		})
	}
	t := table.New(table.WithColumns(columns), table.WithRows(rows), table.WithHeight(10))
	t.SetStyles(table.Styles{
		Header: lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("205")),
		Cell:   lipgloss.NewStyle().Foreground(lipgloss.Color("255")),
	})
	return t
}
