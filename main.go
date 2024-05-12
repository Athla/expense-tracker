package main

import (
	"flag"
	"fmt"
	"log/slog"
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

var (
	debug = flag.Bool("Debug", false, "Sets the debug mode")
)
var choices = []string{"Manage expenses", "Generate report"}

type model struct {
	choice string
	cursor int
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q", "esc":
			return m, tea.Quit
		case "enter":
			m.choice = choices[m.cursor]
			return m, tea.Quit

		case "down", "j":
			m.cursor++
			if m.cursor >= len(choices) {
				m.cursor = 0
			}

		case "up", "k":
			m.cursor--
			if m.cursor < 0 {
				m.cursor = len(choices) - 1
			}
		}
	}

	return m, nil
}

func (m model) View() string {
	s := strings.Builder{}
	s.WriteString("What do you want to do?")

	for i := 0; i < len(choices); i++ {
		if m.cursor == i {
			s.WriteString("[x]")
		} else {
			s.WriteString("[ ]")
		}

		s.WriteString(choices[i])
		s.WriteString("\n")
	}

	s.WriteString("\n(Press q to quit)\n")

	return s.String()
}
func main() {
	p := tea.NewProgram(model{})

	m, err := p.Run()
	if err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}
	if m, ok := m.(model); ok && m.choice != "" {
		fmt.Printf("\n---\nYou choose %s!\n", m.choice)
	}
}
