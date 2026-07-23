package tui

import (
	"bytes"
	"testing"

	tea "charm.land/bubbletea/v2"
	"github.com/charmbracelet/colorprofile"
	teatest "github.com/charmbracelet/x/exp/teatest/v2"
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
