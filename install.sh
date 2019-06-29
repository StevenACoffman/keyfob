#!/usr/bin/env bash

# Script will snarf secrets from 2fa file
if [ ! -x "$(command -v keyfob)" ]; then
  echo "keyfob is not installed, so I'm going to go grab the mac one for you"
  wget -O - https://github.com/StevenACoffman/keyfob/releases/download/v0.1.0/keyfob_0.1.0_Darwin_x86_64.tar.gz | tar xzvf
  mkdir -p /usr/local/bin
  mv keyfob /usr/local/bin
fi

filename="${HOME}/.2fa"

if [ -f $filename ]; then
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
