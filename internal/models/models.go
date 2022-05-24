package models

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
