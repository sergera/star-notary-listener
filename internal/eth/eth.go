package eth

import (
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/sergera/star-notary-listener/internal/env"
	"github.com/sergera/star-notary-listener/internal/gocontracts/starnotary"
	"github.com/sergera/star-notary-listener/internal/logger"
)

var Client *ethclient.Client
var Contract *starnotary.Starnotary
var ABI *abi.ABI

func Setup() {
	setClient()
	setContract()
	setABI()
}

func setClient() {
	client, err := ethclient.Dial(env.InfuraWebsocketURL)
	if err != nil {
		logger.Panic("Could not dial eth client", logger.String("error", err.Error()))
	}

	Client = client
}

func setContract() {
	contractAddress := common.HexToAddress(env.ContractAddress)
	starNotary, err := starnotary.NewStarnotary(contractAddress, Client)
	if err != nil {
		logger.Panic("Could not instance go contract", logger.String("error", err.Error()))
	}

	Contract = starNotary
}

func setABI() {
	starnotaryABI, err := abi.JSON(strings.NewReader(string(starnotary.StarnotaryMetaData.ABI)))
	if err != nil {
		logger.Panic("Could not read contract ABI", logger.String("error", err.Error()))
	}

	ABI = &starnotaryABI
}
