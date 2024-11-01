package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/govote-sh/govote/internal/api"
)

func formatElectionAdministration(admin api.ElectionAdministrationBody, render *lipgloss.Renderer) string {
	var sections []string

	sectionTitleStyle := render.NewStyle().
		Foreground(lipgloss.Color("205")).
		Bold(true).
		Render
	fieldLabelStyle := render.NewStyle().
		Foreground(lipgloss.Color("240")).
		Bold(true).
		Render
	fieldValueStyle := render.NewStyle().
		Foreground(lipgloss.Color("63")).
		Render

	// Title for Election Administration section
	sections = append(sections, sectionTitleStyle("Election Administration"))
	sections = append(sections, fieldValueStyle(admin.Name))

	// Append URLs if they exist
	urlFields := []struct {
		label string
		value string
	}{
		{"Election Info", admin.ElectionInfoUrl},
		{"Registration URL", admin.ElectionRegistrationUrl},
		{"Confirmation URL", admin.ElectionRegistrationConfirmationUrl},
		{"Absentee Voting Info", admin.AbsenteeVotingInfoUrl},
		{"Location Finder", admin.VotingLocationFinderUrl},
		{"Ballot Info", admin.BallotInfoUrl},
		{"Election Rules", admin.ElectionRulesUrl},
	}
	for _, field := range urlFields {
		if field.value != "" {
			sections = append(sections, fmt.Sprintf("%s: %s", fieldLabelStyle(field.label), fieldValueStyle(field.value)))
		}
	}

	// Append Hours of Operation if they exist
	if admin.HoursOfOperation != "" {
		sections = append(sections, fmt.Sprintf("%s: %s", fieldLabelStyle("Hours of Operation"), fieldValueStyle(admin.HoursOfOperation)))
	}

	// Voter Services if they exist
	if len(admin.VoterServices) > 0 {
		sections = append(sections, sectionTitleStyle("Voter Services"))
		sections = append(sections, fieldValueStyle(strings.Join(admin.VoterServices, ", ")))
	}

	// Correspondence Address
	if admin.CorrespondenceAddress != (api.Address{}) {
		sections = append(sections, sectionTitleStyle("Correspondence Address"))
		sections = append(sections, fieldValueStyle(admin.CorrespondenceAddress.String()))
	}

	// Physical Address
	if admin.PhysicalAddress != (api.Address{}) {
		sections = append(sections, sectionTitleStyle("Physical Address"))
		sections = append(sections, fieldValueStyle(admin.PhysicalAddress.String()))
	}

	// Election Officials
	if len(admin.ElectionOfficials) > 0 {
		sections = append(sections, sectionTitleStyle("Election Officials"))
		for _, official := range admin.ElectionOfficials {
			officialInfo := []string{fieldValueStyle(official.Name)}
			if official.Title != "" {
				officialInfo = append(officialInfo, fmt.Sprintf("Title: %s", fieldValueStyle(official.Title)))
			}
			if official.OfficePhoneNumber != "" {
				officialInfo = append(officialInfo, fmt.Sprintf("Office Phone: %s", fieldValueStyle(official.OfficePhoneNumber)))
			}
			if official.EmailAddress != "" {
				officialInfo = append(officialInfo, fmt.Sprintf("Email: %s", fieldValueStyle(official.EmailAddress)))
			}
			sections = append(sections, strings.Join(officialInfo, ", "))
		}
	}

	return strings.Join(sections, "\n")
}

func formatStateResource(state api.State, render *lipgloss.Renderer) string {
	var stateDisplay []string

	// Main header for State
	mainHeaderStyle := render.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("205")).
		Background(lipgloss.Color("63")).
		Padding(0, 1).
		Render
	stateDisplay = append(stateDisplay, mainHeaderStyle(fmt.Sprintf("Register in %s", state.Name)))

	stateDisplay = append(stateDisplay, formatElectionAdministration(state.ElectionAdministrationBody, render))

	if state.LocalJurisdiction != nil {
		stateDisplay = append(stateDisplay, render.NewStyle().Bold(true).Foreground(lipgloss.Color("205")).Render("Local Jurisdiction: "+state.LocalJurisdiction.Name))
		stateDisplay = append(stateDisplay, formatElectionAdministration(state.LocalJurisdiction.ElectionAdministrationBody, render))
	}

	return strings.Join(stateDisplay, "\n\n")
}

func (m model) viewRegister() string {
	if len(m.electionData.State) == 0 {
		return "No registration information available."
	}

	stateInfo := formatStateResource(m.electionData.State[0], m.render)
	return m.render.NewStyle().Margin(1, 2).MaxWidth(m.width).MaxHeight(m.height).Render(
		lipgloss.JoinVertical(
			lipgloss.Top,
			m.HeaderView(),
			stateInfo,
		),
	)
}
