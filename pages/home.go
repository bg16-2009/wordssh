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
		width:  pty.Window.Width,
		height: pty.Window.Height,

		renderer: renderer,
		pty:      pty,
		db:       db,
		user:     user,
	}
}

type homeModel struct {
	width  int
	height int

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
		case "n":
			return rootScreenModel{}.switchScreen(GameScreen(m.renderer, m.pty, m.db, m.user))
		case "l":
			return rootScreenModel{}.switchScreen(LeaderboardScreen(m.user.Username, m.renderer, m.pty, m.db))
		case "/":
			return rootScreenModel{}.switchScreen(SearchScreen(m.user.Username, m.renderer, m.pty, m.db))
		}
	}
	return m, nil
}

func (m homeModel) View() string {
	newGameButton := `
┌────────────┐
│ n New Game │
└────────────┘`[1:]
	leaderboardButton := `
┌───────────────┐
│ l Leaderboard │
└───────────────┘`[1:]
	searchButton := `
┌─────────────────────────────┐
│ / Search for a user's score │
└─────────────────────────────┘`[1:]
	greeting := fmt.Sprintf("\nHi %s", m.user.Username)
	return lipgloss.JoinVertical(
		lipgloss.Center,
		greeting,
		m.renderer.Place(
			m.width, m.height-2, lipgloss.Center, lipgloss.Center,
			lipgloss.JoinVertical(
				lipgloss.Center,
				fmt.Sprintf("Your score is %d\n\n", m.user.Score),
				newGameButton,
				leaderboardButton,
				searchButton,
				m.renderer.NewStyle().Foreground(lipgloss.Color("8")).Render("\nPress 'q' to quit"),
			),
		),
	)
}
