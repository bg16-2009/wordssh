package pages

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/ssh"
)

func RootScreen(renderer *lipgloss.Renderer, pty ssh.Pty) rootScreenModel{
    return rootScreenModel{
        renderer: renderer,
        pty: pty,
        currentScreen: Welcome(renderer, pty),
    }
}

type rootScreenModel struct {
	currentScreen tea.Model
    renderer *lipgloss.Renderer
    pty ssh.Pty
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

func (m rootScreenModel) SwitchScreen(f func(*lipgloss.Renderer, ssh.Pty) gameModel) (tea.Model, tea.Cmd) {
    m.currentScreen = f(m.renderer, m.pty)
    return m.currentScreen, m.currentScreen.Init()
}
