package eth

import (
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/params"
	"github.com/sergera/star-notary-listener/internal/conf"
	"github.com/sergera/star-notary-listener/internal/gocontracts/starnotary"
	"github.com/sergera/star-notary-listener/internal/logger"
)

var Client *ethclient.Client
var Contract *starnotary.Starnotary
var ABI *abi.ABI

func Setup() {
	setClient()
	setContract()
	setABI()
}

func setClient() {
	client, err := ethclient.Dial(conf.RPCProviderWebsocketURL)
	if err != nil {
		logger.Panic("Could not dial eth client", logger.String("error", err.Error()))
	}

	Client = client
}

func setContract() {
	contractAddress := common.HexToAddress(conf.ContractAddress)
	starNotary, err := starnotary.NewStarnotary(contractAddress, Client)
	if err != nil {
		logger.Panic("Could not instance go contract", logger.String("error", err.Error()))
	}

	Contract = starNotary
}

func setABI() {
	starnotaryABI, err := abi.JSON(strings.NewReader(string(starnotary.StarnotaryMetaData.ABI)))
	if err != nil {
		logger.Panic("Could not read contract ABI", logger.String("error", err.Error()))
	}

	ABI = &starnotaryABI
}

func WeiToEther(wei *big.Int) *big.Float {
	f := new(big.Float)
	f.SetPrec(236) //  IEEE 754 octuple-precision binary floating-point format: binary256
	f.SetMode(big.ToNearestEven)
	fWei := new(big.Float)
	fWei.SetPrec(236) //  IEEE 754 octuple-precision binary floating-point format: binary256
	fWei.SetMode(big.ToNearestEven)
	return f.Quo(fWei.SetInt(wei), big.NewFloat(params.Ether))
}
