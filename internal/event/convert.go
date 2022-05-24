package event

import (
	"log"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/sergera/star-notary-listener/internal/eth"
	"github.com/sergera/star-notary-listener/internal/gocontracts/starnotary"
	"github.com/sergera/star-notary-listener/internal/models"
)

func eventToCreatedEvent(event genericEvent) models.CreatedEvent {
	return models.CreatedEvent{
		Owner:       event.sender,
		Name:        event.name,
		TokenId:     event.tokenId,
		Coordinates: event.coordinates,
	}
}

func eventToChangedNameEvent(event genericEvent) models.ChangedNameEvent {
	return models.ChangedNameEvent{
		NewName: event.name,
		TokenId: event.tokenId,
	}
}

func eventToPutForSaleEvent(event genericEvent) models.PutForSaleEvent {
	return models.PutForSaleEvent{
		TokenId:    event.tokenId,
		PriceInWei: event.priceInWei,
	}
}

func eventToRemovedFromSaleEvent(event genericEvent) models.RemovedFromSaleEvent {
	return models.RemovedFromSaleEvent{
		TokenId: event.tokenId,
	}
}

func eventToSoldEvent(event genericEvent) models.SoldEvent {
	return models.SoldEvent{
		NewOwner: event.sender,
		TokenId:  event.tokenId,
	}
}

func contractCreatedToEvent(subscribedEvent starnotary.StarnotaryCreated) genericEvent {
	return genericEvent{
		sender:      common.Address.Hex(subscribedEvent.Owner),
		tokenId:     subscribedEvent.TokenId.Text(10),
		coordinates: string(subscribedEvent.Coordinates[:]),
		name:        string(subscribedEvent.Name),
		eventType:   eventSignatureToType[subscribedEvent.Raw.Topics[0].Hex()],

		contractHash: subscribedEvent.Raw.Address.Hex(),
		topics:       subscribedEvent.Raw.Topics,
		data:         subscribedEvent.Raw.Data,
		blockNumber:  subscribedEvent.Raw.BlockNumber,
		txHash:       subscribedEvent.Raw.TxHash.Hex(),
		txIndex:      subscribedEvent.Raw.TxIndex,
		blockHash:    subscribedEvent.Raw.BlockHash.Hex(),
		logIndex:     subscribedEvent.Raw.Index,
		removed:      subscribedEvent.Raw.Removed,
	}
}

func contractChangedNameToEvent(subscribedEvent starnotary.StarnotaryChangedName) genericEvent {
	return genericEvent{
		sender:    common.Address.Hex(subscribedEvent.Owner),
		name:      string(subscribedEvent.NewName),
		tokenId:   subscribedEvent.TokenId.Text(10),
		eventType: eventSignatureToType[subscribedEvent.Raw.Topics[0].Hex()],

		contractHash: subscribedEvent.Raw.Address.Hex(),
		topics:       subscribedEvent.Raw.Topics,
		data:         subscribedEvent.Raw.Data,
		blockNumber:  subscribedEvent.Raw.BlockNumber,
		txHash:       subscribedEvent.Raw.TxHash.Hex(),
		txIndex:      subscribedEvent.Raw.TxIndex,
		blockHash:    subscribedEvent.Raw.BlockHash.Hex(),
		logIndex:     subscribedEvent.Raw.Index,
		removed:      subscribedEvent.Raw.Removed,
	}
}

func contractPutForSaleToEvent(subscribedEvent starnotary.StarnotaryPutForSale) genericEvent {
	return genericEvent{
		sender:     common.Address.Hex(subscribedEvent.Owner),
		tokenId:    subscribedEvent.TokenId.Text(10),
		priceInWei: subscribedEvent.PriceInWei.Text(10),
		eventType:  eventSignatureToType[subscribedEvent.Raw.Topics[0].Hex()],

		contractHash: subscribedEvent.Raw.Address.Hex(),
		topics:       subscribedEvent.Raw.Topics,
		data:         subscribedEvent.Raw.Data,
		blockNumber:  subscribedEvent.Raw.BlockNumber,
		txHash:       subscribedEvent.Raw.TxHash.Hex(),
		txIndex:      subscribedEvent.Raw.TxIndex,
		blockHash:    subscribedEvent.Raw.BlockHash.Hex(),
		logIndex:     subscribedEvent.Raw.Index,
		removed:      subscribedEvent.Raw.Removed,
	}
}

