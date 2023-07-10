package config

import (
	"command"
	"gh"
	"out"
	"tools"
	"ui"

	"io/ioutil"
	"fmt"
	"os"
	"github.com/AlecAivazis/survey/v2"
	"gopkg.in/yaml.v2"
)

type FileStruct struct {
    Os string `yaml:"os"`
    Installer string `yaml:"installer"`
    UnInstaller string `yaml:"uninstaller"`
    Repo string `yaml:"repo"`
}

type AnyConfig struct {
    Os string 
    Installer string 
    UnInstaller string 
    Repo string 
    User string
    HomeDir string
    Debug bool
}

func Init() AnyConfig {
  // load config if exists
  if tools.CheckExist("/etc/anyconfig/anyconfig.yml") == true {
    yamlFile, _ := ioutil.ReadFile("/etc/anyconfig/anyconfig.yml")
    fileConfig := FileStruct{}
    if err := yaml.Unmarshal(yamlFile, &fileConfig); err != nil {
      fmt.Printf(out.ErrorString("Error while reading config: ") + "%v \n", err)
    }
    config := AnyConfig{
      Os: fileConfig.Os,
      Installer: fileConfig.Installer,
      UnInstaller: fileConfig.UnInstaller,
      Repo: fileConfig.Repo,
      User: tools.GetUser(),
      HomeDir: tools.GetHomeDir(),
      Debug: false,
    }
    return config
  } else {
  // ask to create config
    input := ""
    prompt := &survey.Select{
      Message: "No configuration found, want to create it now? ",
      Options: []string{"Yes", "No"},
    }
    survey.AskOne(prompt, &input)
    if input == "No" {
      out.Error("Then go write it yourself!")
      os.Exit(0)
    }
    if tools.GetUser() != "root"{
      command.Cmd("sudo true", false, true)
    }
    // create config
    fileConfig, err := createConfig()
    if err != nil {
      out.Error(err)
      os.Exit(1)
    }
    config := AnyConfig{
      Os: fileConfig.Os,
      Installer: fileConfig.Installer,
      UnInstaller: fileConfig.UnInstaller,
      Repo: fileConfig.Repo,
      User: tools.GetUser(),
      HomeDir: tools.GetHomeDir(),
      Debug: false,
    }
    return config
  }
}

func createConfig() (FileStruct, error) {
  // create directory
  if tools.CheckExist("/etc/anyconfig") == false {
    err, _, _ := command.Cmd("sudo mkdir /etc/anyconfig", false, false)
    if err != nil {
      out.Error("Could not create /etc/anyconfig")
      os.Exit(1)
    }
  }
  // get Info
  osType := tools.GetOS()
  installer := ""
  uninstaller := ""
  if osType == "arch" {
    installer = "sudo pacman --needed --noconfirm -Sy "
    uninstaller = "sudo pacman --noconfirm -R "
  } else if osType == "debian" {
    installer = "sudo apt update && sudo apt install -y "
    uninstaller = "sudo apt remove -y "
  } else {
    fmt.Println("Did not detect OS, exiting...")
    os.Exit(1)
  }
  // get repo
  input := ""
  prompt := &survey.Select{
    Message: "No dotfiles Repository configured, what next?",
    Options: []string{
      out.Style(3, false, "link") + " existing Repository", 
      out.Style(3, false, "create") + " new Repository (needs github authentication)", 
      out.Style(3, false, "clone") + " existing Repository (needs github authentication)",
      "exit",
    },
  }
  repo := ""
  survey.AskOne(prompt, &input)
  if input == "Exit" {
    os.Exit(0)
  } else if input == out.Style(3, false, "link") + " existing Repository" {
    repo = ui.FilePicker("Select Repository:", "/etc")
  } else if input == out.Style(3, false, "create") + " new Repository (needs github authentication)"{
    gh.Config()
    gh.Create()
    repo = gh.Clone()
  } else if input == out.Style(3, false, "clone") + " existing Repository (needs github authentication)"{
    gh.Config()
    repo = gh.Clone()
  }

  // check if repo is configured for anyconfig
  if tools.CheckExist(repo + "/.anyconfig") == false {
    // create .anyconfig directory
    err, _, _ := command.Cmd("mkdir" + repo + "/.anyconfig", false, false)
    if err != nil {
      out.Error("Could not create .anyconfig directory")
      os.Exit(1)
    }
  }
  anyconfig := FileStruct{
    Os: osType,
    Installer: installer,
    UnInstaller: uninstaller,
    Repo: repo,
  }

  yamlData, err := yaml.Marshal(&anyconfig)
  if err != nil {
    out.Error(err)
    os.Exit(0)
  }
  fileName := "/tmp/anyconfig.yml"
  err = ioutil.WriteFile(fileName, yamlData, 0644)
  if err != nil {
    fmt.Printf(out.ErrorString("Error while writing Config: ") + "%v \n", err)
    os.Exit(1)
  }
  command.Cmd("sudo mv /tmp/anyconfig.yml /etc/anyconfig", false, false)
  return anyconfig, nil
}
