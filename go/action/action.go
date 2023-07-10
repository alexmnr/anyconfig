package action

import (
  "command"
  "out"
  "config"
  "ui"
  "task"
  "tools"

  "strings"
  "os"
)

func GetActions(file string, config config.AnyConfig) []ui.Action {
  actions := []ui.Action{}
  taskDescription := task.GetTask(config.Repo + "/.anyconfig/" + file)
  // check dependencies
  check := true
  for _, actionDescription := range taskDescription.Dependencies{
    if actionDescription.Name == "noDir" {
      for _, input := range actionDescription.Args{
        input := fillArg(input, config)
        ret := tools.CheckExist(input)
        if ret == true  {
          check = false
        }
      }
    } else if actionDescription.Name == "noCommand" {
      for _, input := range actionDescription.Args{
        ret := tools.CommandExists(input)
        if ret == true  {
          check = false
        }
      }
    } else if actionDescription.Name == "user" {
      if actionDescription.Args[0] == "noroot" {
        if tools.GetUser() == "root" {
          out.Error("Can't run task " + file + " as root")
          os.Exit(1)
          return actions
        }
      }
    } else if actionDescription.Name == "os" {
      if tools.GetOS() != actionDescription.Args[0] {
          out.Error("Can't run task " + file + " with this operating system")
          os.Exit(1)
          return actions
      }
    }
  }
  if check == false {
    return actions
  }
  // create actions from task
  for _, actionDescription := range taskDescription.Install{
    // Package Install
    if actionDescription.Name == "pkg" {
      buffer := createAction_pkg(actionDescription, config)
      actions = append(actions, buffer...)
    // apt Install
    } else if actionDescription.Name == "apt" {
      buffer := createAction_apt(actionDescription, config)
      actions = append(actions, buffer...)
    // pacman Install
    } else if actionDescription.Name == "pacman" {
      buffer := createAction_pacman(actionDescription, config)
      actions = append(actions, buffer...)
    // yay Install
    } else if actionDescription.Name == "yay" {
      buffer := createAction_yay(actionDescription, config)
      actions = append(actions, buffer...)
    // custom command
    } else if actionDescription.Name == "cmd" {
      buffer := createAction_cmd(actionDescription, config)
      actions = append(actions, buffer...)
    // create dir
    } else if actionDescription.Name == "mkdir" {
      buffer := createAction_mkdir(actionDescription, config)
      actions = append(actions, buffer...)
    // backup
    } else if actionDescription.Name == "backup" {
      buffer :=createAction_backup(actionDescription, config)
      actions = append(actions, buffer...)
    // linking
    } else if actionDescription.Name == "ln" {
      buffer := createAction_ln(actionDescription, config)
      actions = append(actions, buffer...)
    // copying
    } else if actionDescription.Name == "cp" {
      buffer := createAction_cp(actionDescription, config)
      actions = append(actions, buffer...)
    // moving
    } else if actionDescription.Name == "mv" {
      buffer := createAction_mv(actionDescription, config)
      actions = append(actions, buffer...)
    // set env
    } else if actionDescription.Name == "env" {
      buffer := createAction_env(actionDescription, config)
      actions = append(actions, buffer...)
    } else {
      out.Error("Error: Did not recognise Action-Type: " + actionDescription.Name)
    }
  }
  return actions
}

// general package manager
func createAction_pkg(actionDescription task.ActionDescription, config config.AnyConfig) []ui.Action {
  actions := []ui.Action{}
  for _, input := range actionDescription.Args{
    name, arg := splitArg(input)
    arg = fillArg(arg, config)
    if name == "" {
      name = "Installing: " + arg
    }
    action := ui.Action{
      Name: name,
      Cmd: func() error {
        return command.PkgInstall(config.Installer, arg, config.Debug)
      },
    }
    actions = append(actions, action)
  }
  return actions
}
// apt install
func createAction_apt(actionDescription task.ActionDescription, config config.AnyConfig) []ui.Action {
  actions := []ui.Action{}
  if tools.CommandExists("apt") == false {
    return actions
  }
  for _, input := range actionDescription.Args{
    name, arg := splitArg(input)
    arg = fillArg(arg, config)
    if name == "" {
      name = "Installing: " + arg
    }
    string := "sudo apt update && sudo apt install -y " + arg
    action := ui.Action{
      Name: name,
      Cmd: func() error {
        err, output, error := command.Cmd(string, false, config.Debug)
        if err != nil {
          out.CommandError(string, err, output, error)
        }
        return err
      },
    }
    actions = append(actions, action)
  }
  return actions
}
// pacman install
func createAction_pacman(actionDescription task.ActionDescription, config config.AnyConfig) []ui.Action {
  actions := []ui.Action{}
  if tools.CommandExists("pacman") == false {
    return actions
  }
  for _, input := range actionDescription.Args{
    name, arg := splitArg(input)
    arg = fillArg(arg, config)
    if name == "" {
      name = "Installing: " + arg
    }
    string := "sudo pacman -Sy " + arg + " --needed --noconfirm"
    action := ui.Action{
      Name: name,
      Cmd: func() error {
        err, output, error := command.Cmd(string, false, config.Debug)
        if err != nil {
          out.CommandError(string, err, output, error)
        }
        return err
      },
    }
    actions = append(actions, action)
  }
  return actions
}
// yay install
func createAction_yay(actionDescription task.ActionDescription, config config.AnyConfig) []ui.Action {
  actions := []ui.Action{}
  for _, input := range actionDescription.Args{
    name, arg := splitArg(input)
    arg = fillArg(arg, config)
    if name == "" {
      name = "Installing: " + arg
    }
    string := "yay -Sy " + arg + " --needed --noconfirm"
    action := ui.Action{
      Name: name,
      Cmd: func() error {
        err, output, error := command.Cmd(string, false, config.Debug)
        if err != nil {
          out.CommandError(string, err, output, error)
        }
        return err
      },
    }
    actions = append(actions, action)
  }
  return actions
}
func createAction_cmd(actionDescription task.ActionDescription, config config.AnyConfig) []ui.Action {
  actions := []ui.Action{}
  for _, input := range actionDescription.Args {
    name, arg := splitArg(input)
    arg = fillArg(arg, config)
    if name == "" {
      name = "Command: " + arg
    }
    string := arg
    action := ui.Action{
      Name: name,
      Cmd: func() error {
        err, output, error := command.Cmd(string, false, config.Debug)
        if err != nil {
          out.CommandError(string, err, output, error)
        }
        return err
      },
    }
    actions = append(actions, action)
  }
  return actions
}
func createAction_backup(actionDescription task.ActionDescription, config config.AnyConfig) []ui.Action {
  actions := []ui.Action{}
  for _, input := range actionDescription.Args {
    name, arg := splitArg(input)
    arg = fillArg(arg, config)
    if name == "" {
      name = "Backup: " + arg + " to ~/.old"
    }
    action := ui.Action{
      Name: name,
      Cmd: func() error {
        err := command.Backup(arg)
        return err
      },
    }
    actions = append(actions, action)
  }
  return actions
}

