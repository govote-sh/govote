package tui

import (
	"testing"
	"time"

	tea "charm.land/bubbletea/v2"
	teatest "github.com/charmbracelet/x/exp/teatest/v2"
	"github.com/govote-sh/govote/internal/secrets"
)

func TestProgramBootsAndQuits(t *testing.T) {
	tm := newTestProgram(t, newModel(80, 24))

	teatest.WaitFor(t, tm.Output(), containsBytes("Welcome to govote.sh!"),
		teatest.WithDuration(3*time.Second))

	// The form accepts typing.
	tm.Type("1234 W Broad St")
	// The cell-diff renderer only emits changed cells per keystroke, so the
	// typed text never appears contiguously in the stream. A resize forces a
	// full repaint of the settled frame.
	tm.Send(tea.WindowSizeMsg{Width: 100, Height: 30})
	teatest.WaitFor(t, tm.Output(), containsBytes("1234 W Broad St"),
		teatest.WithDuration(3*time.Second))

	// ctrl+c aborts the huh form, which returns tea.Quit.
	tm.Send(tea.KeyPressMsg{Code: 'c', Mod: tea.ModCtrl})
	tm.WaitFinished(t, teatest.WithFinalTimeout(2*time.Second))
}

// pressEnter advances the huh form one field (submits on the last field).
func pressEnter(tm *teatest.TestModel) {
	tm.Send(tea.KeyPressMsg{Code: tea.KeyEnter})
}

func TestEmptySubmitShowsErrorThenRecovers(t *testing.T) {
	tm := newTestProgram(t, newModel(80, 24))
	teatest.WaitFor(t, tm.Output(), containsBytes("Welcome to govote.sh!"),
		teatest.WithDuration(3*time.Second))

	// Submit all four fields empty: street, city, state, postal code.
	for range 4 {
		pressEnter(tm)
	}

	teatest.WaitFor(t, tm.Output(),
		containsBytes("at least one address field is required"),
		teatest.WithDuration(3*time.Second))

	// Any key returns to a fresh form.
	tm.Type("x")
	teatest.WaitFor(t, tm.Output(), containsBytes("Welcome to govote.sh!"),
		teatest.WithDuration(3*time.Second))
}

func TestSubmitTransportErrorShowsFriendlyMessage(t *testing.T) {
	// CheckServer needs a key to get as far as the HTTP request, and the
	// dead proxy makes that request fail at connect time — no network I/O.
	// Gotcha: net/http caches proxy env per process on first use, so every
	// test in this binary that triggers HTTP must use this same proxy value.
	t.Setenv("API_KEY", "test-key-not-real")
	t.Setenv("HTTPS_PROXY", "http://127.0.0.1:9")
	if err := secrets.SetupSecrets(); err != nil {
		t.Fatalf("SetupSecrets: %v", err)
	}

	tm := newTestProgram(t, newModel(80, 24))
	teatest.WaitFor(t, tm.Output(), containsBytes("Welcome to govote.sh!"),
		teatest.WithDuration(3*time.Second))

	tm.Type("1234 W Broad St")
	pressEnter(tm) // -> city
	tm.Type("Richmond")
	pressEnter(tm) // -> state
	tm.Type("VA")
	pressEnter(tm) // -> postal code
	pressEnter(tm) // submit (postal code empty is allowed)

	// The user-facing message must be the generic one — never the raw
	// *url.Error (which would embed the request URL).
	teatest.WaitFor(t, tm.Output(),
		containsBytes("could not reach the election information service"),
		teatest.WithDuration(3*time.Second))
}
