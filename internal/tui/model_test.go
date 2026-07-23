package tui

import "testing"

func TestNewModelInitialState(t *testing.T) {
	m := newModel(80, 24)

	if m.currPage != inputPage {
		t.Errorf("currPage = %v, want inputPage", m.currPage)
	}
	if m.form == nil {
		t.Error("form is nil, want initialized address form")
	}
	if m.hasMenu {
		t.Error("hasMenu = true, want false before election data arrives")
	}
	if m.width != 80 || m.height != 24 {
		t.Errorf("size = %dx%d, want 80x24", m.width, m.height)
	}
}
