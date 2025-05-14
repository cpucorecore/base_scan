package uniswap_v2

import (
	"base_scan/parser/protocol2"
	"base_scan/types"
	"base_scan/types/orm"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestSwapEvent_TokensOrdered_Buy(t *testing.T) {
	// https://basescan.org/tx/0x18115c0256f29b56e07e84fc580602a4b9dd0e562c2988c0ff662257b4d1a1e7#eventlog#186
	ethLogGetter, pairService := protocol2.PrepareTest()
	receiptLog := ethLogGetter.GetEthLog("0x18115c0256f29b56e07e84fc580602a4b9dd0e562c2988c0ff662257b4d1a1e7", 5)

	event, pErr := EventParserSwap.Parse(receiptLog)
	require.NoError(t, pErr)
	require.False(t, event.CanGetPair())

	pairWrap := pairService.GetPairAndTokens(event.GetPairAddress(), event.GetPossibleProtocolIds())
	event.SetPair(pairWrap.Pair)

	tx := event.GetTx(decimal.NewFromFloat(604.0751501742621019))
	expectAmt0, _ := decimal.NewFromString("765555165000000000000")
	expectAmt1, _ := decimal.NewFromString("363398")
	expectTx := &orm.Tx{
		TxHash:        "0x18115c0256f29b56e07e84fc580602a4b9dd0e562c2988c0ff662257b4d1a1e7",
		Event:         "buy",
		Token0Amount:  expectAmt0,
		Token1Amount:  expectAmt1,
		Maker:         "0x43517018e4fdCb3026E08C7F9F9a61DB143735FE",
		Token0Address: "0x962C8A85F500519266269f77DFfBA4CEa0B46Da1",
		Token1Address: "0xC02aaA39b223FE8D0A0e5C4F27eAD9083C756Cc2",
		Block:         22465490,
		BlockIndex:    36,
		TxIndex:       31,
		PairAddress:   "0x8192D5254284a14d85a58dEF7cef5B91Bf247cd9",
		Program:       types.ProtocolNameUniswapV2,
	}
	require.True(t, tx.Equal(expectTx))
}

func TestSwapEvent_TokensOrdered_Sell(t *testing.T) {
	// https://etherscan.io/tx/0x82f84df12f550c29d2d6e48951edc859b53831b3cd07fc5ebc28971488a412cb#eventlog#228
	ethLogGetter, pairService := protocol2.PrepareTest()
	receiptLog := ethLogGetter.GetEthLog("0x82f84df12f550c29d2d6e48951edc859b53831b3cd07fc5ebc28971488a412cb", 4)

	event, pErr := EventParserSwap.Parse(receiptLog)
	require.NoError(t, pErr)
	require.False(t, event.CanGetPair())

	pairWrap := pairService.GetPairAndTokens(event.GetPairAddress(), event.GetPossibleProtocolIds())
	event.SetPair(pairWrap.Pair)

	tx := event.GetTx(decimal.NewFromFloat(605.3452528642224193))
	expectAmt0, _ := decimal.NewFromString("13144.598878441")
	expectAmt1, _ := decimal.NewFromString("0.058020489364277553")
	expectTx := &orm.Tx{
		TxHash:        "0x82f84df12f550c29d2d6e48951edc859b53831b3cd07fc5ebc28971488a412cb",
		Event:         "sell",
		Token0Amount:  expectAmt0,
		Token1Amount:  expectAmt1,
		Maker:         "0xD8Db782c4A3ffCb14e4928f43f01F3bC3FBe2AB0",
		Token0Address: "0x962C8A85F500519266269f77DFfBA4CEa0B46Da1",
		Token1Address: "0xC02aaA39b223FE8D0A0e5C4F27eAD9083C756Cc2",
		Block:         22466303,
		BlockIndex:    39,
		TxIndex:       228,
		PairAddress:   "0x8192D5254284a14d85a58dEF7cef5B91Bf247cd9",
		Program:       types.ProtocolNameUniswapV2,
	}
	require.True(t, tx.Equal(expectTx), "expect: %v, actual: %v", expectTx, tx)
}

func TestSwapEvent_TokensNotOrdered_Buy(t *testing.T) {
	// https://etherscan.io/tx/0xd644390df8be285cccfe0c9680208677e3fddaf2416fe63692f815065a40c6be#eventlog#861
	ethLogGetter, pairService := protocol2.PrepareTest()
	receiptLog := ethLogGetter.GetEthLog("0xd644390df8be285cccfe0c9680208677e3fddaf2416fe63692f815065a40c6be", 4)

	event, pErr := EventParserSwap.Parse(receiptLog)
	require.NoError(t, pErr)
	require.False(t, event.CanGetPair())

	pairWrap := pairService.GetPairAndTokens(event.GetPairAddress(), event.GetPossibleProtocolIds())
	event.SetPair(pairWrap.Pair)

	tx := event.GetTx(decimal.NewFromFloat(605.3414860237204781))
	expectAmt0, _ := decimal.NewFromString("38.33809097")
	expectAmt1, _ := decimal.NewFromString("0.012866000286028690")
	expectTx := &orm.Tx{
		TxHash:        "0xd644390df8be285cccfe0c9680208677e3fddaf2416fe63692f815065a40c6be",
		Event:         "buy",
		Token0Amount:  expectAmt0,
		Token1Amount:  expectAmt1,
		Maker:         "0xa45FDf8fb671Dc138603E9F180dae920B45FABB0",
		Token0Address: "0xE0f63A424a4439cBE457D80E4f4b51aD25b2c56C",
		Token1Address: "0xC02aaA39b223FE8D0A0e5C4F27eAD9083C756Cc2",
		Block:         22466272,
		BlockIndex:    369,
		TxIndex:       861,
		PairAddress:   "0x52c77b0CB827aFbAD022E6d6CAF2C44452eDbc39",
		Program:       types.ProtocolNameUniswapV2,
	}
	require.True(t, tx.Equal(expectTx), "expect: %v, actual: %v", expectTx, tx)
}

func TestSwapEvent_TokensNotOrdered_Sell(t *testing.T) {
	// https://etherscan.io/tx/0x398141d7ac7639b3cc49cea0d9e1cc6f194396b5ee474a9836756c1fe314bc02#eventlog#19
	ethLogGetter, pairService := protocol2.PrepareTest()
	receiptLog := ethLogGetter.GetEthLog("0x398141d7ac7639b3cc49cea0d9e1cc6f194396b5ee474a9836756c1fe314bc02", 13)

	event, pErr := EventParserSwap.Parse(receiptLog)
	require.NoError(t, pErr)
	require.False(t, event.CanGetPair())

	pairWrap := pairService.GetPairAndTokens(event.GetPairAddress(), event.GetPossibleProtocolIds())
	event.SetPair(pairWrap.Pair)

	tx := event.GetTx(decimal.NewFromFloat(605.3783974833613165))
	expectAmt0, _ := decimal.NewFromString("5254.60418965")
	expectAmt1, _ := decimal.NewFromString("1.743230050558972947")
	expectTx := &orm.Tx{
		TxHash:        "0x398141d7ac7639b3cc49cea0d9e1cc6f194396b5ee474a9836756c1fe314bc02",
		Event:         "sell",
		Token0Amount:  expectAmt0,
		Token1Amount:  expectAmt1,
		Maker:         "0x01FD0484014142E12598cC4Db20db6AC7E07703a",
		Token0Address: "0xE0f63A424a4439cBE457D80E4f4b51aD25b2c56C",
		Token1Address: "0xC02aaA39b223FE8D0A0e5C4F27eAD9083C756Cc2",
		Block:         22466334,
		BlockIndex:    1,
		TxIndex:       19,
		PairAddress:   "0x52c77b0CB827aFbAD022E6d6CAF2C44452eDbc39",
		Program:       types.ProtocolNameUniswapV2,
	}
	require.True(t, tx.Equal(expectTx), "expect: %v, actual: %v", expectTx, tx)
}
