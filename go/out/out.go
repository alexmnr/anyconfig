package out

import (
  "github.com/charmbracelet/lipgloss"
  "fmt"
)

var color [5]string = [5] string {
  "#FA4453",
  "#6EFA73",
  "#F270F5",
  "#96AEF9",
  "#d0d0d0",
}
var info_color = "#6EFA73"
var warning_color = "#FA4453"

func Style(style int, bold bool, input string) string {
  if style < len(color) {
    return lipgloss.NewStyle().Foreground(lipgloss.Color(color[style])).Bold(bold).Render(input)
  } else {
    return input
  }
}

func ErrorString(input string) string {
  return lipgloss.NewStyle().Foreground(lipgloss.Color(color[0])).Bold(true).Render(input)
}
func InfoString(input string) string {
  return lipgloss.NewStyle().Foreground(lipgloss.Color(color[3])).Bold(true).Render(input)
}

func Error(error interface{}) {
  string := fmt.Sprint(error)
  fmt.Println(Style(0, true, "Error: ") + string)
}
func Info(info interface{}) {
  string := fmt.Sprint(info)
  fmt.Println(Style(1, true, "Info: ") + string)
}

func CommandError(command string, err error, out string, error string) {
  fmt.Println(lipgloss.NewStyle().Foreground(lipgloss.Color(color[0])).Bold(true).PaddingLeft(0).Render("Error running Command: ") + command) 
  fmt.Println(lipgloss.NewStyle().Foreground(lipgloss.Color(color[2])).Bold(false).PaddingLeft(1).Render("Error Code: ") + fmt.Sprint(err)) 
  fmt.Println(lipgloss.NewStyle().Foreground(lipgloss.Color(color[2])).Bold(false).PaddingLeft(1).Render("Command stdout: ") + out) 
  fmt.Println(lipgloss.NewStyle().Foreground(lipgloss.Color(color[2])).Bold(false).PaddingLeft(1).Render("Command stderr: ") + error) 
}
