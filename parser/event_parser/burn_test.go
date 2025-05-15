package event_parser

import (
	"base_scan/parser/event_parser/common"
	"base_scan/types"
	"base_scan/types/orm"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestBurn_Aerodrome(t *testing.T) {
	// https://basescan.org/tx/0xeb50cf26b45a8ca72d9343dab433f5eeecaeb3eabc2716a7dd80ddef966948b1#eventlog#213
	ethLogGetter, pairService := common.PrepareTest()
	receiptLog := ethLogGetter.GetEthLog("0xeb50cf26b45a8ca72d9343dab433f5eeecaeb3eabc2716a7dd80ddef966948b1", 6)

	event, pErr := Topic2EventParser[receiptLog.Topics[0]].Parse(receiptLog)
	require.NoError(t, pErr)

	pairWrap := pairService.GetPairAndTokens(event.GetPairAddress(), event.GetPossibleProtocolIds())
	event.SetPair(pairWrap.Pair)

	tx := event.GetTx(common.MockNativeTokenPrice)
	expectAmt0, _ := decimal.NewFromString("4047.640731408680145311")
	expectAmt1, _ := decimal.NewFromString("9.88229999999999995")
	expectTx := &orm.Tx{
		TxHash:        "0xeb50cf26b45a8ca72d9343dab433f5eeecaeb3eabc2716a7dd80ddef966948b1",
		Event:         common.EventNameRemove,
		Token0Amount:  expectAmt0,
		Token1Amount:  expectAmt1,
		Token0Address: "0x1E50309675d5C41D38Ba14133B4DB5b9f44FfBCd",
		Token1Address: types.WETH,
		Block:         30233109,
		BlockIndex:    66,
		TxIndex:       213,
		PairAddress:   "0xC09F68906B1DC60F1BA5771Ec6625cA947031Aaf",
		Program:       types.ProtocolNameAerodrome,
	}
	require.True(t, tx.Equal(expectTx), "expect: %v, actual: %v", expectTx, tx)
}

func TestBurn_UniswapV2(t *testing.T) {
	// TODO
}

func TestBurn_PancakeV2(t *testing.T) {
	// TODO
}
