package models

import (
	"github.com/sergera/star-notary-listener/internal/logger"
)

type CreatedEvent struct {
	Owner       string `json:"owner"`
	TokenId     string `json:"token_id"`
	Coordinates string `json:"coordinates"`
	Name        string `json:"name"`
	Date        string `json:"date"`
}

func (e *CreatedEvent) MarshalLogObject(enc logger.ObjectEncoder) error {
	enc.AddString("Owner", e.Owner)
	enc.AddString("TokenId", e.TokenId)
	enc.AddString("Coordinates", e.Coordinates)
	enc.AddString("Name", e.Name)
	enc.AddString("Date", e.Date)
	return nil
}

type ChangedNameEvent struct {
	Owner   string `json:"owner"`
	TokenId string `json:"token_id"`
	NewName string `json:"name"`
	Date    string `json:"date"`
}

func (e *ChangedNameEvent) MarshalLogObject(enc logger.ObjectEncoder) error {
	enc.AddString("Owner", e.Owner)
	enc.AddString("TokenId", e.TokenId)
	enc.AddString("NewName", e.NewName)
	enc.AddString("Date", e.Date)
	return nil
}

type PutForSaleEvent struct {
	Owner        string `json:"owner"`
	TokenId      string `json:"token_id"`
	PriceInEther string `json:"price"`
	Date         string `json:"date"`
}

func (e *PutForSaleEvent) MarshalLogObject(enc logger.ObjectEncoder) error {
	enc.AddString("Owner", e.Owner)
	enc.AddString("TokenId", e.TokenId)
	enc.AddString("PriceInEther", e.PriceInEther)
	enc.AddString("Date", e.Date)
	return nil
}

type RemovedFromSaleEvent struct {
	Owner   string `json:"owner"`
	TokenId string `json:"token_id"`
	Date    string `json:"date"`
}

func (e *RemovedFromSaleEvent) MarshalLogObject(enc logger.ObjectEncoder) error {
	enc.AddString("Owner", e.Owner)
	enc.AddString("TokenId", e.TokenId)
	enc.AddString("Date", e.Date)
	return nil
}

type SoldEvent struct {
	NewOwner string `json:"owner"`
	TokenId  string `json:"token_id"`
	Date     string `json:"date"`
}

func (e *SoldEvent) MarshalLogObject(enc logger.ObjectEncoder) error {
	enc.AddString("NewOwner", e.NewOwner)
	enc.AddString("TokenId", e.TokenId)
	enc.AddString("Date", e.Date)
	return nil
}
