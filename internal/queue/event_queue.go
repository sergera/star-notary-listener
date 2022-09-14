package queue

import (
	"math/big"
	"sort"
	"sync"

	"github.com/sergera/star-notary-listener/internal/conf"
	"github.com/sergera/star-notary-listener/internal/domain"
	"github.com/sergera/star-notary-listener/internal/logger"
	"github.com/sergera/star-notary-listener/pkg/slc"
)

type EventQueue struct {
	queue []domain.GenericEvent
	lock  *sync.Mutex
}

func NewEventQueue() *EventQueue {
	return &EventQueue{
		queue: []domain.GenericEvent{},
		lock:  &sync.Mutex{},
	}
}

func (q *EventQueue) Length() int {
	return len(q.queue)
}

func (q *EventQueue) FirstEventBlockNumber() *big.Int {
	if q.Length() > 0 {
		return q.queue[0].BlockNumber
	}

	return new(big.Int)
}

func (q *EventQueue) InsertEventByBlockNumber(event domain.GenericEvent) {
	q.lock.Lock()
	q.queue = append(q.queue, event)
	q.lock.Unlock()
	q.sortListByBlockNumber()
}

func (q *EventQueue) sortListByBlockNumber() {
	q.lock.Lock()
	defer q.lock.Unlock()
	sort.SliceStable(q.queue, func(i, j int) bool {
		/* return i < j */
		return q.queue[i].BlockNumber.Cmp(q.queue[j].BlockNumber) == -1
	})
}

func (q *EventQueue) RemoveEventsLike(event domain.GenericEvent) {
	q.lock.Lock()
	defer q.lock.Unlock()
	q.queue = slc.Filter(q.queue, func(duplicate domain.GenericEvent) bool {
		if event.IsDuplicate(&duplicate) {
			logger.Info("removing event", logger.Object("event", &event))
		}
		return !event.IsDuplicate(&duplicate)
	})
}

func (q *EventQueue) RemoveLeftoverEvents(latestBlock *big.Int) {
	conf := conf.GetConf()
	q.lock.Lock()
	defer q.lock.Unlock()
	q.queue = slc.Filter(q.queue, func(orphan domain.GenericEvent) bool {
		if big.NewInt(0).Sub(latestBlock, orphan.BlockNumber).Cmp(big.NewInt(int64(conf.ConfirmationBlocks()))) == 1 {
			/* if latestBlock - orphanBlockNumber > confirmationBlocks */
			logger.Info("removing leftover event", logger.Object("event", &orphan))
		}
		/* return latestBlock - orphanBlockNumber <= confirmationBlocks */
		return big.NewInt(0).Sub(latestBlock, orphan.BlockNumber).Cmp(big.NewInt(int64(conf.ConfirmationBlocks()))) != 1
	})
}

func (q *EventQueue) IsEventInList(event domain.GenericEvent) bool {
	_, exists := slc.Find(q.queue, func(duplicate domain.GenericEvent) bool {
		return event.IsDuplicate(&duplicate)
	})
	return exists
}
