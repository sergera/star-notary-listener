package event

import (
	"bytes"

	"github.com/sergera/star-notary-listener/internal/domain"
	"github.com/sergera/star-notary-listener/pkg/slc"
)

func isDuplicateEvent(event domain.GenericEvent, duplicate domain.GenericEvent) bool {
	if !slc.ShallowEqual(event.Topics, duplicate.Topics) ||
		event.BlockNumber.String() != duplicate.BlockNumber.String() ||
		event.BlockHash != duplicate.BlockHash ||
		event.LogIndex != duplicate.LogIndex ||
		event.TxIndex != duplicate.TxIndex ||
		event.TxHash != duplicate.TxHash ||
		!bytes.Equal(event.Data, duplicate.Data) ||
		event.Sender != duplicate.Sender ||
		event.TokenId != duplicate.TokenId ||
		event.Name != duplicate.Name ||
		event.Coordinates != duplicate.Coordinates ||
		event.PriceInEther.String() != duplicate.PriceInEther.String() ||
		event.ContractHash != duplicate.ContractHash {
		return false
	}
	return true
}
