package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/sergera/star-notary-listener/starnotary"
)

func main() {

	confirmationsThresholdString, exists := os.LookupEnv("CONFIRMATIONS_THRESHOLD")
	if !exists {
		log.Fatal("Confirmations threshold environment variable not found!")
	}

	confirmationsThreshold, err := strconv.ParseUint(confirmationsThresholdString, 10, 64)
	if err != nil {
		log.Fatal("Could not convert onfirmations threshold environment variable to int!")
	}

	infuraProjectId, exists := os.LookupEnv("INFURA_PROJECT_ID")
	if !exists {
		log.Fatal("Infura project id environment variable not found!")
	}

	client, err := ethclient.Dial("wss://rinkeby.infura.io/ws/v3/" + infuraProjectId)
	if err != nil {
		log.Fatal(err)
	}

	contractAddress := common.HexToAddress("0xf130D6968587fb69DE2DC1249293860446fB3823")

	starNotary, err := starnotary.NewStarnotary(contractAddress, client)
	if err != nil {
		log.Fatal(err)
	}

	createdEventQueue := []CreatedEvent{}
	changedNameEventQueue := []ChangedNameEvent{}
	putForSaleEventQueue := []PutForSaleEvent{}
	removedFromSaleEventQueue := []RemovedFromSaleEvent{}
	boughtEventQueue := []BoughtEvent{}

	createdResChan := make(chan *starnotary.StarnotaryCreated)
	changedNameResChan := make(chan *starnotary.StarnotaryChangedName)
	putForSaleResChan := make(chan *starnotary.StarnotaryPutForSale)
	removedFromSaleResChan := make(chan *starnotary.StarnotaryRemovedFromSale)
	boughtResChan := make(chan *starnotary.StarnotaryBought)

	starNotary.WatchCreated(&bind.WatchOpts{Start: nil, Context: context.Background()}, createdResChan)
	starNotary.WatchChangedName(&bind.WatchOpts{Start: nil, Context: context.Background()}, changedNameResChan)
	starNotary.WatchPutForSale(&bind.WatchOpts{Start: nil, Context: context.Background()}, putForSaleResChan)
	starNotary.WatchRemovedFromSale(&bind.WatchOpts{Start: nil, Context: context.Background()}, removedFromSaleResChan)
	starNotary.WatchBought(&bind.WatchOpts{Start: nil, Context: context.Background()}, boughtResChan)

	for {
		select {
		case createdLog := <-createdResChan:
			createdEvent := CreatedEvent{
				owner:       common.Address.Hex(createdLog.Owner),
				tokenId:     createdLog.TokenId.Text(10),
				coordinates: string(createdLog.Coordinates[:]),
				name:        string(createdLog.Name),
				blockNumber: createdLog.Raw.BlockNumber,
				eventIndex:  createdLog.Raw.Index,
			}
			createdEventQueue = append(createdEventQueue, createdEvent)
			fmt.Printf("Created Event To Queue: %+v\n", createdEvent)
		case changedNameLog := <-changedNameResChan:
			changedNameEvent := ChangedNameEvent{
				tokenId:     changedNameLog.TokenId.Text(10),
				newName:     string(changedNameLog.NewName),
				blockNumber: changedNameLog.Raw.BlockNumber,
				eventIndex:  changedNameLog.Raw.Index,
			}
			changedNameEventQueue = append(changedNameEventQueue, changedNameEvent)
			fmt.Printf("Changed Name Event To Queue: %+v\n", changedNameEvent)
		case putForSaleLog := <-putForSaleResChan:
			putForSaleEvent := PutForSaleEvent{
				tokenId:     putForSaleLog.TokenId.Text(10),
				priceInWei:  putForSaleLog.PriceInWei.Text(10),
				blockNumber: putForSaleLog.Raw.BlockNumber,
				eventIndex:  putForSaleLog.Raw.Index,
			}
			putForSaleEventQueue = append(putForSaleEventQueue, putForSaleEvent)
			fmt.Printf("Put For Sale Event To Queue: %+v\n", putForSaleEvent)
		case removedFromSaleLog := <-removedFromSaleResChan:
			removedFromSaleEvent := RemovedFromSaleEvent{
				tokenId:     removedFromSaleLog.TokenId.Text(10),
				blockNumber: removedFromSaleLog.Raw.BlockNumber,
				eventIndex:  removedFromSaleLog.Raw.Index,
			}
			removedFromSaleEventQueue = append(removedFromSaleEventQueue, removedFromSaleEvent)
			fmt.Printf("Removed From Sale Event To Queue: %+v\n", removedFromSaleEvent)
		case boughtLog := <-boughtResChan:
			boughtEvent := BoughtEvent{
				tokenId:     boughtLog.TokenId.Text(10),
				newOwner:    common.Address.Hex(boughtLog.NewOwner),
				blockNumber: boughtLog.Raw.BlockNumber,
				eventIndex:  boughtLog.Raw.Index,
			}
			boughtEventQueue = append(boughtEventQueue, boughtEvent)
			fmt.Printf("Bought Event To Queue: %+v\n", boughtEvent)
		default:
			currentBlockNumber, err := client.BlockNumber(context.Background())
			if err == nil {
				checkConfirmations(currentBlockNumber, confirmationsThreshold, &createdEventQueue, &changedNameEventQueue, &putForSaleEventQueue, &removedFromSaleEventQueue, &boughtEventQueue)
			}
		}
	}
}

