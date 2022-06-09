SHELL:=/bin/bash

.PHONY: help

KERNEL_NAME := $(shell uname -s)
ifeq ($(KERNEL_NAME),Linux)
    OPEN := xdg-open
else ifeq ($(KERNEL_NAME),Darwin)
    OPEN := open
else
    $(error unsupported system: $(KERNEL_NAME))
endif

TRUFFLE_PROJECT_ROOT_PATH=
SOLIDITY_VERSION=0.8.11
CONTRACT_NAME=StarNotary
CONTRACT_PACKAGE_NAME=$(shell echo "$(CONTRACT_NAME)" | tr 'A-Z' 'a-z')
PROJECT_ROOT_PATH:=$(shell dirname $(realpath $(lastword $(MAKEFILE_LIST))))

help: ## Print this help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

install: ## Install dependencies to go modules cache
	@go mod tidy

run: ## Start the application with go run
	@source ./scripts/env.bash && go run cmd/app/*.go

contract: ## Generate go contract file into internal/gocontracts/CONTRACT_PACKAGE_NAME
	@./scripts/contract/install_abigen.bash
	@./scripts/contract/generate_abi.bash $(SOLIDITY_VERSION) $(CONTRACT_NAME) $(CONTRACT_PACKAGE_NAME) $(TRUFFLE_PROJECT_ROOT_PATH)
	@./scripts/contract/generate_bytecode.bash $(SOLIDITY_VERSION) $(CONTRACT_NAME) $(CONTRACT_PACKAGE_NAME) $(TRUFFLE_PROJECT_ROOT_PATH)
	@./scripts/contract/generate_contract.bash $(CONTRACT_NAME) $(CONTRACT_PACKAGE_NAME)
