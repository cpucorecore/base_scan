package service

import (
	"base_scan/cache"
	"base_scan/config"
	"base_scan/service/contract_caller"
	"github.com/ethereum/go-ethereum/ethclient"
	"math/big"
	"testing"
)

func TestPriceService_GetBNBPrice(t *testing.T) {
	c := cache.MockCache{}

	ethClient, err := ethclient.Dial(config.G.Chain.EndpointArchive)
	if err != nil {
		t.Fatal(err)
	}

	cc := contract_caller.NewContractCaller(ethClient, config.G.ContractCaller.Retry.GetRetryParams())

	ps := NewPriceService(c, cc, ethClient, 0)
	price, err := ps.GetNativeTokenPrice(big.NewInt(22466005))
	if err != nil {
		t.Fatal(err)
	}
	t.Log(price)
}
