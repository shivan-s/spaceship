package main

import (
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	lipgloss "github.com/charmbracelet/lipgloss"
	"os"
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
		case "ctrl+x", "q":
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
	var style = lipgloss.NewStyle().BorderStyle(lipgloss.NormalBorder()).BorderForeground(lipgloss.Color("63")).Width(24).Height(32)
	s := "[]>"
	return style.Render(s)
}

func main() {
	p := tea.NewProgram(initialModel(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error: %v", err)
		os.Exit(1)
	}
}
