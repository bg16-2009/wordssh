package pages

import (
	"errors"
	"fmt"

	"github.com/bg16-2009/wordssh/models"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/ssh"
	"gorm.io/gorm"
)

func SearchScreen(username string, renderer *lipgloss.Renderer, pty ssh.Pty, db *gorm.DB) searchModel {
	ti := textinput.New()
	ti.Placeholder = "Enter a username"
	ti.Focus()
	ti.CharLimit = 16
	ti.Width = 20
	return searchModel{
        width: pty.Window.Width,
        height: pty.Window.Height,

		textInput:      ti,
		renderer:       renderer,
		pty:            pty,
		username:       username,
		db:             db,
		searchComplete: false,
	}
}

type searchModel struct {
	textInput textinput.Model
	width     int
	height    int

	renderer       *lipgloss.Renderer
	pty            ssh.Pty
	username       string
	db             *gorm.DB
	searchComplete bool
	message        string
}

func (m searchModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m searchModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() != "q" && msg.String() != "ctrl+c" && m.searchComplete {
			return rootScreenModel{}.switchScreen(HomeScreen(m.username, m.renderer, m.pty, m.db))
		}
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		case "enter":
			var user models.User
            m.textInput.Blur()
			result := m.db.First(&user, "username = ?", m.textInput.Value())
			if errors.Is(result.Error, gorm.ErrRecordNotFound) {
				m.message = lipgloss.JoinVertical(
					lipgloss.Center,
					"User not found",
                    m.renderer.NewStyle().Foreground(lipgloss.Color("8")).Render("\nPress any key to go back"),
				)
			}else{
				m.message = lipgloss.JoinVertical(
					lipgloss.Center,
					fmt.Sprintf("User %s has %d points.", user.Username, user.Score),
                    m.renderer.NewStyle().Foreground(lipgloss.Color("8")).Render("\nPress any key to go back"),
				)
            }
            m.searchComplete = true
		}
	}
	var cmd tea.Cmd
	m.textInput, cmd = m.textInput.Update(msg)
	return m, cmd
}

func (m searchModel) View() string {
	return m.renderer.Place(
		m.width, m.height, lipgloss.Center, lipgloss.Center,
		lipgloss.JoinVertical(
			lipgloss.Center,
			m.textInput.View(),
            m.message,
		),
	)
}
