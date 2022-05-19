package eth

import (
	"log"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/sergera/star-notary-listener/internal/env"
	"github.com/sergera/star-notary-listener/internal/gocontracts/starnotary"
)

var Client *ethclient.Client
var Contract *starnotary.Starnotary
var ABI *abi.ABI

func Init() {
	setClient()
	setContract()
	setABI()
}

func setClient() {
	client, err := ethclient.Dial(env.InfuraWebsocketURL)
	if err != nil {
		log.Fatal(err)
	}

	Client = client
}

func setContract() {
	contractAddress := common.HexToAddress(env.ContractAddress)
	starNotary, err := starnotary.NewStarnotary(contractAddress, Client)
	if err != nil {
		log.Fatal(err)
	}

	Contract = starNotary
}

func setABI() {
	starnotaryABI, err := abi.JSON(strings.NewReader(string(starnotary.StarnotaryMetaData.ABI)))
	if err != nil {
		log.Fatal("Could not read contract ABI")
	}

	ABI = &starnotaryABI
}
