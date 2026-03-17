package tui

import (
	"charm.land/bubbles/v2/list"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/govote-sh/govote/internal/api"
)

func (m model) InitContestsList() *list.Model {
	items := []list.Item{}
	for _, contest := range m.electionData.Contests {
		items = append(items, list.Item(contest))
	}
	model := list.New(items, list.NewDefaultDelegate(), m.width, m.height)
	model.SetHeight(m.height - 4)
	return &model
}

func (m model) updateContests(msg tea.Msg) (tea.Model, tea.Cmd) {
	if m.contestsList != nil {
		var cmd tea.Cmd
		contestsList, cmd := m.contestsList.Update(msg)
		m.contestsList = &contestsList
		if cmd != nil {
			return m, cmd
		}
	}

	if keyMsg, ok := msg.(tea.KeyPressMsg); ok {
		switch keyMsg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		case "enter":
			_, ok := m.contestsList.SelectedItem().(api.Contest)
			if ok {
				m.currPage = contestContentPage
			}
			return m, nil
		}
	}
	return m, nil
}

func (m model) viewContests() string {
	if m.contestsList == nil || len(m.electionData.Contests) == 0 {
		return lipgloss.NewStyle().Margin(1, 1).MaxWidth(m.width).MaxHeight(m.height).Render(
			lipgloss.JoinVertical(
				lipgloss.Top,
				m.HeaderView(),
				m.RenderErrorBox("No contests available..."),
			),
		)
	}
	return lipgloss.NewStyle().Margin(1, 1).MaxWidth(m.width).MaxHeight(m.height).Render(lipgloss.JoinVertical(
		lipgloss.Top,
		m.HeaderView(),
		m.contestsList.View(),
	))
}
