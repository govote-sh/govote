package tui

import (
	"strings"
	"testing"

	"github.com/charmbracelet/x/ansi"
	"github.com/govote-sh/govote/internal/api"
)

// Real Civic API payloads (e.g. Virginia) often list officials with no name,
// which used to render as ", Office Phone: ...".
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
