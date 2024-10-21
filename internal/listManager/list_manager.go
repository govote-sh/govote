package listManager

import (
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

type ListManager struct {
	lists       []list.Model
	activeIndex int
}

// AddList allows adding a new list with items and a title using builder syntax
func (lm *ListManager) AddList(items []list.Item, title string, width, height int) *ListManager {
	newList := list.New(items, list.NewDefaultDelegate(), width, height)
	newList.Title = title
	lm.lists = append(lm.lists, newList)
	return lm
}

// Init initializes the lists if the ListManager is nil
func InitListManager(data [][]list.Item, titles []string, width, height int) *ListManager {
	lm := &ListManager{}
	for i := 0; i < len(data); i++ {
		lm.AddList(data[i], titles[i], width, height)
	}
	return lm
}

// SetSize updates the width and height for all lists
func (lm *ListManager) SetSize(width, height int) {
	for i := range lm.lists {
		lm.lists[i].SetWidth(width)
		lm.lists[i].SetHeight(height)
	}
}

// ActiveList returns the current active list
func (lm *ListManager) ActiveList() list.Model {
	return lm.lists[lm.activeIndex]
}

// CycleNext cycles to the next list
func (lm *ListManager) CycleNext() {
	lm.activeIndex = (lm.activeIndex + 1) % len(lm.lists)
}

// CyclePrev cycles to the previous list
func (lm *ListManager) CyclePrev() {
	lm.activeIndex = (lm.activeIndex - 1 + len(lm.lists)) % len(lm.lists)
}

// UpdateActiveList updates the active list based on tea.Msg
func (lm *ListManager) UpdateActiveList(msg tea.Msg) (*ListManager, tea.Cmd) {
	var cmd tea.Cmd
	lm.lists[lm.activeIndex], cmd = lm.lists[lm.activeIndex].Update(msg)
	return lm, cmd
}

func (lm *ListManager) SettingFilter() bool {
	return lm.lists[lm.activeIndex].SettingFilter()
}

func (lm *ListManager) IsFiltered() bool {
	return lm.lists[lm.activeIndex].IsFiltered()
}

func (lm *ListManager) SelectedItem() list.Item {
	if lm.activeIndex < 0 || lm.activeIndex >= len(lm.lists) {
		return nil
	}
	item := lm.lists[lm.activeIndex].SelectedItem()
	return item
}
