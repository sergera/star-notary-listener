package event

import (
	"math/big"
	"strconv"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/sergera/star-notary-listener/internal/eth"
	"github.com/sergera/star-notary-listener/internal/gocontracts/starnotary"
	"github.com/sergera/star-notary-listener/internal/logger"
	"github.com/sergera/star-notary-listener/internal/models"
)

func genericToCreatedModel(event genericEvent) models.CreatedEvent {
	return models.CreatedEvent{
		Owner:       event.sender,
		Name:        event.name,
		TokenId:     event.tokenId,
		Coordinates: event.coordinates,
		Date:        event.date,
	}
}

func genericToChangedNameModel(event genericEvent) models.ChangedNameEvent {
	return models.ChangedNameEvent{
		NewName: event.name,
		TokenId: event.tokenId,
		Date:    event.date,
	}
}

func genericToPutForSaleModel(event genericEvent) models.PutForSaleEvent {
	return models.PutForSaleEvent{
		TokenId:      event.tokenId,
		PriceInEther: event.priceInEther,
		Date:         event.date,
	}
}

func genericToRemovedFromSaleModel(event genericEvent) models.RemovedFromSaleEvent {
	return models.RemovedFromSaleEvent{
		TokenId: event.tokenId,
		Date:    event.date,
	}
}

func genericToSoldModel(event genericEvent) models.SoldEvent {
	return models.SoldEvent{
		NewOwner: event.sender,
		TokenId:  event.tokenId,
		Date:     event.date,
	}
}

func createdToGeneric(subscribedEvent starnotary.StarnotaryCreated) genericEvent {
	blockNumberBig, _ := big.NewInt(0).SetString(strconv.FormatUint(subscribedEvent.Raw.BlockNumber, 10), 10)
	return genericEvent{
		sender:       common.Address.Hex(subscribedEvent.Owner),
		tokenId:      subscribedEvent.TokenId.Text(10),
		coordinates:  string(subscribedEvent.Coordinates[:]),
		name:         string(subscribedEvent.Name),
		priceInEther: big.NewFloat(0),
		eventType:    eventSignatureToType[subscribedEvent.Raw.Topics[0].Hex()],

		contractHash: subscribedEvent.Raw.Address.Hex(),
		topics:       subscribedEvent.Raw.Topics,
		data:         subscribedEvent.Raw.Data,
		blockNumber:  blockNumberBig,
		txHash:       subscribedEvent.Raw.TxHash.Hex(),
		txIndex:      subscribedEvent.Raw.TxIndex,
		blockHash:    subscribedEvent.Raw.BlockHash.Hex(),
		logIndex:     subscribedEvent.Raw.Index,
		removed:      subscribedEvent.Raw.Removed,
	}
}

func changedNameToGeneric(subscribedEvent starnotary.StarnotaryChangedName) genericEvent {
	blockNumberBig, _ := big.NewInt(0).SetString(strconv.FormatUint(subscribedEvent.Raw.BlockNumber, 10), 10)
	return genericEvent{
		sender:       common.Address.Hex(subscribedEvent.Owner),
		name:         string(subscribedEvent.NewName),
		tokenId:      subscribedEvent.TokenId.Text(10),
		priceInEther: big.NewFloat(0),
		eventType:    eventSignatureToType[subscribedEvent.Raw.Topics[0].Hex()],

		contractHash: subscribedEvent.Raw.Address.Hex(),
		topics:       subscribedEvent.Raw.Topics,
		data:         subscribedEvent.Raw.Data,
		blockNumber:  blockNumberBig,
		txHash:       subscribedEvent.Raw.TxHash.Hex(),
		txIndex:      subscribedEvent.Raw.TxIndex,
		blockHash:    subscribedEvent.Raw.BlockHash.Hex(),
		logIndex:     subscribedEvent.Raw.Index,
		removed:      subscribedEvent.Raw.Removed,
	}
}

func putForSaleToGeneric(subscribedEvent starnotary.StarnotaryPutForSale) genericEvent {
	blockNumberBig, _ := big.NewInt(0).SetString(strconv.FormatUint(subscribedEvent.Raw.BlockNumber, 10), 10)
	return genericEvent{
		sender:       common.Address.Hex(subscribedEvent.Owner),
		tokenId:      subscribedEvent.TokenId.Text(10),
		priceInEther: eth.WeiToEther(subscribedEvent.PriceInWei),
		eventType:    eventSignatureToType[subscribedEvent.Raw.Topics[0].Hex()],

		contractHash: subscribedEvent.Raw.Address.Hex(),
		topics:       subscribedEvent.Raw.Topics,
		data:         subscribedEvent.Raw.Data,
		blockNumber:  blockNumberBig,
		txHash:       subscribedEvent.Raw.TxHash.Hex(),
		txIndex:      subscribedEvent.Raw.TxIndex,
		blockHash:    subscribedEvent.Raw.BlockHash.Hex(),
		logIndex:     subscribedEvent.Raw.Index,
		removed:      subscribedEvent.Raw.Removed,
	}
}

