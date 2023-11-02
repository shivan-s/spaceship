package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	lipgloss "github.com/charmbracelet/lipgloss"
)

const (
	maxHeight = 24
	maxWidth  = 64
)

type pos struct {
	x int
	y int
}

type model struct {
	firing   bool
	ship     pos
	asteroid pos
}

type TickMsg time.Time

func createTick() tea.Cmd {
	return tea.Tick(time.Second, func(t time.Time) tea.Msg {
		return TickMsg(t)
	})
}

func initialModel() model {
	return model{
		ship: pos{x: 0, y: 0},

		asteroid: pos{x: 0, y: 0},
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "up", "k":
			if m.ship.y < 1 {
				m.ship.y = maxHeight
			} else {
				m.ship.y--
			}
		case "down", "j":
			if m.ship.y > maxHeight-1 {
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
	ship := "[]>"
	var titleScreen = lipgloss.NewStyle().
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("20")).
		Width(maxWidth).
		Align(lipgloss.Center)
	var gameScreen = lipgloss.NewStyle().
		BorderStyle(lipgloss.ThickBorder()).
		BorderForeground(lipgloss.Color("63")).
		Height(maxHeight).
		Width(maxWidth)
	var scoreScreen = lipgloss.NewStyle().
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("50")).
		Width(maxWidth).
		Align(lipgloss.Center)

	var screenArr []string
	for y := 0; y < maxHeight+1; y++ {
		var innerArr []string
		for x := 0; x < maxWidth; x++ {
			if y == m.ship.y {
				if x == 0 {
					innerArr = append(innerArr, ship)
				} else if x > len(ship)-1 {
					if m.firing == true {
						innerArr = append(innerArr, "~")
					} else {
						innerArr = append(innerArr, "·")
					}
				}
			} else {
				innerArr = append(innerArr, "·")
			}
		}
		innerArr = append(innerArr, "\n")
		screenArr = append(screenArr, strings.Join(innerArr, ""))
	}
	return lipgloss.JoinVertical(lipgloss.Center,
		titleScreen.Render("Spaceship by Shivan"),
		gameScreen.Render(strings.Join(screenArr, "")),
		scoreScreen.Render(strconv.Itoa(m.ship.y)))
}

func main() {
	p := tea.NewProgram(initialModel(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error: %v", err)
		os.Exit(1)
	}
}
