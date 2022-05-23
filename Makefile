.PHONY: help

KERNEL_NAME := $(shell uname -s)
ifeq ($(KERNEL_NAME),Linux)
    OPEN := xdg-open
else ifeq ($(KERNEL_NAME),Darwin)
    OPEN := open
else
    $(error unsupported system: $(KERNEL_NAME))
endif

GO_PATH=~/go
GO_ROOT=/usr/local/go

CONTRACT_NAME=StarNotary
SOLIDITY_VERSION=0.8.11
DEPLOYED_NETWORK=rinkeby
CONTRACT_PACKAGE_NAME=$(shell echo "$(CONTRACT_NAME)" | tr 'A-Z' 'a-z')
PROJECT_ROOT_PATH:=$(shell dirname $(realpath $(lastword $(MAKEFILE_LIST))))
TRUFFLE_PROJECT_ROOT_PATH=~/code/star-notary

# ENVIRONMENT
INFURA_WEBSOCKET_URL=
CONTRACT_ADDRESS=0x623D6e2B1BB45Fb21b96b7CB3AaeE7C627Cd32C9
CONFIRMED_THRESHOLD=12
ORPHANED_THRESHOLD=24
SLEEP_INTERVAL_SECONDS=5

help: ## Print this help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

start: ## Start the application with go run
	CONTRACT_ADDRESS=$(CONTRACT_ADDRESS) INFURA_WEBSOCKET_URL=$(INFURA_WEBSOCKET_URL) CONFIRMED_THRESHOLD=$(CONFIRMED_THRESHOLD) ORPHANED_THRESHOLD=$(ORPHANED_THRESHOLD) SLEEP_INTERVAL_SECONDS=$(SLEEP_INTERVAL_SECONDS) go run cmd/app/main.go

install-eth-tools: ## Run go-ethereum make files to install abigen
	go get github.com/ethereum/go-ethereum && cd $(GO_PATH)/pkg/mod/github.com/ethereum/go-ethereum* && sudo -E env "PATH=$$PATH" make && sudo -E env "PATH=$$PATH" make devtools

generate-abi: ## Compile contract ABI using solc docker image and write to /build/contracts/CONTRACT_NAME/abi
	mkdir -p build && mkdir -p build/contracts && mkdir -p build/contracts/$(CONTRACT_PACKAGE_NAME) && mkdir -p build/contracts/$(CONTRACT_PACKAGE_NAME)/abi && sudo docker pull ethereum/solc:$(SOLIDITY_VERSION) && sudo docker run --rm -v $(PROJECT_ROOT_PATH):/root -v $(TRUFFLE_PROJECT_ROOT_PATH):/truffle ethereum/solc:$(SOLIDITY_VERSION) openzeppelin-solidity=/truffle/node_modules/openzeppelin-solidity --abi /truffle/contracts/$(CONTRACT_NAME).sol --overwrite -o /root/build/contracts/$(CONTRACT_PACKAGE_NAME)/abi

generate-bytecode: ## Compile EVM bytecode using solc docker image and write to /build/contracts/CONTRACT_NAME/bin
	mkdir -p build && mkdir -p build/contracts && mkdir -p build/contracts/$(CONTRACT_PACKAGE_NAME) && mkdir -p build/contracts/$(CONTRACT_PACKAGE_NAME)/bin && sudo docker pull ethereum/solc:$(SOLIDITY_VERSION) && sudo docker run --rm -v $(PROJECT_ROOT_PATH):/root -v $(TRUFFLE_PROJECT_ROOT_PATH):/truffle ethereum/solc:$(SOLIDITY_VERSION) openzeppelin-solidity=/truffle/node_modules/openzeppelin-solidity --bin /truffle/contracts/$(CONTRACT_NAME).sol --overwrite -o /root/build/contracts/$(CONTRACT_PACKAGE_NAME)/bin

generate-go-contract: ## Make Go contract file into /build/contracts/CONTRACT_NAME
	mkdir -p internal && mkdir -p internal/$(CONTRACT_PACKAGE_NAME) && $(GO_PATH)/bin/abigen -bin=build/contracts/$(CONTRACT_PACKAGE_NAME)/bin/$(CONTRACT_NAME).bin --abi=build/contracts/$(CONTRACT_PACKAGE_NAME)/abi/$(CONTRACT_NAME).abi --pkg=$(CONTRACT_PACKAGE_NAME) --out=internal/gocontracts/$(CONTRACT_PACKAGE_NAME)/$(CONTRACT_PACKAGE_NAME).go

generate: generate-abi generate-bytecode generate-go-contract ## Generate ABI, EVM bytecode, and Go contract file
