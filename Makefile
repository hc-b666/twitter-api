GOLANGCI_LINT := $(shell command -v golangci-lint 2> /dev/null)
CUSTOM_LINT := $(shell command -v ./bin/custom-gcl 2> /dev/null)

setup:
	if [ -z $(GOLANGCI_LINT) ]; then \
		go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@latest; \
	fi

custom-gcl-build: setup
	if [ -z $(CUSTOM_LINT) ]; then \
		$$(which golangci-lint) custom; \
	fi

lint: custom-gcl-build
	bin/custom-gcl run ./...

watch-lint:
	$$(which nodemon) -w . -x "bin/custom-gcl run ./..." -e "go"

.PHONY: run custom-gcl-build lint

.DEFAULT_GOAL:=run
