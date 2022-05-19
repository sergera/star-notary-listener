package models

import (
	"github.com/ethereum/go-ethereum/common"
)

type Event struct {
	ContractHash string
	EventName    string
	Topics       []common.Hash
	Data         []byte
	BlockNumber  uint64
	TxHash       string
	TxIndex      uint
	BlockHash    string
	LogIndex     uint
	Removed      bool
	/* specific event fields */
	Coordinates string
	Sender      string
	PriceInWei  string
	TokenId     string
	Name        string
}

type CreatedEvent struct {
	Owner       string
	TokenId     string
	Coordinates string
	Name        string
}

type ChangedNameEvent struct {
	Owner   string
	TokenId string
	NewName string
}

type PutForSaleEvent struct {
	Owner      string
	TokenId    string
	PriceInWei string
}

type RemovedFromSaleEvent struct {
	Owner   string
	TokenId string
}

type SoldEvent struct {
	NewOwner string
	TokenId  string
}
