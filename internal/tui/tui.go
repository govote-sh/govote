package tui

import (
	"fmt"

	huh "github.com/charmbracelet/huh"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/ssh"
	"github.com/charmbracelet/wish/bubbletea"
	handler "github.com/govote-sh/govote/internal/http"
	"github.com/govote-sh/govote/internal/utils"
)

type model struct {
	// Input
	form *huh.Form

	// Style
	style lipgloss.Style

	// State
	state appState

	// Response
	electionData *handler.VoterInfoResponse
	err          error
}

type appState int

const (
	inputState appState = iota
	loadingState
	reinputConfirmationState
	resultState
)

func TeaHandler(s ssh.Session) (tea.Model, []tea.ProgramOption) {
	// Set up the huh form for user input
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().Title("Address").Key("address"),
		),
	)

	r := bubbletea.MakeRenderer(s)
	style := r.NewStyle().
		Border(lipgloss.NormalBorder()).
		Padding(1, 2).
		BorderForeground(lipgloss.Color("#444444")).
		Foreground(lipgloss.Color("#7571F9"))

	// Create the model with the form and style
	m := model{form: form, style: style, state: inputState}
	return m, []tea.ProgramOption{tea.WithAltScreen()}
}

func (m model) Init() tea.Cmd {
	if m.form == nil {
		return nil
	}
	return m.form.Init()
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	// Update the form and handle form completion or exit
	if m.form != nil {
		f, cmd := m.form.Update(msg)
		m.form = f.(*huh.Form)
		cmds = append(cmds, cmd)
	}

	switch m.state {
	case inputState:
		if m.form.State == huh.StateCompleted {
			// Get the user input and switch to loading state
			address := m.form.GetString("address")
			m.state = loadingState

			// Wrap the CheckServer call in a function to make it a tea.Cmd
			return m, func() tea.Msg {
				return handler.CheckServer(address)
			}
		} else if m.form.State == huh.StateAborted {
			return m, tea.Quit
		}

	case loadingState:
		// Handle the server response
		switch msg := msg.(type) {
		case handler.VoterInfoResponse:
			// Save the response and move to the result state
			m.electionData = &msg
			m.state = resultState
			return m, nil

		case utils.ErrMsg:
			// Capture the error and transition to the reinputConfirmationState
			m.err = msg.Err
			m.state = reinputConfirmationState
			return m, nil
		}

	case reinputConfirmationState:
		// Wait for any key press to continue
		if _, ok := msg.(tea.KeyMsg); ok {
			// Reset the form and return to the input state
			m.form = huh.NewForm(
				huh.NewGroup(
					huh.NewInput().Title("Address").Key("address"),
				),
			)
			m.form.Init()
			m.err = nil
			m.state = inputState
			return m, nil
		}

	case resultState:
		// Allow the user to exit by pressing "q" or "ctrl+c"
		if keyMsg, ok := msg.(tea.KeyMsg); ok {
			if keyMsg.String() == "q" || keyMsg.Type == tea.KeyCtrlC {
				return m, tea.Quit
			}
		}
	}

	return m, tea.Batch(cmds...)
}

func (m model) View() string {
	switch m.state {
	case inputState:
		return m.form.View()
	case loadingState:
		return "Loading election information, please wait..."
	case reinputConfirmationState:
		return fmt.Sprintf("Error: %v\nPress any key to continue...", m.err)
	case resultState:
		if m.electionData != nil {
			return fmt.Sprintf("\nElection Day: %s\n\n", m.electionData.Election.ElectionDay)
		}
		return "No election data available."
	}
	return ""
}