func removedFromSaleToGeneric(subscribedEvent starnotary.StarnotaryRemovedFromSale) genericEvent {
	blockNumberBig, _ := big.NewInt(0).SetString(strconv.FormatUint(subscribedEvent.Raw.BlockNumber, 10), 10)
	return genericEvent{
		sender:       common.Address.Hex(subscribedEvent.Owner),
		tokenId:      subscribedEvent.TokenId.Text(10),
		priceInEther: big.NewFloat(0),
		eventType:    eventSignatureToType[subscribedEvent.Raw.Topics[0].Hex()],

		contractHash: subscribedEvent.Raw.Address.Hex(),
		topics:       subscribedEvent.Raw.Topics,
		data:         subscribedEvent.Raw.Data,
		blockNumber:  blockNumberBig,
		txHash:       subscribedEvent.Raw.TxHash.Hex(),
		txIndex:      subscribedEvent.Raw.TxIndex,
		blockHash:    subscribedEvent.Raw.BlockHash.Hex(),
		logIndex:     subscribedEvent.Raw.Index,
		removed:      subscribedEvent.Raw.Removed,
	}
}

func soldToGeneric(subscribedEvent starnotary.StarnotarySold) genericEvent {
	blockNumberBig, _ := big.NewInt(0).SetString(strconv.FormatUint(subscribedEvent.Raw.BlockNumber, 10), 10)
	return genericEvent{
		sender:       common.Address.Hex(subscribedEvent.NewOwner),
		tokenId:      subscribedEvent.TokenId.Text(10),
		priceInEther: big.NewFloat(0),
		eventType:    eventSignatureToType[subscribedEvent.Raw.Topics[0].Hex()],

		contractHash: subscribedEvent.Raw.Address.Hex(),
		topics:       subscribedEvent.Raw.Topics,
		data:         subscribedEvent.Raw.Data,
		blockNumber:  blockNumberBig,
		txHash:       subscribedEvent.Raw.TxHash.Hex(),
		txIndex:      subscribedEvent.Raw.TxIndex,
		blockHash:    subscribedEvent.Raw.BlockHash.Hex(),
		logIndex:     subscribedEvent.Raw.Index,
		removed:      subscribedEvent.Raw.Removed,
	}
}

func scrappedToGeneric(logEvent types.Log) (event genericEvent) {
	eventSignature := logEvent.Topics[0].Hex()
	eventType := eventSignatureToType[eventSignature]

	switch eventType {
	case "Created":
		event = scrappedCreatedToGeneric(logEvent)
	case "ChangedName":
		event = scrappedChangedNameToGeneric(logEvent)
	case "PutForSale":
		event = scrappedPutForSaleToGeneric(logEvent)
	case "RemovedFromSale":
		event = scrappedRemovedFromSaleToGeneric(logEvent)
	case "Sold":
		event = scrappedSoldToGeneric(logEvent)
	default:
		logger.Error("Tried to parse a scrapped non listened event", logger.String("signature", eventSignature))
	}

	return
}

func scrappedCreatedToGeneric(logEvent types.Log) genericEvent {
	eth := eth.GetEth()
	parsedCreated, err := eth.Contract.ParseCreated(logEvent)
	if err != nil {
		logger.Error("Could not parse scrapped created event", logger.String("error", err.Error()))
	}
	return createdToGeneric(*parsedCreated)
}

func scrappedChangedNameToGeneric(logEvent types.Log) genericEvent {
	eth := eth.GetEth()
	parsedChangedName, err := eth.Contract.ParseChangedName(logEvent)
	if err != nil {
		logger.Error("Could not parse scrapped changed name event", logger.String("error", err.Error()))
	}
	return changedNameToGeneric(*parsedChangedName)
}

func scrappedPutForSaleToGeneric(logEvent types.Log) genericEvent {
	eth := eth.GetEth()
	parsedPutForSale, err := eth.Contract.ParsePutForSale(logEvent)
	if err != nil {
		logger.Error("Could not parse scrapped put for sale event", logger.String("error", err.Error()))
	}
	return putForSaleToGeneric(*parsedPutForSale)
}

func scrappedRemovedFromSaleToGeneric(logEvent types.Log) genericEvent {
	eth := eth.GetEth()
	parsedRemovedFromSale, err := eth.Contract.ParseRemovedFromSale(logEvent)
	if err != nil {
		logger.Error("Could not parse scrapped removed from sale event", logger.String("error", err.Error()))
	}
	return removedFromSaleToGeneric(*parsedRemovedFromSale)
}

func scrappedSoldToGeneric(logEvent types.Log) genericEvent {
	eth := eth.GetEth()
	parsedSold, err := eth.Contract.ParseSold(logEvent)
	if err != nil {
		logger.Error("Could not parse scrapped sold event", logger.String("error", err.Error()))
	}
	return soldToGeneric(*parsedSold)
}
