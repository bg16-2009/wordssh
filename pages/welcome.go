package pages

import (
	"github.com/bg16-2009/wordssh/models"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/ssh"
	"gorm.io/gorm"
)

func WelcomeScreen(newUserPublicKey []byte, renderer *lipgloss.Renderer, pty ssh.Pty, db *gorm.DB) welcomeModel {
	ti := textinput.New()
	ti.Placeholder = "Choose a username"
	ti.Focus()
	ti.CharLimit = 16
	ti.Width = 20
	return welcomeModel{
		width:  pty.Window.Width,
		height: pty.Window.Height,

		renderer:         renderer,
		pty:              pty,
		newUserPublicKey: newUserPublicKey,
		db:               db,
		textInput:        ti,
	}
}

type welcomeModel struct {
	width  int
	height int

	renderer         *lipgloss.Renderer
	pty              ssh.Pty
	textInput        textinput.Model
	newUserPublicKey []byte
	db               *gorm.DB
}

func (m welcomeModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m welcomeModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		case "enter":
			m.textInput.Blur()
			m.db.Create(&models.User{
				Username:  m.textInput.Value(),
				PublicKey: m.newUserPublicKey,
				Score:     0,
			})
			return rootScreenModel{}.switchScreen(HomeScreen(m.textInput.Value(), m.renderer, m.pty, m.db))
		}
	}
	var cmd tea.Cmd
	m.textInput, cmd = m.textInput.Update(msg)
	return m, cmd
}

func (m welcomeModel) View() string {
	return m.renderer.Place(
		m.width, m.height, lipgloss.Center, lipgloss.Center,
		lipgloss.JoinVertical(
			lipgloss.Center,
			"Hello there new user!",
			"Welcome to WordSSH\n",
			m.textInput.View(),
		),
	)
}
