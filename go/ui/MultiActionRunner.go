package ui

import (
	"command"
	"out"
	"tools"

	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/bubbles/progress"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Action struct{
  Name string
  Cmd func() error
}

type installedPkgMsg string
var (
	currentPkgNameStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#6CD0D4"))
	doneStyle           = lipgloss.NewStyle().Margin(1, 2)
	checkMark           = lipgloss.NewStyle().Foreground(lipgloss.Color("#4EF465")).SetString("✓")
)

////////// multi action runner ////////////
type multi_action_model struct {
	actions    []Action
	index    int
	width    int
	height   int
	spinner  spinner.Model
	progress progress.Model
	done     bool
  debug    bool
}

func new_multi_model(actions []Action, debug bool) multi_action_model {
	p := progress.New(
		progress.WithScaledGradient("#6CD0D4", "#C45CFA"),
		progress.WithWidth(90),
		progress.WithoutPercentage(),
	)
	s := spinner.New()
	s.Spinner = spinner.MiniDot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("63"))
	return multi_action_model{
	  actions: actions,
    index: 0,
		spinner:  s,
		progress: p,
		debug: debug,
	}
}

func (m multi_action_model) Init() tea.Cmd {
	return tea.Batch(
    m.spinner.Tick,
  )
}

func (m multi_action_model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width, m.height = msg.Width, msg.Height
    m.progress.Width = msg.Width - 20
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc", "q":
			return m, tea.Quit
		}
	case installedPkgMsg:
		// Update progress bar
		progressCmd := m.progress.SetPercent(float64(m.index + 1) / float64(len(m.actions )))
		if m.index >= len(m.actions) - 1 {
			m.done = true
			return m, tea.Sequence(
        tea.Printf("%s %s", checkMark, m.actions[m.index].Name), // print success message above our program
        tea.Quit,
      )
		}

    m.index++

		return m, tea.Batch(
			progressCmd,
			tea.Printf("%s %s", checkMark, m.actions[m.index - 1].Name), // print success message above our program
		)
	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	case progress.FrameMsg:
		newModel, cmd := m.progress.Update(msg)
		if newModel, ok := newModel.(progress.Model); ok {
			m.progress = newModel
		}
		return m, cmd
	}

	return m, nil
}

func (m multi_action_model) View() string {
	n := len(m.actions)
	w := lipgloss.Width(fmt.Sprintf("%d", n))

	if m.done {
		return doneStyle.Render(fmt.Sprintf("Done! Ran %d actions.\n", n))
	}

	pkgCount := fmt.Sprintf(" %*d/%*d", w, m.index, w, n)

	spin := m.spinner.View() + " "
	prog := m.progress.View()
	cellsAvail := max(0, m.width-lipgloss.Width(spin))

	pkgName := currentPkgNameStyle.Render(m.actions[m.index].Name)
	info := lipgloss.NewStyle().MaxWidth(cellsAvail).Render(pkgName + "\n")

	cellsRemaining := max(0, m.width-lipgloss.Width(prog+pkgCount))
	gap := strings.Repeat(" ", cellsRemaining / 2)

  if m.debug == true {
    return ""
  }

	return spin + info + "\n" + gap + prog + pkgCount + "\n"
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func RunActions(actions []Action, debug bool) {
  if len(actions) == 0 {
    out.Info("nothing to do")
    return
  }
  // get sudo rights
  if tools.GetUser() != "root"{
    command.Cmd("sudo true", false, true)
  }
  // create model
  model := new_multi_model(actions, debug)
  p := tea.NewProgram(model)
  // run actions
  go func(){
    for _, action := range actions {
      err := action.Cmd()
      if err != nil {
        p.Kill()
        p.Quit()
      }
      p.Send(installedPkgMsg(action.Name))
    }
  }()
  // start manager
	if _, err := p.Run(); err != nil {
    // printError(err)
		// fmt.Println("Error running package_manager:", err)
		os.Exit(0)
	}
}

