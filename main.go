package main

import (
	"fmt"
	"os"

	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var baseStyle = lipgloss.
	NewStyle().
	BorderStyle(lipgloss.NormalBorder()).
	BorderForeground(lipgloss.Color("420"))

type Model struct {
	options    []string
	cursor     int
	choice     string
	table      table.Model
	inputField textinput.Model
	chosen     bool
	quitting   bool
}

type Expense struct {
	Name  string
	Value float32
	Tag   string
	Type  string
}

func (m Model) Init() tea.Cmd { return nil }

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if msg, ok := msg.(tea.KeyMsg); ok {
		k := msg.String()
		if k == "q" || k == "esc" || k == "ctrl+c" {
			m.quitting = true
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
			if m.cursor > 3 {
				m.cursor = 3
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
	var cmd tea.Cmd
	switch m.choice {
	case "Add":
		m.inputField.Focus()
		m.inputField, cmd = m.inputField.Update(msg)
		return m, cmd
	case "See":
		// return the current table of expenses
	case "Manage":
		// call manage function
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
			cursor = ">>"
		}
		checked := " "
		if m.choice == m.options[i] {
			checked = "x"
		}
		s += fmt.Sprintf("\n%s [%s] %s", cursor, checked, opt)

		if m.quitting {
			return "\n\nCYA!!!!\n\n"
		}
	}
	return s
}

func chosenView(m Model, s string) string {
	switch m.choice {
	case "Add":
		s = m.inputField.View()
	case "See":
		s = m.table.View()
	case "Manage":
		s = "This is the managing function"
	}
	return s
}

func initModel() Model {
	t := table.New()
	t.SetColumns(columns)
	t.SetRows(rows)

	ti := textinput.New()
	ti.Placeholder = "\n\nWhat did you expend on?\n\n"
	ti.CharLimit = 156
	ti.Width = 20
	return Model{
		options:    []string{"Add", "Manage", "See"},
		inputField: ti,
		table:      t,
	}
}

func main() {
	p := tea.NewProgram(initModel(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error! %s", err)
		os.Exit(1)
	}
}
