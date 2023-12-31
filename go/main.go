package main

import (
	"action"
	"config"
	"gh"
	"out"
	"task"
	"tools"
	"ui"
  // "command"

	"fmt"
	"os"

	"github.com/AlecAivazis/survey/v2"
)


func main() {
  // config
  config := config.Init()
  config.Debug = false
  directInstall := false
  silent := false
  var files []string
  //////// Arguments ///////
  args := os.Args
  for _, arg := range args {
    if arg == "-d" {
      config.Debug = true
    } else if arg == "-s" {
      silent = true
    } else if arg == "-i" {
      directInstall = true
    } else if directInstall == true {
      files = append(files, arg)
    }
  }
  // direct install
  if len(files) > 0 {
    // sort them in the right order
    taskFiles := task.SortFiles(files, config)
    // Create actions from selected files
    actions := []ui.Action{}
    for _, file := range taskFiles {
      actions = append(actions, action.GetActions(file, config)...)
    }
    fmt.Println()
    // run actions
    if silent == true {
      for _, action := range actions {
        out.Info("Running action with name: " + action.Name)
        action.Cmd()
      }
    } else {
      ui.RunActions(actions, config.Debug)
    }
    os.Exit(0)
  }
  // check if update is possible
  anyconfig_update := false
  repo_update := false
  if tools.CheckExist("/tmp/anyconfig_update") == true {
    anyconfig_update = true
  } 
  // check if repo has updates
  if tools.CheckExist("/tmp/repo_update") == true {
    repo_update = true
  } 
  // check what to do
  input := ""
  options := []string{
    out.Style(3, false, "Install") + " existing tasks in repo", 
    out.Style(3, false, "Create") + " new task in repo", 
  }
  if anyconfig_update == true {
    options = append(options, out.Style(2, false, "Update") + " anyconfig")
  }
  if repo_update == true {
    options = append(options, out.Style(2, false, "Update") + " Repository")
  }
  options = append(options, out.Style(0, true, "Exit"))
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
  } else if input == out.Style(2, false, "Update") + " anyconfig" {
    fmt.Println()
    gh.UpdateAnyconfig()
  } else if input == out.Style(2, false, "Update") + " Repository" {
    fmt.Println()
    gh.UpdateRepo(config.Repo)
  } else {
    fmt.Println("Bye!")
    os.Exit(0)
  }
}
