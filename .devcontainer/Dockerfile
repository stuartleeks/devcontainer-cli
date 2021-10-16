#-------------------------------------------------------------------------------------------------------------
# Copyright (c) Microsoft Corporation. All rights reserved.
# Licensed under the MIT License. See https://go.microsoft.com/fwlink/?linkid=2090316 for license information.
#-------------------------------------------------------------------------------------------------------------

FROM golang:1.15-stretch

# Avoid warnings by switching to noninteractive
ENV DEBIAN_FRONTEND=noninteractive

# Configure apt, install packages and tools
RUN apt-get update \
    && apt-get -y install --no-install-recommends apt-utils dialog nano sudo bsdmainutils \
    #
    # Verify git, process tools, lsb-release (common in install instructions for CLIs) installed
    && apt-get -y install git iproute2 procps lsb-release build-essential \
    # Install Release Tools
    #
    # --> RPM used by goreleaser
    && apt install -y rpm 

# This Dockerfile adds a non-root user with sudo access. Use the "remoteUser"
# property in devcontainer.json to use it. On Linux, the container user's GID/UIDs
# will be updated to match your local UID/GID (when using the dockerFile property).
# See https://aka.ms/vscode-remote/containers/non-root-user for details.
ARG USERNAME=vscode
ARG USER_UID=1000
ARG USER_GID=$USER_UID

# Create the user
RUN groupadd --gid $USER_GID $USERNAME \
    && useradd --uid $USER_UID --gid $USER_GID -m $USERNAME \
    && apt-get update \
    && apt-get install -y sudo \
    && echo $USERNAME ALL=\(root\) NOPASSWD:ALL > /etc/sudoers.d/$USERNAME \
    && chmod 0440 /etc/sudoers.d/$USERNAME

# Set default user
USER $USERNAME
RUN mkdir -p ~/.local/bin
ENV PATH /home/${USERNAME}/.local/bin:$PATH

# Set env for tracking that we're running in a devcontainer
ENV DEVCONTAINER=true

# Enable go modules
ENV GO111MODULE=on

# Install Go tools
ARG GO_PLS_VERSION=0.7.2
ARG DLV_VERSION=1.7.2
ARG GO_RELEASER_VERSION=0.180.3
ARG GOLANGCI_LINT_VERSION=1.42.1
RUN \
    # --> Delve for debugging
    go get github.com/go-delve/delve/cmd/dlv@v${DLV_VERSION}\
    # --> Go language server
    && go get golang.org/x/tools/gopls@v${GO_PLS_VERSION} \
    # --> Go symbols and outline for go to symbol support and test support 
    && go get github.com/acroca/go-symbols@v0.1.1 && go get github.com/ramya-rao-a/go-outline@7182a932836a71948db4a81991a494751eccfe77 \
    # --> GolangCI-lint
    && curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin v${GOLANGCI_LINT_VERSION} \
    # --> Go releaser 
    && curl -sfL https://install.goreleaser.com/github.com/goreleaser/goreleaser.sh | sh -s -- v${GO_RELEASER_VERSION} \
    # --> Install junit converter
    && go get github.com/jstemmer/go-junit-report@v0.9.1 \
    && sudo rm -rf /go/src/ \
    && sudo rm -rf /go/pkg

# Switch back to dialog for any ad-hoc use of apt-get
ENV DEBIAN_FRONTEND=dialog

# gh
COPY scripts/gh.sh /tmp/
RUN /tmp/gh.sh

# symlink gh config folder
RUN echo 'if [[ ! -d /home/vscode/.config/gh ]]; then mkdir -p /home/vscode/.config; ln -s /config/gh /home/vscode/.config/gh; fi ' >> ~/.bashrc

ARG DOCKER_GROUP_ID

# docker-from-docker
COPY scripts/docker-client.sh /tmp/
RUN /tmp/docker-client.sh
