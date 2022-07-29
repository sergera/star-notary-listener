package event

import (
	"math/big"
	"sort"

	"github.com/ethereum/go-ethereum/common"
	"github.com/sergera/star-notary-listener/internal/conf"
	"github.com/sergera/star-notary-listener/internal/logger"
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
	coordinates  string
	sender       string
	priceInEther big.Float
	tokenId      string
	name         string
}

func (e *genericEvent) MarshalLogObject(enc logger.ObjectEncoder) error {
	enc.AddString("contractHash", e.contractHash)
	enc.AddString("eventType", e.eventType)
	enc.AddString("data", string(e.data))
	enc.AddUint64("blockNumber", e.blockNumber)
	enc.AddString("txHash", e.txHash)
	enc.AddUint("txIndex", e.txIndex)
	enc.AddString("blockHash", e.blockHash)
	enc.AddUint("logIndex", e.logIndex)
	enc.AddBool("removed", e.removed)
	enc.AddString("coordinates", e.coordinates)
	enc.AddString("sender", e.sender)
	enc.AddString("priceInEther", e.priceInEther.String())
	enc.AddString("tokenId", e.tokenId)
	enc.AddString("name", e.name)
	return nil
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
			logger.Info("Removing event", logger.Object("event", &event))
		}
		return !isDuplicateEvent(event, duplicate)
	})
}

func removeLeftoverEvents(currentBlock uint64) {
	unconfirmedEventsList = slc.Filter(unconfirmedEventsList, func(orphan genericEvent) bool {
		if currentBlock-orphan.blockNumber > conf.ConfirmationBlocks {
			logger.Info("Removing leftover event", logger.Object("event", &orphan))
		}
		return currentBlock-orphan.blockNumber <= conf.ConfirmationBlocks
	})
}

func isEventInList(event genericEvent) bool {
	_, exists := slc.Find(unconfirmedEventsList, func(duplicate genericEvent) bool {
		return isDuplicateEvent(event, duplicate)
	})
	return exists
}
