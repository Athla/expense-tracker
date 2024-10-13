package main

import (
	"fmt"
	"os"
	"time"

	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var baseStyle = lipgloss.
	NewStyle().
	BorderStyle(lipgloss.NormalBorder()).
	BorderForeground(lipgloss.Color("420"))

const (
	hotPink  = lipgloss.Color("#FF06B7")
	darkGray = lipgloss.Color("#767676")
)

var (
	inputStyle    = lipgloss.NewStyle().Foreground(hotPink)
	continueStyle = lipgloss.NewStyle().Foreground(darkGray)
)

type Model struct {
	options  []string
	cursor   int
	choice   string
	table    table.Model
	inputs   []textinput.Model
	submit   bool
	chosen   bool
	quitting bool
	err      error
}

type Expense struct {
	Name  string
	Value float32
	Tag   string
	Type  string
}

func (m Model) Init() tea.Cmd { return textinput.Blink }

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if msg, ok := msg.(tea.KeyMsg); ok {
		k := msg.String()
		if k == "q" || k == "esc" || k == "ctrl+c" {
			m.quitting = true
			time.Sleep(time.Duration(200) * time.Millisecond)
			return m, tea.Quit
		}
		if k == "c" {
			m.chosen = true
		} else if k == "v" {
			m.chosen = false
		}
	}
	if !m.chosen {
		return updateChoices(msg, m) //sub func for the update logic in the choices area
	}
	return updateChosen(msg, m) //sub func for the update logic in the chosen area
}

func updateChoices(msg tea.Msg, m Model) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "j", "down":
			m.cursor++
			if m.cursor > len(m.options)-1 {
				m.cursor = len(m.options) - 1
			}
		case "k", "up":
			m.cursor--
			if m.cursor < 0 {
				m.cursor = 0
			}
		case "enter":
			m.choice = m.options[m.cursor]
			m.chosen = true
			return m, nil
		}
		return m, nil
	}
	return m, nil
}

func updateChosen(msg tea.Msg, m Model) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd = make([]tea.Cmd, len(m.inputs))
	switch m.choice {
	case "Add":
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.Type {
			case tea.KeyTab:
				m.cursor++
				if m.cursor > len(m.inputs)-1 {
					m.cursor = len(m.inputs) - 1
				}
			case tea.KeyShiftTab:
				m.cursor--
				if m.cursor < 0 {
					m.cursor = 0
				}
			case tea.KeyEnter:
				m.submit = true
			}
		}

		for i := range m.inputs {
			m.inputs[i].Blur()
		}
		m.inputs[m.cursor].Focus()

		if m.submit == true {
			rows = append(rows, table.Row{
				m.inputs[name].Value(),
				m.inputs[val].Value(),
				m.inputs[tag].Value(),
				m.inputs[etype].Value(),
			})
			m.table.SetRows(rows)
			m.chosen = false
		}
		for i := range m.inputs {
			m.inputs[i], cmds[i] = m.inputs[i].Update(msg)
		}
		return m, tea.Batch(cmds...)
	case "See & Manage":
		m.table.Focus()
		switch msg := msg.(type) {
		case tea.KeyMsg:
			size := len(m.table.Rows()) - 1
			switch msg.String() {
			case "j", "down":
				m.cursor++
				if m.cursor > size {
					m.cursor = size
				}
			case "k", "up":
				m.cursor--
				if m.cursor < 0 {
					m.cursor = 0
				}
			case "d":
				rows := m.table.Rows()
				rows = append(rows[:m.cursor], rows[m.cursor+1:]...)
				m.table.SetRows(rows)
			}
			return m, nil
		}
	}
	return m, nil
}

func (m Model) View() string {
	var s string
	if !m.chosen {
		s = choicesView(m, s)
	} else {
		s = chosenView(m, s)
	}
	return s
}

func choicesView(m Model, s string) string {
	s = "\n\tWhat to do?\n\n"
	for i, opt := range m.options {
		cursor := " "
		if m.cursor == i {
			cursor = "► "
		}
		checked := "◇ "
		if m.choice == m.options[i] {
			checked = "◆ "
		}
		s += fmt.Sprintf("\n%s %s %s", cursor, checked, opt)
	}
	return s
}

func chosenView(m Model, s string) string {
	switch m.choice {
	case "Add":
		s = "What are we adding today?\n"
		for i, ip := range m.inputs {
			cursor := " "
			if m.cursor == i {
				cursor = "► "
			}
			s += fmt.Sprintf("\n%v%v", cursor, ip.View())
		}
		s += "\n\n\nPress 'Tab' to traverse field.\nPress 'Enter' to submit."
	case "See & Manage":
		for i, opt := range m.table.Rows() {
			cursor := " "
			if m.cursor == i {
				cursor = "►  "
			}
			s += fmt.Sprintf("\n%s %s", cursor, opt)
		}
	}
	return s
}

const (
	name = iota
	val
	tag
	etype
)

func initModel() Model {
	var inputs []textinput.Model = make([]textinput.Model, 4)

	inputs[name] = textinput.New()
	inputs[name].Placeholder = "\n\nName of your expense\n\n"
	inputs[name].CharLimit = 156
	inputs[name].Width = 55
	inputs[name].Prompt = ""
	inputs[name].Focus()

	inputs[val] = textinput.New()
	inputs[val].Placeholder = "\n\nValue of your expense\n\n"
	inputs[val].CharLimit = 156
	inputs[val].Width = 55
	inputs[val].Prompt = ""

	inputs[tag] = textinput.New()
	inputs[tag].Placeholder = "\n\nTag your expense\n\n"
	inputs[tag].CharLimit = 156
	inputs[tag].Width = 55
	inputs[tag].Prompt = ""

	inputs[etype] = textinput.New()
	inputs[etype].Placeholder = "\n\nWhat's the type of your expense\n\n"
	inputs[etype].CharLimit = 156
	inputs[etype].Width = 55
	inputs[etype].Prompt = ""

	t := table.New()
	t.SetColumns(columns)
	t.SetRows(rows)

	return Model{
		options: []string{"Add", "See & Manage"},
		inputs:  inputs,
		table:   t,
		submit:  false,
	}
}

func main() {
	p := tea.NewProgram(initModel(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error! %s", err)
		os.Exit(1)
	}
} // nextInput focuses the next input field
