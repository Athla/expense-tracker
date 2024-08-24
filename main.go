package main

import (
	"fmt"
	"log"
	"log/slog"
	"os"
	"strconv"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/fogleman/ease"
	"github.com/lucasb-eyer/go-colorful"
)

const (
	progBarWidth  = 75
	progBarFull   = "█"
	progEmptyChar = "░"
	dot           = " • "
)

var (
	keywordStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("211"))
	subtleStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("211"))
	ticksStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("211"))
	checkboxStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("211"))
	progressEmpty = subtleStyle.Render(progEmptyChar)
	dotStyle      = lipgloss.NewStyle().Foreground(lipgloss.Color("288")).Render(dot)
	mainStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("211"))
)

var choices = []string{"Manage expenses", "Generate report"}
var expensesChoices = []string{"Add expense", "Delete expense"}

type (
	tickMsg  struct{}
	frameMsg struct{}
)

type model struct {
	Choice   int
	Cursor   int
	Ticks    int
	Frames   int
	Progress float64
	Chosen   bool
	Loaded   bool
	Quitting bool
}

func tick() tea.Cmd {
	return tea.Tick(time.Second, func(t time.Time) tea.Msg {
		return tickMsg{}
	})
}
func frame() tea.Cmd {
	return tea.Tick(time.Second/60, func(t time.Time) tea.Msg {
		return tickMsg{}
	})
}

func main() {
	log.Println("INCOMPLETE -- ADD FUNCTIONALITY TO THE CHOICES -- KEEP BROKEN UNTIL OK")
	p := tea.NewProgram(model{})

	m, err := p.Run()
	if err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}
	if m, ok := m.(model); ok {
		fmt.Printf("\n---\nYou choose %v\n", m.Choice)
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if msg, ok := msg.(tea.KeyMsg); ok {
		k := msg.String()
		if k == "q" || k == "esc" || k == "ctrl+c" {
			m.Quitting = true
			return m, tea.Quit
		}
	}

	if !m.Chosen {
		return updateChoices(msg, m)
	}
	return updateChosen(msg, m)
}

func (m model) View() string {
	var s string
	if m.Quitting {
		return "\n\t See you later!\n\n"
	}

	if !m.Chosen {
		s = choicesView(m)
	} else {
		s = chosenView(m)
	}

	return mainStyle.Render("\n" + s + "\n\n")
}

// Updating funcs

func updateChoices(msg tea.Msg, m model) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "j", "down":
			m.Choice++
			if m.Choice > 2 {
				m.Choice = 2
			}
		case "k", "up":
			m.Choice--
			if m.Choice < 0 {
				m.Choice = 0
			}
		case "enter":
			m.Chosen = true
			return m, frame()

		}
	case tickMsg:
		if m.Ticks == 0 {
			m.Quitting = true
			return m, tea.Quit
		}
		m.Ticks--
		return m, tick()
	}
	return m, nil
}

func updateChosen(msg tea.Msg, m model) (tea.Model, tea.Cmd) {
	switch msg.(type) {
	case frameMsg:
		if !m.Loaded {
			m.Frames++
			m.Progress = ease.OutBounce(float64(m.Frames) / float64(100))
			if m.Progress >= 1 {
				m.Progress = 1
				m.Loaded = true
				m.Ticks = 2
				return m, tick()
			}
			return m, frame()
		}
	case tickMsg:
		if m.Loaded {
			if m.Ticks == 0 {
				m.Quitting = true
				return m, tea.Quit
			}
			m.Ticks--
			return m, tick()
		}
	}

	return m, nil
}

// views

func choicesView(m model) string {
	c := m.Choice
	tpl := "What do you want to do?\n\n"
	tpl += "%s\n\n"
	tpl += "Program quits in %s"
	tpl += subtleStyle.Render("j/k, up/down: select") + dotStyle +
		subtleStyle.Render("enter: choose") + dotStyle +
		subtleStyle.Render("q/esc: exit")

	choices := fmt.Sprintf(
		"%s\n%s\n%s\n",
		checkbox("Add expense", c == 0),
		checkbox("See expenses", c == 1),
		checkbox("Manage expenses", c == 2),
	)

	return fmt.Sprintf(tpl, choices, ticksStyle.Render(strconv.Itoa(m.Ticks)))
}

func chosenView(m model) string {
	var msg string

	switch m.Choice {
	case 0:
		msg = fmt.Sprintf("Adding expenses?")
	case 1:
		msg = fmt.Sprintf("Seeing expenses?")
	case 2:
		msg = fmt.Sprintf("Managing expenses?")
	}

	label := "Loading..."
	if m.Loaded {
		label = fmt.Sprintf("Loaded. Exiting in %ss...", ticksStyle.Render(strconv.Itoa(m.Ticks)))
	}

	return msg + "\n\n" + label + "\n"
}

func checkbox(label string, checked bool) string {
	if checked {
		return checkboxStyle.Render("[x] " + label)
	}
	return fmt.Sprintf("[ ] %s", label)
}

// Generate a blend of colors.
func makeRampStyles(colorA, colorB string, steps float64) (s []lipgloss.Style) {
	cA, _ := colorful.Hex(colorA)
	cB, _ := colorful.Hex(colorB)

	for i := 0.0; i < steps; i++ {
		c := cA.BlendLuv(cB, i/steps)
		s = append(s, lipgloss.NewStyle().Foreground(lipgloss.Color(colorToHex(c))))
	}
	return
}

// Convert a colorful.Color to a hexadecimal format.
func colorToHex(c colorful.Color) string {
	return fmt.Sprintf("#%s%s%s", colorFloatToHex(c.R), colorFloatToHex(c.G), colorFloatToHex(c.B))
}

// Helper function for converting colors to hex. Assumes a value between 0 and
func colorFloatToHex(f float64) (s string) {
	s = strconv.FormatInt(int64(f*255), 16)
	if len(s) == 1 {
		s = "0" + s
	}
	return
}