func checkConfirmations(blockNumber uint64, confirmationsThreshold uint64, createdEventQueue *[]CreatedEvent, changedNameEventQueue *[]ChangedNameEvent, putForSaleEventQueue *[]PutForSaleEvent, removedFromSaleEventQueue *[]RemovedFromSaleEvent, boughtEventQueue *[]BoughtEvent) {
	consumeConfirmedEvents(createdEventQueue, blockNumber, confirmationsThreshold)
	consumeConfirmedEvents(changedNameEventQueue, blockNumber, confirmationsThreshold)
	consumeConfirmedEvents(putForSaleEventQueue, blockNumber, confirmationsThreshold)
	consumeConfirmedEvents(removedFromSaleEventQueue, blockNumber, confirmationsThreshold)
	consumeConfirmedEvents(boughtEventQueue, blockNumber, confirmationsThreshold)
}

type CreatedEvent struct {
	owner       string
	tokenId     string
	coordinates string
	name        string
	blockNumber uint64
	eventIndex  uint
}

type ChangedNameEvent struct {
	tokenId     string
	newName     string
	blockNumber uint64
	eventIndex  uint
}

type PutForSaleEvent struct {
	tokenId     string
	priceInWei  string
	blockNumber uint64
	eventIndex  uint
}

type RemovedFromSaleEvent struct {
	tokenId     string
	blockNumber uint64
	eventIndex  uint
}

type BoughtEvent struct {
	newOwner    string
	tokenId     string
	blockNumber uint64
	eventIndex  uint
}

type SharedEventFieldsInterface interface {
	getBlockNumber() uint64
	getEventIndex() uint
}

func (e CreatedEvent) getBlockNumber() uint64         { return e.blockNumber }
func (e CreatedEvent) getEventIndex() uint            { return e.eventIndex }
func (e ChangedNameEvent) getBlockNumber() uint64     { return e.blockNumber }
func (e ChangedNameEvent) getEventIndex() uint        { return e.eventIndex }
func (e PutForSaleEvent) getBlockNumber() uint64      { return e.blockNumber }
func (e PutForSaleEvent) getEventIndex() uint         { return e.eventIndex }
func (e RemovedFromSaleEvent) getBlockNumber() uint64 { return e.blockNumber }
func (e RemovedFromSaleEvent) getEventIndex() uint    { return e.eventIndex }
func (e BoughtEvent) getBlockNumber() uint64          { return e.blockNumber }
func (e BoughtEvent) getEventIndex() uint             { return e.eventIndex }

func consumeConfirmedEvents[E SharedEventFieldsInterface](queue *[]E, currentBlockNumber uint64, confirmationsThreshold uint64) {
	for _, event := range *queue {
		if currentBlockNumber-event.getBlockNumber() > confirmationsThreshold {
			fmt.Printf("Consuming event: %+v\n", event)
			removeFromQueue(queue, event.getEventIndex())
		}
	}
}

func removeFromQueue[E SharedEventFieldsInterface](queue *[]E, eventIndex uint) {
	for index, event := range *queue {
		if event.getEventIndex() == eventIndex {
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

// type Event struct {
// 	contractHash   string
// 	eventSignature string
// 	topics         []string
// 	data           []string
// 	blockNumber    int
// 	txHash         string
// 	txIndex        int
// 	blockHash      string
// 	logIndex       int
// 	removed        bool
// }
