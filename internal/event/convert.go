package event

import (
	"math/big"
	"strconv"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/sergera/star-notary-listener/internal/domain"
	"github.com/sergera/star-notary-listener/internal/eth"
	"github.com/sergera/star-notary-listener/internal/gocontracts/starnotary"
	"github.com/sergera/star-notary-listener/internal/logger"
)

func genericToCreatedModel(event domain.GenericEvent) domain.CreatedEvent {
	return domain.CreatedEvent{
		Owner:       event.Sender,
		Name:        event.Name,
		TokenId:     event.TokenId,
		Coordinates: event.Coordinates,
		Date:        event.Date,
	}
}

func genericToChangedNameModel(event domain.GenericEvent) domain.ChangedNameEvent {
	return domain.ChangedNameEvent{
		NewName: event.Name,
		TokenId: event.TokenId,
		Date:    event.Date,
	}
}

func genericToPutForSaleModel(event domain.GenericEvent) domain.PutForSaleEvent {
	return domain.PutForSaleEvent{
		TokenId:      event.TokenId,
		PriceInEther: strings.TrimRight(event.PriceInEther.Text('f', 18), ".0"),
		Date:         event.Date,
	}
}

func genericToRemovedFromSaleModel(event domain.GenericEvent) domain.RemovedFromSaleEvent {
	return domain.RemovedFromSaleEvent{
		TokenId: event.TokenId,
		Date:    event.Date,
	}
}

func genericToSoldModel(event domain.GenericEvent) domain.SoldEvent {
	return domain.SoldEvent{
		NewOwner: event.Sender,
		TokenId:  event.TokenId,
		Date:     event.Date,
	}
}

func createdToGeneric(subscribedEvent starnotary.StarnotaryCreated) domain.GenericEvent {
	blockNumberBig, _ := big.NewInt(0).SetString(strconv.FormatUint(subscribedEvent.Raw.BlockNumber, 10), 10)
	return domain.GenericEvent{
		Sender:       common.Address.Hex(subscribedEvent.Owner),
		TokenId:      subscribedEvent.TokenId.Text(10),
		Coordinates:  string(subscribedEvent.Coordinates[:]),
		Name:         string(subscribedEvent.Name),
		PriceInEther: big.NewFloat(0),
		EventType:    eventSignatureToType[subscribedEvent.Raw.Topics[0].Hex()],

		ContractHash: subscribedEvent.Raw.Address.Hex(),
		Topics:       subscribedEvent.Raw.Topics,
		Data:         subscribedEvent.Raw.Data,
		BlockNumber:  blockNumberBig,
		TxHash:       subscribedEvent.Raw.TxHash.Hex(),
		TxIndex:      subscribedEvent.Raw.TxIndex,
		BlockHash:    subscribedEvent.Raw.BlockHash.Hex(),
		LogIndex:     subscribedEvent.Raw.Index,
		Removed:      subscribedEvent.Raw.Removed,
	}
}

func changedNameToGeneric(subscribedEvent starnotary.StarnotaryChangedName) domain.GenericEvent {
	blockNumberBig, _ := big.NewInt(0).SetString(strconv.FormatUint(subscribedEvent.Raw.BlockNumber, 10), 10)
	return domain.GenericEvent{
		Sender:       common.Address.Hex(subscribedEvent.Owner),
		Name:         string(subscribedEvent.NewName),
		TokenId:      subscribedEvent.TokenId.Text(10),
		PriceInEther: big.NewFloat(0),
		EventType:    eventSignatureToType[subscribedEvent.Raw.Topics[0].Hex()],

		ContractHash: subscribedEvent.Raw.Address.Hex(),
		Topics:       subscribedEvent.Raw.Topics,
		Data:         subscribedEvent.Raw.Data,
		BlockNumber:  blockNumberBig,
		TxHash:       subscribedEvent.Raw.TxHash.Hex(),
		TxIndex:      subscribedEvent.Raw.TxIndex,
		BlockHash:    subscribedEvent.Raw.BlockHash.Hex(),
		LogIndex:     subscribedEvent.Raw.Index,
		Removed:      subscribedEvent.Raw.Removed,
	}
}

func putForSaleToGeneric(subscribedEvent starnotary.StarnotaryPutForSale) domain.GenericEvent {
	blockNumberBig, _ := big.NewInt(0).SetString(strconv.FormatUint(subscribedEvent.Raw.BlockNumber, 10), 10)
	return domain.GenericEvent{
		Sender:       common.Address.Hex(subscribedEvent.Owner),
		TokenId:      subscribedEvent.TokenId.Text(10),
		PriceInEther: eth.WeiToEther(subscribedEvent.PriceInWei),
		EventType:    eventSignatureToType[subscribedEvent.Raw.Topics[0].Hex()],

		ContractHash: subscribedEvent.Raw.Address.Hex(),
		Topics:       subscribedEvent.Raw.Topics,
		Data:         subscribedEvent.Raw.Data,
		BlockNumber:  blockNumberBig,
		TxHash:       subscribedEvent.Raw.TxHash.Hex(),
		TxIndex:      subscribedEvent.Raw.TxIndex,
		BlockHash:    subscribedEvent.Raw.BlockHash.Hex(),
		LogIndex:     subscribedEvent.Raw.Index,
		Removed:      subscribedEvent.Raw.Removed,
	}
}

