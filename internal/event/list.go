package event

import (
	"log"
	"sort"

	"github.com/ethereum/go-ethereum/common"
	"github.com/sergera/star-notary-listener/internal/env"
	"github.com/sergera/star-notary-listener/pkg/slc"
)

var unconfirmedEventsList []genericEvent = []genericEvent{}

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
	sort.SliceStable(unconfirmedEventsList, func(i, j int) bool {
		return unconfirmedEventsList[i].blockNumber < unconfirmedEventsList[j].blockNumber
	})
}

func insertEventByBlockNumber(event genericEvent) {
	unconfirmedEventsList = append(unconfirmedEventsList, event)
	sortListByBlockNumber()
}

func removeEvents(event genericEvent) {
	unconfirmedEventsList = slc.Filter(unconfirmedEventsList, func(duplicate genericEvent) bool {
		if isDuplicateEvent(event, duplicate) {
			log.Printf("Removing duplicate event: %+v\n\n", event)
		}
		return !isDuplicateEvent(event, duplicate)
	})
}

func removeOrphanedEvents(currentBlock uint64) {
	unconfirmedEventsList = slc.Filter(unconfirmedEventsList, func(orphan genericEvent) bool {
		if currentBlock-orphan.blockNumber >= env.OrphanedThreshold {
			log.Printf("Removing orphan event: %+v\n\n", orphan)
		}
		return currentBlock-orphan.blockNumber < env.OrphanedThreshold
	})
}

func isEventInList(event genericEvent) bool {
	_, exists := slc.Find(unconfirmedEventsList, func(duplicate genericEvent) bool {
		return isDuplicateEvent(event, duplicate)
	})
	return exists
}
