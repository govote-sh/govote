package tui

import (
	"strings"
	"testing"

	"github.com/charmbracelet/x/exp/golden"
	"github.com/govote-sh/govote/internal/utils"
)

// requireGoldenView snapshots a page render. View is a pure function of the
// model, so this is deterministic — no program, no timing. ESC bytes are
// escaped to "\x1b" before writing so the golden files stay diffable text
// (upstream's RequireEqualEscape no longer escapes; it is a deprecated
// forwarder to RequireEqual, which writes raw bytes).
// Re-record with: go test ./internal/tui/ -update
func requireGoldenView(t *testing.T, m model) {
	t.Helper()
	escaped := strings.ReplaceAll(m.View().Content, "\x1b", `\x1b`)
	golden.RequireEqual(t, []byte(escaped))
}

func TestGoldenVotePage(t *testing.T) {
	requireGoldenView(t, newVotePageModel(80, 24))
}

func TestGoldenContestsPage(t *testing.T) {
	m := newVotePageModel(80, 24)
	m.currPage = contestsPage
	requireGoldenView(t, m)
}

func TestGoldenContestDetailPage(t *testing.T) {
	m := newVotePageModel(80, 24)
	m.currPage = contestContentPage
	requireGoldenView(t, m)
}

func TestGoldenRegisterPage(t *testing.T) {
	m := newVotePageModel(80, 24)
	m.currPage = registerPage
	requireGoldenView(t, m)
}

func TestGoldenPollingPlacePage(t *testing.T) {
	m := newVotePageModel(80, 24)
	m.currPage = pollingPlacePage
	requireGoldenView(t, m)
}

func TestGoldenErrorPage(t *testing.T) {
	m := newModel(80, 24)
	m.currPage = reinputConfirmationPage
	m.err = &utils.ErrMsg{HTTPStatusCode: 400}
	requireGoldenView(t, m)
}
