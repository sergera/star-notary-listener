package eth

import (
	"log"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/sergera/star-notary-listener/internal/env"
	"github.com/sergera/star-notary-listener/internal/starnotary"
)

var Client *ethclient.Client
var Contract *starnotary.Starnotary

func Init() {
	setClient()
	setContract()
}

func setClient() {
	client, err := ethclient.Dial("wss://rinkeby.infura.io/ws/v3/" + env.InfuraProjectID)
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
