package uniswap_v2

import (
	"base_scan/parser/protocol2"
	"base_scan/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestSyncEvent_TokenOrdered(t *testing.T) {
	// https://etherscan.io/tx/0x15a1e1638edf8d2154df3a2bbed898e99451861a172bb968ff765b02a7c5f00d#eventlog#30
	ethLogGetter, pairService := protocol2.PrepareTest()
	receiptLog := ethLogGetter.GetEthLog("0x15a1e1638edf8d2154df3a2bbed898e99451861a172bb968ff765b02a7c5f00d", 4)

	event, pErr := EventParserSync.Parse(receiptLog)
	require.NoError(t, pErr)
	require.False(t, event.CanGetPair())

	pairWrap := pairService.GetPairAndTokens(event.GetPairAddress(), event.GetPossibleProtocolIds())
	event.SetPair(pairWrap.Pair)

	require.True(t, event.CanGetPoolUpdate())
	poolUpdate := event.GetPoolUpdate()
	expectAmt0, _ := decimal.NewFromString("12333258.060059299")
	expectAmt1, _ := decimal.NewFromString("65.610598598622848807")
	expectPoolUpdate := &types.PoolUpdate{
		Program:       program,
		LogIndex:      30,
		Address:       common.HexToAddress("0x8192D5254284a14d85a58dEF7cef5B91Bf247cd9"),
		Token0Address: common.HexToAddress("0x962C8A85F500519266269f77DFfBA4CEa0B46Da1"),
		Token1Address: common.HexToAddress("0xC02aaA39b223FE8D0A0e5C4F27eAD9083C756Cc2"),
		Token0Amount:  expectAmt0,
		Token1Amount:  expectAmt1,
	}
	require.True(t, expectPoolUpdate.Equal(poolUpdate))
}

func TestSyncEvent_TokenNotOrdered(t *testing.T) {
	// https://etherscan.io/tx/0xbd905a7b53b5e1649bde10101a86f6898c37257ead60cf933690c2e391f718ec#eventlog#197
	ethLogGetter, pairService := protocol2.PrepareTest()
	receiptLog := ethLogGetter.GetEthLog("0xbd905a7b53b5e1649bde10101a86f6898c37257ead60cf933690c2e391f718ec", 4)

	event, pErr := EventParserSync.Parse(receiptLog)
	require.NoError(t, pErr)
	require.False(t, event.CanGetPair())

	pairWrap := pairService.GetPairAndTokens(event.GetPairAddress(), event.GetPossibleProtocolIds())
	event.SetPair(pairWrap.Pair)

	require.True(t, event.CanGetPoolUpdate())
	poolUpdate := event.GetPoolUpdate()
	expectAmt0, _ := decimal.NewFromString("8275276.45401660")
	expectAmt1, _ := decimal.NewFromString("2766.702481325104089903")
	expectPoolUpdate := &types.PoolUpdate{
		Program:       program,
		LogIndex:      197,
		Address:       common.HexToAddress("0x52c77b0CB827aFbAD022E6d6CAF2C44452eDbc39"),
		Token0Address: common.HexToAddress("0xE0f63A424a4439cBE457D80E4f4b51aD25b2c56C"),
		Token1Address: common.HexToAddress("0xC02aaA39b223FE8D0A0e5C4F27eAD9083C756Cc2"),
		Token0Amount:  expectAmt0,
		Token1Amount:  expectAmt1,
	}
	require.True(t, expectPoolUpdate.Equal(poolUpdate), "expect: %v, actual: %v", expectPoolUpdate, poolUpdate)
}
