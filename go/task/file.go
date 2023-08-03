package task

import (
	"config"
	"out"
	"tools"

	"os"
	"strings"

	"github.com/AlecAivazis/survey/v2"
)


func SelectFiles(config config.AnyConfig) []string {
    repo := config.Repo
    configPath := repo + "/.anyconfig/"
    fileNames := getNamesOfFiles(configPath, true)
    // ged rid if files that dont meet requirements
    filter := filterFiles(fileNames, config)
    fileNames = filter
    fileNames_colored :=  []string{}
    for _, filename := range fileNames {
      _, err := os.Open(configPath + filename)
      if err != nil {
        fileNames_colored = append(fileNames_colored, filename)
      } else {
        fileNames_colored = append(fileNames_colored, out.Style(3, false, filename))
      }
    }
    // create choices from files
    var actionFiles []string
    selectedFileNames := []string{}
    prompt2 := &survey.MultiSelect{
        Message: "Select Tasks to run ",
        Options: fileNames_colored,
    }
    survey.AskOne(prompt2, &selectedFileNames)
    // loop through selected choices
    for _, filename_colored  := range selectedFileNames {
      // get index of filenames
      filename := ""
      for n, k := range fileNames_colored {
        if k == filename_colored {
          filename = fileNames[n]
        }
      }
      _, err := os.Open(configPath + filename)
      if err != nil {
        // if a file, add to list
        actionFiles = append(actionFiles, filename + ".yml")
      } else {
        // if directory, loop through files
        fileNames2 := getNamesOfFiles(configPath + filename, false)
        temp := []string{}
        for _, k := range fileNames2 {
          if k == filename {
            actionFiles = append(actionFiles, filename + "/" + k + ".yml")
          } else {
            temp = append(temp, filename + "/" + k)
          }
        }
        fileNames2 = temp
        // ged rid if files that dont meet requirements
        filter := filterFiles(fileNames2, config)
        fileNames2 = filter
        selectedFileNamesSubDir := []string{}
        prompt2 := &survey.MultiSelect{
          Message: "Select Tasks to run in subDir: " + filename,
            Options: fileNames2,
        }
        survey.AskOne(prompt2, &selectedFileNamesSubDir)
        for _, filename2 := range selectedFileNamesSubDir {
          if filename == filename2 {
            continue
          }
          path := filename2 + ".yml"
          _, err := os.Open(configPath + path)
          if err != nil {
            out.Error("Could not open file: " + configPath + path)
            os.Exit(1)
          } else {
            actionFiles = append(actionFiles, path)
          }
        }
      }
    }
  return actionFiles
}

func SortFiles(files []string, config config.AnyConfig) []string {
  buffer := []string{}
  for _, k := range files {
    task := GetTask(config.Repo + "/.anyconfig/" + k)[0]
    for _, dep := range task.Dependencies{
      if dep.Name == "task" {
        for _, arg := range dep.Args {
          buffer = append(buffer, (arg + ".yml"))
        }
      } else if dep.Name == "user" {
        if dep.Args[0] == "noroot" {
          if tools.GetUser() == "root" {
            out.Error("Can't run task " + k + " as root")
            os.Exit(1)
          }
        }
      }
    }
    buffer = append(buffer, k) 
  }
  // filter buffer
  buffer2 := []string{}
  for _, k := range buffer {
    found := 0
    for _, k2 := range buffer2 {
      if k == k2 {
        found = 1
      }
    }
    if found == 0 {
      buffer2 = append(buffer2, k)
    }
  }
  return buffer2
}

func getNamesOfFiles(path string, allowDir bool) []string {
  items, _ := os.ReadDir(path)
  var action_names []string
  for _, item := range items {
    if item.IsDir() {
      if allowDir == true {
        action_names = append(action_names, item.Name())
      } else {
        out.Error("directorys are not allowed")
        os.Exit(1)
      }
    } else {
      if !strings.Contains(item.Name(), ".yml") {
        out.Error("File " + item.Name() +" is not of type yml")
        os.Exit(1)
      }
      name := strings.TrimSuffix(item.Name(), ".yml")
      action_names = append(action_names, name)
    }
  }
  return action_names
}

func filterFiles(files []string, config config.AnyConfig) []string {
  configPath := config.Repo + "/.anyconfig/"
  temp := []string{}
  for _, k := range files {
    if k == "template" {
      continue
    }
    keep := false
    found := false
    groups := GetTask(configPath + k + ".yml")
    if len(groups) == 0 {
      groups = GetTask(configPath + k + "/" + k + ".yml")
    }
    for _, group := range groups {
      for _, actionDescription := range group.Dependencies{
        if actionDescription.Name == "os" {
          found = true
          if actionDescription.Args[0] == config.Os {
            keep = true
          }
        }
      }
    }
    if keep == true {
      temp = append(temp, k)
    } else if found == false {
      temp = append(temp, k)
    }
  }
  return temp
}
