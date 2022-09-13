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
	StarNotaryAPIHost        string
	StarNotaryAPIPort        string
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
	c.setStarNotaryAPIHost()
	c.setStarNotaryAPIPort()
	c.setLogPath()
}

func (c *conf) setConfig() {
	hocon, err := hocon.ParseResource("application.conf")
	if err != nil {
		log.Panic("could not parse configuration file: ", err.Error())
	}

	log.Printf("configurations: %+v", *hocon)

	c.hocon = hocon
}

func (c *conf) setRPCProviderWebsocketURL() {
	rpcProviderWebsocketURL := c.hocon.GetString("rpc-provider.websocket-url")
	if len(rpcProviderWebsocketURL) == 0 {
		log.Panic("infura websocket url environment variable not found")
	}

	c.RPCProviderWebsocketURL = rpcProviderWebsocketURL
}

func (c *conf) setContractAddress() {
	contractAddress := c.hocon.GetString("contract.address")
	if len(contractAddress) == 0 {
		log.Panic("contract address environment variable not found")
	}

	c.ContractAddress = contractAddress
}

func (c *conf) setConfirmationBlocks() {
	confirmationBlocksString := c.hocon.GetString("confirmation.blocks")
	if len(confirmationBlocksString) == 0 {
		log.Panic("confirmed threshold environment variable not found")
	}

	confirmationBlocks, err := strconv.ParseUint(confirmationBlocksString, 10, 64)
	if err != nil {
		log.Panic("could not convert confirmed threshold environment variable to uint: ", err.Error())
	}

	c.ConfirmationBlocks = confirmationBlocks
}

func (c *conf) setConfirmationSleepSeconds() {
	confirmationSleepSecondsString := c.hocon.GetString("confirmation.sleep-seconds")
	if len(confirmationSleepSecondsString) == 0 {
		log.Panic("sleep interval seconds environment variable not found")
	}

	confirmationSleepSeconds, err := strconv.ParseUint(confirmationSleepSecondsString, 10, 64)
	if err != nil {
		log.Panic("could not convert sleep interval seconds environment variable to uint: ", err.Error())
	}

	c.ConfirmationSleepSeconds = confirmationSleepSeconds
}

func (c *conf) setStarNotaryAPIHost() {
	starNotaryAPIHost := c.hocon.GetString("star-notary-api.host")
	if len(starNotaryAPIHost) == 0 {
		log.Panic("star notary api host environment variable not found")
	}

	c.StarNotaryAPIHost = starNotaryAPIHost
}

func (c *conf) setStarNotaryAPIPort() {
	starNotaryAPIPort := c.hocon.GetString("star-notary-api.port")
	if len(starNotaryAPIPort) == 0 {
		log.Panic("star notary api port environment variable not found")
	}

	c.StarNotaryAPIPort = starNotaryAPIPort
}

func (c *conf) setLogPath() {
	c.LogPath = c.hocon.GetString("log.path")
}
