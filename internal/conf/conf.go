package conf

import (
	"log"
	"strconv"

	"github.com/gurkankaymak/hocon"
)

var RPCProviderWebsocketURL string
var ContractAddress string
var ConfirmationBlocks uint64
var ConfirmationSleepSeconds uint64
var LogPath string

var config *hocon.Config

func Setup() {
	setConfig()
	setRPCProviderWebsocketURL()
	setContractAddress()
	setConfirmationBlocks()
	setConfirmationSleepSeconds()
	setLogPath()
}

func setConfig() {
	conf, err := hocon.ParseResource("application.conf")
	if err != nil {
		log.Panic("error while parsing configuration file: ", err)
	}

	log.Printf("all configuration: %+v", *conf)

	config = conf
}

func setRPCProviderWebsocketURL() {
	rpcProviderWebsocketURL := config.GetString("rpc-provider.websocket-url")
	if len(rpcProviderWebsocketURL) == 0 {
		log.Panic("Infura websocket URL environment variable not found")
	}

	RPCProviderWebsocketURL = rpcProviderWebsocketURL
}

func setContractAddress() {
	contractAddress := config.GetString("contract.address")
	if len(contractAddress) == 0 {
		log.Panic("Contract address environment variable not found")
	}

	ContractAddress = contractAddress
}

func setConfirmationBlocks() {
	confirmationBlocksString := config.GetString("confirmation.blocks")
	if len(confirmationBlocksString) == 0 {
		log.Panic("Confirmed threshold environment variable not found")
	}

	confirmationBlocks, err := strconv.ParseUint(confirmationBlocksString, 10, 64)
	if err != nil {
		log.Panicf("Could not convert confirmed threshold environment variable to uint: %+v\n", err)
	}

	ConfirmationBlocks = confirmationBlocks
}

func setConfirmationSleepSeconds() {
	confirmationSleepSecondsString := config.GetString("confirmation.sleep-seconds")
	if len(confirmationSleepSecondsString) == 0 {
		log.Panic("Sleep interval seconds environment variable not found")
	}

	confirmationSleepSeconds, err := strconv.ParseUint(confirmationSleepSecondsString, 10, 64)
	if err != nil {
		log.Panicf("Could not convert sleep interval seconds environment variable to uint: %+v\n", err)
	}

	ConfirmationSleepSeconds = confirmationSleepSeconds
}

func setLogPath() {
	LogPath = config.GetString("log-path")
}
