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

type Model struct {
	options    []string
	cursor     int
	choice     string
	table      table.Model
	nameInput  textinput.Model
	valueInput textinput.Model
	tagInput   textinput.Model
	typeInput  textinput.Model
	submit     bool
	chosen     bool
	quitting   bool
	err        error
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
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.Type {
			case tea.KeyTab:
				if m.nameInput.Focused() {
					m.nameInput.Blur()
					m.valueInput.Focus()
				} else if m.valueInput.Focused() {
					m.valueInput.Blur()
					m.tagInput.Focus()
				} else if m.tagInput.Focused() {
					m.tagInput.Blur()
					m.typeInput.Focus()
				} else {
					m.typeInput.Blur()
					m.nameInput.Focus()
				}

				return m, nil
			case tea.KeyEnter:
				m.submit = true
			}
		}

		m.nameInput, cmd = m.nameInput.Update(msg)
		m.valueInput, _ = m.valueInput.Update(msg)
		m.tagInput, _ = m.tagInput.Update(msg)
		m.typeInput, _ = m.typeInput.Update(msg)

		if m.submit == true {
			rows = append(rows, table.Row{
				m.nameInput.Value(),
				m.valueInput.Value(),
				m.tagInput.Value(),
				m.tagInput.Value(),
			})
			m.table.SetRows(rows)
		}
		return m, cmd
	case "See":
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
			}
			return m, nil
		}
	case "Manage":
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
			cursor = "►"
		}
		checked := "◇"
		if m.choice == m.options[i] {
			checked = "◆"
		}
		s += fmt.Sprintf("\n%s %s %s", cursor, checked, opt)
	}
	return s
}

func chosenView(m Model, s string) string {
	switch m.choice {
	case "Add":
		s = "What are we adding today?\n"
		s += fmt.Sprintf(
			"%s\n%s\n%s\n%s\n",
			m.nameInput.View(),
			m.valueInput.View(),
			m.tagInput.View(),
			m.typeInput.View(),
		)

		s += "\n\n\nPress 'Tab' to traverse field.\nPress 'Enter' to submit."
	case "See":
		for i, opt := range m.table.Rows() {
			cursor := " "
			if m.cursor == i {
				cursor = "►  "
			}
			s += fmt.Sprintf("\n%s %s", cursor, opt)
		}
	case "Manage":
		s = "This is the managing function"
	}
	return s
}

func initModel() Model {
	t := table.New()
	t.SetColumns(columns)
	t.SetRows(rows)

	nameInput := textinput.New()
	nameInput.Placeholder = "\n\nName of your expense\n\n"
	nameInput.CharLimit = 156
	nameInput.Width = 55
	nameInput.Focus()

	valueInput := textinput.New()
	valueInput.Placeholder = "\n\nValue of your expense\n\n"
	valueInput.CharLimit = 156
	valueInput.Width = 55

	tagInput := textinput.New()
	tagInput.Placeholder = "\n\nTag your expense\n\n"
	tagInput.CharLimit = 156
	tagInput.Width = 55

	typeInput := textinput.New()
	typeInput.Placeholder = "\n\nWhat's the type of your expense\n\n"
	typeInput.CharLimit = 156
	typeInput.Width = 55

	return Model{
		options:    []string{"Add", "Manage", "See"},
		nameInput:  nameInput,
		valueInput: valueInput,
		tagInput:   tagInput,
		typeInput:  typeInput,
		table:      t,
		submit:     false,
	}
}

func main() {
	p := tea.NewProgram(initModel(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error! %s", err)
		os.Exit(1)
	}
}
