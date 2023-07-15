package gh

import (
	"out"
	"tools"
  "ui"
  "command"

	"os/exec"
  "os"
	"github.com/AlecAivazis/survey/v2"
	"strings"
  "encoding/json"

)

type FileStruct struct {
    Os string `yaml:"os"`
    Installer string `yaml:"installer"`
    UnInstaller string `yaml:"uninstaller"`
    Repo string `yaml:"repo"`
}

func Config() {
  // check if github-cli is installed
  check := tools.CommandExists("gh")
  if check == false {
    if tools.GetOS() == "arch" {
      // install github-cli on arch
      installer := tools.GetInstaller()
      command := func() error {
        err, _, _ := command.Cmd(installer + "github-cli git openssh", false, false)
        return err
      }
      action := ui.Action{
        Name: "Installing github-cli",
        Cmd: command,
      }
      ui.RunAction(action)
    } else if tools.GetOS() == "debian" {
      // install github-cli on debian
      installer := tools.GetInstaller()
      actions := []ui.Action{}
      command_string := func() error {
        err, _, _ := command.Cmd(installer + "git openssh-client", false, false)
        return err
      }
      action := ui.Action{Name: "Installing git and ssh", Cmd: command_string} 
      actions = append(actions, action)

      command_string = func() error {
        err, _, _ := command.Cmd("sudo bash /opt/anyconfig/etc/github-cli.sh", false, false)
        return err
      }
      action = ui.Action{Name: "Installing github-cli", Cmd: command_string} 
      actions = append(actions, action)
      ui.RunActions(actions, false)
    }
  }
  // check if github-cli is logged in
  // command := exec.Command("gh", "auth", "status")
  err, _, _ := command.Cmd("gh auth status", false, false)
  if err != nil {
    err, _, _ := command.Cmd("gh auth login", false, true)
    if err != nil {
      out.Error(err)
      os.Exit(1)
    }
  }
}

func Create() {
  err, _, _ := command.Cmd("gh repo create", false, true)
  if err != nil {
      out.Error(err)
      os.Exit(1)
  }
}

func Clone() string {
  repos := getRepos()
  input := ""
  prompt := &survey.Select{
    Message: "Select Repository:",
    Options: repos,
    PageSize: 10,
  }
  survey.AskOne(prompt, &input)
  path := ui.TextInput("Select Directory to clone into:", "path")
  err, _, _ := command.Cmd("gh repo clone " + input + " " + path + "/" + input, false, true)
  if err != nil {
    out.Error(err)
    os.Exit(1)
}
  return path + "/" + input
}

func getRepos() []string {
  cmd := exec.Command("gh", "repo", "list", "--json", "name")
  out, _ := cmd.Output()
  out = out[:len(out)-1]
  out = out[1:]
  var repos []string
  for _, k := range strings.Split(string(out[:]), ",") {
    var result map[string]string
    json.Unmarshal([]byte(k), &result)
    repos = append(repos, result["name"])
  }
  return repos
}

