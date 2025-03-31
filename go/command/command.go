package command

import (
  "errors"
	"bytes"
	"os"
	"os/exec"
	"out"
	"strings"
	"tools"
)

func PkgInstall(installer string, pkgs string, debug bool) error {
  string := installer + " " + pkgs
  err, output, error := Cmd(string, false, debug)
  if err != nil {
    out.CommandError(string, err, output, error)
  }
  return err
}

func Cmd(input string, ignore bool, debug bool) (error, string, string) {
  if strings.Contains(input, "{ignore}") {
    input = strings.Replace(input, "{ignore}", "", -1)
    ignore = true
  }
  args := ""
  if strings.Contains(input, " && ") {
    split := strings.Split(input, " && ")
    for _, part := range split {
      args = args + part + "; "
    }
  } else {
    args = input
  }
  cmd := exec.Command("/bin/bash", "-c", args)
  if debug == true {
    cmd.Stdin = os.Stdin
    cmd.Stdout = os.Stdout
    cmd.Stderr = os.Stderr
    err := cmd.Run()
    return err, "", ""
  } else {
    if ignore == false {
      var outb, errb bytes.Buffer
      cmd.Stdout = &outb
      cmd.Stderr = &errb
      err := cmd.Run()
      return err, outb.String(), errb.String()
    } else {
      cmd.Run()
      return nil, "", ""
    }
  }
}

func Ln(org string, des string, backup bool) error {
  split := strings.Split(org, "/")
  filename := split[len(split) - 1]
  // check if Origin exists
  if tools.CheckExist(org) == false {
    out.Error("Linking Origin " + org + " does not exist")
    err := errors.New("Error while running linking function")
    return err
  }
  
  var check string
  if string(des[len(des) - 1]) == "/" {
    check = des + filename
  } else {
    check = des + "/" + filename
  }
  // check if Destination exists
  if tools.CheckExist(check) == true {
    if backup == true {
      err := Backup(check)
      if err != nil {
        out.Error("Could not backup directory")
        return err
      }
    } else {
      return nil
    }
  }
  // linking
  string := "ln -s " + org + " " + des
  err, _, _ := Cmd(string, false, false)
  if err != nil {
    string := "sudo ln -s " + org + " " + des
    err, output, error := Cmd(string, false, false)
    if err != nil {
      out.CommandError(string, err, output, error)
      return err
    }
  }
  return nil
}

func Mkdir(input string, backup bool) error {
  if tools.CheckExist(input) == true {
    if backup == true {
      // Backup
      err := Backup(input)
      if err != nil {
        out.Error("Could not backup directory")
        return err
      }
    } else {
      return nil
    }
  }
  string := "mkdir " + input
  err, output, error := Cmd(string, false, false)
  if err != nil {
    out.CommandError(string, err, output, error)
  }
  return err
}

func Backup(input string) error {
  if tools.CheckExist(input) == false {
    // remove file if not already gone
		removeFile := strings.TrimSuffix(input, "/")
    string := "sudo rm -rf " + removeFile
    err, output, error := Cmd(string, false, false)
    if err != nil {
      out.CommandError(string, err, output, error)
    }
    return nil
  }
  oldDir := tools.GetHomeDir() + "/.old/"
  // create .old Dir
  if tools.CheckExist(oldDir)  == false {
    err := Mkdir(oldDir, false) 
    if err != nil {
      out.Error("Could not create .old Directory")
      return err
    }
  }
  // check if backuped file already exists
	strippedInput := strings.TrimSuffix(input, "/")
	strippedInput = strings.TrimPrefix(strippedInput, "/")
  split := strings.Split(input, "/")
  filename := split[len(split) - 1]
  if tools.CheckExist(oldDir + "/" + filename)  == true {
    string := "rm -rf " + oldDir + "/" + filename
    err, output, error := Cmd(string, false, false)
    if err != nil {
      out.CommandError(string, err, output, error)
      return err
    }
  }
  // backup file
  string := "mv " + input + " " + oldDir 
  err, _, _ := Cmd(string, false, false)
  if err != nil {
    string := "sudo mv " + input + " " + oldDir 
    err, output, error := Cmd(string, false, false)
    if err != nil {
      out.CommandError(string, err, output, error)
      return err
    }
  }
  return nil
}

func Cp(org string, des string, backup bool) error {
  split := strings.Split(org, "/")
  filename := split[len(split) - 1]
  // check if Origin exists
  if tools.CheckExist(org) == false {
    out.Error("Copying Origin " + org + " does not exist")
    err := errors.New("Error while running Copy function")
    return err
  }
  var check string
  if string(des[len(des) - 1]) == "/" {
    check = des + filename
  } else {
    check = des + "/" + filename
  }
  // check if Destination exists
  if tools.CheckExist(check) == true {
    if backup == true {
      err := Backup(check)
      if err != nil {
        out.Error("Could not backup file or directory")
        return err
      }
    } else {
      return nil
    }
  }
  // copying
  string := "cp -r " + org + " " + des
  err, _, _ := Cmd(string, false, false)
  if err != nil {
    string := "sudo cp -r " + org + " " + des
    err, output, error := Cmd(string, false, false)
    if err != nil {
      out.CommandError(string, err, output, error)
      return err
    }
  }
  return nil
}

func Mv(org string, des string, backup bool) error {
  split := strings.Split(org, "/")
  filename := split[len(split) - 1]
  // check if Origin exists
  if tools.CheckExist(org) == false {
    out.Error("Moving Origin " + org + " does not exist")
    err := errors.New("Error while running Moving function")
    return err
  }
  var check string
  if string(des[len(des) - 1]) == "/" {
    check = des + filename
  } else {
    check = des + "/" + filename
  }
  // check if Destination exists
  if tools.CheckExist(check) == true {
    if backup == true {
      err := Backup(check)
      if err != nil {
        out.Error("Could not backup file or directory")
        return err
      }
    } else {
      return nil
    }
  }
  // moving
  string := "mv -f " + org + " " + des
  err, _, _ := Cmd(string, false, false)
  if err != nil {
    string := "sudo mv -f " + org + " " + des
    err, output, error := Cmd(string, false, false)
    if err != nil {
      out.CommandError(string, err, output, error)
      return err
    }
  }
  return nil
}
