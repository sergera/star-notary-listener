package event

import (
	"math/big"
	"sort"

	"github.com/sergera/star-notary-listener/internal/conf"
	"github.com/sergera/star-notary-listener/internal/domain"
	"github.com/sergera/star-notary-listener/internal/logger"
	"github.com/sergera/star-notary-listener/pkg/slc"
)

var unconfirmedEventsList []domain.GenericEvent = []domain.GenericEvent{}

func sortListByBlockNumber() {
	sort.SliceStable(unconfirmedEventsList, func(i, j int) bool {
		/* return i < j */
		return unconfirmedEventsList[i].BlockNumber.Cmp(unconfirmedEventsList[j].BlockNumber) == -1
	})
}

func insertEventByBlockNumber(event domain.GenericEvent) {
	unconfirmedEventsList = append(unconfirmedEventsList, event)
	sortListByBlockNumber()
}

func removeEvents(event domain.GenericEvent) {
	unconfirmedEventsList = slc.Filter(unconfirmedEventsList, func(duplicate domain.GenericEvent) bool {
		if isDuplicateEvent(event, duplicate) {
			logger.Info("Removing event", logger.Object("event", &event))
		}
		return !isDuplicateEvent(event, duplicate)
	})
}

func removeLeftoverEvents(latestBlock *big.Int) {
	conf := conf.GetConf()
	unconfirmedEventsList = slc.Filter(unconfirmedEventsList, func(orphan domain.GenericEvent) bool {
		if big.NewInt(0).Sub(latestBlock, orphan.BlockNumber).Cmp(big.NewInt(int64(conf.ConfirmationBlocks))) == 1 {
			/* if latestBlock - orphanBlockNumber > confirmationBlocks */
			logger.Info("Removing leftover event", logger.Object("event", &orphan))
		}
		/* return latestBlock - orphanBlockNumber <= confirmationBlocks */
		return big.NewInt(0).Sub(latestBlock, orphan.BlockNumber).Cmp(big.NewInt(int64(conf.ConfirmationBlocks))) != 1
	})
}

func isEventInList(event domain.GenericEvent) bool {
	_, exists := slc.Find(unconfirmedEventsList, func(duplicate domain.GenericEvent) bool {
		return isDuplicateEvent(event, duplicate)
	})
	return exists
}
