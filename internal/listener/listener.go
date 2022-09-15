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
	"0x645f015284b54bc389e86c58c9e6fa711dabfe2a1fd7d7fdf9440458b0074587": "Create",
	"0xf744b02e2965631679dba1d6ccc015286348dc6452ab377764b9bce7c65aec0c": "ChangeName",
	"0xeef8701c784dcc5b12eb5ce2687a9e42d1d94b6e81f660dcb84b51554c37f082": "PutForSale",
	"0xbfbf7e7677a0c423106146f1ee86ac042526b53581de06ba54c51e8acfeac746": "RemoveFromSale",
	"0x2499a5330ab0979cc612135e7883ebc3cd5c9f7a8508f042540c34723348f632": "Purchase",
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

	createResChan := make(chan *starnotary.StarnotaryCreate)
	changeNameResChan := make(chan *starnotary.StarnotaryChangeName)
	putForSaleResChan := make(chan *starnotary.StarnotaryPutForSale)
	removeFromSaleResChan := make(chan *starnotary.StarnotaryRemoveFromSale)
	purchaseResChan := make(chan *starnotary.StarnotaryPurchase)

	defer close(createResChan)
	defer close(changeNameResChan)
	defer close(putForSaleResChan)
	defer close(removeFromSaleResChan)
	defer close(purchaseResChan)

	eth.Contract.WatchCreate(&bind.WatchOpts{Start: nil, Context: context.Background()}, createResChan)
	eth.Contract.WatchChangeName(&bind.WatchOpts{Start: nil, Context: context.Background()}, changeNameResChan)
	eth.Contract.WatchPutForSale(&bind.WatchOpts{Start: nil, Context: context.Background()}, putForSaleResChan)
	eth.Contract.WatchRemoveFromSale(&bind.WatchOpts{Start: nil, Context: context.Background()}, removeFromSaleResChan)
	eth.Contract.WatchPurchase(&bind.WatchOpts{Start: nil, Context: context.Background()}, purchaseResChan)

	for {
		select {
		case createEvent := <-createResChan:
			genericCreate := createToGeneric(*createEvent)
			l.queue.InsertEventByBlockNumber(genericCreate)
			logger.Info("create event to list", logger.Object("event", &genericCreate))
		case changeNameEvent := <-changeNameResChan:
			genericChangeName := changeNameToGeneric(*changeNameEvent)
			l.queue.InsertEventByBlockNumber(genericChangeName)
			logger.Info("changed name event to list", logger.Object("event", &genericChangeName))
		case putForSaleEvent := <-putForSaleResChan:
			genericPutForSale := putForSaleToGeneric(*putForSaleEvent)
			l.queue.InsertEventByBlockNumber(genericPutForSale)
			logger.Info("put for sale event to list", logger.Object("event", &genericPutForSale))
		case removeFromSaleEvent := <-removeFromSaleResChan:
			genericRemoveFromSale := removeFromSaleToGeneric(*removeFromSaleEvent)
			l.queue.InsertEventByBlockNumber(genericRemoveFromSale)
			logger.Info("removed from sale event to list", logger.Object("event", &genericRemoveFromSale))
		case purchaseEvent := <-purchaseResChan:
			genericPurchase := purchaseToGeneric(*purchaseEvent)
			l.queue.InsertEventByBlockNumber(genericPurchase)
			logger.Info("purchase event to list", logger.Object("event", &genericPurchase))
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
	case "Create":
		createModel := generic.ToCreateEvent()
		logger.Info("consuming create event", logger.Object("event", &createModel))
		l.api.CreateStar(createModel)
	case "ChangeName":
		changeNameModel := generic.ToChangeNameEvent()
		logger.Info("consuming changed name event", logger.Object("event", &changeNameModel))
		l.api.ChangeName(changeNameModel)
	case "PutForSale":
		putForSaleModel := generic.ToPutForSaleEvent()
		logger.Info("consuming put for sale event", logger.Object("event", &putForSaleModel))
		l.api.PutForSale(putForSaleModel)
	case "RemoveFromSale":
		removeFromSaleModel := generic.ToRemoveFromSaleEvent()
		logger.Info("consuming removed from sale event", logger.Object("event", &removeFromSaleModel))
		l.api.RemoveFromSale(removeFromSaleModel)
	case "Purchase":
		purchaseModel := generic.ToPurchaseEvent()
		logger.Info("consuming purchase event", logger.Object("event", &purchaseModel))
		l.api.Purchase(purchaseModel)
	}
}
