.PHONY: help

KERNEL_NAME := $(shell uname -s)
ifeq ($(KERNEL_NAME),Linux)
    OPEN := xdg-open
else ifeq ($(KERNEL_NAME),Darwin)
    OPEN := open
else
    $(error unsupported system: $(KERNEL_NAME))
endif

CONTRACT_NAME=StarNotary
SOLIDITY_VERSION=0.8.11
DEPLOYED_NETWORK=rinkeby
PACKAGE_NAME=$(shell echo "$(CONTRACT_NAME)" | tr 'A-Z' 'a-z')
PROJECT_ROOT_PATH:=$(shell dirname $(realpath $(lastword $(MAKEFILE_LIST))))
TRUFFLE_PROJECT_ROOT_PATH=~/code/star-notary
GO_PATH=~/go
GO_ROOT=/usr/local/go

help: ## Print this help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

install-eth-tools: ## Run go-ethereum make files to install abigen
	go get github.com/ethereum/go-ethereum && cd $(GO_PATH)/pkg/mod/github.com/ethereum/go-ethereum* && sudo -E env "PATH=$$PATH" make && sudo -E env "PATH=$$PATH" make devtools

generate-abi: ## Compile contract ABI using solc docker image and write to /build/abis
	docker pull ethereum/solc:$(SOLIDITY_VERSION) && docker run --rm -v $(PROJECT_ROOT_PATH):/root -v $(TRUFFLE_PROJECT_ROOT_PATH):/truffle ethereum/solc:$(SOLIDITY_VERSION) openzeppelin-solidity=/truffle/node_modules/openzeppelin-solidity --abi /truffle/contracts/$(CONTRACT_NAME).sol --overwrite -o /root/build/

generate-bytecode: ## Make EVM bytecode into /build
	docker pull ethereum/solc:$(SOLIDITY_VERSION) && docker run --rm -v $(PROJECT_ROOT_PATH):/root -v $(TRUFFLE_PROJECT_ROOT_PATH):/truffle ethereum/solc:$(SOLIDITY_VERSION) openzeppelin-solidity=/truffle/node_modules/openzeppelin-solidity --bin /truffle/contracts/$(CONTRACT_NAME).sol --overwrite -o /root/build/

generate-go-contract: ## Make Go contract file into root directory 
	mkdir -p $(PACKAGE_NAME) && $(GO_PATH)/bin/abigen -bin=./build/${CONTRACT_NAME}.bin --abi=./build/$(CONTRACT_NAME).abi --pkg=$(PACKAGE_NAME) --out=$(PACKAGE_NAME)/$(CONTRACT_NAME).go

generate: generate-abi generate-bytecode generate-go-contract ## Generate ABI, EVM bytecode, and Go contract file
