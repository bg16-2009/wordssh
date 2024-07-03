package pages

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/ssh"
)

func Welcome(renderer *lipgloss.Renderer, pty ssh.Pty) welcomeModel {
	return welcomeModel{}
}

type welcomeModel struct{}

func (m welcomeModel) Init() tea.Cmd {
	return nil
}

func (m welcomeModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type){
    case tea.KeyMsg:
        switch msg.String(){
        case "q", "ctrl+c":
            return m, tea.Quit
        default:
            // TODO: Add page switch
            return m, nil
        }
    }
    return m, nil
}

func (m welcomeModel) View() string{
    return "Welcome to WordSSH\nPress any key to play......"
}
