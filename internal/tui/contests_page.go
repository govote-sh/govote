package tui

import (
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

func (m model) InitContestsList() *list.Model {
	items := []list.Item{}
	for _, contest := range m.electionData.Contests {
		items = append(items, list.Item(contest))
	}
	model := list.New(items, list.NewDefaultDelegate(), m.width, m.height)
	return &model
}

func (m model) UpdateContests(msg tea.Msg) (tea.Model, tea.Cmd) {
	if m.contestsList == nil {
		return m, nil
	}

	var cmd tea.Cmd
	contestsList, cmd := m.contestsList.Update(msg)
	m.contestsList = &contestsList
	if cmd != nil {
		return m, cmd
	}

	return m, nil
}

func (m model) ViewContests() string {
	if m.contestsList == nil || len(m.electionData.Contests) == 0 {
		return m.render.NewStyle().Margin(1, 1).MaxWidth(m.width).MaxHeight(m.height).Render(
			lipgloss.JoinVertical(
				lipgloss.Top,
				m.HeaderView(),
				m.RenderErrorBox("No contests available..."),
			),
		)
	}
	return m.render.NewStyle().Margin(1, 1).MaxWidth(m.width).MaxHeight(m.height).Render(lipgloss.JoinVertical(
		lipgloss.Top,
		m.HeaderView(),
		m.contestsList.View(),
	))
}
