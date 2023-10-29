package main

import (
	"fmt"
	"os"
	"strconv"

	tea "github.com/charmbracelet/bubbletea"
	lipgloss "github.com/charmbracelet/lipgloss"
)

const (
	maxWidth  = 40
	maxHeight = 10
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
			m.ship.y--
		case "down", "j":
			m.ship.y++
		case " ":
			m.firing = true
		}
	}
	return m, nil
}

func (m model) View() string {
	var screen = lipgloss.NewStyle().
		BorderStyle(lipgloss.ThickBorder()).
		BorderForeground(lipgloss.Color("63")).
		Width(48).
		Height(24)

	ship := "[]>"
	w, h := lipgloss.Size(screen.Render())

	var screenArr [][]string
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			if y == m.ship.y {
				screenArr[y][0] = "S"
			}
			if x == len(w) {
				screenArr[y][x]
			}
		}
	}

	return screen.Render(s, strconv.Itoa(m.ship.y))
}

func main() {
	p := tea.NewProgram(initialModel(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error: %v", err)
		os.Exit(1)
	}
}
