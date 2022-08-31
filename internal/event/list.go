package event

import (
	"math/big"
	"sort"
	"time"

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
	blockNumber  *big.Int
	txHash       string
	txIndex      uint
	blockHash    string
	logIndex     uint
	removed      bool
	date         time.Time
	/* specific event fields */
	coordinates  string
	sender       string
	priceInEther *big.Float
	tokenId      string
	name         string
}

func (e *genericEvent) MarshalLogObject(enc logger.ObjectEncoder) error {
	enc.AddString("contractHash", e.contractHash)
	enc.AddString("eventType", e.eventType)
	enc.AddString("data", string(e.data))
	enc.AddString("blockNumber", e.blockNumber.String())
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
	enc.AddString("date", e.date.String())
	return nil
}

func sortListByBlockNumber() {
	sort.SliceStable(unconfirmedEventsList, func(i, j int) bool {
		/* return i < j */
		return unconfirmedEventsList[i].blockNumber.Cmp(unconfirmedEventsList[j].blockNumber) == -1
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

func removeLeftoverEvents(latestBlock *big.Int) {
	conf := conf.GetConf()
	unconfirmedEventsList = slc.Filter(unconfirmedEventsList, func(orphan genericEvent) bool {
		if big.NewInt(0).Sub(latestBlock, orphan.blockNumber).Cmp(big.NewInt(int64(conf.ConfirmationBlocks))) == 1 {
			/* if latestBlock - orphanBlockNumber > confirmationBlocks */
			logger.Info("Removing leftover event", logger.Object("event", &orphan))
		}
		/* return latestBlock - orphanBlockNumber <= confirmationBlocks */
		return big.NewInt(0).Sub(latestBlock, orphan.blockNumber).Cmp(big.NewInt(int64(conf.ConfirmationBlocks))) != 1
	})
}

func isEventInList(event genericEvent) bool {
	_, exists := slc.Find(unconfirmedEventsList, func(duplicate genericEvent) bool {
		return isDuplicateEvent(event, duplicate)
	})
	return exists
}
