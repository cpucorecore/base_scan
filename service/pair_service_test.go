package service

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestPairService_GetPair_UniswapV2(t *testing.T) {
	tc := GetTestContext()
	pair := pairUniswapV2

	pw := tc.PairService.GetPair(pair.address, possibleProtocolIds)
	require.False(t, pw.Pair.Filtered, "pair should not be filtered")
	require.Equal(t, true, pw.NewPair)
	require.Equal(t, true, pw.NewToken0)
	require.Equal(t, true, pw.NewToken1)
	require.True(t, pw.Pair.Equal(pair.GetPair()), "pair should be equal", pw.Pair, pair.GetPair())

	pw = tc.PairService.GetPair(pair.address, possibleProtocolIds)
	require.False(t, pw.Pair.Filtered, "pair should not be filtered")
	require.Equal(t, false, pw.NewPair)
	require.Equal(t, false, pw.NewToken0)
	require.Equal(t, false, pw.NewToken1)
	require.True(t, pw.Pair.Equal(pair.GetPair()), "pair should be equal", pw.Pair, pair.GetPair())
}

func TestPairService_GetPair_UniswapV3(t *testing.T) {
	tc := GetTestContext()
	pair := pairUniswapV3

	pw := tc.PairService.GetPair(pair.address, possibleProtocolIds)
	require.False(t, pw.Pair.Filtered, "pair should not be filtered")
	require.Equal(t, true, pw.NewPair)
	require.Equal(t, true, pw.NewToken0)
	require.Equal(t, true, pw.NewToken1)
	require.True(t, pw.Pair.Equal(pair.GetPair()), "pair should be equal", pw.Pair, pair.GetPair())

	pw = tc.PairService.GetPair(pair.address, possibleProtocolIds)
	require.False(t, pw.Pair.Filtered, "pair should not be filtered")
	require.Equal(t, false, pw.NewPair)
	require.Equal(t, false, pw.NewToken0)
	require.Equal(t, false, pw.NewToken1)
	require.True(t, pw.Pair.Equal(pair.GetPair()), "pair should be equal", pw.Pair, pair.GetPair())
}

func TestPairService_GetPair_PancakeV2(t *testing.T) {
	tc := GetTestContext()
	pair := pairPancakeV2

	pw := tc.PairService.GetPair(pair.address, possibleProtocolIds)
	require.False(t, pw.Pair.Filtered, "pair should not be filtered")
	require.Equal(t, true, pw.NewPair)
	require.Equal(t, true, pw.NewToken0)
	require.Equal(t, true, pw.NewToken1)
	require.True(t, pw.Pair.Equal(pair.GetPair()), "pair should be equal", pw.Pair, pair.GetPair())

	pw = tc.PairService.GetPair(pair.address, possibleProtocolIds)
	require.False(t, pw.Pair.Filtered, "pair should not be filtered")
	require.Equal(t, false, pw.NewPair)
	require.Equal(t, false, pw.NewToken0)
	require.Equal(t, false, pw.NewToken1)
	require.True(t, pw.Pair.Equal(pair.GetPair()), "pair should be equal", pw.Pair, pair.GetPair())
}

func TestPairService_GetPair_PancakeV3(t *testing.T) {
	tc := GetTestContext()
	pair := pairPancakeV3

	pw := tc.PairService.GetPair(pair.address, possibleProtocolIds)
	require.False(t, pw.Pair.Filtered, "pair should not be filtered")
	require.Equal(t, true, pw.NewPair)
	require.Equal(t, true, pw.NewToken0)
	require.Equal(t, true, pw.NewToken1)
	require.True(t, pw.Pair.Equal(pair.GetPair()), "pair should be equal", pw.Pair, pair.GetPair())

	pw = tc.PairService.GetPair(pair.address, possibleProtocolIds)
	require.False(t, pw.Pair.Filtered, "pair should not be filtered")
	require.Equal(t, false, pw.NewPair)
	require.Equal(t, false, pw.NewToken0)
	require.Equal(t, false, pw.NewToken1)
	require.True(t, pw.Pair.Equal(pair.GetPair()), "pair should be equal", pw.Pair, pair.GetPair())
}

func TestPairService_GetPair_Aerodrome(t *testing.T) {
	tc := GetTestContext()
	pair := pairAerodrome

	pw := tc.PairService.GetPair(pair.address, possibleProtocolIds)
	require.False(t, pw.Pair.Filtered, "pair should not be filtered")
	require.Equal(t, true, pw.NewPair)
	require.Equal(t, true, pw.NewToken0)
	require.Equal(t, true, pw.NewToken1)
	require.True(t, pw.Pair.Equal(pair.GetPair()), "pair should be equal", pw.Pair, pair.GetPair())

	pw = tc.PairService.GetPair(pair.address, possibleProtocolIds)
	require.False(t, pw.Pair.Filtered, "pair should not be filtered")
	require.Equal(t, false, pw.NewPair)
	require.Equal(t, false, pw.NewToken0)
	require.Equal(t, false, pw.NewToken1)
	require.True(t, pw.Pair.Equal(pair.GetPair()), "pair should be equal", pw.Pair, pair.GetPair())
}

