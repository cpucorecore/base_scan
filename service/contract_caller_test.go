package service

import (
	pancakev2 "base_scan/abi/pancake/v2"
	pancakev3 "base_scan/abi/pancake/v3"
	"base_scan/config"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/stretchr/testify/require"
	"math/big"
	"testing"
)

func TestContractCaller_ParseToken(t *testing.T) {
	addresses := []common.Address{
		common.HexToAddress("0x6636F7B89f64202208f608DEFFa71293EEF7b466"),
		common.HexToAddress("0x4811d87B7Ab45F380Af38e5830Ab3D8A03B2F4Df"),
		common.HexToAddress("0x1736Eceea489e9afb0342612453CB4661a0Ad887"),
	}

	ethClient, err := ethclient.Dial(config.G.Chain.Endpoint)
	if err != nil {
		t.Fatal(err)
	}
	cc := NewContractCaller(ethClient, config.G.ContractCaller.Retry.GetRetryParams())

	for _, address := range addresses {
		name, callNameErr := cc.CallName(&address)
		if callNameErr != nil {
			t.Fatal(callNameErr)
		}
		t.Log(address, name)

		symbol, callSymbolErr := cc.CallSymbol(&address)
		if callSymbolErr != nil {
			t.Fatal(callSymbolErr)
		}
		t.Log(address, symbol)

		decimals, callDecimalsErr := cc.CallDecimals(&address)
		if callDecimalsErr != nil {
			t.Fatal(callDecimalsErr)
		}
		t.Log(address, decimals)

		totalSupply, callTotalSupplyErr := cc.CallTotalSupply(&address)
		if callTotalSupplyErr != nil {
			t.Fatal(callTotalSupplyErr)
		}
		t.Log(address, totalSupply)
	}
}

func TestContractCaller_ParsePairPancakeV2(t *testing.T) {
	tests := []struct {
		pairAddress    common.Address
		expectedToken0 common.Address
		expectedToken1 common.Address
	}{
		{
			pairAddress:    common.HexToAddress("0x00cc8c4549ad70d515d5AA64afd0ac99562d010d"),
			expectedToken0: common.HexToAddress("0xdcc342647a84d25220e0C4b4cDE3eD39Ff68f099"),
			expectedToken1: common.HexToAddress("0xF415bec722DF2F14A9F12f357b930529FC6166B2"),
		},
	}

	ethClient, err := ethclient.Dial(config.G.Chain.Endpoint)
	if err != nil {
		t.Fatal(err)
	}
	cc := NewContractCaller(ethClient, config.G.ContractCaller.Retry.GetRetryParams())

	for _, test := range tests {
		token0, callToken0Err := cc.CallToken0(&test.pairAddress)
		if callToken0Err != nil {
			require.Nil(t, callToken0Err)
		}
		require.Equal(t, test.expectedToken0, token0)

		token1, callToken1Err := cc.CallToken1(&test.pairAddress)
		if callToken1Err != nil {
			require.Nil(t, callToken1Err)
		}
		require.Equal(t, test.expectedToken1, token1)

		pairAddress, callErr := cc.CallGetPair(&pancakev2.FactoryAddress, &token0, &token1)
		if callErr != nil {
			require.Nil(t, callErr)
		}
		require.Equal(t, test.pairAddress, pairAddress)
	}
}

func TestContractCaller_ParsePairPancakeV3(t *testing.T) {
	tests := []struct {
		pairAddress    common.Address
		expectedToken0 common.Address
		expectedToken1 common.Address
		fee            *big.Int
	}{
		{
			pairAddress:    common.HexToAddress("0x2b303C32c2e8E5B2FA82B69AF1D263C1EBc9ed22"),
			expectedToken0: common.HexToAddress("0x8519EA49c997f50cefFa444d240fB655e89248Aa"),
			expectedToken1: common.HexToAddress("0xe9e7CEA3DedcA5984780Bafc599bD69ADd087D56"),
			fee:            big.NewInt(2500),
		},
	}

	tc := GetTestContext()

	for _, test := range tests {
		token0, callToken0Err := tc.ContractCaller.CallToken0(&test.pairAddress)
		if callToken0Err != nil {
			require.Nil(t, callToken0Err)
		}
		require.Equal(t, test.expectedToken0, token0)

		token1, callToken1Err := tc.ContractCaller.CallToken1(&test.pairAddress)
		if callToken1Err != nil {
			require.Nil(t, callToken1Err)
		}
		require.Equal(t, test.expectedToken1, token1)

		fee, callFeeErr := tc.ContractCaller.CallFee(&test.pairAddress)
		if callFeeErr != nil {
			require.Nil(t, callFeeErr)
		}
		require.Equal(t, test.fee, fee)

		pairAddress, callErr := tc.ContractCaller.CallGetPool(&pancakev3.FactoryAddress, &token0, &token1, big.NewInt(2500))
		if callErr != nil {
			require.Nil(t, callErr)
		}
		require.Equal(t, test.pairAddress, pairAddress)
	}
}

func TestContractCaller_GetBnbPrice(t *testing.T) {
	t.Skip()
	ethClient, err := ethclient.Dial(config.G.Chain.Endpoint)
	if err != nil {
		t.Fatal(err)
	}
	cc := NewContractCaller(ethClient, config.G.ContractCaller.Retry.GetRetryParams())

	price, err := cc.GetNativeTokenPriceByBlockNumber(big.NewInt(48433894))
	if err != nil {
		t.Fatal(err)
	}
	t.Log(price)
}
