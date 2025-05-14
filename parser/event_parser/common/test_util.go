package common

import (
	"base_scan/cache"
	"base_scan/config"
	"base_scan/service"
	"base_scan/service/contract_caller"
	"context"
	"github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"math/big"
)

type EthLogGetter struct {
	ethClient      *ethclient.Client
	cache          cache.Cache
	contractCaller *contract_caller.ContractCaller
	pairService    service.PairService
}

func PrepareTest() (*EthLogGetter, service.PairService) {
	ethClient, err := ethclient.Dial("https://base-rpc.publicnode.com")
	if err != nil {
		panic(err)
	}

	contractCaller := contract_caller.NewContractCaller(ethClient, config.G.ContractCaller.Retry.GetRetryParams())
	cache := cache.MockCache{}
	pairService := service.NewPairService(cache, contractCaller)

	return &EthLogGetter{
		ethClient:      ethClient,
		cache:          cache,
		contractCaller: contractCaller,
		pairService:    pairService,
	}, pairService
}

func (g *EthLogGetter) GetEthLog(txHashStr string, logIndex int) *ethtypes.Log {
	txHash := common.HexToHash(txHashStr)
	txReceipt, apiErr := g.ethClient.TransactionReceipt(context.Background(), txHash)
	if apiErr != nil {
		panic(apiErr)
	}

	return txReceipt.Logs[logIndex]
}

func (g *EthLogGetter) GetBlockTimestamp(blockNumber uint64) uint64 {
	blockHeader, err := g.ethClient.HeaderByNumber(context.Background(), big.NewInt(int64(blockNumber)))
	if err != nil {
		panic(err)
	}
	return blockHeader.Time
}