func contractRemovedFromSaleToEvent(subscribedEvent starnotary.StarnotaryRemovedFromSale) genericEvent {
	return genericEvent{
		sender:    common.Address.Hex(subscribedEvent.Owner),
		tokenId:   subscribedEvent.TokenId.Text(10),
		eventType: eventSignatureToType[subscribedEvent.Raw.Topics[0].Hex()],

		contractHash: subscribedEvent.Raw.Address.Hex(),
		topics:       subscribedEvent.Raw.Topics,
		data:         subscribedEvent.Raw.Data,
		blockNumber:  subscribedEvent.Raw.BlockNumber,
		txHash:       subscribedEvent.Raw.TxHash.Hex(),
		txIndex:      subscribedEvent.Raw.TxIndex,
		blockHash:    subscribedEvent.Raw.BlockHash.Hex(),
		logIndex:     subscribedEvent.Raw.Index,
		removed:      subscribedEvent.Raw.Removed,
	}
}

func contractSoldToEvent(subscribedEvent starnotary.StarnotarySold) genericEvent {
	return genericEvent{
		sender:    common.Address.Hex(subscribedEvent.NewOwner),
		tokenId:   subscribedEvent.TokenId.Text(10),
		eventType: eventSignatureToType[subscribedEvent.Raw.Topics[0].Hex()],

		contractHash: subscribedEvent.Raw.Address.Hex(),
		topics:       subscribedEvent.Raw.Topics,
		data:         subscribedEvent.Raw.Data,
		blockNumber:  subscribedEvent.Raw.BlockNumber,
		txHash:       subscribedEvent.Raw.TxHash.Hex(),
		txIndex:      subscribedEvent.Raw.TxIndex,
		blockHash:    subscribedEvent.Raw.BlockHash.Hex(),
		logIndex:     subscribedEvent.Raw.Index,
		removed:      subscribedEvent.Raw.Removed,
	}
}

func logToEvent(logEvent types.Log) (event genericEvent) {
	eventSignature := logEvent.Topics[0].Hex()
	eventType := eventSignatureToType[eventSignature]

	switch eventType {
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
		log.Printf("Tried to parse a non listened event: %+v\n\n", eventType)
	}

	return
}

func logCreatedToEvent(logEvent types.Log) genericEvent {
	parsedCreated, err := eth.Contract.ParseCreated(logEvent)
	if err != nil {
		log.Printf("Could not parse log event: %+v\n\n", err)
	}
	return contractCreatedToEvent(*parsedCreated)
}

func logChangedNameToEvent(logEvent types.Log) genericEvent {
	parsedChangedName, err := eth.Contract.ParseChangedName(logEvent)
	if err != nil {
		log.Printf("Could not parse log event: %+v\n\n", err)
	}
	return contractChangedNameToEvent(*parsedChangedName)
}

func logPutForSaleToEvent(logEvent types.Log) genericEvent {
	parsedPutForSale, err := eth.Contract.ParsePutForSale(logEvent)
	if err != nil {
		log.Printf("Could not parse log event: %+v\n\n", err)
	}
	return contractPutForSaleToEvent(*parsedPutForSale)
}

func logRemovedFromSaleToEvent(logEvent types.Log) genericEvent {
	parsedRemovedFromSale, err := eth.Contract.ParseRemovedFromSale(logEvent)
	if err != nil {
		log.Printf("Could not parse log event: %+v\n\n", err)
	}
	return contractRemovedFromSaleToEvent(*parsedRemovedFromSale)
}

func logSoldToEvent(logEvent types.Log) genericEvent {
	parsedSold, err := eth.Contract.ParseSold(logEvent)
	if err != nil {
		log.Printf("Could not parse log event: %+v\n\n", err)
	}
	return contractSoldToEvent(*parsedSold)
}
