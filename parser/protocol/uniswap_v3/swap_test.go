package uniswap_v3

import (
	"base_scan/parser/protocol"
	"base_scan/types"
	"base_scan/types/orm"
	"encoding/json"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestSwapEvent_TokensOrdered_Buy(t *testing.T) {
	// https://etherscan.io/tx/0xbff012bb00982626ade612b2f40e08f46281166f5b049b7c7c28847ea1018f92#eventlog#362
	ethLogGetter, pairService := protocol.PrepareTest()
	receiptLog := ethLogGetter.GetEthLog("0xbff012bb00982626ade612b2f40e08f46281166f5b049b7c7c28847ea1018f92", 8)

	event, pErr := EventParserSwap.Parse(receiptLog)
	require.NoError(t, pErr)
	require.False(t, event.CanGetPair())

	pairWrap := pairService.GetPairAndTokens(event.GetPairAddress(), event.GetProtocolId())
	event.SetPair(pairWrap.Pair)

	tx := event.GetTx(decimal.NewFromFloat(613.8234269349539585))
	expectAmt0, _ := decimal.NewFromString("46.983825885337328850")
	expectAmt1, _ := decimal.NewFromString("46.993548957558595625")
	expectTx := &orm.Tx{
		TxHash:        "0xbff012bb00982626ade612b2f40e08f46281166f5b049b7c7c28847ea1018f92",
		Event:         "buy",
		Token0Amount:  expectAmt0,
		Token1Amount:  expectAmt1,
		Maker:         "0x3686ED59A65952c31c8C578142acB83757ae6cf8",
		Token0Address: "0x97Ad75064b20fb2B2447feD4fa953bF7F007a706",
		Token1Address: "0xC02aaA39b223FE8D0A0e5C4F27eAD9083C756Cc2",
		Block:         22464846,
		BlockIndex:    131,
		TxIndex:       362,
		PairAddress:   "0x6dcba3657EE750A51A13A235B4Ed081317dA3066",
		Program:       types.ProtocolNameUniswapV3,
	}
	require.True(t, tx.Equal(expectTx))
}

func TestSwapEvent_TokensOrdered_Sell(t *testing.T) {
	// https://etherscan.io/tx/0xe782222eac634349d65a7e8209d91f85b65fea095430518550513377f83bd385#eventlog#681
	ethLogGetter, pairService := protocol.PrepareTest()
	receiptLog := ethLogGetter.GetEthLog("0xe782222eac634349d65a7e8209d91f85b65fea095430518550513377f83bd385", 2)

	event, pErr := EventParserSwap.Parse(receiptLog)
	require.NoError(t, pErr)
	require.False(t, event.CanGetPair())

	pairWrap := pairService.GetPairAndTokens(event.GetPairAddress(), event.GetProtocolId())
	event.SetPair(pairWrap.Pair)

	tx := event.GetTx(decimal.NewFromFloat(614.5725867943508445))
	txBytes, _ := json.Marshal(tx)
	t.Log(string(txBytes))
	expectAmt0, _ := decimal.NewFromString("0.207696511163875253")
	expectAmt1, _ := decimal.NewFromString("0.207587092155862667")
	expectTx := &orm.Tx{
		TxHash:        "0xe782222eac634349d65a7e8209d91f85b65fea095430518550513377f83bd385",
		Event:         "sell",
		Token0Amount:  expectAmt0,
		Token1Amount:  expectAmt1,
		Maker:         "0x06F17315A7CfFA2756Eec18a6Cc2cADef2c1Bb64",
		Token0Address: "0x97Ad75064b20fb2B2447feD4fa953bF7F007a706",
		Token1Address: "0xC02aaA39b223FE8D0A0e5C4F27eAD9083C756Cc2",
		Block:         22466456,
		BlockIndex:    144,
		TxIndex:       681,
		PairAddress:   "0x6dcba3657EE750A51A13A235B4Ed081317dA3066",
		Program:       program,
	}
	require.True(t, tx.Equal(expectTx), "expect: %v, actual: %v", expectTx, tx)
}

func TestSwapEvent_TokensNotOrdered_Buy(t *testing.T) {
	// https://etherscan.io/tx/0x176f0d565cfcc79f6e0ee22621c41ce379850855a3f31e1b05419dfa2adb1067#eventlog#34
	ethLogGetter, pairService := protocol.PrepareTest()
	receiptLog := ethLogGetter.GetEthLog("0x176f0d565cfcc79f6e0ee22621c41ce379850855a3f31e1b05419dfa2adb1067", 6)

	event, pErr := EventParserSwap.Parse(receiptLog)
	require.NoError(t, pErr)
	require.False(t, event.CanGetPair())

	pairWrap := pairService.GetPairAndTokens(event.GetPairAddress(), event.GetProtocolId())
	event.SetPair(pairWrap.Pair)

	tx := event.GetTx(decimal.NewFromFloat(613.7696023929051877))
	txBytes, _ := json.Marshal(tx)
	t.Log(string(txBytes))
	expectAmt0, _ := decimal.NewFromString("6203005.174892979")
	expectAmt1, _ := decimal.NewFromString("0.129069133597499153")
	expectTx := &orm.Tx{
		TxHash:        "0x176f0d565cfcc79f6e0ee22621c41ce379850855a3f31e1b05419dfa2adb1067",
		Event:         "buy",
		Token0Amount:  expectAmt0,
		Token1Amount:  expectAmt1,
		Maker:         "0x6a29D13603A9C6Ad5c204aE4DeAC682f5605049F",
		Token0Address: "0xf816507E690f5Aa4E29d164885EB5fa7a5627860",
		Token1Address: "0xC02aaA39b223FE8D0A0e5C4F27eAD9083C756Cc2",
		Block:         22466430,
		BlockIndex:    13,
		TxIndex:       34,
		PairAddress:   "0x41b5b06Ccf883FC1652E8d4f73A444f6Bb75e384",
		Program:       types.ProtocolNameUniswapV3,
	}
	require.True(t, tx.Equal(expectTx), "expect: %v, actual: %v", expectTx, tx)
}

func TestSwapEvent_TokensNotOrdered_Sell(t *testing.T) {
	// https://etherscan.io/tx/0x650dfbcaeec93063626375d6347f62c74539b7c5b521b2ef16263374aef28e94#eventlog#152
	ethLogGetter, pairService := protocol.PrepareTest()
	receiptLog := ethLogGetter.GetEthLog("0x650dfbcaeec93063626375d6347f62c74539b7c5b521b2ef16263374aef28e94", 4)

	event, pErr := EventParserSwap.Parse(receiptLog)
	require.NoError(t, pErr)
	require.False(t, event.CanGetPair())

	pairWrap := pairService.GetPairAndTokens(event.GetPairAddress(), event.GetProtocolId())
	event.SetPair(pairWrap.Pair)

	tx := event.GetTx(decimal.NewFromFloat(613.7696023929051877))
	txBytes, _ := json.Marshal(tx)
	t.Log(string(txBytes))
	expectAmt0, _ := decimal.NewFromString("21000000")
	expectAmt1, _ := decimal.NewFromString("0.429161850896409879")
	expectTx := &orm.Tx{
		TxHash:        "0x650dfbcaeec93063626375d6347f62c74539b7c5b521b2ef16263374aef28e94",
		Event:         "sell",
		Token0Amount:  expectAmt0,
		Token1Amount:  expectAmt1,
		Maker:         "0x6a29D13603A9C6Ad5c204aE4DeAC682f5605049F",
		Token0Address: "0xf816507E690f5Aa4E29d164885EB5fa7a5627860",
		Token1Address: "0xC02aaA39b223FE8D0A0e5C4F27eAD9083C756Cc2",
		Block:         22466428,
		BlockIndex:    57,
		TxIndex:       152,
		PairAddress:   "0x41b5b06Ccf883FC1652E8d4f73A444f6Bb75e384",
		Program:       types.ProtocolNameUniswapV3,
	}
	require.True(t, tx.Equal(expectTx), "expect: %v, actual: %v", expectTx, tx)
}
