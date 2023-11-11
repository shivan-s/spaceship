package main

import (
	"fmt"
	"math"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	lipgloss "github.com/charmbracelet/lipgloss"
)

const (
	H = 24
	W = 96
)
const SHIP = "[]>"
const MS_TICK = 250

type position struct {
	x int
	y int
	s int
}

type asteroid struct {
	x         int
	y         int
	s         int
	destroyed bool
}

type model struct {
	paused    bool
	gameOver  bool
	score     int
	firing    bool
	ship      position
	asteroids []asteroid
}

func (m model) activeAsteroids() int {
	var count int
	for _, a := range m.asteroids {
		if a.destroyed == false {
			count++
		}
	}
	return count
}

type TickMsg time.Time

func doTick() tea.Cmd {
	tickDuration := time.Duration(MS_TICK * time.Millisecond)
	return tea.Tick(tickDuration, func(t time.Time) tea.Msg {
		return TickMsg(t)
	})
}

func initialModel() model {
	return model{
		paused:   false,
		gameOver: false,
		score:    0,
		ship:     position{x: 0, y: rand.Intn(H), s: 0},
	}
}

func (m model) Init() tea.Cmd {
	return doTick()
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case TickMsg:
		if m.paused == false {
			for i, asteroid := range m.asteroids {
				if asteroid.destroyed == false {
					m.asteroids[i].x = asteroid.x - asteroid.s
					if asteroid.x < len(SHIP) {
						m.gameOver = true
					} else if m.firing == true && m.ship.y == asteroid.y && asteroid.destroyed == false {
						m.asteroids[i].destroyed = true
						m.score++
					}
				}
			}
			if rand.Intn(3) == 0 && m.gameOver == false {
				f := int(math.Round(float64(m.score)/100)) + 1
				s := rand.Intn(2*f+1) + 1*f
				newAsteroid := asteroid{x: W - 2, y: rand.Intn(H), s: s, destroyed: false}
				m.asteroids = append(m.asteroids, newAsteroid)
			}
			m.firing = false
		}
		return m, doTick()
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			if m.gameOver == true || m.paused == true {
				return m, tea.Quit
			}
		case "r":
			if m.gameOver == true || m.paused == true {
				m.paused = false
				m.score = 0
				m.asteroids = make([]asteroid, 0)
			}
		case "p":
			if m.gameOver == false {
				m.paused = !m.paused
			}
		case "up", "k":
			if m.paused == false {
				if m.ship.y < 1 {
					m.ship.y = H
				} else {
					m.ship.y--
				}
			}
		case "down", "j":
			if m.paused == false {
				if m.ship.y > H-1 {
					m.ship.y = 0
				} else {
					m.ship.y++
				}
			}
		case " ":
			if m.paused == false {
				m.firing = true
			}
		}
	}
	return m, nil
}

func offset(w int, t string) int {
	s := strings.Split(t, "")
	if len(s) > w {
		s = s[:w]
	}
	return int(math.Round(float64(w)/2) - math.Round(float64(len(s)/2)))
}

func (m model) View() string {
	col := "1"
	if m.gameOver == true || m.paused == true {
		col = strconv.Itoa(rand.Intn(255))
	}
	var titleScreen = lipgloss.NewStyle().
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("1")).
		Width(W).
		Align(lipgloss.Center)
	var mainScreen = lipgloss.NewStyle().
		BorderStyle(lipgloss.ThickBorder()).
		BorderForeground(lipgloss.Color(col)).
		Height(H).
		Width(W)
	var scoreScreen = lipgloss.NewStyle().
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("1")).
		Width(W).
		Align(lipgloss.Center)
	screen := make([][]string, H+1)
	hY := int(math.Round(H / 2))
	for y := range screen {
		line := make([]string, W)
		for x := range line {
			if x == len(line)-1 {
				line[x] = "\n"
			} else {
				line[x] = " "
			}
		}
		screen[y] = line
	}
	if m.gameOver == true {
		{
			msg := "Game Over!"
			offsetX := offset(W, msg)
			for x, c := range strings.Split(msg, "") {
				screen[hY][x+offsetX] = c
			}
		}
		{
			msg := "Press 'q' to quit or 'r' to play again."
			offsetX := offset(W, msg)
			for x, c := range strings.Split(msg, "") {
				screen[hY+1][x+offsetX] = c
			}
		}

	} else if m.paused == true {
		{
			msg := "Paused!"
			offsetX := offset(W, msg)
			for x, c := range strings.Split(msg, "") {
				screen[hY+1][x+offsetX] = c
			}
		}
		{
			msg := "Press 'p' to resume, 'q' to quit, or 'r' to restart game."
			offsetX := offset(W, msg)
			for x, c := range strings.Split(msg, "") {
				screen[hY+1][x+offsetX] = c
			}
		}
	} else {
		for x, c := range strings.Split(SHIP, "") {
			screen[m.ship.y][x] = c
		}
		for _, a := range m.asteroids {
			if a.destroyed == false {
				screen[a.y][a.x] = "*"
			}
		}
		if m.firing {
			for x := range screen[m.ship.y] {
				if x > len(SHIP) {
					screen[m.ship.y][x-1] = "-"
				}

			}
		}
	}
	gameScreen := make([]string, len(screen))
	for y, line := range screen {
		gameScreen[y] = strings.Join(line, "")
	}
	gameScreenContent := strings.Join(gameScreen, "")
	return lipgloss.JoinVertical(
		lipgloss.Center,
		titleScreen.Render("Spaceship by Shivan"),
		mainScreen.Render(gameScreenContent),
		scoreScreen.Render(
			"Score: ",
			strconv.Itoa(m.score),
			"Asteroids: ",
			strconv.Itoa(m.activeAsteroids()),
		),
	)
}

func main() {
	p := tea.NewProgram(initialModel(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error: %v", err)
		os.Exit(1)
	}
}
