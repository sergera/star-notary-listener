package event

import (
	"log"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/sergera/star-notary-listener/internal/eth"
	"github.com/sergera/star-notary-listener/internal/models"
	"github.com/sergera/star-notary-listener/internal/starnotary"
)

func eventToCreatedEvent(event models.Event) models.CreatedEvent {
	return models.CreatedEvent{
		Owner:       event.Sender,
		Name:        event.Name,
		TokenId:     event.TokenId,
		Coordinates: event.Coordinates,
	}
}

func eventToChangedNameEvent(event models.Event) models.ChangedNameEvent {
	return models.ChangedNameEvent{
		NewName: event.Name,
		TokenId: event.TokenId,
	}
}

func eventToPutForSaleEvent(event models.Event) models.PutForSaleEvent {
	return models.PutForSaleEvent{
		TokenId:    event.TokenId,
		PriceInWei: event.PriceInWei,
	}
}

func eventToRemovedFromSaleEvent(event models.Event) models.RemovedFromSaleEvent {
	return models.RemovedFromSaleEvent{
		TokenId: event.TokenId,
	}
}

func eventToSoldEvent(event models.Event) models.SoldEvent {
	return models.SoldEvent{
		NewOwner: event.Sender,
		TokenId:  event.TokenId,
	}
}

func contractCreatedToEvent(subscribedEvent starnotary.StarnotaryCreated) models.Event {
	return models.Event{
		Sender:      common.Address.Hex(subscribedEvent.Owner),
		TokenId:     subscribedEvent.TokenId.Text(10),
		Coordinates: string(subscribedEvent.Coordinates[:]),
		Name:        string(subscribedEvent.Name),
		EventName:   eventSignatureToName[subscribedEvent.Raw.Topics[0].Hex()],

		ContractHash: subscribedEvent.Raw.Address.Hex(),
		Topics:       subscribedEvent.Raw.Topics,
		Data:         subscribedEvent.Raw.Data,
		BlockNumber:  subscribedEvent.Raw.BlockNumber,
		TxHash:       subscribedEvent.Raw.TxHash.Hex(),
		TxIndex:      subscribedEvent.Raw.TxIndex,
		BlockHash:    subscribedEvent.Raw.BlockHash.Hex(),
		LogIndex:     subscribedEvent.Raw.Index,
		Removed:      subscribedEvent.Raw.Removed,
	}
}

func contractChangedNameToEvent(subscribedEvent starnotary.StarnotaryChangedName) models.Event {
	return models.Event{
		Sender:    common.Address.Hex(subscribedEvent.Owner),
		Name:      string(subscribedEvent.NewName),
		TokenId:   subscribedEvent.TokenId.Text(10),
		EventName: eventSignatureToName[subscribedEvent.Raw.Topics[0].Hex()],

		ContractHash: subscribedEvent.Raw.Address.Hex(),
		Topics:       subscribedEvent.Raw.Topics,
		Data:         subscribedEvent.Raw.Data,
		BlockNumber:  subscribedEvent.Raw.BlockNumber,
		TxHash:       subscribedEvent.Raw.TxHash.Hex(),
		TxIndex:      subscribedEvent.Raw.TxIndex,
		BlockHash:    subscribedEvent.Raw.BlockHash.Hex(),
		LogIndex:     subscribedEvent.Raw.Index,
		Removed:      subscribedEvent.Raw.Removed,
	}
}

func contractPutForSaleToEvent(subscribedEvent starnotary.StarnotaryPutForSale) models.Event {
	return models.Event{
		Sender:     common.Address.Hex(subscribedEvent.Owner),
		TokenId:    subscribedEvent.TokenId.Text(10),
		PriceInWei: subscribedEvent.PriceInWei.Text(10),
		EventName:  eventSignatureToName[subscribedEvent.Raw.Topics[0].Hex()],

		ContractHash: subscribedEvent.Raw.Address.Hex(),
		Topics:       subscribedEvent.Raw.Topics,
		Data:         subscribedEvent.Raw.Data,
		BlockNumber:  subscribedEvent.Raw.BlockNumber,
		TxHash:       subscribedEvent.Raw.TxHash.Hex(),
		TxIndex:      subscribedEvent.Raw.TxIndex,
		BlockHash:    subscribedEvent.Raw.BlockHash.Hex(),
		LogIndex:     subscribedEvent.Raw.Index,
		Removed:      subscribedEvent.Raw.Removed,
	}
}

