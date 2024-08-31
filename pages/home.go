package pages

import (
	"fmt"

	"github.com/bg16-2009/wordssh/models"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/ssh"
	"gorm.io/gorm"
)

func HomeScreen(username string, renderer *lipgloss.Renderer, pty ssh.Pty, db *gorm.DB) homeModel {
	var user models.User
	db.First(&user, "username = ?", username)
	return homeModel{
		renderer: renderer,
		pty:      pty,
		db:       db,
		user:     user,
	}
}

type homeModel struct {
	renderer *lipgloss.Renderer
	pty      ssh.Pty
	db       *gorm.DB
	user     models.User
}

func (m homeModel) Init() tea.Cmd {
	return nil
}

func (m homeModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m homeModel) View() string {
	return fmt.Sprintf("Hi %s\nPress any key to play......", m.user.Username)
}
