package pages

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/ssh"
	"gorm.io/gorm"
)

type baseScreen struct {
	renderer *lipgloss.Renderer
	pty      ssh.Pty
}

func RootScreen(username string, newUserPublicKey []byte, isNewUser bool, renderer *lipgloss.Renderer, pty ssh.Pty, db *gorm.DB) rootScreenModel {
	var currentScreen tea.Model
	if isNewUser {
		currentScreen = WelcomeScreen(newUserPublicKey, renderer, pty, db)
	} else {
		currentScreen = HomeScreen(username, renderer, pty, db)
	}
	return rootScreenModel{
		currentScreen: currentScreen,
	}
}

type rootScreenModel struct {
	currentScreen tea.Model
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
