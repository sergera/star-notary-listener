package event

import (
	"log"
	"sort"

	"github.com/ethereum/go-ethereum/common"
	"github.com/sergera/star-notary-listener/internal/env"
	"github.com/sergera/star-notary-listener/pkg/slc"
)

var subscribedEventsList []genericEvent = []genericEvent{}

type genericEvent struct {
	contractHash string
	eventType    string
	topics       []common.Hash
	data         []byte
	blockNumber  uint64
	txHash       string
	txIndex      uint
	blockHash    string
	logIndex     uint
	removed      bool
	/* specific event fields */
	coordinates string
	sender      string
	priceInWei  string
	tokenId     string
	name        string
}

func sortListByBlockNumber() {
	sort.SliceStable(subscribedEventsList, func(i, j int) bool {
		return subscribedEventsList[i].blockNumber < subscribedEventsList[j].blockNumber
	})
}

func insertEvent(event genericEvent) {
	subscribedEventsList = append(subscribedEventsList, event)
	sortListByBlockNumber()
}

func removeDuplicateEvents(event genericEvent) {
	subscribedEventsList = slc.Filter(subscribedEventsList, func(duplicate genericEvent) bool {
		if isDuplicateEvent(event, duplicate) {
			log.Printf("Removing duplicate event: %+v\n\n", event)
		}
		return !isDuplicateEvent(event, duplicate)
	})
}

func removeOrphanedEvents(currentBlock uint64) {
	subscribedEventsList = slc.Filter(subscribedEventsList, func(event genericEvent) bool {
		if currentBlock-event.blockNumber >= env.OrphanedThreshold {
			log.Printf("Removing orphan event: %+v\n\n", event)
		}
		return currentBlock-event.blockNumber < env.OrphanedThreshold
	})
}

func isEventInList(event genericEvent) bool {
	_, exists := slc.Find(subscribedEventsList, func(duplicate genericEvent) bool {
		return isDuplicateEvent(event, duplicate)
	})
	return exists
}
