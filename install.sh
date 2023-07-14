#!/bin/bash

if command -v pacman &> /dev/null; then
  os="arch"
  echo "INFO: Detected arch based system"
elif command -v apt &> /dev/null; then
  os="debian"
  echo "INFO: Detected debian based system"
else
  echo "ERROR: OS not supported"
  exit 1
fi

if [ ! -z "$(lscpu | grep 'aarch64')" ]; then
  arc="aarch64"
  echo "INFO: Detected aarch64 architecture"
elif [ ! -z "$(lscpu | grep 'x86_64')" ]; then
  arc="x86_64"
  echo "INFO: Detected x86_64 architecture"
else
  echo "ERROR: architecture not detected"
  exit 1
fi

SCRIPT_DIR=$( cd -- "$( dirname -- "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )

if [ "$SCRIPT_DIR" = "/opt/anyconfig" ]; then
  sudo chown $USER:$USER /opt/anyconfig -R
  cd /opt/anyconfig
else
  sudo rm -rf /opt/anyconfig &>/dev/null
  sudo cp -r $SCRIPT_DIR /opt/anyconfig
  sudo chown $USER:$USER /opt/anyconfig -R
  cd /opt/anyconfig
fi

# check if dependencies are met
dep=true
if ! command -v git &> /dev/null; then
  dep=false
fi
if ! command -v go &> /dev/null; then
  dep=false
fi
if ! command -v wget &> /dev/null; then
  dep=false
fi

if [ "$dep" = "false" ]; then
  # arch install
  if [ "$os" = "arch" ]; then
    if [ "$arc" = "x86_64" ]; then
      sudo pacman -Sy wget go git --needed --noconfirm
    else
      echo "ERROR: anyconfig can't automatically install dependencies on your system, you need to do it manually"
      exit 1
    fi
  # debian install
  elif [ "$os" = "debian" ]; then
    if [ "$arc" = "x86_64" ]; then
      sudo apt update && sudo apt install wget git -y
      sudo rm -rf /usr/local/go &> /dev/null
      cd /tmp
      wget https://go.dev/dl/go1.20.6.linux-amd64.tar.gz
      sudo tar -C /usr/local -xzf go1.20.6.linux-amd64.tar.gz
      if [ ! -z "$(echo $PATH | grep '/usr/local/go/bin')" ]; then
        echo 'PATH="\$PATH:/usr/local/go/bin"' | sudo tee -a /etc/profile
        PATH="$PATH:/usr/local/go/bin"
        export PATH
      fi
    elif [ "$arc" = "aarch64" ]; then
      sudo apt update && sudo apt install wget git -y
      sudo rm -rf /usr/local/go &> /dev/null
      cd /tmp
      wget https://go.dev/dl/go1.20.6.linux-armv6l.tar.gz
      sudo tar -C /usr/local -xzf go1.20.6.linux-armv6l.tar.gz
      if [ ! -z "$(echo $PATH | grep '/usr/local/go/bin')" ]; then
        echo 'PATH="\$PATH:/usr/local/go/bin"' | sudo tee -a /etc/profile
        PATH="$PATH:/usr/local/go/bin"
        export PATH
      fi
    else
      echo "ERROR: anyconfig can't automatically install dependencies on your system, you need to do it manually"
      exit 1
    fi
  fi
else
  echo "INFO: All Dependencies found"
fi

cd /opt/anyconfig/go
go build .

sudo rm -f /usr/bin/anyconfig &>/dev/null
sudo ln -s /opt/anyconfig/go/anyconfig /usr/bin

echo ""
echo "INFO: Installation complete"
