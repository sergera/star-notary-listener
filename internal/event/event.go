package event

import (
	"context"
	"log"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/sergera/star-notary-listener/internal/env"
	"github.com/sergera/star-notary-listener/internal/eth"
	"github.com/sergera/star-notary-listener/internal/gocontracts/starnotary"
)

var eventSignatureToType = map[string]string{
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
		case createdEvent := <-createdResChan:
			genericCreated := createdToGeneric(*createdEvent)
			insertEventByBlockNumber(genericCreated)
			log.Printf("Created Event To List: %+v\n\n", genericCreated)
		case changedNameEvent := <-changedNameResChan:
			genericChangedName := changedNameToGeneric(*changedNameEvent)
			insertEventByBlockNumber(genericChangedName)
			log.Printf("Changed Name Event To List: %+v\n\n", genericChangedName)
		case putForSaleEvent := <-putForSaleResChan:
			genericPutForSale := putForSaleToGeneric(*putForSaleEvent)
			insertEventByBlockNumber(genericPutForSale)
			log.Printf("Put For Sale Event To List: %+v\n\n", genericPutForSale)
		case removedFromSaleEvent := <-removedFromSaleResChan:
			genericRemovedFromSale := removedFromSaleToGeneric(*removedFromSaleEvent)
			insertEventByBlockNumber(genericRemovedFromSale)
			log.Printf("Removed From Sale Event To List: %+v\n\n", genericRemovedFromSale)
		case soldEvent := <-soldResChan:
			genericSold := soldToGeneric(*soldEvent)
			insertEventByBlockNumber(genericSold)
			log.Printf("Sold Event To List: %+v\n\n", genericSold)
		default:
			if len(unconfirmedEventsList) > 0 {
				currentBlock, err := eth.Client.BlockNumber(context.Background())
				if err != nil {
					log.Printf("Could not update current block number: %+v\n\n", err)
				}
				scrapAndConfirm(currentBlock)
				removeOrphanedEvents(currentBlock)
				time.Sleep(time.Duration(env.SleepIntervalSeconds) * time.Second)
			}
		}
	}
}

func scrapAndConfirm(currentBlock uint64) {
	query := ethereum.FilterQuery{
		FromBlock: big.NewInt(int64(unconfirmedEventsList[0].blockNumber)),
		ToBlock:   nil, /* nil will query to latest block */
		Addresses: []common.Address{
			common.HexToAddress(env.ContractAddress),
		},
	}

	logs, err := eth.Client.FilterLogs(context.Background(), query)
	if err != nil {
		log.Printf("Could not query contract logs: %+v\n\n", err)
	}

	for _, scrappedEvent := range logs {
		listenedEventType := eventSignatureToType[scrappedEvent.Topics[0].Hex()]
		if len(listenedEventType) == 0 {
			/* if event is not listened to, ignore it */
			continue
		}
		event := scrappedToGeneric(scrappedEvent)
		if event.removed {
			/* if event was removed, remove it from list */
			removeEvents(event)
			continue
		}
		if currentBlock-event.blockNumber < env.ConfirmedThreshold {
			/* if event is not yet confirmed, ignore it */
			continue
		}
		if !isEventInList(event) {
			/* if event is not in list, ignore it */
			/* subscribed events might arrive after being added to logs */
			/* which would make the event be consumed again upon arrival */
			continue
		}
		consume(event)
		removeEvents(event)
	}
}

func consume(event genericEvent) {
	switch event.eventType {
	case "Created":
		createdModel := genericToCreatedModel(event)
		log.Printf("Consuming created event: %+v\n\n", createdModel)
	case "ChangedName":
		changedNameModel := genericToChangedNameModel(event)
		log.Printf("Consuming changed name event: %+v\n\n", changedNameModel)
	case "PutForSale":
		putForSaleModel := genericToPutForSaleModel(event)
		log.Printf("Consuming put for sale event: %+v\n\n", putForSaleModel)
	case "RemovedFromSale":
		removedFromSaleModel := genericToRemovedFromSaleModel(event)
		log.Printf("Consuming removed from sale event: %+v\n\n", removedFromSaleModel)
	case "Sold":
		soldModel := genericToSoldModel(event)
		log.Printf("Consuming sold event: %+v\n\n", soldModel)
	}
}
