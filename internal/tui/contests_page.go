package tui

import (
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/govote-sh/govote/internal/api"
)

func (m model) InitContestsList() *list.Model {
	items := []list.Item{}
	for _, contest := range m.electionData.Contests {
		items = append(items, list.Item(contest))
	}
	model := list.New(items, list.NewDefaultDelegate(), m.width, m.height)
	return &model
}

func (m model) updateContests(msg tea.Msg) (tea.Model, tea.Cmd) {
	if keyMsg, ok := msg.(tea.KeyMsg); ok {
		switch keyMsg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		case "enter":
			if m.contestsList != nil && m.contestsList.SettingFilter() {
				break
			}
			_, ok := m.contestsList.SelectedItem().(api.Contest)
			if ok {
				m.currPage = contestContentPage
			}
			return m, nil
		}
	}

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

func (m model) viewContests() string {
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
