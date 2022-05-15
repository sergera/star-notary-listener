package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

func main() {

	infuraProjectId, exists := os.LookupEnv("INFURA_PROJECT_ID")
	if !exists {
		log.Fatal("Infura project id environment variable not found!")
	}

	client, err := ethclient.Dial("wss://rinkeby.infura.io/ws/v3/" + infuraProjectId)
	if err != nil {
		log.Fatal(err)
	}

	contractAddress := common.HexToAddress("0xf130D6968587fb69DE2DC1249293860446fB3823")
	query := ethereum.FilterQuery{
		Addresses: []common.Address{contractAddress},
	}

	logs := make(chan types.Log)

	sub, err := client.SubscribeFilterLogs(context.Background(), query, logs)
	if err != nil {
		log.Fatal(err)
	}

	for {
		select {
		case err := <-sub.Err():
			log.Fatal(err)
		case vLog := <-logs:
			fmt.Println(vLog) // pointer to event log
		}
	}
}
