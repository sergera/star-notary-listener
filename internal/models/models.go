package models

import (
	"math/big"
	"time"

	"github.com/sergera/star-notary-listener/internal/logger"
)

type CreatedEvent struct {
	Owner       string    `json:"owner"`
	TokenId     string    `json:"tokenid"`
	Coordinates string    `json:"coordinates"`
	Name        string    `json:"name"`
	Date        time.Time `json:"date"`
}

func (e *CreatedEvent) MarshalLogObject(enc logger.ObjectEncoder) error {
	enc.AddString("Owner", e.Owner)
	enc.AddString("TokenId", e.TokenId)
	enc.AddString("Coordinates", e.Coordinates)
	enc.AddString("Name", e.Name)
	enc.AddString("Date", e.Date.String())
	return nil
}

type ChangedNameEvent struct {
	Owner   string    `json:"owner"`
	TokenId string    `json:"tokenid"`
	NewName string    `json:"name"`
	Date    time.Time `json:"date"`
}

func (e *ChangedNameEvent) MarshalLogObject(enc logger.ObjectEncoder) error {
	enc.AddString("Owner", e.Owner)
	enc.AddString("TokenId", e.TokenId)
	enc.AddString("NewName", e.NewName)
	enc.AddString("Date", e.Date.String())
	return nil
}

type PutForSaleEvent struct {
	Owner        string     `json:"owner"`
	TokenId      string     `json:"tokenid"`
	PriceInEther *big.Float `json:"price"`
	Date         time.Time  `json:"date"`
}

func (e *PutForSaleEvent) MarshalLogObject(enc logger.ObjectEncoder) error {
	enc.AddString("Owner", e.Owner)
	enc.AddString("TokenId", e.TokenId)
	enc.AddString("PriceInEther", e.PriceInEther.String())
	enc.AddString("Date", e.Date.String())
	return nil
}

type RemovedFromSaleEvent struct {
	Owner   string    `json:"owner"`
	TokenId string    `json:"tokenid"`
	Date    time.Time `json:"date"`
}

func (e *RemovedFromSaleEvent) MarshalLogObject(enc logger.ObjectEncoder) error {
	enc.AddString("Owner", e.Owner)
	enc.AddString("TokenId", e.TokenId)
	enc.AddString("Date", e.Date.String())
	return nil
}

type SoldEvent struct {
	NewOwner string    `json:"owner"`
	TokenId  string    `json:"tokenid"`
	Date     time.Time `json:"date"`
}

func (e *SoldEvent) MarshalLogObject(enc logger.ObjectEncoder) error {
	enc.AddString("NewOwner", e.NewOwner)
	enc.AddString("TokenId", e.TokenId)
	enc.AddString("Date", e.Date.String())
	return nil
}
