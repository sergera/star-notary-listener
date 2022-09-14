package eth

import (
	"context"
	"errors"
	"math/big"
	"strings"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/params"
	"github.com/sergera/star-notary-listener/internal/conf"
	"github.com/sergera/star-notary-listener/internal/gocontracts/starnotary"
	"github.com/sergera/star-notary-listener/internal/logger"
)

var once sync.Once
var instance *eth

type eth struct {
	Client   *ethclient.Client
	Contract *starnotary.Starnotary
	ABI      *abi.ABI
}

func GetEth() *eth {
	once.Do(func() {
		var e *eth = &eth{}
		e.setup()
		instance = e
	})
	return instance
}

func (e *eth) setup() {
	e.setClient()
	e.setContract()
	e.setABI()
	go e.avoidProviderTimeout()
}

func (e *eth) setClient() {
	conf := conf.GetConf()
	client, err := ethclient.Dial(conf.RPCProviderWebsocketURL())
	if err != nil {
		logger.Panic("could not dial eth client", logger.String("message", err.Error()))
	}

	e.Client = client
}

func (e *eth) avoidProviderTimeout() {
	for {
		_, err := e.Client.BlockNumber(context.Background())
		if err != nil {
			logger.Error("disconnected: ", logger.String("message", err.Error()))
		}
		time.Sleep(30 * time.Second)
	}
}

func (e *eth) setContract() {
	conf := conf.GetConf()
	contractAddress := common.HexToAddress(conf.ContractAddress())
	starNotary, err := starnotary.NewStarnotary(contractAddress, e.Client)
	if err != nil {
		logger.Panic("could not instance go contract", logger.String("message", err.Error()))
	}

	e.Contract = starNotary
}

func (e *eth) setABI() {
	starnotaryABI, err := abi.JSON(strings.NewReader(string(starnotary.StarnotaryMetaData.ABI)))
	if err != nil {
		logger.Panic("could not read contract ABI", logger.String("message", err.Error()))
	}

	e.ABI = &starnotaryABI
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
