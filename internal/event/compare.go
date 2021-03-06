package event

import (
	"bytes"

	"github.com/sergera/star-notary-listener/pkg/slc"
)

func isDuplicateEvent(event genericEvent, duplicate genericEvent) bool {
	if !slc.ShallowEqual(event.topics, duplicate.topics) ||
		event.blockNumber.String() != duplicate.blockNumber.String() ||
		event.blockHash != duplicate.blockHash ||
		event.logIndex != duplicate.logIndex ||
		event.txIndex != duplicate.txIndex ||
		event.txHash != duplicate.txHash ||
		!bytes.Equal(event.data, duplicate.data) ||
		event.sender != duplicate.sender ||
		event.tokenId != duplicate.tokenId ||
		event.name != duplicate.name ||
		event.coordinates != duplicate.coordinates ||
		event.priceInEther.String() != duplicate.priceInEther.String() ||
		event.contractHash != duplicate.contractHash {
		return false
	}
	return true
}
