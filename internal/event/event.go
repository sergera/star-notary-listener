package event

import (
	"context"
	"log"
	"math/big"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/sergera/star-notary-listener/internal/env"
	"github.com/sergera/star-notary-listener/internal/eth"
	"github.com/sergera/star-notary-listener/internal/gocontracts/starnotary"
	"github.com/sergera/star-notary-listener/internal/models"
)

var eventSignatureToName = map[string]string{
	"0xc6c75e3dff0d786834a52d041dd27162b7b18821da7d44f32eda867409aff50b": "Created",
	"0x5fcff9840e521b42af6784f574cc1da0706118211c0a01958aeea8b03cb2382b": "ChangedName",
	"0xeef8701c784dcc5b12eb5ce2687a9e42d1d94b6e81f660dcb84b51554c37f082": "PutForSale",
	"0x288d34a7c93145667176da74757455994ea9d6ab6fe10c04f3c9c5d5ba78df78": "RemovedFromSale",
	"0xae92ab4b6f8f401ead768d3273e6bb937a13e39827d19c6376e8fd4512a05d9a": "Sold",
}

func Listen() {
	createdResChan := make(chan *starnotary.StarnotaryCreated)
	changedNameResChan := make(chan *starnotary.StarnotaryChangedName)
	putForSaleResChan := make(chan *starnotary.StarnotaryPutForSale)
	removedFromSaleResChan := make(chan *starnotary.StarnotaryRemovedFromSale)
	soldResChan := make(chan *starnotary.StarnotarySold)

	defer close(createdResChan)
	defer close(changedNameResChan)
	defer close(putForSaleResChan)
	defer close(removedFromSaleResChan)
	defer close(soldResChan)

	eth.Contract.WatchCreated(&bind.WatchOpts{Start: nil, Context: context.Background()}, createdResChan)
	eth.Contract.WatchChangedName(&bind.WatchOpts{Start: nil, Context: context.Background()}, changedNameResChan)
	eth.Contract.WatchPutForSale(&bind.WatchOpts{Start: nil, Context: context.Background()}, putForSaleResChan)
	eth.Contract.WatchRemovedFromSale(&bind.WatchOpts{Start: nil, Context: context.Background()}, removedFromSaleResChan)
	eth.Contract.WatchSold(&bind.WatchOpts{Start: nil, Context: context.Background()}, soldResChan)

	for {
		select {
		case createdLog := <-createdResChan:
			createdEvent := contractCreatedToEvent(*createdLog)
			insertEvent(createdEvent)
			log.Printf("Created Event To List: %+v\n\n", createdEvent)
		case changedNameLog := <-changedNameResChan:
			changedNameEvent := contractChangedNameToEvent(*changedNameLog)
			insertEvent(changedNameEvent)
			log.Printf("Changed Name Event To List: %+v\n\n", changedNameEvent)
		case putForSaleLog := <-putForSaleResChan:
			putForSaleEvent := contractPutForSaleToEvent(*putForSaleLog)
			insertEvent(putForSaleEvent)
			log.Printf("Put For Sale Event To List: %+v\n\n", putForSaleEvent)
		case removedFromSaleLog := <-removedFromSaleResChan:
			removedFromSaleEvent := contractRemovedFromSaleToEvent(*removedFromSaleLog)
			insertEvent(removedFromSaleEvent)
			log.Printf("Removed From Sale Event To List: %+v\n\n", removedFromSaleEvent)
		case soldLog := <-soldResChan:
			soldEvent := contractSoldToEvent(*soldLog)
			insertEvent(soldEvent)
			log.Printf("Sold Event To List: %+v\n\n", soldEvent)
		default:
			if len(subscribedEventsList) > 0 {
				currentBlock, err := eth.Client.BlockNumber(context.Background())
				if err != nil {
					log.Printf("Could not update current block number: %+v\n\n", err)
				}
				scrapAndConsume(currentBlock)
				removeOrphanedEvents(currentBlock)
			}
		}
	}
}

func scrapAndConsume(currentBlock uint64) {
	query := ethereum.FilterQuery{
		FromBlock: big.NewInt(int64(subscribedEventsList[0].BlockNumber)),
		ToBlock:   nil, /* nil will query to latest block */
		Addresses: []common.Address{
			common.HexToAddress(env.ContractAddress),
		},
	}

	logs, err := eth.Client.FilterLogs(context.Background(), query)
	if err != nil {
		log.Printf("Could not query contract logs: %+v\n\n", err)
	}

	for _, vLog := range logs {
		listenedEvent := eventSignatureToName[vLog.Topics[0].Hex()]
		if len(listenedEvent) == 0 {
			continue
		}
		event := logToEvent(vLog)
		if event.Removed {
			removeDuplicateEvents(event)
			continue
		}
		if currentBlock-event.BlockNumber < env.ConfirmedThreshold {
			continue
		}
		consume(event)
		removeDuplicateEvents(event)
	}
}

func consume(event models.Event) {
	switch event.EventName {
	case "Created":
		createdModel := eventToCreatedEvent(event)
		log.Printf("Consuming created event: %+v\n\n", createdModel)
	case "ChangedName":
		changedNameModel := eventToChangedNameEvent(event)
		log.Printf("Consuming changed name event: %+v\n\n", changedNameModel)
	case "PutForSale":
		putForSaleModel := eventToPutForSaleEvent(event)
		log.Printf("Consuming put for sale event: %+v\n\n", putForSaleModel)
	case "RemovedFromSale":
		removedFromSaleModel := eventToRemovedFromSaleEvent(event)
		log.Printf("Consuming removed from sale event: %+v\n\n", removedFromSaleModel)
	case "Sold":
		soldModel := eventToSoldEvent(event)
		log.Printf("Consuming sold event: %+v\n\n", soldModel)
	}
}
