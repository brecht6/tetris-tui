package main

import (
	"fmt"
	"os"
	"time"

	tea "charm.land/bubbletea/v2"
	lipgloss "charm.land/lipgloss/v2"
)

func initialModel() model {
	return model{
		blocks:       [boardHeight][boardWidth]string{},
		currentBlock: getRandomPiece(),
		score:        0,
		isGameOver:   false,
		width:        0,
		height:       0,
	}
}

func tick(score int) tea.Cmd {
	return tea.Tick(time.Second/time.Duration(2+min(score, 10000)/500), func(t time.Time) tea.Msg {
		return "tick"
	})
}

func (m model) Init() tea.Cmd {
	return tick(m.score)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	// update window size
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

	case string:
		// move down
		if msg == "tick" {
			if hitBottom(m.currentBlock, m.blocks) {
				blocks, points := addBlockToMatrix(m.blocks, m.currentBlock, true)
				m.score += points
				m.blocks = blocks
				m.currentBlock = getRandomPiece()

				if isGameOver(m.blocks, m.currentBlock) {
					m.isGameOver = true
					return m, nil // stop the ticks
				}
			} else {
				m.currentBlock.positions = moveBlockDown(m.currentBlock, m.blocks)
			}
			return m, tick(m.score)
		}

	case tea.KeyPressMsg:

		switch msg.String() {

		// stop
		case "ctrl+c", "q":
			return m, tea.Quit

		// restart
		case "r":
			m.blocks = [boardHeight][boardWidth]string{}
			m.currentBlock = getRandomPiece()
			m.score = 0

			ended := m.isGameOver

			m.isGameOver = false

			if ended {
				return m, tick(m.score)
			}

		// movement
		case "left":
			m.currentBlock.positions = moveBlockLeft(m.currentBlock, m.blocks)

		case "right":
			m.currentBlock.positions = moveBlockRight(m.currentBlock, m.blocks)

		case "down":
			m.currentBlock.positions = moveBlockDown(m.currentBlock, m.blocks)

		case "up":
			m.currentBlock.positions = rotateBlock(m.currentBlock, m.blocks)

		case "space":
			for !hitBottom(m.currentBlock, m.blocks) {
				m.currentBlock.positions = moveBlockDown(m.currentBlock, m.blocks)
			}
			blocks, points := addBlockToMatrix(m.blocks, m.currentBlock, true)
			m.score += points
			m.blocks = blocks

			m.currentBlock = getRandomPiece()

			if isGameOver(m.blocks, m.currentBlock) {
				m.isGameOver = true
			}
		}
	}

	return m, nil
}

// also show score in a window to the right of the tetris field
func (m model) View() tea.View {
	if m.height < boardHeight+2 || m.width < boardWidth*2+2 {
		view := tea.NewView("Please make the window bigger to play the game!")
		view.AltScreen = true
		return view
	}

	if m.isGameOver {
		gameOver := `
 ██████╗  █████╗ ███╗   ███╗███████╗     ██████╗ ██╗   ██╗███████╗██████╗ 
██╔════╝ ██╔══██╗████╗ ████║██╔════╝    ██╔═══██╗██║   ██║██╔════╝██╔══██╗
██║  ███╗███████║██╔████╔██║█████╗      ██║   ██║██║   ██║█████╗  ██████╔╝
██║   ██║██╔══██║██║╚██╔╝██║██╔══╝      ██║   ██║╚██╗ ██╔╝██╔══╝  ██╔══██╗
╚██████╔╝██║  ██║██║ ╚═╝ ██║███████╗    ╚██████╔╝ ╚████╔╝ ███████╗██║  ██║
 ╚═════╝ ╚═╝  ╚═╝╚═╝     ╚═╝╚══════╝     ╚═════╝   ╚═══╝  ╚══════╝╚═╝  ╚═╝`

		s := lipgloss.NewStyle().Foreground(lipgloss.Color("160")).Render(gameOver)
		score := fmt.Sprintf("Your score: %d", m.score)
		s += "\n\n" + lipgloss.NewStyle().Foreground(lipgloss.Color("255")).Render(score)
		restart := lipgloss.NewStyle().Foreground(lipgloss.Color("255")).Render("Press 'r' to restart or 'q' to quit")
		combined := lipgloss.JoinVertical(lipgloss.Center, s, restart)
		centered := lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, combined)
		view := tea.NewView(centered)
		view.AltScreen = true
		return view
	}

	var s string

	matrix, _ := addBlockToMatrix(m.blocks, m.currentBlock, false)

	for row := 0; row < boardHeight; row++ {
		for col := 0; col < boardWidth; col++ {
			color := matrix[row][col]
			if color == "" {
				s += "  "
			} else {
				s += lipgloss.NewStyle().Foreground(lipgloss.Color(color)).Render("██")
			}
		}
		if row < boardHeight-1 {
			s += "\n"
		}
	}

	border := lipgloss.Border{
		Top:         "▄",
		Bottom:      "▀",
		Left:        "█",
		Right:       "█",
		TopLeft:     "▄",
		TopRight:    "▄",
		BottomLeft:  "▀",
		BottomRight: "▀",
	}

	style := lipgloss.NewStyle().
		Border(border).
		BorderForeground(lipgloss.Color("239"))

	borderWrapped := style.Render(s)

	scoreContent := fmt.Sprintf("Score\n%d", m.score)
	scoreView := lipgloss.NewStyle().Border(border).BorderForeground(lipgloss.Color("239")).Padding(0, 1).Render(scoreContent)

	combined := lipgloss.JoinHorizontal(lipgloss.Top, borderWrapped, " ", scoreView)

	centeredUI := lipgloss.Place(
		m.width,
		m.height,
		lipgloss.Center,
		lipgloss.Center,
		combined,
	)

	// return the view
	view := tea.NewView(centeredUI)
	view.AltScreen = true

	return view
}

func main() {
	p := tea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
