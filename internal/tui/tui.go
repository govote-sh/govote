package tui

import (
	"fmt"

	spinner "github.com/charmbracelet/bubbles/spinner"
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

	// Style & Bubbles
	style   lipgloss.Style
	spinner spinner.Model

	// State
	state appState

	// Response
	electionData *handler.VoterInfoResponse
	err          error

	// Header and subtitle styles
	headerStyle   lipgloss.Style
	subtitleStyle lipgloss.Style
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

	// Define the styles for the header and subtitle
	headerStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("205")).
		Align(lipgloss.Center).
		Bold(true).
		Padding(0, 1)

	subtitleStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("240")).
		Align(lipgloss.Center).
		Padding(0, 1)

	r := bubbletea.MakeRenderer(s)
	style := r.NewStyle().
		Border(lipgloss.NormalBorder()).
		Padding(1, 2).
		BorderForeground(lipgloss.Color("#444444")).
		Foreground(lipgloss.Color("#7571F9"))

	spin := spinner.New()
	spin.Spinner = spinner.Dot
	spin.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))

	// Create the model with the form, style, and spinner
	m := model{
		form:          form,
		style:         style,
		spinner:       spin,
		state:         inputState,
		headerStyle:   headerStyle,   // Assign the header style
		subtitleStyle: subtitleStyle, // Assign the subtitle style
	}
	return m, []tea.ProgramOption{tea.WithAltScreen()}
}

func (m model) Init() tea.Cmd {
	if m.form == nil {
		return nil
	}
	return tea.Batch(
		m.form.Init(),  // Initialize the form
		m.spinner.Tick, // Pass the command to start the spinner ticking
	)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch m.state {
	case inputState:
		if m.form != nil {
			f, cmd := m.form.Update(msg)
			m.form = f.(*huh.Form)
			cmds = append(cmds, cmd)
		}
		if m.form.State == huh.StateCompleted {
			// Get the user input and switch to loading state
			address := m.form.GetString("address")
			m.state = loadingState

			// Return the CheckServer call as a tea.Cmd
			return m, tea.Batch(
				m.spinner.Tick, // Start the spinner ticking
				func() tea.Msg {
					return handler.CheckServer(address)
				},
			)
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
		case spinner.TickMsg:
			var cmd tea.Cmd
			m.spinner, cmd = m.spinner.Update(msg)
			return m, cmd
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
		return m.viewInput()
	case loadingState:
		return fmt.Sprintf("%s Loading election information, please wait...\n\n", m.spinner.View())
	case reinputConfirmationState:
		return fmt.Sprintf("Error: %v\nPress any key to continue...", m.err)
	case resultState:
		return m.viewResult()
	}
	return ""
}

func (m model) viewInput() string {
	header := m.headerStyle.Render("Welcome to govote.sh!")
	subtitle := m.subtitleStyle.Render("Please enter your address to get election information from the Voting Information Project")
	return fmt.Sprintf("%s\n%s\n\n%s", header, subtitle, m.form.View())
}

func (m model) viewResult() string {
	headerText := fmt.Sprintf("Upcoming %s on %s", m.electionData.Election.Name, m.electionData.Election.ElectionDay)
	header := m.headerStyle.Render(headerText)
	subtitleText := fmt.Sprintf("Results for: %s", m.electionData.NormalizedInput.String())
	subtitle := m.headerStyle.Render(subtitleText)
	if m.electionData != nil {
		return fmt.Sprintf("%s\n%s\nUpcoming election: %s on %s\n\n", header, subtitle, m.electionData.Election.Name, m.electionData.Election.ElectionDay)
	}
	return "No election data available." // Transition to Reinput confirmation state? or catch error?
}
