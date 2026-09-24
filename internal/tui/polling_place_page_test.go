package tui

import (
	"strings"
	"testing"

	"github.com/charmbracelet/x/ansi"
)

func TestFormatPollingDates(t *testing.T) {
	tests := []struct {
		name, start, end, want string
	}{
		{"neither", "", "", ""},
		{"same day", "2026-11-03", "2026-11-03", "Date: 2026-11-03"},
		{"range", "2026-10-20", "2026-11-01", "Available Dates: 2026-10-20 → 2026-11-01"},
		{"start only", "2026-10-20", "", "Available From: 2026-10-20"},
		{"end only", "", "2026-11-01", "Available Until: 2026-11-01"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ansi.Strip(formatPollingDates(tt.start, tt.end)); got != tt.want {
				t.Errorf("formatPollingDates(%q, %q) = %q, want %q", tt.start, tt.end, got, tt.want)
			}
		})
	}
}

// The range values used to skip fieldValueStyle, unlike every other field.
func TestPollingDateRangeValuesAreStyled(t *testing.T) {
	got := formatPollingDates("2026-10-20", "2026-11-01")
	if want := fieldValueStyle("2026-10-20"); !strings.Contains(got, want) {
		t.Errorf("start date not rendered with fieldValueStyle: %q", got)
	}
}
