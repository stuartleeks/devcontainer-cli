#!/bin/bash
set -e

# Expects
# OS=linux or darwin
# ARCH=amd64 or 386

if [[ -z "$OS" ]]; then
echo "$OS"
    echo "OS not set: expected linux or darwin"
    exit 1
fi
if [[ -z "$ARCH" ]]; then
    echo "ARCH not set: expected amd64 or 386"
    exit 1
fi


LATEST_VERSION=$(curl --silent "https://api.github.com/repos/stuartleeks/devcontainer-cli/releases/latest" | grep -Po '"tag_name": "\K.*?(?=")')
echo $LATEST_VERSION
mkdir -p ~/bin
wget https://github.com/stuartleeks/devcontainer-cli/releases/download/${LATEST_VERSION}/devcontainer-cli_${OS}_${ARCH}.tar.gz
tar -C ~/bin -zxvf devcontainer-cli_${OS}_${ARCH}.tar.gz devcontainerx
chmod +x ~/bin/devcontainerx
