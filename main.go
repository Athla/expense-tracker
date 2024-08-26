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

var cmd tea.Cmd

type Model struct {
	options    []string
	cursor     int
	current    string
	table      table.Model
	inputField textinput.Model
	checked    bool
}

func (m Model) Init() tea.Cmd { return nil }

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	//  options: []string{"Add", "Manage", "See"},

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q", "esc":
			return m, tea.Quit
		case "k", "up":
			m.cursor--
			if m.cursor < 0 {
				m.cursor = 0
			}
		case "j", "down":
			m.cursor++
			if m.cursor > len(m.options)-1 {
				m.cursor = len(m.options) - 1
			}
		case "enter", " ":
			m.current = m.options[m.cursor]
			return m, nil
		}
	}
	switch m.current {
	case "Add":
		m.inputField, cmd = m.inputField.Update(msg)
		return m, cmd
	}
	return m, nil
}

func (m Model) View() string {
	s := "\n\tWhat to do?\n\n"

	for i, opt := range m.options {
		cursor := " "
		if m.cursor == i {
			cursor = ">>"
		}
		checked := " "
		if m.current == m.options[i] {
			checked = "x"
		}
		s += fmt.Sprintf("\n%s [%s] %s", cursor, checked, opt)
	}
	switch m.current {
	case "Add":
		s += m.inputField.View()
	case "See":
		cs := ""
		for _, v := range expenses {
			cs += fmt.Sprintf("\n\t%+v", v)
		}
		s += cs
	case "Manage":
		s += ""
	}
	return s
}

func initModel() Model {
	ti := textinput.New()
	ti.Placeholder = "\n\nWhat did you expend?"
	ti.Focus()
	ti.CharLimit = 156
	ti.Width = 20
	return Model{
		options:    []string{"Add", "Manage", "See"},
		inputField: ti,
	}
}

var expenses = []Expense{
	{Name: "Netflix", Value: 55.90, Tag: "Entretenimento", Type: "Mensal"},
	{Name: "Gympass", Value: 49.90, Tag: "Entretenimento", Type: "Mensal"},
	{Name: "Espetinho", Value: 35.00, Tag: "Comida", Type: "Avulso"},
}

func main() {
	p := tea.NewProgram(initModel(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error! %s", err)
		os.Exit(1)
	}
}
