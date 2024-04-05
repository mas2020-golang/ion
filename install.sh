#!/usr/bin/env bash

#########################
# Repo specific content #
#########################

export OWNER=mas2020-golang
export REPO=ion
export BINLOCATION="/usr/local/bin"
export SUCCESS_CMD="$BINLOCATION/$REPO version"

# -- COLORS
export STOP_COLOR="\e[0m"
# color for a main activity
export ACTIVITY="\e[38;5;184m"
# color for a sub activity
export SUB_ACT="\n\n➜"
# export SUB_ACT="\n\n\e[1;34m➜\e[0m"
export DONE="\e[1;32m✔︎\e[0m"
export ERROR="\e[1;31mERROR\e[0m:"
export WARNING="\e[38;5;216mWARNING\e[0m:"

###############################
# Content common across repos #
###############################
#set -x
printf "${ACTIVITY}%s${STOP_COLOR}" "installing the $REPO application..."
version=$(curl -sI https://github.com/$OWNER/$REPO/releases/latest | grep -i "location:" | awk -F"/" '{ printf "%s", $NF }' | tr -d '\r')
#set -x
#version=${version:1}
printf "\n✔︎ selected version for %s is '%q'" $REPO $version
if [ ! $version ]; then
  echo "Failed while attempting to install $REPO. Please manually install:"
  echo ""
  echo "1. Open your web browser and go to https://github.com/$OWNER/$REPO/releases"
  echo "2. Download the latest release for your platform. Extract it and call it '$REPO'."
  echo "3. chmod +x ./$REPO"
  echo "4. mv ./$REPO $BINLOCATION"
  exit 1
fi

hasCli() {
  hasCurl=$(which curl)
  if [ "$?" = "1" ]; then
    echo "You need curl to use this script."
    exit 1
  fi
}

getPackage() {
  uname=$(uname)
  userid=$(id -u)

  suffix=""
  case $uname in
  "Darwin")
    arch=$(uname -m)
    case $arch in
    "x86_64")
      suffix="Darwin_x86_64"
      ;;
    esac
    case $arch in
    "arm64")
      suffix="Darwin_arm64"
      ;;
    esac
    ;;

  "MINGW"*)
    echo
    echo "================================================================================================"
    echo "  The script doesn't provide an automatic Windows installation."
    echo "  Pls connect to the official release page (https://github.com/mas2020-golang/ion/releases)"
    echo "  and download the version compatible with you Windows operating system."
    echo "  Unzip the file and place it into a directory that is in the Windows CLASSPATH."
    echo "================================================================================================"
    echo
    exit 0

    ;;
  "Linux")
    arch=$(uname -m)
    case $arch in
    "aarch64")
      suffix="Linux_arm64"
      ;;
    esac
    case $arch in
    "x86_64")
      suffix="Linux_x86_64"
      ;;
    esac
    ;;
  esac
  targetFile="/tmp/$REPO_$version_$suffix.tar.gz"
  downloadFile="${REPO}_${suffix}.tar.gz"
  printf "\n✔︎ the file to download is '%q'" "${downloadFile}"

  if [ "$userid" != "0" ]; then
    targetFile="$(pwd)/$REPO_$version_$suffix.tar.gz"
  fi

  if [ -e "$targetFile" ]; then
    rm "$targetFile"
  fi

  url="https://github.com/$OWNER/$REPO/releases/download/$version/${downloadFile}"
  printf "${SUB_ACT} %s ${STOP_COLOR}" "downloading package $url as $targetFile..."

  http_code=$(curl -sSL $url -w '%{http_code}\n' --output "$targetFile")

  # check the file not found
  if [ "$?" != "0" ] || [ ${http_code} -eq 404 ] ; then
    printf "\n${ERROR} no file as a target download has been found"
    exit 1
  fi

  chmod +x "$targetFile"
  printf "\n${DONE} download complete"
  
  printf "${SUB_ACT} %s\n" "extracting the file $targetFile..."
  tar -xvf $targetFile

  # is the location writable?
  if [ ! -w "$BINLOCATION" ]; then
    echo
    echo "============================================================"
    echo "  The script was run as a user who is unable to write"
    echo "  to $BINLOCATION. To complete the installation the"
    echo "  following commands may need to be run manually."
    echo "============================================================"
    echo
    echo "$ sudo mv $REPO $BINLOCATION/$REPO"

    # final operations
    greets $targetFile
  else
    printf "${SUB_ACT} %s" "moving $REPO to $BINLOCATION..."

    if [ ! -w "$BINLOCATION/$REPO" ] && [ -f "$BINLOCATION/$REPO" ]; then
      echo
      echo "================================================================"
      echo "  $BINLOCATION/$REPO already exists and is not writeable"
      echo "  by the current user.  Please adjust the binary ownership"
      echo "  or run sh/bash with sudo."
      echo "================================================================"
      echo
      exit 1
    fi
    mv $REPO $BINLOCATION/$REPO

    if [ "$?" = "0" ]; then
      printf "\n${DONE} new version of $REPO installed to $BINLOCATION"
    fi

    printf "${SUB_ACT} checking the application...\n"
    sleep 0.5
    ${SUCCESS_CMD}
    
    # final operations
    greets $targetFile
  fi

}

greets() {
  if [ -e "$1" ]; then
    rm "$1"
  fi
  printf "${DONE} take a look at these files for further information: README.md, LICENSE, CHANGELOG.md\n"
}

hasCli
getPackage