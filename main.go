package main

import (
	"context"
	"errors"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"
	"unicode"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/log"
	"github.com/charmbracelet/ssh"
	"github.com/charmbracelet/wish"
	"github.com/charmbracelet/wish/activeterm"
	"github.com/charmbracelet/wish/bubbletea"
	"github.com/charmbracelet/wish/logging"
)

const (
	host = "localhost"
	port = "23234"
)

func teaHandler(s ssh.Session) (tea.Model, []tea.ProgramOption) {
	pty, _, _ := s.Pty()

	renderer := bubbletea.MakeRenderer(s)
	txtStyle := renderer.NewStyle().Foreground(lipgloss.Color("10"))
	quitStyle := renderer.NewStyle().Foreground(lipgloss.Color("8"))

	bg := "light"
	if renderer.HasDarkBackground() {
		bg = "dark"
	}

	m := model{
		term:           pty.Term,
		width:          pty.Window.Width,
		height:         pty.Window.Height,
		bg:             bg,
		txtStyle:       txtStyle,
		quitStyle:      quitStyle,
		currentAttempt: 2,
		currentChar:    1,
		attempts:       6,
		wordLenght:     5,
	}

    // This is temporary
	m.gameState = [][]string{
		{"a", "a", "a", "a", "a"},
		{"a", "a", "a", "a", "a"},
		{"a", "a", "a", "a", "a"},
		{"a", "a", "a", "a", "a"},
		{"a", "a", "a", "a", "a"},
		{"a", "a", "a", "a", "a"},
	}

	return m, []tea.ProgramOption{tea.WithAltScreen()}
}

type model struct {
	currentAttempt int
	currentChar    int
	attempts       int
	wordLenght     int
	term           string
	gameState     [][]string
	width          int
	height         int
	bg             string
	txtStyle       lipgloss.Style
	quitStyle      lipgloss.Style
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.height = msg.Height
		m.width = msg.Width
	case tea.KeyMsg:
		if msg.String() == "q" || msg.String() == "ctrl+c" {
			return m, tea.Quit
		}
		if len(msg.String()) == 1 && unicode.IsLower([]rune(msg.String())[0]) && m.currentChar < m.wordLenght {
			m.gameState[m.currentAttempt][m.currentChar] = msg.String()
			m.currentChar++
		}
		if msg.String() == "backspace" && m.currentChar > 0 {
			m.currentChar--
		}
	}
	return m, nil
}

func (m model) View() string {
	s := "┌───┬───┬───┬───┬───┐\n"
	for i := 0; i < m.attempts; i++ {
		if i < m.currentAttempt {
			s += "│ "
			for j := 0; j < m.wordLenght; j++ {
				s += m.gameState[i][j] + " │ "
			}
		}
		if i == m.currentAttempt {
			s += "│ "
			for j := 0; j < m.currentChar; j++ {
				s += m.gameState[i][j] + " │ "
			}
			if m.currentChar != m.wordLenght {
				s += "_ │ "
			}
			for j := m.currentChar + 2; j <= m.wordLenght; j++ {
				s += "  │ "
			}
		}
		if i > m.currentAttempt {
			s += "│   │   │   │   │   │"
		}
		s += "\n"

		if i != m.attempts-1 {
			s += "├───┼───┼───┼───┼───┤\n"
		}
	}

	s += "└───┴───┴───┴───┴───┘\n"
	return m.txtStyle.Render(s) + "\n\n" + m.quitStyle.Render("Press 'q' to quit\n")
}

func main() {
	s, err := wish.NewServer(
		wish.WithAddress(net.JoinHostPort(host, port)),
		wish.WithHostKeyPath(".ssh/id_ed25519"),
		wish.WithMiddleware(
			bubbletea.Middleware(teaHandler),
			activeterm.Middleware(), // Bubble Tea apps usually require a PTY.
			logging.Middleware(),
		),
	)
	if err != nil {
		log.Error("Could not start server", "error", err)
	}

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	log.Info("Starting SSH server", "host", host, "port", port)
	go func() {
		if err = s.ListenAndServe(); err != nil && !errors.Is(err, ssh.ErrServerClosed) {
			log.Error("Could not start server", "error", err)
			done <- nil
		}
	}()

	<-done
	log.Info("Stopping SSH server")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer func() { cancel() }()
	if err := s.Shutdown(ctx); err != nil && !errors.Is(err, ssh.ErrServerClosed) {
		log.Error("Could not stop server", "error", err)
	}
}
