package gh

import (
	"command"
	"out"
	"tools"
	"ui"

	"encoding/json"
	"os"
	"os/exec"
	"strings"

	"github.com/AlecAivazis/survey/v2"
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
      ui.RunAction(action, false)
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

func CheckAnyconfigUpdate() {
  if tools.CheckExist("/tmp/anyconfig_update") == true {
    return
  }
  update_needed := false
  // fetch remote
  command_string := "git -C /opt/anyconfig remote update"
  cmd := exec.Command(command_string)
  cmd.Run()

  // check if remote is ahead
  command_string = "cd /opt/anyconfig && git status"
  err, output, _ := command.Cmd(command_string, false, false)
  if err == nil {
    if strings.Contains(output, "behind") == true {
      update_needed = true
    }
  }

  // create file to signal that a update is available
  if update_needed == true {
    command_string = "touch /tmp/anyconfig_update" 
    err, _, _ := command.Cmd(command_string, false, false)
    if err != nil {
      out.Error("Could not create file")
      return
    }
  }
}

func CheckRepoUpdate(repo string) {
  if tools.CheckExist("/tmp/repo_update") == true {
    return
  }
  update_needed := false
  // fetch remote
  command_string := "git -C " + repo + " remote update"
  cmd := exec.Command(command_string)
  cmd.Run()

  // check if remote is ahead
  command_string = "cd " + repo + " && git status"
  err, output, _ := command.Cmd(command_string, false, false)
  if err == nil {
    if strings.Contains(output, "behind") == true {
      update_needed = true
    }
  }

  // create file to signal that a update is available
  if update_needed == true {
    command_string = "touch /tmp/repo_update" 
    err, _, _ := command.Cmd(command_string, false, false)
    if err != nil {
      out.Error("Could not create file")
      return
    }
  }
}

func UpdateRepo(repo string) {
  command_string := "cd " + repo + " && git pull"
  exe := func() error {
    err, output, error := command.Cmd(command_string, false, false)
    if err != nil {
      out.CommandError(command_string, err, output, error)
      os.Exit(1)
    }
    return err
  }
  action := ui.Action{
    Name: "Pulling latest commits",
    Cmd: exe,
  }
  ui.RunAction(action, false)

  command_string = "rm -f /tmp/repo_update"
  command.Cmd(command_string, false, false)
}

func UpdateAnyconfig() {
  // Pulling repo
  command_string := "cd /opt/anyconfig && git pull"
  exe := func() error {
    err, output, error := command.Cmd(command_string, false, false)
    if err != nil {
      out.CommandError(command_string, err, output, error)
      os.Exit(1)
    }
    return err
  }
  action := ui.Action{
    Name: "Pulling latest commits",
    Cmd: exe,
  }
  ui.RunAction(action, false)

  // Building anyconfig
  command_string = "cd /opt/anyconfig/go && go build ."
  exe = func() error {
    err, output, error := command.Cmd(command_string, false, false)
    if err != nil {
      out.CommandError(command_string, err, output, error)
      os.Exit(1)
    }
    return err
  }
  action = ui.Action{
    Name: "Building anyconfig",
    Cmd: exe,
  }
  ui.RunAction(action, false)

  command_string = "rm -f /tmp/anyconfig_update"
  command.Cmd(command_string, false, false)
}
