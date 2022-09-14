package listener

import (
	"context"
	"math/big"
	"strconv"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/sergera/star-notary-listener/internal/conf"
	"github.com/sergera/star-notary-listener/internal/domain"
	"github.com/sergera/star-notary-listener/internal/eth"
	"github.com/sergera/star-notary-listener/internal/gocontracts/starnotary"
	"github.com/sergera/star-notary-listener/internal/logger"
	"github.com/sergera/star-notary-listener/internal/queue"
	"github.com/sergera/star-notary-listener/internal/service"
)

var eventSignatureToType = map[string]string{
	"0xc6c75e3dff0d786834a52d041dd27162b7b18821da7d44f32eda867409aff50b": "Created",
	"0x5fcff9840e521b42af6784f574cc1da0706118211c0a01958aeea8b03cb2382b": "ChangedName",
	"0xeef8701c784dcc5b12eb5ce2687a9e42d1d94b6e81f660dcb84b51554c37f082": "PutForSale",
	"0x288d34a7c93145667176da74757455994ea9d6ab6fe10c04f3c9c5d5ba78df78": "RemovedFromSale",
	"0xae92ab4b6f8f401ead768d3273e6bb937a13e39827d19c6376e8fd4512a05d9a": "Sold",
}

type Listener struct {
	queue           *queue.EventQueue
	api             *service.StarNotaryAPIService
	contractAddress string
	confirmDelay    uint64
	confirmBlocks   uint64
}

func NewListener() *Listener {
	conf := conf.GetConf()
	return &Listener{
		queue:           queue.NewEventQueue(),
		api:             service.NewStarNotaryAPIService(),
		contractAddress: conf.ContractAddress(),
		confirmDelay:    conf.ConfirmationSleepSeconds(),
		confirmBlocks:   conf.ConfirmationBlocks(),
	}
}

func (l *Listener) Listen() {
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
			l.queue.InsertEventByBlockNumber(genericCreated)
			logger.Info("created event to list", logger.Object("event", &genericCreated))
		case changedNameEvent := <-changedNameResChan:
			genericChangedName := changedNameToGeneric(*changedNameEvent)
			l.queue.InsertEventByBlockNumber(genericChangedName)
			logger.Info("changed name event to list", logger.Object("event", &genericChangedName))
		case putForSaleEvent := <-putForSaleResChan:
			genericPutForSale := putForSaleToGeneric(*putForSaleEvent)
			l.queue.InsertEventByBlockNumber(genericPutForSale)
			logger.Info("put for sale event to list", logger.Object("event", &genericPutForSale))
		case removedFromSaleEvent := <-removedFromSaleResChan:
			genericRemovedFromSale := removedFromSaleToGeneric(*removedFromSaleEvent)
			l.queue.InsertEventByBlockNumber(genericRemovedFromSale)
			logger.Info("removed from sale event to list", logger.Object("event", &genericRemovedFromSale))
		case soldEvent := <-soldResChan:
			genericSold := soldToGeneric(*soldEvent)
			l.queue.InsertEventByBlockNumber(genericSold)
			logger.Info("sold event to list", logger.Object("event", &genericSold))
		default:
			if l.queue.Length() > 0 {
				latestBlock, err := eth.Client.BlockNumber(context.Background())
				if err != nil {
					logger.Error("could not update current block number", logger.String("message", err.Error()))
				}
				latestBlockBig, _ := big.NewInt(0).SetString(strconv.FormatUint(latestBlock, 10), 10)
				l.scrapAndConfirm(latestBlockBig)
				l.queue.RemoveLeftoverEvents(latestBlockBig)
				time.Sleep(time.Duration(l.confirmDelay) * time.Second)
			}
		}
	}
}

func (l *Listener) scrapAndConfirm(latestBlock *big.Int) {
	eth := eth.GetEth()

	query := ethereum.FilterQuery{
		FromBlock: l.queue.FirstEventBlockNumber(),
		ToBlock:   nil, /* nil will query to latest block */
		Addresses: []common.Address{
			common.HexToAddress(l.contractAddress),
		},
	}

	logs, err := eth.Client.FilterLogs(context.Background(), query)
	if err != nil {
		logger.Error("could not query contract logs", logger.String("message", err.Error()))
	}

	for _, scrappedEvent := range logs {
		listenedEventType := eventSignatureToType[scrappedEvent.Topics[0].Hex()]
		if len(listenedEventType) == 0 {
			/* if event is not listened to, ignore it */
			continue
		}
		event := scrappedToGeneric(scrappedEvent)
		if event.Removed {
			/* if event was removed, remove it and duplicates from list */
			l.queue.RemoveEventsLike(event)
			continue
		}
		if big.NewInt(0).Sub(latestBlock, event.BlockNumber).Cmp(new(big.Int).SetUint64(l.confirmBlocks)) == -1 {
			/* if latestBlock - eventBlockNumber < confirmationBlocks */
			/* if event is not yet confirmed, ignore it */
			continue
		}
		if !l.queue.IsEventInList(event) {
			/* if event is not in list, ignore it */
			/* subscribed events might arrive after being added to logs */
			/* which would make the event be consumed again upon arrival */
			continue
		}
		block, err := eth.Client.BlockByNumber(context.Background(), event.BlockNumber)
		if err != nil {
			/* if fail to get block, return to try again */
			logger.Error("failed to get block", logger.String("message", err.Error()))
			return
		}
		event.Date = time.Unix(int64(block.Time()), 0).Format(time.RFC3339)
		l.consume(event)
		l.queue.RemoveEventsLike(event)
	}
}

func (l *Listener) consume(generic domain.GenericEvent) {
	switch generic.EventType {
	case "Created":
		createdModel := generic.ToCreatedEvent()
		logger.Info("consuming created event", logger.Object("event", &createdModel))
		l.api.CreateStar(createdModel)
	case "ChangedName":
		changedNameModel := generic.ToChangedNameEvent()
		logger.Info("consuming changed name event", logger.Object("event", &changedNameModel))
		l.api.ChangeName(changedNameModel)
	case "PutForSale":
		putForSaleModel := generic.ToPutForSaleEvent()
		logger.Info("consuming put for sale event", logger.Object("event", &putForSaleModel))
		l.api.PutForSale(putForSaleModel)
	case "RemovedFromSale":
		removedFromSaleModel := generic.ToRemovedFromSaleEvent()
		logger.Info("consuming removed from sale event", logger.Object("event", &removedFromSaleModel))
		l.api.RemoveFromSale(removedFromSaleModel)
	case "Sold":
		soldModel := generic.ToSoldEvent()
		logger.Info("consuming sold event", logger.Object("event", &soldModel))
		l.api.Sell(soldModel)
	}
}
