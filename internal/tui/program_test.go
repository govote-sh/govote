package tui

import (
	"testing"
	"time"

	tea "charm.land/bubbletea/v2"
	teatest "github.com/charmbracelet/x/exp/teatest/v2"
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
