build:
	go build ./cmd/devcontainer

devcontainer:
	docker build -f ./.devcontainer/Dockerfile ./.devcontainer -t devcontainer-cli

devcontainer-release:
ifdef DEVCONTAINER
	$(error This target can only be run outside of the devcontainer as it mounts files and this fails within a devcontainer. Don't worry all it needs is docker)
endif
	@docker run -v ${PWD}:${PWD} \
		--entrypoint /bin/bash \
		--workdir "${PWD}" \
		devcontainer-cli \
		-c "${PWD}/scripts/ci_release.sh"