func contractRemovedFromSaleToEvent(subscribedEvent starnotary.StarnotaryRemovedFromSale) models.Event {
	return models.Event{
		Sender:    common.Address.Hex(subscribedEvent.Owner),
		TokenId:   subscribedEvent.TokenId.Text(10),
		EventName: eventSignatureToName[subscribedEvent.Raw.Topics[0].Hex()],

		ContractHash: subscribedEvent.Raw.Address.Hex(),
		Topics:       subscribedEvent.Raw.Topics,
		Data:         subscribedEvent.Raw.Data,
		BlockNumber:  subscribedEvent.Raw.BlockNumber,
		TxHash:       subscribedEvent.Raw.TxHash.Hex(),
		TxIndex:      subscribedEvent.Raw.TxIndex,
		BlockHash:    subscribedEvent.Raw.BlockHash.Hex(),
		LogIndex:     subscribedEvent.Raw.Index,
		Removed:      subscribedEvent.Raw.Removed,
	}
}

func contractSoldToEvent(subscribedEvent starnotary.StarnotarySold) models.Event {
	return models.Event{
		Sender:    common.Address.Hex(subscribedEvent.NewOwner),
		TokenId:   subscribedEvent.TokenId.Text(10),
		EventName: eventSignatureToName[subscribedEvent.Raw.Topics[0].Hex()],

		ContractHash: subscribedEvent.Raw.Address.Hex(),
		Topics:       subscribedEvent.Raw.Topics,
		Data:         subscribedEvent.Raw.Data,
		BlockNumber:  subscribedEvent.Raw.BlockNumber,
		TxHash:       subscribedEvent.Raw.TxHash.Hex(),
		TxIndex:      subscribedEvent.Raw.TxIndex,
		BlockHash:    subscribedEvent.Raw.BlockHash.Hex(),
		LogIndex:     subscribedEvent.Raw.Index,
		Removed:      subscribedEvent.Raw.Removed,
	}
}

func logToEvent(logEvent types.Log) (event models.Event) {
	eventSignature := logEvent.Topics[0].Hex()
	eventName := eventSignatureToName[eventSignature]

	switch eventName {
	case "Created":
		event = logCreatedToEvent(logEvent)
	case "ChangedName":
		event = logChangedNameToEvent(logEvent)
	case "PutForSale":
		event = logPutForSaleToEvent(logEvent)
	case "RemovedFromSale":
		event = logRemovedFromSaleToEvent(logEvent)
	case "Sold":
		event = logSoldToEvent(logEvent)
	default:
		log.Fatal("Tried to parse a non listened event")
	}

	return
}

func logCreatedToEvent(logEvent types.Log) models.Event {
	parsedCreated, err := eth.Contract.ParseCreated(logEvent)
	if err != nil {
		log.Fatal("Could not parse log event")
	}
	return contractCreatedToEvent(*parsedCreated)
}

func logChangedNameToEvent(logEvent types.Log) models.Event {
	parsedChangedName, err := eth.Contract.ParseChangedName(logEvent)
	if err != nil {
		log.Fatal("Could not parse log event")
	}
	return contractChangedNameToEvent(*parsedChangedName)
}

func logPutForSaleToEvent(logEvent types.Log) models.Event {
	parsedPutForSale, err := eth.Contract.ParsePutForSale(logEvent)
	if err != nil {
		log.Fatal("Could not parse log event")
	}
	return contractPutForSaleToEvent(*parsedPutForSale)
}

func logRemovedFromSaleToEvent(logEvent types.Log) models.Event {
	parsedRemovedFromSale, err := eth.Contract.ParseRemovedFromSale(logEvent)
	if err != nil {
		log.Fatal("Could not parse log event")
	}
	return contractRemovedFromSaleToEvent(*parsedRemovedFromSale)
}

func logSoldToEvent(logEvent types.Log) models.Event {
	parsedSold, err := eth.Contract.ParseSold(logEvent)
	if err != nil {
		log.Fatal("Could not parse log event")
	}
	return contractSoldToEvent(*parsedSold)
}
