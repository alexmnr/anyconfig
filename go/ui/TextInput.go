package ui

import (
  "out"

	"fmt"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)


type (
	errMsg error
)

type text_input_model struct {
	textInput textinput.Model
  title string
	err       error
  quitting  bool
}

func new_text_input_model(title string, placeholder string) text_input_model {
	ti := textinput.New()
	ti.Placeholder = placeholder
	ti.Focus()
	ti.CharLimit = 156
	ti.Width = 30

	return text_input_model{
		textInput: ti,
		title:     title,
		err:       nil,
		quitting:  false,
	}
}

func (m text_input_model) Init() tea.Cmd {
	return textinput.Blink
}

func (m text_input_model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter, tea.KeyCtrlC, tea.KeyEsc:
      m.quitting = true
			return m, tea.Quit
		}

	// We handle errors just like any other message
	case errMsg:
		m.err = msg
		return m, nil
	}

	m.textInput, cmd = m.textInput.Update(msg)
	return m, cmd
}

func (m text_input_model) View() string {
  if m.quitting == true {
    return fmt.Sprintf(
      out.Style(2, true, "! ") + out.Style(4, true, "Filename: ") + m.textInput.Value() + "\n%s",
    ) 
  } else {
    return fmt.Sprintf(
      out.Style(2, true, "? ") + out.Style(4, true, m.title) + "\n%s\n\n%s",
      m.textInput.View(),
      "(Enter to accept)",
    ) 
  }
}

func TextInput(title string, placeholder string) string {
  m := new_text_input_model(title, placeholder)
	tm, _ := tea.NewProgram(&m).Run()
	mm := tm.(text_input_model)
  return mm.textInput.Value()
}
