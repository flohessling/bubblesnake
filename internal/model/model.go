package model

import (
	"math/rand"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

const (
	Left = 1 + iota
	Right
	Up
	Down
)

const (
	HELP     = "\n\nuse arrow keys, wasd or hjkl to move the snake.\n"
	GAMEOVER = "\n\nGAME OVER\n"
	QUIT     = "press 'q' or 'ctrl + c' to quit.\n"
)

const INTERVAL = 100

type TickMsg time.Time

type Model struct {
	horizontalLine string
	verticalLine   string
	emptySymbol    string
	snakeSymbol    string
	foodSymbol     string
	width          int
	height         int
	arena          [][]string
	snake          snake
	lostGame       bool
	score          int
	food           coord
}

func (m Model) tick() tea.Cmd {
	return tea.Tick(time.Duration(INTERVAL)*time.Millisecond, func(t time.Time) tea.Msg {
		return TickMsg(t)
	})
}

func (m Model) changeDirection(direction int) (tea.Model, tea.Cmd) {
	if m.snake.hitWall(m) {
		m.lostGame = true
		return m, tea.Quit
	}

	opposite := map[int]int{
		Up:    Down,
		Down:  Up,
		Left:  Right,
		Right: Left,
	}

	if m.snake.direction != opposite[direction] {
		m.snake.direction = direction
	}

	return m, nil
}

func (m Model) moveSnake() (tea.Model, tea.Cmd) {
	head := m.snake.getHead()
	pos := coord{x: head.x, y: head.y}

	switch m.snake.direction {
	case Up:
		pos.x--
	case Down:
		pos.x++
	case Left:
		pos.y--
	case Right:
		pos.y++
	}

	if pos.x == m.food.x && pos.y == m.food.y {
		m.snake.length++
		x := rand.Intn(m.height-2) + 1
		y := rand.Intn(m.width-2) + 1

		for {
			if !m.snake.hitSelf(coord{x: x, y: y}) {
				break
			}
		}

		m.food.x = x
		m.food.y = y
	}

	if m.snake.hitSelf(pos) || m.snake.hitWall(m) {
		m.lostGame = true
		return m, tea.Quit
	}

	if len(m.snake.body) < m.snake.length {
		m.snake.body = append(m.snake.body, pos)
		m.score += 10
	} else {
		m.snake.body = append(m.snake.body[1:], pos)
	}

	return m, m.tick()
}

func InitialModel() Model {
	return Model{
		horizontalLine: "+",
		verticalLine:   "+",
		emptySymbol:    " ",
		snakeSymbol:    "o",
		foodSymbol:     "x",
		width:          60,
		height:         20,
		arena:          [][]string{},
		lostGame:       false,
		score:          0,
		food:           coord{x: 10, y: 10},
		snake: snake{
			body: []coord{
				{x: 1, y: 1},
				{x: 1, y: 2},
				{x: 1, y: 3},
			},
			length:    3,
			direction: Right,
		},
	}
}

func (m Model) Init() tea.Cmd {
	var x, y int

	x = rand.Intn(m.height - 1)
	y = rand.Intn(m.width - 1)

	m.food = coord{x: x, y: y}
	return m.tick()
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		case "up", "w", "k":
			return m.changeDirection(Up)
		case "down", "s", "j":
			return m.changeDirection(Down)
		case "left", "a", "h":
			return m.changeDirection(Left)
		case "right", "d", "l":
			return m.changeDirection(Right)
		}
	case TickMsg:
		return m.moveSnake()
	}

	return m, nil
}

func (m Model) View() string {
	var sb strings.Builder
	sb.WriteString(RenderTitle())
	sb.WriteString("\n")

	var stringArena strings.Builder
	RenderArena(&m)
	RenderSnake(&m)
	RenderFood(&m)

	for _, row := range m.arena {
		stringArena.WriteString(strings.Join(row, "") + "\n")
	}

	sb.WriteString(stringArena.String())
	sb.WriteString("\n")
	sb.WriteString(RenderScore(m.score))
	sb.WriteString("\n")

	if m.lostGame {
		sb.WriteString(RenderGameOver())
	}

	sb.WriteString(RenderHelp(HELP))
	sb.WriteString("\n")
	sb.WriteString(RenderHelp(QUIT))

	return sb.String()
}
