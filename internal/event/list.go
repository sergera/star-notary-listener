package event

import (
	"log"
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
			log.Printf("Removing duplicate event: %+v\n\n", event)
		}
		return !isDuplicateEvent(event, duplicate)
	})
}

func removeOrphanedEvents(currentBlock uint64) {
	subscribedEventsList = slc.Filter(subscribedEventsList, func(event models.Event) bool {
		if currentBlock-event.BlockNumber >= env.OrphanedThreshold {
			log.Printf("Removing orphan event: %+v\n\n", event)
		}
		return currentBlock-event.BlockNumber < env.OrphanedThreshold
	})
}

func isEventInList(event models.Event) bool {
	_, exists := slc.Find(subscribedEventsList, func(duplicate models.Event) bool {
		return isDuplicateEvent(event, duplicate)
	})
	return exists
}
