package event_parser

import (
	"base_scan/parser/event_parser/common"
	"base_scan/types"
	"base_scan/types/orm"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestMint_Aerodrome(t *testing.T) {
	// https://basescan.org/tx/0xff85824c89b77fb78641d11d20738817dbc7fdd0dbad9e791b4e8b2ad8f1a4e7#eventlog#72
	txHash := "0xff85824c89b77fb78641d11d20738817dbc7fdd0dbad9e791b4e8b2ad8f1a4e7"
	ethLogGetter, pairService := common.PrepareTest()
	receiptLog := ethLogGetter.GetEthLog(txHash, 6)

	event, pErr := Topic2EventParser[receiptLog.Topics[0]].Parse(receiptLog)
	require.NoError(t, pErr)

	pairWrap := pairService.GetPairAndTokens(event.GetPairAddress(), event.GetPossibleProtocolIds())
	event.SetPair(pairWrap.Pair)

	tx := event.GetTx(common.MockNativeTokenPrice)
	expectAmt0, _ := decimal.NewFromString("10000")
	expectAmt1, _ := decimal.NewFromString("4")
	expectTx := &orm.Tx{
		TxHash:        txHash,
		Event:         common.EventNameAdd,
		Token0Amount:  expectAmt0,
		Token1Amount:  expectAmt1,
		Token0Address: "0x1E50309675d5C41D38Ba14133B4DB5b9f44FfBCd",
		Token1Address: types.WETH,
		Block:         30217994,
		BlockIndex:    186,
		TxIndex:       72,
		PairAddress:   "0xC09F68906B1DC60F1BA5771Ec6625cA947031Aaf",
		Program:       types.ProtocolNameAerodrome,
	}
	require.True(t, tx.Equal(expectTx), "expect: %v, actual: %v", expectTx, tx)
}

func TestMint_UniswapV2(t *testing.T) {
	// TODO
}

func TestMint_PancakeV2(t *testing.T) {
	// TODO
}