func removedFromSaleToGeneric(subscribedEvent starnotary.StarnotaryRemovedFromSale) domain.GenericEvent {
	blockNumberBig, _ := big.NewInt(0).SetString(strconv.FormatUint(subscribedEvent.Raw.BlockNumber, 10), 10)
	return domain.GenericEvent{
		Sender:       common.Address.Hex(subscribedEvent.Owner),
		TokenId:      subscribedEvent.TokenId.Text(10),
		PriceInEther: big.NewFloat(0),
		EventType:    eventSignatureToType[subscribedEvent.Raw.Topics[0].Hex()],

		ContractHash: subscribedEvent.Raw.Address.Hex(),
		Topics:       subscribedEvent.Raw.Topics,
		Data:         subscribedEvent.Raw.Data,
		BlockNumber:  blockNumberBig,
		TxHash:       subscribedEvent.Raw.TxHash.Hex(),
		TxIndex:      subscribedEvent.Raw.TxIndex,
		BlockHash:    subscribedEvent.Raw.BlockHash.Hex(),
		LogIndex:     subscribedEvent.Raw.Index,
		Removed:      subscribedEvent.Raw.Removed,
	}
}

func soldToGeneric(subscribedEvent starnotary.StarnotarySold) domain.GenericEvent {
	blockNumberBig, _ := big.NewInt(0).SetString(strconv.FormatUint(subscribedEvent.Raw.BlockNumber, 10), 10)
	return domain.GenericEvent{
		Sender:       common.Address.Hex(subscribedEvent.NewOwner),
		TokenId:      subscribedEvent.TokenId.Text(10),
		PriceInEther: big.NewFloat(0),
		EventType:    eventSignatureToType[subscribedEvent.Raw.Topics[0].Hex()],

		ContractHash: subscribedEvent.Raw.Address.Hex(),
		Topics:       subscribedEvent.Raw.Topics,
		Data:         subscribedEvent.Raw.Data,
		BlockNumber:  blockNumberBig,
		TxHash:       subscribedEvent.Raw.TxHash.Hex(),
		TxIndex:      subscribedEvent.Raw.TxIndex,
		BlockHash:    subscribedEvent.Raw.BlockHash.Hex(),
		LogIndex:     subscribedEvent.Raw.Index,
		Removed:      subscribedEvent.Raw.Removed,
	}
}

func scrappedToGeneric(logEvent types.Log) (event domain.GenericEvent) {
	eventSignature := logEvent.Topics[0].Hex()
	EventType := eventSignatureToType[eventSignature]

	switch EventType {
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

func scrappedCreatedToGeneric(logEvent types.Log) domain.GenericEvent {
	eth := eth.GetEth()
	parsedCreated, err := eth.Contract.ParseCreated(logEvent)
	if err != nil {
		logger.Error("Could not parse scrapped created event", logger.String("error", err.Error()))
	}
	return createdToGeneric(*parsedCreated)
}

func scrappedChangedNameToGeneric(logEvent types.Log) domain.GenericEvent {
	eth := eth.GetEth()
	parsedChangedName, err := eth.Contract.ParseChangedName(logEvent)
	if err != nil {
		logger.Error("Could not parse scrapped changed Name event", logger.String("error", err.Error()))
	}
	return changedNameToGeneric(*parsedChangedName)
}

func scrappedPutForSaleToGeneric(logEvent types.Log) domain.GenericEvent {
	eth := eth.GetEth()
	parsedPutForSale, err := eth.Contract.ParsePutForSale(logEvent)
	if err != nil {
		logger.Error("Could not parse scrapped put for sale event", logger.String("error", err.Error()))
	}
	return putForSaleToGeneric(*parsedPutForSale)
}

func scrappedRemovedFromSaleToGeneric(logEvent types.Log) domain.GenericEvent {
	eth := eth.GetEth()
	parsedRemovedFromSale, err := eth.Contract.ParseRemovedFromSale(logEvent)
	if err != nil {
		logger.Error("Could not parse scrapped Removed from sale event", logger.String("error", err.Error()))
	}
	return removedFromSaleToGeneric(*parsedRemovedFromSale)
}

func scrappedSoldToGeneric(logEvent types.Log) domain.GenericEvent {
	eth := eth.GetEth()
	parsedSold, err := eth.Contract.ParseSold(logEvent)
	if err != nil {
		logger.Error("Could not parse scrapped sold event", logger.String("error", err.Error()))
	}
	return soldToGeneric(*parsedSold)
}
