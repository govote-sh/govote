package tui

import (
	"bytes"
	"testing"

	tea "charm.land/bubbletea/v2"
	"github.com/charmbracelet/colorprofile"
	teatest "github.com/charmbracelet/x/exp/teatest/v2"
	"github.com/govote-sh/govote/internal/api"
)

// newTestProgram runs m as a full tea.Program against an in-memory 80x24
// terminal. The color profile is pinned so local and CI output are identical.
func newTestProgram(t *testing.T, m model) *teatest.TestModel {
	t.Helper()
	tm := teatest.NewTestModel(t, m,
		teatest.WithInitialTermSize(80, 24),
		teatest.WithProgramOptions(tea.WithColorProfile(colorprofile.ANSI256)),
	)
	t.Cleanup(func() { _ = tm.Quit() })
	return tm
}

// containsBytes adapts a substring check to teatest.WaitFor's condition func.
func containsBytes(s string) func([]byte) bool {
	return func(b []byte) bool { return bytes.Contains(b, []byte(s)) }
}

// fixtureVoterInfo is a minimal but page-complete Civic API response: one
// polling place, one contest with a candidate, one state with registration
// info — enough to render every page.
func fixtureVoterInfo() api.VoterInfoResponse {
	return api.VoterInfoResponse{
		Election: api.Election{
			ID:          "2000",
			Name:        "Test General Election",
			ElectionDay: "2026-11-03",
		},
		PollingLocations: []api.PollingPlace{{
			Address: api.Address{
				LocationName: "Main St Community Center",
				Line1:        "100 Main St",
				City:         "Richmond",
				State:        "VA",
				Zip:          "23220",
			},
			PollingHours: "Tuesday: 6:00 AM - 7:00 PM",
			StartDate:    "2026-11-03",
			EndDate:      "2026-11-03",
		}},
		Contests: []api.Contest{{
			Type:        "General",
			BallotTitle: "Governor",
			Office:      "Governor",
			Candidates: []api.Candidate{
				{Name: "Alex Doe", Party: "Independent"},
			},
		}},
		State: []api.State{{
			Name: "Test State",
			ElectionAdministrationBody: api.ElectionAdministrationBody{
				Name:                    "Test State Board of Elections",
				ElectionRegistrationUrl: "https://vote.example.gov/register",
			},
		}},
	}
}

// newVotePageModel builds a model as if a lookup just succeeded — it must
// mirror the api.VoterInfoResponse branch in Update (tui.go): electionData,
// hasMenu, both lists, currPage. If that branch changes, change this too.
func newVotePageModel(width, height int) model {
	m := newModel(width, height)
	data := fixtureVoterInfo()
	m.electionData = &data
	m.hasMenu = true
	m.currPage = votePage
	m.lm = m.InitVotePageListManager()
	m.contestsList = m.InitContestsList()
	return m
}
