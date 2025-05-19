package event_parser

import (
	"base_scan/repository/orm"
	"base_scan/service"
	"base_scan/types"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestBurn_UniswapV3(t *testing.T) {
	// https://basescan.org/tx/0x91db85460d929bfc664779b1bf7fc23ea47436f84bdd1782e6466cb0bb2962ef#eventlog#797
	tc := service.GetTestContext()
	receiptLog := tc.GetEthLog("0x09d3714d936513bfc2e36b7c96420da3824c0f273c860e463f360de76cc68f75", 0)

	event, pErr := Topic2EventParser[receiptLog.Topics[0]].Parse(receiptLog)
	require.NoError(t, pErr)

	pairWrap := tc.PairService.GetPairAndTokens(event.GetPairAddress(), event.GetPossibleProtocolIds())
	event.SetPair(pairWrap.Pair)

	tx := event.GetTx(decimal.NewFromFloat(582.8979382022061794))
	expectAmt0, _ := decimal.NewFromString("2.000451657793625289")
	expectAmt1, _ := decimal.NewFromString("2.000632094315897124")
	expectTx := &orm.Tx{
		TxHash:        "0x09d3714d936513bfc2e36b7c96420da3824c0f273c860e463f360de76cc68f75",
		Event:         "remove",
		Token0Amount:  expectAmt0,
		Token1Amount:  expectAmt1,
		Maker:         "0x8Cb56CFC374cbDC47d2ae6CdBFD8E54e0C7391B8",
		Token0Address: "0x97Ad75064b20fb2B2447feD4fa953bF7F007a706",
		Token1Address: "0xC02aaA39b223FE8D0A0e5C4F27eAD9083C756Cc2",
		Block:         22465581,
		BlockIndex:    58,
		TxIndex:       146,
		PairAddress:   "0x6dcba3657EE750A51A13A235B4Ed081317dA3066",
		Program:       types.ProtocolNameUniswapV3,
	}
	require.True(t, tx.Equal(expectTx))
}
