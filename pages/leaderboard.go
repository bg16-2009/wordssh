package pages

import (
	"fmt"

	"github.com/bg16-2009/wordssh/models"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/ssh"
	"gorm.io/gorm"
)

func LeaderboardScreen(username string, renderer *lipgloss.Renderer, pty ssh.Pty, db *gorm.DB) leaderboardModel {
	columns := []table.Column{
		{Title: "Rank", Width: 4},
		{Title: "Username", Width: 16},
		{Title: "Score", Width: 10},
	}
	rows := []table.Row{}
	var users []models.User
	result := db.Order("score desc").Limit(10).Find(&users)
	if result.Error != nil {
		panic(fmt.Sprintf("failed to query users: %v", result.Error))
	}
	for i, user := range users {
		rows = append(rows, []string{fmt.Sprintf("%d", i+1), user.Username, fmt.Sprintf("%d", user.Score)})
	}

	t := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(true),
		table.WithHeight(7),
	)

	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240")).
		BorderBottom(true).
		Bold(false)
	s.Selected = s.Selected.
		Foreground(lipgloss.Color("229")).
		Background(lipgloss.Color("57")).
		Bold(false)
	t.SetStyles(s)
    t.Focus()

	return leaderboardModel{
		width:  pty.Window.Width,
		height: pty.Window.Height,

		table: t,
		baseStyle: renderer.NewStyle().
			BorderStyle(lipgloss.NormalBorder()).
			BorderForeground(lipgloss.Color("240")),
		renderer: renderer,
		pty:      pty,
		username: username,
		db:       db,
	}
}

type leaderboardModel struct {
	width  int
	height int

	baseStyle lipgloss.Style
	table     table.Model
	renderer  *lipgloss.Renderer
	pty       ssh.Pty
	username  string
	db        *gorm.DB
}

func (m leaderboardModel) Init() tea.Cmd {
	return nil
}

func (m leaderboardModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		}
	}
    var cmd tea.Cmd
    m.table, cmd = m.table.Update(msg)
	return m, cmd
}

func (m leaderboardModel) View() string {
	return m.renderer.Place(
		m.width, m.height, lipgloss.Center, lipgloss.Center,
		lipgloss.JoinVertical(
			lipgloss.Center,
			m.baseStyle.Render(m.table.View())+"\n  "+m.table.HelpView()+"\n",
			"Press any key to go back",
		),
	)
}
