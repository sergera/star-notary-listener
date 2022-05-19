package event

import (
	"bytes"

	"github.com/sergera/star-notary-listener/internal/models"
	"github.com/sergera/star-notary-listener/pkg/slc"
)

func isDuplicateEvent(event models.Event, duplicate models.Event) bool {
	if !slc.ShallowEqual(event.Topics, duplicate.Topics) ||
		event.BlockNumber != duplicate.BlockNumber ||
		event.BlockHash != duplicate.BlockHash ||
		event.LogIndex != duplicate.LogIndex ||
		event.TxIndex != duplicate.TxIndex ||
		event.TxHash != duplicate.TxHash ||
		!bytes.Equal(event.Data, duplicate.Data) ||
		event.Sender != duplicate.Sender ||
		event.TokenId != duplicate.TokenId ||
		event.Name != duplicate.Name ||
		event.Coordinates != duplicate.Coordinates ||
		event.PriceInWei != duplicate.PriceInWei ||
		event.ContractHash != duplicate.ContractHash {
		return false
	}
	return true
}
