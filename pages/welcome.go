package pages

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/ssh"
)

func WelcomeScreen(renderer *lipgloss.Renderer, pty ssh.Pty) welcomeModel {
	return welcomeModel{
		renderer: renderer,
		pty:      pty,
	}
}

type welcomeModel struct {
	renderer *lipgloss.Renderer
	pty      ssh.Pty
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
	return "Welcome to WordSSH\nPress any key to play......"
}
