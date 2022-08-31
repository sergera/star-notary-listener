package conf

import (
	"log"
	"strconv"
	"sync"

	"github.com/gurkankaymak/hocon"
)

var once sync.Once
var instance *conf

type conf struct {
	hocon                    *hocon.Config
	RPCProviderWebsocketURL  string
	ContractAddress          string
	ConfirmationBlocks       uint64
	ConfirmationSleepSeconds uint64
	BackendHost              string
	BackendPort              string
	LogPath                  string
}

func GetConf() *conf {
	once.Do(func() {
		var c *conf = &conf{}
		c.setup()
		instance = c
	})
	return instance
}

func (c *conf) setup() {
	c.setConfig()
	c.setRPCProviderWebsocketURL()
	c.setContractAddress()
	c.setConfirmationBlocks()
	c.setConfirmationSleepSeconds()
	c.setBackendHost()
	c.setBackendPort()
	c.setLogPath()
}

func (c *conf) setConfig() {
	hocon, err := hocon.ParseResource("application.conf")
	if err != nil {
		log.Panic("error while parsing configuration file: ", err)
	}

	log.Printf("configurations: %+v", *hocon)

	c.hocon = hocon
}

func (c *conf) setRPCProviderWebsocketURL() {
	rpcProviderWebsocketURL := c.hocon.GetString("rpc-provider.websocket-url")
	if len(rpcProviderWebsocketURL) == 0 {
		log.Panic("Infura websocket URL environment variable not found")
	}

	c.RPCProviderWebsocketURL = rpcProviderWebsocketURL
}

func (c *conf) setContractAddress() {
	contractAddress := c.hocon.GetString("contract.address")
	if len(contractAddress) == 0 {
		log.Panic("Contract address environment variable not found")
	}

	c.ContractAddress = contractAddress
}

func (c *conf) setConfirmationBlocks() {
	confirmationBlocksString := c.hocon.GetString("confirmation.blocks")
	if len(confirmationBlocksString) == 0 {
		log.Panic("Confirmed threshold environment variable not found")
	}

	confirmationBlocks, err := strconv.ParseUint(confirmationBlocksString, 10, 64)
	if err != nil {
		log.Panicf("Could not convert confirmed threshold environment variable to uint: %+v\n", err)
	}

	c.ConfirmationBlocks = confirmationBlocks
}

func (c *conf) setConfirmationSleepSeconds() {
	confirmationSleepSecondsString := c.hocon.GetString("confirmation.sleep-seconds")
	if len(confirmationSleepSecondsString) == 0 {
		log.Panic("Sleep interval seconds environment variable not found")
	}

	confirmationSleepSeconds, err := strconv.ParseUint(confirmationSleepSecondsString, 10, 64)
	if err != nil {
		log.Panicf("Could not convert sleep interval seconds environment variable to uint: %+v\n", err)
	}

	c.ConfirmationSleepSeconds = confirmationSleepSeconds
}

func (c *conf) setBackendHost() {
	backendHost := c.hocon.GetString("backend.host")
	if len(backendHost) == 0 {
		log.Panic("Backend host environment variable not found")
	}

	c.BackendHost = backendHost
}

func (c *conf) setBackendPort() {
	backendPort := c.hocon.GetString("backend.port")
	if len(backendPort) == 0 {
		log.Panic("Backend port environment variable not found")
	}

	c.BackendPort = backendPort
}

func (c *conf) setLogPath() {
	c.LogPath = c.hocon.GetString("log-path")
}
