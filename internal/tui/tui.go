package tui

import (
	"fmt"

	"github.com/charmbracelet/bubbles/list"
	spinner "github.com/charmbracelet/bubbles/spinner"
	huh "github.com/charmbracelet/huh"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/ssh"
	"github.com/charmbracelet/wish/bubbletea"
	"github.com/govote-sh/govote/internal/api"
	"github.com/govote-sh/govote/internal/listManager"
	"github.com/govote-sh/govote/internal/utils"
)

type model struct {
	// Input
	form *huh.Form

	// Style & Bubbles
	render  *lipgloss.Renderer
	spinner spinner.Model

	// Page
	currPage page

	// Response
	electionData *api.VoterInfoResponse
	err          *utils.ErrMsg

	// Header and subtitle styles
	headerStyle   lipgloss.Style
	subtitleStyle lipgloss.Style

	// Lists
	selectedPollingPlace *api.PollingPlace
	lm                   *listManager.ListManager // List manager for the vote page
	contestsList         *list.Model

	hasMenu bool

	// Track window size
	width, height int
}

type page int

const (
	inputPage page = iota
	loadingPage
	reinputConfirmationPage
	votePage
	// earlyVotePage
	// ballotDropOffPage
	contestsPage
	registerPage
	pollingPlacePage
)

func TeaHandler(s ssh.Session) (tea.Model, []tea.ProgramOption) {
	pty, _, _ := s.Pty()

	// Set up the huh form for user input
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().Title("Address").Key("address"),
		),
	)

	r := bubbletea.MakeRenderer(s)

	// Define the styles for the header and subtitle
	headerStyle := r.NewStyle().
		Foreground(lipgloss.Color("205")).
		Align(lipgloss.Center).
		Bold(true).
		Padding(0, 1)

	subtitleStyle := r.NewStyle().
		Foreground(lipgloss.Color("240")).
		Align(lipgloss.Center).
		Padding(0, 1)

	spin := spinner.New()
	spin.Spinner = spinner.Dot
	spin.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))

	// Create the model with the form, style, and spinner
	m := model{
		form:          form,
		spinner:       spin,
		currPage:      inputPage,
		headerStyle:   headerStyle,   // Assign the header style
		subtitleStyle: subtitleStyle, // Assign the subtitle style
		width:         pty.Window.Width,
		height:        pty.Window.Height,
		render:        r,
		hasMenu:       false,
	}
	return m, []tea.ProgramOption{tea.WithAltScreen()}
}

func (m model) Init() tea.Cmd {
	if m.form == nil {
		return nil
	}
	return m.form.Init()
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var headerCmd tea.Cmd
	m, headerCmd = m.HeaderUpdate(msg)
	cmds := []tea.Cmd{headerCmd}

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		// Capture the window size
		m.width = msg.Width
		m.height = msg.Height

		// If the list is created, adjust its size accordingly
		if m.lm != nil {
			m.lm.SetSize(m.width, m.height)
		}
		return m, nil
	}

	switch m.currPage {
	case inputPage:
		if m.form != nil {
			// Update the form and handle form completion or exit
			f, cmd := m.form.Update(msg)
			m.form = f.(*huh.Form)
			cmds = append(cmds, cmd)
		}
		if m.form.State == huh.StateCompleted {
			// Get the user input and switch to loading state
			address := m.form.GetString("address")
			m.currPage = loadingPage

			// Return the CheckServer call as a tea.Cmd
			return m, tea.Batch(
				m.spinner.Tick, // Start the spinner ticking
				func() tea.Msg {
					return api.CheckServer(address)
				},
			)
		} else if m.form.State == huh.StateAborted {
			return m, tea.Quit
		}

	case loadingPage:
		// Handle the server response
		switch msg := msg.(type) {
		case api.VoterInfoResponse:
			// Save the response and move to the votePage
			m.electionData = &msg
			m.currPage = votePage
			m.hasMenu = true
			m.lm = m.InitVotePageListManager()
			m.contestsList = m.InitContestsList()

			return m, nil

		case spinner.TickMsg:
			var cmd tea.Cmd
			m.spinner, cmd = m.spinner.Update(msg)
			cmds = append(cmds, cmd)
			return m, tea.Batch(cmds...)

		case utils.ErrMsg:
			// Capture the error and transition to reinputConfirmationState
			m.err = &msg
			m.currPage = reinputConfirmationPage
			return m, nil
		}
	case reinputConfirmationPage:
		// Wait for any key press to continue
		if _, ok := msg.(tea.KeyMsg); ok {
			// Reset the form and return to input state
			m.form = huh.NewForm(
				huh.NewGroup(
					huh.NewInput().Title("Address").Key("address"),
				),
			)
			m.form.Init()
			m.err = nil
			m.currPage = inputPage
			return m, nil
		}
	case votePage:
		return m.UpdateVote(msg)
	case contestsPage:
		return m.UpdateContests(msg)
	case registerPage:
		return m, nil
	case pollingPlacePage:
		return m.updatePollingPlace(msg)
	}

	return m, tea.Batch(cmds...)
}

func (m model) View() string {
	switch m.currPage {
	case inputPage:
		return m.viewInput()
	case loadingPage:
		return fmt.Sprintf("%s Loading election information, please wait...\n\n", m.spinner.View())
	case reinputConfirmationPage:
		return m.viewReinputConfirmation()
	case votePage:
		return m.viewVote()
	case contestsPage:
		return m.ViewContests()
	case registerPage:
		return m.viewRegister()
	case pollingPlacePage:
		return m.viewPollingPlace()
	}
	return ""
}

func (m model) viewReinputConfirmation() string {
	var errorMsg string
	if m.err == nil {
		errorMsg = "Error: unknown error"
	} else if m.err.HTTPStatusCode >= 400 && m.err.HTTPStatusCode < 500 { // Client error
		errorMsg = fmt.Sprintf("Error: Client error (code: %d): This is likely due to an invalid address\nor the voter information project not being up to date\nPlease check https://all.votinginfotool.org", m.err.HTTPStatusCode)
	} else if m.err.HTTPStatusCode >= 500 && m.err.HTTPStatusCode < 600 { // Server error
		errorMsg = fmt.Sprintf("Error: Server error (code: %d): This is likely due to the API being down\nPlease check https://all.votinginfotool.org to make sure", m.err.HTTPStatusCode)
	} else {
		errorMsg = fmt.Sprintf("Error: %v", m.err.Err.Error())
	}

	return m.render.NewStyle().Margin(1, 1).MaxWidth(m.width).MaxHeight(m.height).Render(lipgloss.JoinVertical(
		lipgloss.Top,
		errorMsg,
		"Press any key to continue...",
	))
}

func (m model) viewInput() string {
	header := m.headerStyle.Render("Welcome to govote.sh!")
	subtitle := m.subtitleStyle.Render("Please enter your address to get election information from the Voting Information Project")
	return fmt.Sprintf("%s\n%s\n\n%s", header, subtitle, m.form.View())
}

// func (m model) viewContests() string {
// 	return m.render.NewStyle().Margin(1, 1).MaxWidth(m.width).MaxHeight(m.height).Render(lipgloss.JoinVertical(
// 		lipgloss.Top,
// 		m.HeaderView(),
// 		"Contests",
// 	))
// }
