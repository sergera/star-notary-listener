package event

import (
	"fmt"
	"sort"

	"github.com/sergera/star-notary-listener/internal/env"
	"github.com/sergera/star-notary-listener/internal/models"
	"github.com/sergera/star-notary-listener/pkg/slc"
)

var subscribedEventsList []models.Event = []models.Event{}

func sortListByBlockNumber() {
	sort.SliceStable(subscribedEventsList, func(i, j int) bool {
		return subscribedEventsList[i].BlockNumber < subscribedEventsList[j].BlockNumber
	})
}

func insertEvent(event models.Event) {
	subscribedEventsList = append(subscribedEventsList, event)
	sortListByBlockNumber()
}

func removeDuplicateEvents(event models.Event) {
	subscribedEventsList = slc.Filter(subscribedEventsList, func(duplicate models.Event) bool {
		if isDuplicateEvent(event, duplicate) {
			fmt.Println("Removing duplicate event: ", event)
		}
		return !isDuplicateEvent(event, duplicate)
	})
}

func removeOrphanedEvents(currentBlock uint64) {
	subscribedEventsList = slc.Filter(subscribedEventsList, func(event models.Event) bool {
		if currentBlock-event.BlockNumber > env.ConfirmationsThreshold {
			fmt.Println("Removing orphan event: ", event)
		}
		return currentBlock-event.BlockNumber <= env.ConfirmationsThreshold
	})
}