func createAction_ln(actionDescription task.ActionDescription, config config.AnyConfig) []ui.Action {
  actions := []ui.Action{}
  for _, input := range actionDescription.Args {
    name, arg := splitArg(input)
    arg = fillArg(arg, config)
    org := ""
    des := ""
    if strings.Contains(arg, " > ") {
      split := strings.Split(arg, " > ")
      org = split[0]
      des = split[1]
    } else {
      out.Error("Invalid ln command")
      os.Exit(1)
    }
    if name == "" {
      name = "Linking " + org + " to " + des
    }
    action := ui.Action{
      Name: name,
      Cmd: func() error {
        err := command.Ln(org, des, true)
        return err
      },

    }
    actions = append(actions, action)
  }
  return actions
}

func createAction_cp(actionDescription task.ActionDescription, config config.AnyConfig) []ui.Action {
  actions := []ui.Action{}
  for _, input := range actionDescription.Args {
    name, arg := splitArg(input)
    arg = fillArg(arg, config)
    org := ""
    des := ""
    if strings.Contains(arg, " > ") {
      split := strings.Split(arg, " > ")
      org = split[0]
      des = split[1]
    } else {
      out.Error("Invalid cp command")
      os.Exit(1)
    }
    if name == "" {
      name = "Copying " + org + " to " + des
    }
    action := ui.Action{
      Name: name,
      Cmd: func() error {
        err := command.Cp(org, des, true)
        return err
      },
    }
    actions = append(actions, action)
  }
  return actions
}

func createAction_mv(actionDescription task.ActionDescription, config config.AnyConfig) []ui.Action {
  actions := []ui.Action{}
  for _, input := range actionDescription.Args {
    name, arg := splitArg(input)
    arg = fillArg(arg, config)
    org := ""
    des := ""
    if strings.Contains(arg, " > ") {
      split := strings.Split(arg, " > ")
      org = split[0]
      des = split[1]
    } else {
      out.Error("Invalid mv command")
      os.Exit(1)
    }
    if name == "" {
      name = "Moving " + org + " to " + des
    }
    action := ui.Action{
      Name: name,
      Cmd: func() error {
        err := command.Mv(org, des, true)
        return err
      },
    }
    actions = append(actions, action)
  }
  return actions
}

func createAction_mkdir(actionDescription task.ActionDescription, config config.AnyConfig) []ui.Action {
  actions := []ui.Action{}
  for _, input := range actionDescription.Args {
    name, arg := splitArg(input)
    arg = fillArg(arg, config)
    if name == "" {
      name = "Creating dir: " + arg
    }
    // check if dir should be overwritten
    backup := false
    if strings.Contains(arg, "{backup}") {
      arg = strings.Replace(arg, "{backup}", "", -1)
      backup = true
    } else {
      backup = false
    }
    action := ui.Action{
      Name: name,
      Cmd: func() error {
        err := command.Mkdir(arg, backup)
        return err
      },
    }
    actions = append(actions, action)
  }
  return actions
}

func createAction_env(actionDescription task.ActionDescription, config config.AnyConfig) []ui.Action {
  actions := []ui.Action{}
  for _, input := range actionDescription.Args {
    name, arg := splitArg(input)
    arg = fillArg(arg, config)
    env := ""
    key := ""
    if strings.Contains(arg, " = ") {
      split := strings.Split(arg, " = ")
      env = split[0]
      key = split[1]
    } else {
      out.Error("Invalid env command")
      os.Exit(1)
    }
    if name == "" {
      name = "Setting env " + env + " to " + key
    }
    action := ui.Action{
      Name: name,
      Cmd: func() error {
      return os.Setenv(env, key)
    },
    }
    actions = append(actions, action)
  }
  return actions
}

func splitArg(input string) (string, string) {
  name := ""
  args := ""
  if strings.Contains(input, "|") {
    split := strings.Split(input, " | ")
    name = split[0]
    args = split[1]
  } else {
    name = ""
    args = input
  }
  return name, args
}
func fillArg(input string, config config.AnyConfig) string {
  string := input
  // fill in Username
  string = strings.Replace(string, "{user}", config.User, -1)
  // fill in HomeDir
  string = strings.Replace(string, "{home}", config.HomeDir, -1)
  // fill in Repo
  string = strings.Replace(string, "{repo}", config.Repo, -1)
  return string
}


