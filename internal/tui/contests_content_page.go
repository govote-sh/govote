package tui

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/govote-sh/govote/internal/api"
)

func (m model) updateContestContent(msg tea.Msg) (tea.Model, tea.Cmd) {
	if keyMsg, ok := msg.(tea.KeyMsg); ok {
		switch keyMsg.String() {
		case "esc":
			if m.contestsList != nil && !m.contestsList.SettingFilter() {
				m.currPage = contestsPage
			}
			return m, nil
		}
	}
	return m, nil
}

func (m model) viewContestContent() string {
	selectedContest := m.contestsList.SelectedItem().(api.Contest)

	// Title and bold styles
	titleStyle := m.render.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("205")).
		Background(lipgloss.Color("63")).
		Padding(0, 1).
		Render
	boldStyle := m.render.NewStyle().Bold(true).Render

	title := titleStyle("Contest Details")

	// Contest basic info
	basicInfo := fmt.Sprintf(
		"%s: %s\n%s: %s\n%s: %s\n%s: %s",
		boldStyle("Ballot Title"), selectedContest.BallotTitle,
		boldStyle("Office"), selectedContest.Office,
		boldStyle("Number Elected"), selectedContest.NumberElected,
		boldStyle("Ballot Placement"), selectedContest.BallotPlacement,
	)

	// Electorate Specifications, if any
	var electorateSpecs string
	if selectedContest.ElectorateSpecifications != "" {
		electorateSpecs = fmt.Sprintf("%s: %s", boldStyle("Electorate Specifications"), selectedContest.ElectorateSpecifications)
	}

	// Referendum Information, if applicable
	var referendumInfo string
	if selectedContest.ReferendumTitle != "" {
		referendumInfo = fmt.Sprintf(
			"%s\n%s\n%s: %s\n%s: %s\n%s: %s\n%s: %s\n%s: %s\n%s: %s\n",
			boldStyle("Referendum Details"),
			strings.Repeat("-", m.width/2),
			boldStyle("Title"), selectedContest.ReferendumTitle,
			boldStyle("Subtitle"), selectedContest.ReferendumSubtitle,
			boldStyle("Description"), selectedContest.ReferendumBrief,
			boldStyle("Pro Statement"), selectedContest.ReferendumProStatement,
			boldStyle("Con Statement"), selectedContest.ReferendumConStatement,
			boldStyle("URL"), selectedContest.ReferendumUrl,
		)
	}

	// Candidate Information
	var candidateInfo []string
	if len(selectedContest.Candidates) > 0 {
		candidateInfo = append(candidateInfo, boldStyle("Candidates:"))
		for _, candidate := range selectedContest.Candidates {
			candidateInfo = append(candidateInfo, fmt.Sprintf(
				"%s (Party: %s)", candidate.Name, candidate.Party,
			))
		}
	}

	return m.render.NewStyle().Margin(1, 1).MaxWidth(m.width).MaxHeight(m.height).Render(
		lipgloss.JoinVertical(
			lipgloss.Top,
			m.HeaderView(),
			title,
			basicInfo,
			electorateSpecs,
			referendumInfo,
			strings.Join(candidateInfo, "\n"),
		),
	)
}
