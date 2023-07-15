package main

import (
	"action"
	"config"
	"out"
	"task"
	"ui"

	"flag"
	"fmt"
	"os"

	"github.com/AlecAivazis/survey/v2"
)

var debugFlag = flag.Bool("debug", false, "Activate Debugging")

func main() {
  // config
  config := config.Init()
  config.Debug = false
  // flags and arguments
	flag.Parse()
  if *debugFlag == true {
    config.Debug = true
  }
  // check what to do
  input := ""
  options := []string{
    out.Style(3, false, "Install") + " existing tasks in repo", 
    out.Style(3, false, "Create") + " new task in repo", 
    out.Style(0, true, "Exit"), 
  }
  prompt := &survey.Select{
    Message: "What do you want to do? ",
    Options: options,
  }
  survey.AskOne(prompt, &input)

  ////// Install
  if input == options[0] {
    // Select Task files to run
    taskFiles := task.SelectFiles(config)
    if len(taskFiles) == 0 {
      out.Error("You must select something!")
      os.Exit(0)
    }
    // sort them in the right order
    for i := 0; i < 4; i++ {
      taskFiles = task.SortFiles(taskFiles, config)
    }
    // Create actions from selected files
    actions := []ui.Action{}
    for _, file := range taskFiles {
      actions = append(actions, action.GetActions(file, config)...)
    }
    fmt.Println()
    // run actions
    ui.RunActions(actions, config.Debug)

  ////// Create
  } else if input == options[1] {
    task.CreateFile(config)
  } else {
    fmt.Println("Bye!")
    os.Exit(0)
  }
}
