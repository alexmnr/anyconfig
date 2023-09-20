#!/bin/bash
# check if dependencies are met
SCRIPT_DIR=$( cd -- "$( dirname -- "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )
version="1.21.1"
dep=true
missing=""
if ! command -v git &> /dev/null; then
  dep=false
  missing=$(echo "git $missing")
fi
if ! command -v wget &> /dev/null; then
  dep=false
  missing=$(echo "wget $missing")
fi
if ! command -v curl &> /dev/null; then
  dep=false
  missing=$(echo "curl $missing")
fi
if ! command -v sudo &> /dev/null; then
  dep=false
  missing=$(echo "sudo $missing")
fi

if [ "$dep" = "false" ]; then
  echo "INFO: Not all dependencies were met, please install following packages: $missing"
  exit
fi

# check if go is installed
if ! command -v go &> /dev/null; then
  go=false
  echo "INFO: go not found, trying to install..."
else
  go=true
fi

# install go if necessary
if [ "$go" = "false" ]; then
  # get info about cpu
  if [ ! -z "$(lscpu | grep 'aarch64')" ]; then
    arc="aarch64"
    echo "INFO: Detected aarch64 architecture"
  elif [ ! -z "$(lscpu | grep 'armv7l')" ]; then
    arc="armv7l"
    echo "INFO: Detected armv7l architecture"
  elif [ ! -z "$(lscpu | grep 'x86_64')" ]; then
    arc="x86_64"
    echo "INFO: Detected x86_64 architecture"
  else
    echo "ERROR: architecture not detected"
    exit 1
  fi

  # install go from source
  if [ "$arc" = "x86_64" ]; then
    sudo rm -rf /usr/local/go /usr/bin/go /usr/bin/gofmt &> /dev/null
    cd /tmp
    link="https://go.dev/dl/go$version.linux-amd64.tar.gz"
    wget $link
    sudo tar -C /usr/local -xzf go$version.linux-amd64.tar.gz
    sudo rm -f /usr/bin/go
    sudo rm -f /usr/bin/gofmt
    sudo ln -s /usr/local/go/bin/go /usr/bin
    sudo ln -s /usr/local/go/bin/gofmt /usr/bin
  elif [ "$arc" = "armv7l" ]; then
    sudo rm -rf /usr/local/go /usr/bin/go /usr/bin/gofmt &> /dev/null
    cd /tmp
    link="https://go.dev/dl/go$version.linux-armv6l.tar.gz"
    wget $link
    sudo tar -C /usr/local -xzf go$version.linux-armv6l.tar.gz
    sudo rm -f /usr/bin/go
    sudo rm -f /usr/bin/gofmt
    sudo ln -s /usr/local/go/bin/go /usr/bin
    sudo ln -s /usr/local/go/bin/gofmt /usr/bin
  elif [ "$arc" = "aarch64" ]; then
    sudo rm -rf /usr/local/go /usr/bin/go /usr/bin/gofmt &> /dev/null
    cd /tmp
    link="https://go.dev/dl/go$version.linux-arm64.tar.gz"
    wget $link
    sudo tar -C /usr/local -xzf go$version.linux-arm64.tar.gz
    sudo rm -f /usr/bin/go
    sudo rm -f /usr/bin/gofmt
    sudo ln -s /usr/local/go/bin/go /usr/bin
    sudo ln -s /usr/local/go/bin/gofmt /usr/bin
  else
    echo "ERROR: Can't automatically install go on your system, you need to do it manually"
    exit 1
  fi

fi

echo "INFO: All dependencies were met!"

# move current dir to /opt/anyconfig
echo "INFO: Installing anyconfig in /opt/anyconfig..."
if [ "$SCRIPT_DIR" = "/opt/anyconfig" ]; then
  sudo chown $USER:$USER /opt/anyconfig -R
  cd /opt/anyconfig
else
  sudo rm -rf /opt/anyconfig &>/dev/null
  sudo cp -r $SCRIPT_DIR /opt/anyconfig
  sudo chown $USER:$USER /opt/anyconfig -R
  cd /opt/anyconfig
fi

echo "INFO: Building anyconfig..."

cd /opt/anyconfig/go
go build .

sudo rm -f /usr/bin/anyconfig &>/dev/null
sudo ln -s /opt/anyconfig/go/anyconfig /usr/bin

echo "INFO: Linking anyconfig-update service..."
sudo rm -f /usr/lib/systemd/user/anyconfig-update.service
sudo ln -s /opt/anyconfig/etc/anyconfig-update.service /usr/lib/systemd/user/
systemctl --user daemon-reload
systemctl --user enable anyconfig-update
systemctl --user start anyconfig-update

echo "INFO: done!"
