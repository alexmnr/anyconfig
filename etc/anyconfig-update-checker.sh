#!/bin/bash

repo="$(cat /etc/anyconfig/anyconfig.yml | grep "repo:")"
repo=${repo//'repo: '/}

while : ; do
  # anyconfig
  if test -f "/tmp/anyconfig_update"; then
    echo "anyconfig has an available update"
  else
    # check if update needed
    cd /opt/anyconfig
    git remote update
    if git status | grep -q "behind"; then
      touch /tmp/anyconfig_update
    else
      echo "anyconfig is up to date"
    fi
  fi

  # repo
  if test -f "/tmp/repo_update"; then
    echo "repo has an available update"
  else
    # check if update needed
    cd $repo
    git remote update
    if git status | grep -q "behind"; then
      touch /tmp/repo_update
    else
      echo "repo is up to date"
    fi
  fi

  sleep 10
done
