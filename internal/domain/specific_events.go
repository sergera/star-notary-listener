package domain

import (
	"github.com/sergera/star-notary-listener/internal/logger"
)

type CreateEvent struct {
	Owner       string `json:"owner"`
	TokenId     string `json:"token_id"`
	Coordinates string `json:"coordinates"`
	Name        string `json:"name"`
	Date        string `json:"date"`
}

func (e *CreateEvent) MarshalLogObject(enc logger.ObjectEncoder) error {
	enc.AddString("Owner", e.Owner)
	enc.AddString("TokenId", e.TokenId)
	enc.AddString("Coordinates", e.Coordinates)
	enc.AddString("Name", e.Name)
	enc.AddString("Date", e.Date)
	return nil
}

type ChangeNameEvent struct {
	Owner   string `json:"owner"`
	TokenId string `json:"token_id"`
	NewName string `json:"name"`
	Date    string `json:"date"`
}

func (e *ChangeNameEvent) MarshalLogObject(enc logger.ObjectEncoder) error {
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

type RemoveFromSaleEvent struct {
	Owner   string `json:"owner"`
	TokenId string `json:"token_id"`
	Date    string `json:"date"`
}

func (e *RemoveFromSaleEvent) MarshalLogObject(enc logger.ObjectEncoder) error {
	enc.AddString("Owner", e.Owner)
	enc.AddString("TokenId", e.TokenId)
	enc.AddString("Date", e.Date)
	return nil
}

type PurchaseEvent struct {
	NewOwner string `json:"owner"`
	TokenId  string `json:"token_id"`
	Date     string `json:"date"`
}

func (e *PurchaseEvent) MarshalLogObject(enc logger.ObjectEncoder) error {
	enc.AddString("NewOwner", e.NewOwner)
	enc.AddString("TokenId", e.TokenId)
	enc.AddString("Date", e.Date)
	return nil
}
