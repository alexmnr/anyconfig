package ui

import (
	"os"
	"strings"
	"github.com/charmbracelet/bubbles/filepicker"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
)

type file_picker_model struct {
	filepicker   filepicker.Model
	currentDir   string
  keys         keyMap
  help         help.Model
	quitting     bool
  message      string
}

type keyMap struct {
	Up    key.Binding
	Down  key.Binding
	In    key.Binding
	Out   key.Binding
	Select key.Binding
	Quit  key.Binding
}
func (k keyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Up, k.Down, k.In, k.Out, k.Select, k.Quit}
}

func (k keyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{{k.Up, k.Down, k.In, k.Out, k.Select, k.Quit},}
}

var keys = keyMap{
	Up: key.NewBinding(
		key.WithKeys("up", "k"),
		key.WithHelp("↑/k", "move up"),
	),
	Down: key.NewBinding(
		key.WithKeys("down", "j"),
		key.WithHelp("↓/j", "move down"),
	),
	Out: key.NewBinding(
		key.WithKeys("left", "h"),
		key.WithHelp("←/h", "parent dir"),
	),
	In: key.NewBinding(
		key.WithKeys("right", "l"),
		key.WithHelp("→/l", "into dir"),
	),
	Select: key.NewBinding(
		key.WithKeys("q"),
		key.WithHelp("q", "select"),
	),
	Quit: key.NewBinding(
		key.WithKeys("ctrl-c"),
		key.WithHelp("ctrl-c", "quit"),
	),
}

func (m file_picker_model) Init() tea.Cmd {
	return m.filepicker.Init()
}

func (m file_picker_model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			m.quitting = true
      m.currentDir = ""
			return m, tea.Quit
		case "q":
			m.quitting = true
			return m, tea.Quit
    }	
  case tea.WindowSizeMsg:
		m.help.Width = msg.Width
  }

	var cmd tea.Cmd
	m.filepicker, cmd = m.filepicker.Update(msg)

  m.currentDir = m.filepicker.CurrentDirectory
	return m, cmd
}

func (m file_picker_model) View() string {
	if m.quitting {
		return ""
	}
	var s strings.Builder
	s.WriteString("\n  ")
  s.WriteString(m.message + " " + m.currentDir)
	s.WriteString("\n\n" + m.filepicker.View() + "\n")

	helpView := m.help.View(m.keys)

	return s.String() + helpView
}

func FilePicker(message string, startDir string) string {
	fp := filepicker.New()
	fp.CurrentDirectory = "/opt/dotfiles"
	// fp.CurrentDirectory = startDir
	// fp.CurrentDirectory, _ = os.UserHomeDir()

	m := file_picker_model{
		filepicker: fp,
    help:       help.New(),
    keys:       keys,
    message:    message,
	}
	tm, _ := tea.NewProgram(&m, tea.WithOutput(os.Stderr)).Run()
	mm := tm.(file_picker_model)
	// fmt.Println("\n  You selected: " + mm.currentDir + "\n")
  return mm.currentDir
}
