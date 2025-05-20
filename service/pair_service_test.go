package service

import (
	"base_scan/cache"
	"base_scan/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
	"testing"
)

func TearDown(c cache.Cache, pairAddress, token0Address, token1Address common.Address) {
	c.DelPair(pairAddress)
	c.DelToken(token0Address)
	c.DelToken(token1Address)
}

func TestPairService_GetPairAndTokens_PancakeV2_TokenOrdered(t *testing.T) {
	tc := GetTestContext()
	pairAddress := common.HexToAddress("0x41610B9024bd46e7991c274dbBF9Fc02D36567f2")
	token0Address := common.HexToAddress("0x0b3e328455c4059EEb9e3f84b5543F74E24e7E1b")
	token1Address := common.HexToAddress("0x4200000000000000000000000000000000000006")
	protocolId := types.ProtocolIdPancakeV2
	protocolIds := []int{protocolId}
	expectToken0 := &types.TokenCore{
		Address:  token0Address,
		Symbol:   "VIRTUAL",
		Decimals: 18,
	}
	expectToken1 := &types.TokenCore{
		Address:  token1Address,
		Symbol:   "WETH",
		Decimals: 18,
	}
	expectPair := &types.Pair{
		Address:        pairAddress,
		TokensReversed: false,
		Token0Core:     expectToken0,
		Token1Core:     expectToken1,
		ProtocolId:     protocolId,
		Filtered:       false,
	}

	defer TearDown(tc.Cache, pairAddress, token0Address, token1Address)

	pairWrap := tc.PairService.GetPair(pairAddress, protocolIds)
	require.Equal(t, true, pairWrap.NewPair)
	require.Equal(t, true, pairWrap.NewToken0)
	require.Equal(t, true, pairWrap.NewToken1)
	require.True(t, pairWrap.Pair.Equal(expectPair), "pair equal failed", pairWrap.Pair, expectPair)
	require.Equal(t, false, pairWrap.Pair.TokensReversed)

	pairWrapFromCache := tc.PairService.GetPair(pairAddress, protocolIds)
	require.Equal(t, false, pairWrapFromCache.NewPair)
	require.Equal(t, false, pairWrapFromCache.NewToken0)
	require.Equal(t, false, pairWrapFromCache.NewToken1)
	require.True(t, pairWrapFromCache.Pair.Equal(expectPair))
	require.Equal(t, false, pairWrapFromCache.Pair.TokensReversed)
}
