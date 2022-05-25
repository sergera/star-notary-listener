package models

import (
	"github.com/sergera/star-notary-listener/internal/logger"
)

type CreatedEvent struct {
	Owner       string
	TokenId     string
	Coordinates string
	Name        string
}

func (e *CreatedEvent) MarshalLogObject(enc logger.ObjectEncoder) error {
	enc.AddString("Owner", e.Owner)
	enc.AddString("TokenId", e.TokenId)
	enc.AddString("Coordinates", e.Coordinates)
	enc.AddString("Name", e.Name)
	return nil
}

type ChangedNameEvent struct {
	Owner   string
	TokenId string
	NewName string
}

func (e *ChangedNameEvent) MarshalLogObject(enc logger.ObjectEncoder) error {
	enc.AddString("Owner", e.Owner)
	enc.AddString("TokenId", e.TokenId)
	enc.AddString("NewName", e.NewName)
	return nil
}

type PutForSaleEvent struct {
	Owner      string
	TokenId    string
	PriceInWei string
}

func (e *PutForSaleEvent) MarshalLogObject(enc logger.ObjectEncoder) error {
	enc.AddString("Owner", e.Owner)
	enc.AddString("TokenId", e.TokenId)
	enc.AddString("PriceInWei", e.PriceInWei)
	return nil
}

type RemovedFromSaleEvent struct {
	Owner   string
	TokenId string
}

func (e *RemovedFromSaleEvent) MarshalLogObject(enc logger.ObjectEncoder) error {
	enc.AddString("Owner", e.Owner)
	enc.AddString("TokenId", e.TokenId)
	return nil
}

type SoldEvent struct {
	NewOwner string
	TokenId  string
}

func (e *SoldEvent) MarshalLogObject(enc logger.ObjectEncoder) error {
	enc.AddString("NewOwner", e.NewOwner)
	enc.AddString("TokenId", e.TokenId)
	return nil
}
