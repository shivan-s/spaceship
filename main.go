package main

import (
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	lipgloss "github.com/charmbracelet/lipgloss"
)

const (
	H = 12
	W = 48
)

const SHIP = "[]>"

type position struct {
	x int
	y int
	s int
}

type model struct {
	isGameOver bool
	score      int
	firing     bool
	ship       position
	asteroids  []position
}

type TickMsg time.Time

func doTick() tea.Cmd {
	return tea.Tick(time.Second, func(t time.Time) tea.Msg {
		return TickMsg(t)
	})
}

func initialModel() model {
	return model{
		isGameOver: false,
		score:      0,
		ship:       position{x: 0, y: rand.Intn(H), s: 0},
	}
}

func (m model) Init() tea.Cmd {
	return doTick()
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case TickMsg:
		for i, asteroid := range m.asteroids {
			m.asteroids[i].x = asteroid.x - asteroid.s
			if asteroid.x < len(SHIP) {
				m.isGameOver = true
			} else if m.firing == true && m.ship.y == asteroid.y {
				m.score++
				m.asteroids[i] = m.asteroids[len(m.asteroids)-1]
				m.asteroids = m.asteroids[:len(m.asteroids)-1]
			}
		}
		if rand.Intn(2) == 0 {
			newAsteroid := position{x: W - 2, y: rand.Intn(H), s: rand.Intn(3) + 1}
			m.asteroids = append(m.asteroids, newAsteroid)
		}
		m.firing = false
		return m, doTick()
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "up", "k":
			if m.ship.y < 1 {
				m.ship.y = H
			} else {
				m.ship.y--
			}
		case "down", "j":
			if m.ship.y > H-1 {
				m.ship.y = 0
			} else {
				m.ship.y++
			}
		case " ":
			m.firing = true
		}
	}
	return m, nil
}

func (m model) View() string {
	var titleScreen = lipgloss.NewStyle().
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("20")).
		Width(W).
		Align(lipgloss.Center)
	var gameScreen = lipgloss.NewStyle().
		BorderStyle(lipgloss.ThickBorder()).
		BorderForeground(lipgloss.Color("63")).
		Height(H).
		Width(W)
	var scoreScreen = lipgloss.NewStyle().
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("50")).
		Width(W).
		Align(lipgloss.Center)

	var gameScreenContent string
	if m.isGameOver == true {
		gameScreenContent = "Game Over!"
	} else {
		screen := make([][]string, H+1)
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
		for x, c := range strings.Split(SHIP, "") {
			screen[m.ship.y][x] = c
		}
		for _, a := range m.asteroids {
			screen[a.y][a.x] = "*"
		}
		if m.firing {
			for x := range screen[m.ship.y] {
				if x > len(SHIP) {
					screen[m.ship.y][x-1] = "-"
				}

			}
		}
		gameScreen := make([]string, len(screen))
		for y, line := range screen {
			gameScreen[y] = strings.Join(line, "")
		}
		gameScreenContent = strings.Join(gameScreen, "")
	}
	return lipgloss.JoinVertical(
		lipgloss.Center,
		titleScreen.Render("Spaceship by Shivan"),
		gameScreen.Render(gameScreenContent),
		scoreScreen.Render(
			"Score:",
			strconv.Itoa(m.score),
			"Asteroids:",
			strconv.Itoa(len(m.asteroids)),
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