func TestPairService_GetGetPairTokens_UniswapV2(t *testing.T) {
	tc := GetTestContext()
	pair := pairUniswapV2
	expectPair := pair.GetPair()

	pairWithoutTokenInfo := pair.GetPairWithoutTokenInfo()
	pw := tc.PairService.GetPairTokens(pairWithoutTokenInfo)
	require.False(t, pw.Pair.Filtered, "pair should not be filtered")
	require.Equal(t, true, pw.NewPair)
	require.Equal(t, true, pw.NewToken0)
	require.Equal(t, true, pw.NewToken1)
	require.True(t, pw.Pair.Equal(expectPair), "pair should be equal", pw.Pair, expectPair)

	pairWithoutTokenInfo = pair.GetPairWithoutTokenInfo()
	pw = tc.PairService.GetPairTokens(pairWithoutTokenInfo)
	require.False(t, pw.Pair.Filtered, "pair should not be filtered")
	require.Equal(t, false, pw.NewPair)
	require.Equal(t, false, pw.NewToken0)
	require.Equal(t, false, pw.NewToken1)
	require.True(t, pw.Pair.Equal(expectPair), "pair should be equal", pw.Pair, expectPair)
}

func TestPairService_GetGetPairTokens_UniswapV3(t *testing.T) {
	tc := GetTestContext()
	pair := pairUniswapV3
	expectPair := pair.GetPair()

	pairWithoutTokenInfo := pair.GetPairWithoutTokenInfo()
	pw := tc.PairService.GetPairTokens(pairWithoutTokenInfo)
	require.False(t, pw.Pair.Filtered, "pair should not be filtered")
	require.Equal(t, true, pw.NewPair)
	require.Equal(t, true, pw.NewToken0)
	require.Equal(t, true, pw.NewToken1)
	require.True(t, pw.Pair.Equal(expectPair), "pair should be equal", pw.Pair, expectPair)

	pairWithoutTokenInfo = pair.GetPairWithoutTokenInfo()
	pw = tc.PairService.GetPairTokens(pairWithoutTokenInfo)
	require.False(t, pw.Pair.Filtered, "pair should not be filtered")
	require.Equal(t, false, pw.NewPair)
	require.Equal(t, false, pw.NewToken0)
	require.Equal(t, false, pw.NewToken1)
	require.True(t, pw.Pair.Equal(expectPair), "pair should be equal", pw.Pair, expectPair)
}

func TestPairService_GetGetPairTokens_PancakeV2(t *testing.T) {
	tc := GetTestContext()
	pair := pairPancakeV2
	expectPair := pair.GetPair()

	pairWithoutTokenInfo := pair.GetPairWithoutTokenInfo()
	pw := tc.PairService.GetPairTokens(pairWithoutTokenInfo)
	require.False(t, pw.Pair.Filtered, "pair should not be filtered")
	require.Equal(t, true, pw.NewPair)
	require.Equal(t, true, pw.NewToken0)
	require.Equal(t, true, pw.NewToken1)
	require.True(t, pw.Pair.Equal(expectPair), "pair should be equal", pw.Pair, expectPair)

	pairWithoutTokenInfo = pair.GetPairWithoutTokenInfo()
	pw = tc.PairService.GetPairTokens(pairWithoutTokenInfo)
	require.False(t, pw.Pair.Filtered, "pair should not be filtered")
	require.Equal(t, false, pw.NewPair)
	require.Equal(t, false, pw.NewToken0)
	require.Equal(t, false, pw.NewToken1)
	require.True(t, pw.Pair.Equal(expectPair), "pair should be equal", pw.Pair, expectPair)
}

func TestPairService_GetGetPairTokens_PancakeV3(t *testing.T) {
	tc := GetTestContext()
	pair := pairPancakeV3
	expectPair := pair.GetPair()

	pairWithoutTokenInfo := pair.GetPairWithoutTokenInfo()
	pw := tc.PairService.GetPairTokens(pairWithoutTokenInfo)
	require.False(t, pw.Pair.Filtered, "pair should not be filtered")
	require.Equal(t, true, pw.NewPair)
	require.Equal(t, true, pw.NewToken0)
	require.Equal(t, true, pw.NewToken1)
	require.True(t, pw.Pair.Equal(expectPair), "pair should be equal", pw.Pair, expectPair)

	pairWithoutTokenInfo = pair.GetPairWithoutTokenInfo()
	pw = tc.PairService.GetPairTokens(pairWithoutTokenInfo)
	require.False(t, pw.Pair.Filtered, "pair should not be filtered")
	require.Equal(t, false, pw.NewPair)
	require.Equal(t, false, pw.NewToken0)
	require.Equal(t, false, pw.NewToken1)
	require.True(t, pw.Pair.Equal(expectPair), "pair should be equal", pw.Pair, expectPair)
}

func TestPairService_GetGetPairTokens_Aerodrome(t *testing.T) {
	tc := GetTestContext()
	pair := pairAerodrome
	expectPair := pair.GetPair()

	pairWithoutTokenInfo := pair.GetPairWithoutTokenInfo()
	pw := tc.PairService.GetPairTokens(pairWithoutTokenInfo)
	require.False(t, pw.Pair.Filtered, "pair should not be filtered")
	require.Equal(t, true, pw.NewPair)
	require.Equal(t, true, pw.NewToken0)
	require.Equal(t, true, pw.NewToken1)
	require.True(t, pw.Pair.Equal(expectPair), "pair should be equal", pw.Pair, expectPair)

	pairWithoutTokenInfo = pair.GetPairWithoutTokenInfo()
	pw = tc.PairService.GetPairTokens(pairWithoutTokenInfo)
	require.False(t, pw.Pair.Filtered, "pair should not be filtered")
	require.Equal(t, false, pw.NewPair)
	require.Equal(t, false, pw.NewToken0)
	require.Equal(t, false, pw.NewToken1)
	require.True(t, pw.Pair.Equal(expectPair), "pair should be equal", pw.Pair, expectPair)
}
