package pages

import (
	"io"
	"math/rand/v2"
	"os"
	"strings"
	"unicode"

	"github.com/bg16-2009/wordssh/models"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/log"
	"github.com/charmbracelet/ssh"
	"gorm.io/gorm"
)

func GameScreen(renderer *lipgloss.Renderer, pty ssh.Pty, db *gorm.DB, currentUser models.User) gameModel {
	m := gameModel{
		renderer: renderer,
		pty:      pty,

		width:  pty.Window.Width,
		height: pty.Window.Height,

		txtStyle:             renderer.NewStyle().Foreground(lipgloss.Color("10")),
		quitStyle:            renderer.NewStyle().Foreground(lipgloss.Color("8")),
		incorrectLetterStyle: renderer.NewStyle().Foreground(lipgloss.Color("#ff0000")),
		correctLetterStyle:   renderer.NewStyle().Foreground(lipgloss.Color("#00ff00")),
		misplacedLatterStyle: renderer.NewStyle().Foreground(lipgloss.Color("#ffff00")),

		answer:         "",
		answerMap:      make(map[rune]int),
		currentAttempt: 0,
		currentChar:    0,
		attempts:       6,
		wordLenght:     5,
		keyboardState:  make(map[rune]lipgloss.Style),
		wordlistLenght: 14855,

		err:      "",
		gameOver: false,
		win:      false,

		db:          db,
		currentUser: currentUser,
	}
	m.answer = getWordFromWordlist(rand.IntN(m.wordlistLenght), m.wordLenght)
	m.gameState = make([][]string, m.attempts)
	for i := range m.gameState {
		m.gameState[i] = make([]string, m.wordLenght)
	}
	return m
}

type gameModel struct {
	renderer *lipgloss.Renderer
	pty      ssh.Pty

	term   string
	width  int
	height int

	txtStyle             lipgloss.Style
	quitStyle            lipgloss.Style
	incorrectLetterStyle lipgloss.Style
	correctLetterStyle   lipgloss.Style
	misplacedLatterStyle lipgloss.Style

	answer         string
	answerMap      map[rune]int
	currentAttempt int
	currentChar    int
	attempts       int
	wordLenght     int
	keyboardState  map[rune]lipgloss.Style
	gameState      [][]string
	wordlistLenght int

	err      string
	win      bool
	gameOver bool

	db          *gorm.DB
	currentUser models.User
}

func (m gameModel) Init() tea.Cmd {
	for _, c := range m.answer {
		m.answerMap[c]++
	}
	return nil
}

func (m gameModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	m.err = ""
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.height = msg.Height
		m.width = msg.Width
	case tea.KeyMsg:
		if msg.String() == "q" || msg.String() == "ctrl+c" {
			return m, tea.Quit
		}
		if m.gameOver {
			return m, nil
		}
		if len(msg.String()) == 1 && unicode.IsLower([]rune(msg.String())[0]) && m.currentChar < m.wordLenght {
			m.gameState[m.currentAttempt][m.currentChar] = msg.String()
			m.currentChar++
		}
		if msg.String() == "backspace" && m.currentChar > 0 {
			m.currentChar--
		}
		if msg.String() == "enter" {
			if m.currentChar < m.wordLenght {
				m.err = "Word is too short"
				return m, nil
			}
			enteredWordRunes := make([]rune, m.wordLenght)
			for i, c := range m.gameState[m.currentAttempt] {
				enteredWordRunes[i] = []rune(c)[0]
			}
			if !chechIfWordIsInWordlist(string(enteredWordRunes), m.wordlistLenght) {
				m.err = "Invalid word"
				return m, nil
			}

			aTemporaryMap := make(map[rune]int)
			correctLetters := 0
			for i := 0; i < m.wordLenght; i++ {
				if []rune(m.answer)[i] == enteredWordRunes[i] {
					m.gameState[m.currentAttempt][i] = m.correctLetterStyle.Render(m.gameState[m.currentAttempt][i])
					aTemporaryMap[enteredWordRunes[i]]++
					correctLetters++
					m.keyboardState[enteredWordRunes[i]] = m.correctLetterStyle
				}
			}
			if correctLetters == m.wordLenght {
				m.win = true
				m.gameOver = true
				return m, nil
			}
			for i := 0; i < m.wordLenght; i++ {
				if len(m.gameState[m.currentAttempt][i]) == 1 {
					currentRune := []rune(m.gameState[m.currentAttempt][i])[0]
					if m.answerMap[currentRune]-aTemporaryMap[currentRune] > 0 {
						m.gameState[m.currentAttempt][i] = m.misplacedLatterStyle.Render(m.gameState[m.currentAttempt][i])
						aTemporaryMap[currentRune]++
						if _, ok := m.keyboardState[currentRune]; !ok {
							m.keyboardState[currentRune] = m.misplacedLatterStyle
						}
					} else {
						m.gameState[m.currentAttempt][i] = m.incorrectLetterStyle.Render(m.gameState[m.currentAttempt][i])
						if _, ok := m.keyboardState[currentRune]; !ok {
							m.keyboardState[currentRune] = m.incorrectLetterStyle
						}
					}
				}
			}
			m.currentAttempt++
			m.currentChar = 0
			if m.currentAttempt == m.attempts {
				m.gameOver = true
			}
		}
	}
	return m, nil
}

