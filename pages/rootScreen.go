package pages

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/ssh"
)

type baseScreen struct {
	renderer *lipgloss.Renderer
	pty      ssh.Pty
}

func RootScreen(isNewUser bool, renderer *lipgloss.Renderer, pty ssh.Pty) rootScreenModel {
	return rootScreenModel{
		renderer:      renderer,
		pty:           pty,
		currentScreen: WelcomeScreen(isNewUser, renderer, pty),
	}
}

type rootScreenModel struct {
	currentScreen tea.Model
	renderer      *lipgloss.Renderer
	pty           ssh.Pty
}

func (m rootScreenModel) Init() tea.Cmd {
	return m.currentScreen.Init()
}

func (m rootScreenModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return m.currentScreen.Update(msg)
}

func (m rootScreenModel) View() string {
	return m.currentScreen.View()
}

func (m rootScreenModel) switchScreen(model tea.Model) (tea.Model, tea.Cmd) {
	return model, model.Init()
}
