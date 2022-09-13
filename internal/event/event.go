package event

import (
	"context"
	"math/big"
	"strconv"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/sergera/star-notary-listener/internal/backend"
	"github.com/sergera/star-notary-listener/internal/conf"
	"github.com/sergera/star-notary-listener/internal/eth"
	"github.com/sergera/star-notary-listener/internal/gocontracts/starnotary"
	"github.com/sergera/star-notary-listener/internal/logger"
)

var eventSignatureToType = map[string]string{
	"0xc6c75e3dff0d786834a52d041dd27162b7b18821da7d44f32eda867409aff50b": "Created",
	"0x5fcff9840e521b42af6784f574cc1da0706118211c0a01958aeea8b03cb2382b": "ChangedName",
	"0xeef8701c784dcc5b12eb5ce2687a9e42d1d94b6e81f660dcb84b51554c37f082": "PutForSale",
	"0x288d34a7c93145667176da74757455994ea9d6ab6fe10c04f3c9c5d5ba78df78": "RemovedFromSale",
	"0xae92ab4b6f8f401ead768d3273e6bb937a13e39827d19c6376e8fd4512a05d9a": "Sold",
}

func Listen() {
	conf := conf.GetConf()
	eth := eth.GetEth()

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
			logger.Info("Created event to list", logger.Object("event", &genericCreated))
		case changedNameEvent := <-changedNameResChan:
			genericChangedName := changedNameToGeneric(*changedNameEvent)
			insertEventByBlockNumber(genericChangedName)
			logger.Info("Changed name event to list", logger.Object("event", &genericChangedName))
		case putForSaleEvent := <-putForSaleResChan:
			genericPutForSale := putForSaleToGeneric(*putForSaleEvent)
			insertEventByBlockNumber(genericPutForSale)
			logger.Info("Put for sale event to list", logger.Object("event", &genericPutForSale))
		case removedFromSaleEvent := <-removedFromSaleResChan:
			genericRemovedFromSale := removedFromSaleToGeneric(*removedFromSaleEvent)
			insertEventByBlockNumber(genericRemovedFromSale)
			logger.Info("Removed from sale event to list", logger.Object("event", &genericRemovedFromSale))
		case soldEvent := <-soldResChan:
			genericSold := soldToGeneric(*soldEvent)
			insertEventByBlockNumber(genericSold)
			logger.Info("Sold event to list", logger.Object("event", &genericSold))
		default:
			if len(unconfirmedEventsList) > 0 {
				latestBlock, err := eth.Client.BlockNumber(context.Background())
				if err != nil {
					logger.Error("Could not update current block number", logger.String("error", err.Error()))
				}
				latestBlockBig, _ := big.NewInt(0).SetString(strconv.FormatUint(latestBlock, 10), 10)
				scrapAndConfirm(latestBlockBig)
				removeLeftoverEvents(latestBlockBig)
				time.Sleep(time.Duration(conf.ConfirmationSleepSeconds) * time.Second)
			}
		}
	}
}

func scrapAndConfirm(latestBlock *big.Int) {
	conf := conf.GetConf()
	eth := eth.GetEth()

	query := ethereum.FilterQuery{
		FromBlock: unconfirmedEventsList[0].blockNumber,
		ToBlock:   nil, /* nil will query to latest block */
		Addresses: []common.Address{
			common.HexToAddress(conf.ContractAddress),
		},
	}

	logs, err := eth.Client.FilterLogs(context.Background(), query)
	if err != nil {
		logger.Error("Could not query contract logs", logger.String("error", err.Error()))
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
		if big.NewInt(0).Sub(latestBlock, event.blockNumber).Cmp(big.NewInt(int64(conf.ConfirmationBlocks))) == -1 {
			/* if latestBlock - eventBlockNumber < confirmationBlocks */
			/* if event is not yet confirmed, ignore it */
			continue
		}
		if !isEventInList(event) {
			/* if event is not in list, ignore it */
			/* subscribed events might arrive after being added to logs */
			/* which would make the event be consumed again upon arrival */
			continue
		}
		block, err := eth.Client.BlockByNumber(context.Background(), event.blockNumber)
		if err != nil {
			/* if fail to get block, return to try again */
			logger.Error("Failed to get block", logger.String("message", err.Error()))
			return
		}
		event.date = time.Unix(int64(block.Time()), 0).Format(time.RFC3339)
		consume(event)
		removeEvents(event)
	}
}

func consume(event genericEvent) {
	var back *backend.Backend = backend.NewBackend()
	switch event.eventType {
	case "Created":
		createdModel := genericToCreatedModel(event)
		logger.Info("Consuming created event", logger.Object("event", &createdModel))
		back.CreateStar(createdModel)
	case "ChangedName":
		changedNameModel := genericToChangedNameModel(event)
		logger.Info("Consuming changed name event", logger.Object("event", &changedNameModel))
		back.ChangeName(changedNameModel)
	case "PutForSale":
		putForSaleModel := genericToPutForSaleModel(event)
		logger.Info("Consuming put for sale event", logger.Object("event", &putForSaleModel))
		back.PutForSale(putForSaleModel)
	case "RemovedFromSale":
		removedFromSaleModel := genericToRemovedFromSaleModel(event)
		logger.Info("Consuming removed from sale event", logger.Object("event", &removedFromSaleModel))
		back.RemoveFromSale(removedFromSaleModel)
	case "Sold":
		soldModel := genericToSoldModel(event)
		logger.Info("Consuming sold event", logger.Object("event", &soldModel))
		back.Sell(soldModel)
	}
}
