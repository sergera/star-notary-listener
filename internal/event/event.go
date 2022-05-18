package event

import (
	"context"
	"fmt"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/sergera/star-notary-listener/internal/env"
	"github.com/sergera/star-notary-listener/internal/eth"
	"github.com/sergera/star-notary-listener/internal/models"
	"github.com/sergera/star-notary-listener/internal/starnotary"
)

var createdEventQueue []models.CreatedEvent = []models.CreatedEvent{}
var changedNameEventQueue []models.ChangedNameEvent = []models.ChangedNameEvent{}
var putForSaleEventQueue []models.PutForSaleEvent = []models.PutForSaleEvent{}
var removedFromSaleEventQueue []models.RemovedFromSaleEvent = []models.RemovedFromSaleEvent{}
var boughtEventQueue []models.BoughtEvent = []models.BoughtEvent{}

func ListenAndConfirm() {
	createdResChan := make(chan *starnotary.StarnotaryCreated)
	changedNameResChan := make(chan *starnotary.StarnotaryChangedName)
	putForSaleResChan := make(chan *starnotary.StarnotaryPutForSale)
	removedFromSaleResChan := make(chan *starnotary.StarnotaryRemovedFromSale)
	boughtResChan := make(chan *starnotary.StarnotaryBought)

	eth.Contract.WatchCreated(&bind.WatchOpts{Start: nil, Context: context.Background()}, createdResChan)
	eth.Contract.WatchChangedName(&bind.WatchOpts{Start: nil, Context: context.Background()}, changedNameResChan)
	eth.Contract.WatchPutForSale(&bind.WatchOpts{Start: nil, Context: context.Background()}, putForSaleResChan)
	eth.Contract.WatchRemovedFromSale(&bind.WatchOpts{Start: nil, Context: context.Background()}, removedFromSaleResChan)
	eth.Contract.WatchBought(&bind.WatchOpts{Start: nil, Context: context.Background()}, boughtResChan)

	for {
		select {
		case createdLog := <-createdResChan:
			createdEvent := models.CreatedEvent{
				Owner:       common.Address.Hex(createdLog.Owner),
				TokenId:     createdLog.TokenId.Text(10),
				Coordinates: string(createdLog.Coordinates[:]),
				Name:        string(createdLog.Name),
				BlockNumber: createdLog.Raw.BlockNumber,
				EventIndex:  createdLog.Raw.Index,
			}
			createdEventQueue = append(createdEventQueue, createdEvent)
			fmt.Printf("Created Event To Queue: %+v\n", createdEvent)
		case changedNameLog := <-changedNameResChan:
			changedNameEvent := models.ChangedNameEvent{
				TokenId:     changedNameLog.TokenId.Text(10),
				NewName:     string(changedNameLog.NewName),
				BlockNumber: changedNameLog.Raw.BlockNumber,
				EventIndex:  changedNameLog.Raw.Index,
			}
			changedNameEventQueue = append(changedNameEventQueue, changedNameEvent)
			fmt.Printf("Changed Name Event To Queue: %+v\n", changedNameEvent)
		case putForSaleLog := <-putForSaleResChan:
			putForSaleEvent := models.PutForSaleEvent{
				TokenId:     putForSaleLog.TokenId.Text(10),
				PriceInWei:  putForSaleLog.PriceInWei.Text(10),
				BlockNumber: putForSaleLog.Raw.BlockNumber,
				EventIndex:  putForSaleLog.Raw.Index,
			}
			putForSaleEventQueue = append(putForSaleEventQueue, putForSaleEvent)
			fmt.Printf("Put For Sale Event To Queue: %+v\n", putForSaleEvent)
		case removedFromSaleLog := <-removedFromSaleResChan:
			removedFromSaleEvent := models.RemovedFromSaleEvent{
				TokenId:     removedFromSaleLog.TokenId.Text(10),
				BlockNumber: removedFromSaleLog.Raw.BlockNumber,
				EventIndex:  removedFromSaleLog.Raw.Index,
			}
			removedFromSaleEventQueue = append(removedFromSaleEventQueue, removedFromSaleEvent)
			fmt.Printf("Removed From Sale Event To Queue: %+v\n", removedFromSaleEvent)
		case boughtLog := <-boughtResChan:
			boughtEvent := models.BoughtEvent{
				TokenId:     boughtLog.TokenId.Text(10),
				NewOwner:    common.Address.Hex(boughtLog.NewOwner),
				BlockNumber: boughtLog.Raw.BlockNumber,
				EventIndex:  boughtLog.Raw.Index,
			}
			boughtEventQueue = append(boughtEventQueue, boughtEvent)
			fmt.Printf("Bought Event To Queue: %+v\n", boughtEvent)
		default:
			currentBlockNumber, err := eth.Client.BlockNumber(context.Background())
			if err == nil {
				checkConfirmations(currentBlockNumber)
			}
		}
	}
}

func checkConfirmations(blockNumber uint64) {
	consumeConfirmedEvents(&createdEventQueue, blockNumber)
	consumeConfirmedEvents(&changedNameEventQueue, blockNumber)
	consumeConfirmedEvents(&putForSaleEventQueue, blockNumber)
	consumeConfirmedEvents(&removedFromSaleEventQueue, blockNumber)
	consumeConfirmedEvents(&boughtEventQueue, blockNumber)
}

func consumeConfirmedEvents[E models.SharedEventFieldsInterface](queue *[]E, currentBlockNumber uint64) {
	for _, event := range *queue {
		if currentBlockNumber-event.GetBlockNumber() > env.ConfirmationsThreshold {
			fmt.Printf("Consuming event: %+v\n", event)
			removeFromQueue(queue, event.GetEventIndex())
		}
	}
}

func removeFromQueue[E models.SharedEventFieldsInterface](queue *[]E, eventIndex uint) {
	for index, event := range *queue {
		if event.GetEventIndex() == eventIndex {
			fmt.Printf("Removing event: %+v\n", event)
			*queue = remove(*queue, index)
		}
	}
}

func removeFast[T any](slice []T, i int) []T {
	slice[i] = slice[len(slice)-1]
	return slice[:len(slice)-1]
}

func remove[T any](slice []T, i int) []T {
	return append(slice[:i], slice[i+1:]...)
}
