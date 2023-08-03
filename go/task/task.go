package task

import (
	"fmt"
  "os"
	"gopkg.in/yaml.v2"
)

type TaskGroup []TaskDescription

type TaskDescription struct {
  Dependencies []ActionDescription
  Install []ActionDescription
}

type ActionDescription struct {
  Name string
  Args []string
}

func GetTask(file string) []TaskDescription {
  yamlFile, _ := os.ReadFile(file)
  m := yaml.MapSlice{}
  if err := yaml.Unmarshal(yamlFile, &m); err != nil {
    fmt.Println("Couldn't read yaml file: ", err)
  }
  configs := []TaskDescription{}
  dependencies := []ActionDescription{}
  install := []ActionDescription{}
  for _, groupvalue := range m {
    if groupvalue.Key == "dependencies" {
      dependencies = []ActionDescription{}
      switch group := groupvalue.Value.(type) {
      case yaml.MapSlice:
        for _, taskvalue := range group {
          task := ActionDescription{}
          task.Name = taskvalue.Key.(string)
          switch t := taskvalue.Value.(type) {
          case []interface{}:
            for _, p := range t {
              task.Args = append(task.Args, p.(string))
            }
          }
          dependencies = append(dependencies, task)
        }
      }
    } else if groupvalue.Key == "install" {
      install = []ActionDescription{}
      switch group := groupvalue.Value.(type) {
      case yaml.MapSlice:
        for _, taskvalue := range group {
          task := ActionDescription{}
          task.Name = taskvalue.Key.(string)
          switch t := taskvalue.Value.(type) {
          case []interface{}:
            for _, p := range t {
              task.Args = append(task.Args, p.(string))
            }
          }
          install = append(install, task)
        }
      }
      config := TaskDescription{}
      config.Dependencies = dependencies
      config.Install = install
      configs = append(configs, config)
    }
  }
  return configs
}




