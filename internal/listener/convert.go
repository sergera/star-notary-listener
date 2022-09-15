package listener

import (
	"math/big"
	"strconv"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/sergera/star-notary-listener/internal/domain"
	"github.com/sergera/star-notary-listener/internal/eth"
	"github.com/sergera/star-notary-listener/internal/gocontracts/starnotary"
	"github.com/sergera/star-notary-listener/internal/logger"
)

func createToGeneric(subscribedEvent starnotary.StarnotaryCreate) domain.GenericEvent {
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

func changeNameToGeneric(subscribedEvent starnotary.StarnotaryChangeName) domain.GenericEvent {
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

func removeFromSaleToGeneric(subscribedEvent starnotary.StarnotaryRemoveFromSale) domain.GenericEvent {
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

func purchaseToGeneric(subscribedEvent starnotary.StarnotaryPurchase) domain.GenericEvent {
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
	case "Create":
		event = scrappedCreateToGeneric(logEvent)
	case "ChangeName":
		event = scrappedChangeNameToGeneric(logEvent)
	case "PutForSale":
		event = scrappedPutForSaleToGeneric(logEvent)
	case "RemoveFromSale":
		event = scrappedRemoveFromSaleToGeneric(logEvent)
	case "Purchase":
		event = scrappedPurchaseToGeneric(logEvent)
	default:
		logger.Error("tried to parse a scrapped non listened event", logger.String("signature", eventSignature))
	}

	return
}

func scrappedCreateToGeneric(logEvent types.Log) domain.GenericEvent {
	eth := eth.GetEth()
	parsedCreate, err := eth.Contract.ParseCreate(logEvent)
	if err != nil {
		logger.Error("could not parse scrapped create event", logger.String("message", err.Error()))
	}
	return createToGeneric(*parsedCreate)
}

func scrappedChangeNameToGeneric(logEvent types.Log) domain.GenericEvent {
	eth := eth.GetEth()
	parsedChangeName, err := eth.Contract.ParseChangeName(logEvent)
	if err != nil {
		logger.Error("could not parse scrapped changed Name event", logger.String("message", err.Error()))
	}
	return changeNameToGeneric(*parsedChangeName)
}

func scrappedPutForSaleToGeneric(logEvent types.Log) domain.GenericEvent {
	eth := eth.GetEth()
	parsedPutForSale, err := eth.Contract.ParsePutForSale(logEvent)
	if err != nil {
		logger.Error("could not parse scrapped put for sale event", logger.String("message", err.Error()))
	}
	return putForSaleToGeneric(*parsedPutForSale)
}

func scrappedRemoveFromSaleToGeneric(logEvent types.Log) domain.GenericEvent {
	eth := eth.GetEth()
	parsedRemoveFromSale, err := eth.Contract.ParseRemoveFromSale(logEvent)
	if err != nil {
		logger.Error("could not parse scrapped Removed from sale event", logger.String("message", err.Error()))
	}
	return removeFromSaleToGeneric(*parsedRemoveFromSale)
}

func scrappedPurchaseToGeneric(logEvent types.Log) domain.GenericEvent {
	eth := eth.GetEth()
	parsedPurchase, err := eth.Contract.ParsePurchase(logEvent)
	if err != nil {
		logger.Error("could not parse scrapped purchase event", logger.String("message", err.Error()))
	}
	return purchaseToGeneric(*parsedPurchase)
}
