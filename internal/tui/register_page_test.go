package tui

import (
	"strings"
	"testing"

	"github.com/charmbracelet/x/ansi"
	"github.com/govote-sh/govote/internal/api"
)

// Seen in live Virginia data.
func TestElectionOfficialWithoutNameHasNoLeadingComma(t *testing.T) {
	admin := api.ElectionAdministrationBody{
		Name: "Virginia Department of Elections",
		ElectionOfficials: []api.ElectionOfficial{{
			OfficePhoneNumber: "800-552-9745",
			EmailAddress:      "info@elections.virginia.gov",
		}},
	}

	out := ansi.Strip(formatElectionAdministration(admin))

	want := "Office Phone: 800-552-9745, Email: info@elections.virginia.gov"
	for line := range strings.SplitSeq(out, "\n") {
		if strings.Contains(line, "Office Phone") {
			if line != want {
				t.Fatalf("official line = %q, want %q", line, want)
			}
			return
		}
	}
	t.Fatalf("no official line in output:\n%s", out)
}

func assertNoBlankLines(t *testing.T, out string) {
	t.Helper()
	for i, line := range strings.Split(out, "\n") {
		if strings.TrimSpace(line) == "" {
			t.Fatalf("line %d is blank in output:\n%s", i, out)
		}
	}
}

// Delaware's state admin body arrives with "name": null.
func TestElectionAdministrationWithoutNameHasNoBlankLine(t *testing.T) {
	admin := api.ElectionAdministrationBody{ElectionInfoUrl: "https://elections.delaware.gov/"}

	assertNoBlankLines(t, ansi.Strip(formatElectionAdministration(admin)))
}

func TestEmptyElectionOfficialHasNoBlankLine(t *testing.T) {
	admin := api.ElectionAdministrationBody{
		Name:              "Delaware Department of Elections",
		ElectionOfficials: []api.ElectionOfficial{{}},
	}

	out := ansi.Strip(formatElectionAdministration(admin))

	assertNoBlankLines(t, out)
	if strings.Contains(out, "Election Officials") {
		t.Errorf("header rendered for an official with no details:\n%s", out)
	}
}

func TestLocalJurisdictionWithoutNameHasNoDanglingColon(t *testing.T) {
	state := api.State{
		Name: "Delaware",
		LocalJurisdiction: &api.AdministrationRegion{
			ElectionAdministrationBody: api.ElectionAdministrationBody{Name: "New Castle County"},
		},
	}

	out := ansi.Strip(formatStateResource(state))

	for line := range strings.SplitSeq(out, "\n") {
		if strings.HasPrefix(line, "Local Jurisdiction") && line != "Local Jurisdiction" {
			t.Fatalf("jurisdiction header = %q, want %q", line, "Local Jurisdiction")
		}
	}
}
