package tui

import (
	"fmt"
	"strings"
	"unicode"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/govote-sh/govote/internal/api"
	"github.com/govote-sh/govote/internal/utils"
	"github.com/muesli/reflow/wordwrap"
)

func (m model) updateContestContent(msg tea.Msg) (tea.Model, tea.Cmd) {
	if keyMsg, ok := msg.(tea.KeyMsg); ok {
		switch keyMsg.String() {
		case "esc":
			m.currPage = contestsPage
			return m, nil
		}
	}
	return m, nil
}

func (m model) viewContestContent() string {
	selectedContest := m.contestsList.SelectedItem().(api.Contest)

	// Title styling
	title := sectionTitleStyle(m.render, "Contest Details")

	// Contest basic info
	var basicInfo []string
	if selectedContest.BallotTitle != "" {
		basicInfo = append(basicInfo, fmt.Sprintf("%s: %s", fieldLabelStyle(m.render, "Ballot Title"), fieldValueStyle(m.render, selectedContest.BallotTitle)))
	}
	if selectedContest.Office != "" {
		basicInfo = append(basicInfo, fmt.Sprintf("%s: %s", fieldLabelStyle(m.render, "Office"), fieldValueStyle(m.render, selectedContest.Office)))
	}
	if selectedContest.NumberElected != "" {
		basicInfo = append(basicInfo, fmt.Sprintf("%s: %s", fieldLabelStyle(m.render, "Number Elected"), fieldValueStyle(m.render, selectedContest.NumberElected)))
	}
	if selectedContest.BallotPlacement != "" {
		basicInfo = append(basicInfo, fmt.Sprintf("%s: %s", fieldLabelStyle(m.render, "Ballot Placement"), fieldValueStyle(m.render, selectedContest.BallotPlacement)))
	}

	// Electorate Specifications
	var electorateSpecs string
	if selectedContest.ElectorateSpecifications != "" {
		electorateSpecs = fmt.Sprintf("%s: %s", fieldLabelStyle(m.render, "Electorate Specifications"), fieldValueStyle(m.render, selectedContest.ElectorateSpecifications))
	}

	// Referendum Information for ballot-measure contests
	var referendumInfo []string
	if selectedContest.ReferendumTitle != "" {
		referendumInfo = append(referendumInfo, fmt.Sprintf("%s: %s", fieldLabelStyle(m.render, "Referendum Title"), fieldValueStyle(m.render, selectedContest.ReferendumTitle)))
	}
	if selectedContest.ReferendumText != "" {
		referendumInfo = append(referendumInfo, fmt.Sprintf("%s:\n%s", fieldLabelStyle(m.render, "Referendum Text"), fieldValueStyle(m.render, utils.Wrap(selectedContest.ReferendumText, m.width-4))))
	}
	if selectedContest.ReferendumSubtitle != "" {
		referendumInfo = append(referendumInfo, fmt.Sprintf("%s: %s", fieldLabelStyle(m.render, "Subtitle"), fieldValueStyle(m.render, selectedContest.ReferendumSubtitle)))
	}
	if selectedContest.ReferendumBrief != "" {
		referendumInfo = append(referendumInfo, fmt.Sprintf("%s: %s", fieldLabelStyle(m.render, "Description"), fieldValueStyle(m.render, selectedContest.ReferendumBrief)))
	}
	if selectedContest.ReferendumProStatement != "" {
		referendumInfo = append(referendumInfo, fmt.Sprintf("%s: %s", fieldLabelStyle(m.render, "Pro Statement"), fieldValueStyle(m.render, selectedContest.ReferendumProStatement)))
	}
	if selectedContest.ReferendumConStatement != "" {
		referendumInfo = append(referendumInfo, fmt.Sprintf("%s: %s", fieldLabelStyle(m.render, "Con Statement"), fieldValueStyle(m.render, selectedContest.ReferendumConStatement)))
	}
	if selectedContest.ReferendumUrl != "" {
		referendumInfo = append(referendumInfo, fmt.Sprintf("%s: %s", fieldLabelStyle(m.render, "URL"), fieldValueStyle(m.render, selectedContest.ReferendumUrl)))
	}

	// Candidate Information for office contests
	var candidateTable string
	if len(selectedContest.Candidates) > 0 {
		candidateTable = sectionTitleStyle(m.render, "Candidates") + "\n" + newCandidateTable(selectedContest.Candidates).View()
	}

	return m.render.NewStyle().Margin(1, 1).MaxWidth(m.width).MaxHeight(m.height).Render(
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

func wrap(s string, width int) string {
	s = strings.ReplaceAll(s, "\n", " ")
	s = strings.ReplaceAll(s, "\r", " ")
	return wordwrap.String(s, width)
}

// https://stackoverflow.com/questions/59955085/how-can-i-elliptically-truncate-text-in-golang
func ellipticalTruncate(text string, maxLen int) string {
	lastSpaceIx := maxLen
	len := 0
	for i, r := range text {
		if unicode.IsSpace(r) {
			lastSpaceIx = i
		}
		len++
		if len > maxLen {
			return text[:lastSpaceIx] + "..."
		}
	}
	return text
}
