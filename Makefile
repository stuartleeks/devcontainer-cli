build:
	go build ./cmd/devcontainer

devcontainer:
	docker build -f ./.devcontainer/Dockerfile ./.devcontainer -t devcontainer-cli