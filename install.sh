#!/bin/bash

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

# install arch dependencies
if command -v pacman &> /dev/null; then
  if [ "$EUID" -ne 0 ]; then
    sudo pacman -Sy go git --needed --noconfirm
  else
    pacman -Sy sudo go git --needed --noconfirm
  fi
fi
# install debian dependencies
if command -v apt &> /dev/null
then
  if [ "$EUID" -ne 0 ]; then
    sudo apt update && apt install golang-go git -y
  else
    apt update && apt install sudo golang-go git -y
  fi
fi

cd /opt/anyconfig/go
go build .

sudo rm -f /usr/bin/anyconfig &>/dev/null
sudo ln -s /opt/anyconfig/go/anyconfig /usr/bin
