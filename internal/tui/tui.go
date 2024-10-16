package tui

import (
	"fmt"

	spinner "github.com/charmbracelet/bubbles/spinner"
	huh "github.com/charmbracelet/huh"

	"github.com/charmbracelet/bubbles/list"
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
	render  *lipgloss.Renderer
	spinner spinner.Model

	// Page
	page page

	// Response
	electionData *handler.VoterInfoResponse
	err          error

	// Header and subtitle styles
	headerStyle   lipgloss.Style
	subtitleStyle lipgloss.Style

	// Lists
	pollingLocationList        list.Model
	pollingLocationListCreated bool

	hasMenu bool

	// Track window size
	width, height int
}

type page int

const (
	inputPage page = iota
	loadingPage
	reinputConfirmationPage
	pollingLocationPage
	earlyVotePage
	ballotDropOffPage
	contestsPage
	registerPage
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
		form:                       form,
		spinner:                    spin,
		page:                       inputPage,
		headerStyle:                headerStyle,   // Assign the header style
		subtitleStyle:              subtitleStyle, // Assign the subtitle style
		width:                      pty.Window.Width,
		height:                     pty.Window.Height,
		pollingLocationListCreated: false,
		render:                     r,
		hasMenu:                    false,
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
		if m.pollingLocationListCreated {
			m.pollingLocationList.SetWidth(m.width)
			m.pollingLocationList.SetHeight(m.height)
		} else if m.page == pollingLocationPage {
			// Initialize the list if not done yet
			m.initList(m.width, m.height)
		}
		return m, nil
	}

	switch m.page {
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
			m.page = loadingPage

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

	case loadingPage:
		// Handle the server response
		switch msg := msg.(type) {
		case handler.VoterInfoResponse:
			// Save the response and move to the pollingLocationPage state
			m.electionData = &msg
			m.page = pollingLocationPage
			m.hasMenu = true

			// Initialize the list if window size information is available
			if m.width != 0 && m.height != 0 {
				m.initList(m.width, m.height)
			}
			return m, nil

		case spinner.TickMsg:
			var cmd tea.Cmd
			m.spinner, cmd = m.spinner.Update(msg)
			cmds = append(cmds, cmd)
			return m, tea.Batch(cmds...)

		case utils.ErrMsg:
			// Capture the error and transition to reinputConfirmationState
			m.err = msg.Err
			m.page = reinputConfirmationPage
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
			m.page = inputPage
			return m, nil
		}

	case pollingLocationPage:
		// Handle list updates
		if m.pollingLocationListCreated {
			var cmd tea.Cmd
			m.pollingLocationList, cmd = m.pollingLocationList.Update(msg)
			cmds = append(cmds, cmd)
		}

		// Allow the user to exit by pressing "q" or "ctrl+c"
		if keyMsg, ok := msg.(tea.KeyMsg); ok {
			if keyMsg.String() == "q" || keyMsg.Type == tea.KeyCtrlC {
				return m, tea.Quit
			}
		}
	case contestsPage:
		// Handle list updates
	}

	return m, tea.Batch(cmds...)
}

func (m model) View() string {
	switch m.page {
	case inputPage:
		return m.viewInput()
	case loadingPage:
		return fmt.Sprintf("%s Loading election information, please wait...\n\n", m.spinner.View())
	case reinputConfirmationPage:
		return fmt.Sprintf("Error: %v\nPress any key to continue...", m.err)
	case pollingLocationPage:
		return m.viewResult()
	case contestsPage:
		return "Contests"
	case registerPage:
		return "Register"
	}
	return ""
}

func (m model) viewInput() string {
	header := m.headerStyle.Render("Welcome to govote.sh!")
	subtitle := m.subtitleStyle.Render("Please enter your address to get election information from the Voting Information Project")
	return fmt.Sprintf("%s\n%s\n\n%s", header, subtitle, m.form.View())
}

func (m model) viewResult() string {
	// headerText := fmt.Sprintf("Upcoming %s on %s", m.electionData.Election.Name, m.electionData.Election.ElectionDay)
	// header := m.headerStyle.Render(headerText)
	// subtitleText := fmt.Sprintf("Results for: %s", m.electionData.NormalizedInput.String())
	// subtitle := m.headerStyle.MarginBottom(1).Render(subtitleText)
	if !m.pollingLocationListCreated {
		return "building list..."
	}
	return m.render.NewStyle().Margin(1, 1).MaxWidth(m.width).MaxHeight(m.height).Render(lipgloss.JoinVertical(
		lipgloss.Top,
		m.HeaderView(),
		m.pollingLocationList.View(),
	))
}

// func (m model) viewResult() string {
// 	// Prepare the header and subtitle with updated styling
// 	headerText := fmt.Sprintf("Upcoming %s on %s", m.electionData.Election.Name, m.electionData.Election.ElectionDay)
// 	header := m.headerStyle.Bold(true).Foreground(lipgloss.Color("212")).Render(headerText)

// 	subtitleText := fmt.Sprintf("Results for: %s", m.electionData.NormalizedInput.String())
// 	subtitle := m.subtitleStyle.Italic(true).Foreground(lipgloss.Color("240")).MarginBottom(2).Render(subtitleText)

// 	// Ensure the list is ready, otherwise show a loading message
// 	if !m.pollingLocationListCreated {
// 		return m.render.NewStyle().Margin(2, 2).Render("Building list...")
// 	}

// 	// Combine the header, subtitle, and list with adequate padding/margin
// 	return m.render.NewStyle().
// 		Margin(0, 0).
// 		MaxWidth(m.width).
// 		MaxHeight(m.height).
// 		Render(lipgloss.JoinVertical(
// 			lipgloss.Top,
// 			header,
// 			subtitle,
// 			m.pollingLocationList.View(),
// 		))
// }

func (m *model) initList(width, height int) {
	if m.electionData == nil {
		fmt.Println("electionData is nil")
		return
	}

	// Check if PollingLocations is nil or empty
	if m.electionData.PollingLocations == nil || len(m.electionData.PollingLocations) == 0 {
		fmt.Println("PollingLocations is nil or empty")
		return
	}

	// Convert []PollingPlace to []list.Item
	var items []list.Item
	for _, pollingPlace := range m.electionData.EarlyVoteSites {
		items = append(items, pollingPlace)
	}

	m.pollingLocationList = list.New(items, list.NewDefaultDelegate(), width, height)
	m.pollingLocationList.Title = "Polling Locations"
	m.pollingLocationListCreated = true
}
