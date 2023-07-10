package task

import (
	"config"
	"out"
	"ui"
  "command"
  "tools"

	"os"
	"strings"

	"github.com/AlecAivazis/survey/v2"
)

var options []string = []string{
  out.Style(3, false, "Single") + " File Task", 
  out.Style(3, false, "Multi") + " File Task", 
  out.Style(3, false, "Add") + " Task to existing dir", 
  out.Style(0, true, "Exit"), 
}

func CreateFile(config config.AnyConfig) {
  templatePath := "/opt/anyconfig/template/template.yml"
  // check if template exists
  if tools.CheckExist(templatePath) == false {
    out.Error("Could not find template.yml, check if installed correctly!")
    os.Exit(1)
  }
  newFile := ""
  // get type of task
  num := getOption()
  if num == 0 {
    // get name of task
    filename := getName()
    split := strings.Split(filename, ".")
    filename = split[0] + ".yml"
    // check if name is already taken
    newFile = config.Repo + "/.anyconfig/" + filename
    if tools.CheckExist(newFile) == true {
      out.Error("This filename is already taken")
      os.Exit(0)
    }
    // copy template
    cmd_string := "cp " + templatePath + " "+ newFile
    err, output, error := command.Cmd(cmd_string, false, false)
    if err != nil {
      out.CommandError(cmd_string, err, output, error)
      os.Exit(1)
    }
  } else if num == 1 {
    // get name of task
    name := getName()
    split := strings.Split(name, ".")
    name = split[0] 
    // check if name is already taken
    dir := config.Repo + "/.anyconfig/" + name
    if tools.CheckExist(dir) == true {
      out.Error("This name is already taken")
      os.Exit(0)
    }
    command.Mkdir(dir, false)
    // copy template
    cmd_string := "cp " + templatePath + " " + dir + "/" + name + ".yml"
    err, output, error := command.Cmd(cmd_string, false, false)
    if err != nil {
      out.CommandError(cmd_string, err, output, error)
      os.Exit(1)
    }
    newFile = dir + "/" + name + ".yml"
  } else if num == 2 {
    dirs := tools.GetDirs(config.Repo + "/.anyconfig")
    // get dir to create file in
    input := ""
    prompt := &survey.Select{
      Message: "In what directory do you want to add your task?",
      Options: dirs,
    }
    survey.AskOne(prompt, &input)
    // get name of task
    filename := getName()
    split := strings.Split(filename, ".")
    filename = split[0] + ".yml"
    // check if name is already taken
    newFile = config.Repo + "/.anyconfig/" + input + "/"+ filename
    if tools.CheckExist(newFile) == true {
      out.Error("This filename is already taken")
      os.Exit(0)
    }
    // copy template
    cmd_string := "cp " + templatePath + " "+ newFile
    err, output, error := command.Cmd(cmd_string, false, false)
    if err != nil {
      out.CommandError(cmd_string, err, output, error)
      os.Exit(1)
    }
  } else if num == 3 {
    os.Exit(0)
  }
	editor := os.Getenv("EDITOR")
	if editor == "" {
		editor = "vim"
	}
  command.Cmd(editor + " " + newFile, false, true)
  out.Info("Succesfully create task!")
}

func getOption() int {
  input := ""
  // Ask if multi or single task
  prompt := &survey.Select{
    Message: "What kind of task do you want to create?",
    Options: options,
  }
  survey.AskOne(prompt, &input)
  for n, k := range options {
    if input == k {
      return n
    }
  }
  return -1
}

func getName() string {
  string := ui.TextInput("Name of Task:", "filename")
  return string
}
