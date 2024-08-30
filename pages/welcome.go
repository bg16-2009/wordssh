package pages

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/ssh"
)

func WelcomeScreen(isNewUser bool, renderer *lipgloss.Renderer, pty ssh.Pty) welcomeModel {
	return welcomeModel{
		renderer:  renderer,
		pty:       pty,
		isNewUser: isNewUser,
	}
}

type welcomeModel struct {
	renderer  *lipgloss.Renderer
	pty       ssh.Pty
	isNewUser bool
}

func (m welcomeModel) Init() tea.Cmd {
	return nil
}

func (m welcomeModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		default:
			return rootScreenModel{}.switchScreen(GameScreen(m.renderer, m.pty))
		}
	}
	return m, nil
}

func (m welcomeModel) View() string {
	if m.isNewUser {
		return "Hello there new user! Welcome to WordSSH\nPress any key to play......"
	} else {
		return "Welcome back to WordSSH\nPress any key to play......"
	}
}