func getWordFromWordlist(idx int, wordLenght int) string {
	f, err := os.Open("wordlist")
	if err != nil {
		log.Errorf("Error opening wordlist: %v", err)
	}
	defer f.Close()

	_, err = f.Seek(int64((wordLenght+1)*idx), io.SeekStart)
	if err != nil {
		log.Errorf("Error seeking in wordlist: %v", err)
	}
	word := make([]byte, wordLenght)
	io.ReadAtLeast(f, word, wordLenght)
	return string(word)
}
func chechIfWordIsInWordlist(targetWord string, wordlistLength int) bool {
	var (
		left       = 0
		right      = wordlistLength - 1
		middle     = (left + right) / 2
		wordLenght = len(targetWord)
	)
	for {
		if left > right {
			return false
		}
		word := getWordFromWordlist(middle, wordLenght)
		if word == targetWord {
			return true
		}
		if word < targetWord {
			left = middle + 1
		} else {
			right = middle - 1
		}
		middle = (left + right) / 2
	}
}

func (m gameModel) View() string {
	st := m.renderer.NewStyle().Width(m.width).Align(lipgloss.Center)
	var s string
	s += st.Render(m.incorrectLetterStyle.Render(m.err))
	if m.gameOver {
		if m.win {
			s += st.Render(m.correctLetterStyle.Render("\nYou won\n"))
			m.db.Model(&m.currentUser).Update("Score", m.currentUser.Score+(uint(m.attempts-m.currentAttempt)*100))
		} else {
			s += st.Render(m.incorrectLetterStyle.Render("\nYou lost\nThe word was " + m.answer))
		}
	}
	table := "\n ┌───┬───┬───┬───┬───┐ \n"
	for i := 0; i < m.attempts; i++ {
		if i < m.currentAttempt {
			table += " │ "
			for j := 0; j < m.wordLenght; j++ {
				table += m.gameState[i][j] + " │ "
			}
		}
		if i == m.currentAttempt {
			table += " │ "
			for j := 0; j < m.currentChar; j++ {
				table += m.gameState[i][j] + " │ "
			}
			if m.currentChar != m.wordLenght {
				table += "_ │ "
			}
			for j := m.currentChar + 2; j <= m.wordLenght; j++ {
				table += "  │ "
			}
		}
		if i > m.currentAttempt {
			table += " │   │   │   │   │   │ "
		}
		table += "\n"

		if i != m.attempts-1 {
			table += " ├───┼───┼───┼───┼───┤ \n"
		}
	}

	table += " └───┴───┴───┴───┴───┘ \n"

	s += st.Render(table)

	bytes, err := os.ReadFile("keyboard")
	if err != nil {
		log.Errorf("There was a problem opening the keyboard file: %v", err)
	} else {
		keyboardStr := string(bytes)
		for key, style := range m.keyboardState {
			keyboardStr = strings.Replace(keyboardStr, string(unicode.ToUpper(key)), style.Render(string(unicode.ToUpper(key))), -1)
		}
		s += st.Render(keyboardStr)
	}
	s += "\n" + m.quitStyle.Render(st.Render("Press 'q' to quit\n"))
	return m.renderer.PlaceVertical(m.height, lipgloss.Center, s)
}
