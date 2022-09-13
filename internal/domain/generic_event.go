package domain

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/sergera/star-notary-listener/internal/logger"
)

type GenericEvent struct {
	ContractHash string
	EventType    string
	Topics       []common.Hash
	Data         []byte
	BlockNumber  *big.Int
	TxHash       string
	TxIndex      uint
	BlockHash    string
	LogIndex     uint
	Removed      bool
	Date         string
	/* specific event fields */
	Coordinates  string
	Sender       string
	PriceInEther *big.Float
	TokenId      string
	Name         string
}

func (e *GenericEvent) MarshalLogObject(enc logger.ObjectEncoder) error {
	enc.AddString("contractHash", e.ContractHash)
	enc.AddString("eventType", e.EventType)
	enc.AddString("data", string(e.Data))
	enc.AddString("blockNumber", e.BlockNumber.String())
	enc.AddString("txHash", e.TxHash)
	enc.AddUint("txIndex", e.TxIndex)
	enc.AddString("blockHash", e.BlockHash)
	enc.AddUint("logIndex", e.LogIndex)
	enc.AddBool("removed", e.Removed)
	enc.AddString("coordinates", e.Coordinates)
	enc.AddString("sender", e.Sender)
	enc.AddString("priceInEther", e.PriceInEther.String())
	enc.AddString("tokenId", e.TokenId)
	enc.AddString("name", e.Name)
	enc.AddString("date", e.Date)
	return nil
}
