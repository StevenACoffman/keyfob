#!/usr/bin/env bash

# Installing kefob

if [ ! -x "$(command -v keyfob)" ]; then
  echo "keyfob is not installed, so I'm going to go grab the mac one for you"
  if [! -x "(command -v brew)" ]; then
    KEYFOB_RELEASE='0.3.0'
    echo "Homebrew is not installed, so I'm going to grab the v${KEYFOB_RELEASE} current release from github"
    wget -O - "https://github.com/StevenACoffman/keyfob/releases/download/v${KEYFOB_RELEASE}/keyfob_${KEYFOB_RELEASE}_Darwin_x86_64.tar.gz" | tar xzvf
    mkdir -p /usr/local/bin
    mv keyfob /usr/local/bin
  else
    echo "Using homebrew and tapping StevenACoffman/keyfob"
    brew tap StevenACoffman/keyfob
    brew install keyfob
  fi

fi

filename="${HOME}/.2fa"

if [ -f $filename ]; then
  echo "Snarfing secrets from 2fa for you"
  cat $filename | while read line
  do
    SIZE="$(echo $line | awk '{print $2}')"
    KEY="$(echo $line | awk '{print $1}')"
    VALUE="$(echo $line | awk '{print $3}')"
    echo "Processing $KEY" >/dev/tty
    keyfob add "${KEY}" "${VALUE}"
  done
else
  echo "${filename} does not exist so not automatically copying any keys from 2fa"
fi
