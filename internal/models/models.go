package models

type Event struct {
	ContractHash   string
	EventSignature string
	Topics         []string
	Data           []string
	BlockNumber    int
	TxHash         string
	TxIndex        int
	BlockHash      string
	LogIndex       int
	Removed        bool
}

type CreatedEvent struct {
	Owner       string
	TokenId     string
	Coordinates string
	Name        string
	BlockNumber uint64
	EventIndex  uint
}

type ChangedNameEvent struct {
	TokenId     string
	NewName     string
	BlockNumber uint64
	EventIndex  uint
}

type PutForSaleEvent struct {
	TokenId     string
	PriceInWei  string
	BlockNumber uint64
	EventIndex  uint
}

type RemovedFromSaleEvent struct {
	TokenId     string
	BlockNumber uint64
	EventIndex  uint
}

type BoughtEvent struct {
	NewOwner    string
	TokenId     string
	BlockNumber uint64
	EventIndex  uint
}

type SharedEventFieldsInterface interface {
	GetBlockNumber() uint64
	GetEventIndex() uint
}

func (e CreatedEvent) GetBlockNumber() uint64         { return e.BlockNumber }
func (e CreatedEvent) GetEventIndex() uint            { return e.EventIndex }
func (e ChangedNameEvent) GetBlockNumber() uint64     { return e.BlockNumber }
func (e ChangedNameEvent) GetEventIndex() uint        { return e.EventIndex }
func (e PutForSaleEvent) GetBlockNumber() uint64      { return e.BlockNumber }
func (e PutForSaleEvent) GetEventIndex() uint         { return e.EventIndex }
func (e RemovedFromSaleEvent) GetBlockNumber() uint64 { return e.BlockNumber }
func (e RemovedFromSaleEvent) GetEventIndex() uint    { return e.EventIndex }
func (e BoughtEvent) GetBlockNumber() uint64          { return e.BlockNumber }
func (e BoughtEvent) GetEventIndex() uint             { return e.EventIndex }
