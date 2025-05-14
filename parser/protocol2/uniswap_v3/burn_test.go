package uniswap_v3

import (
	"base_scan/parser/protocol2"
	"base_scan/types/orm"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestBurnEvent_TokensOrdered(t *testing.T) {
	// https://etherscan.io/tx/0x09d3714d936513bfc2e36b7c96420da3824c0f273c860e463f360de76cc68f75#eventlog#146
	ethLogGetter, pairService := protocol2.PrepareTest()
	receiptLog := ethLogGetter.GetEthLog("0x09d3714d936513bfc2e36b7c96420da3824c0f273c860e463f360de76cc68f75", 0)

	event, pErr := EventParserBurn.Parse(receiptLog)
	require.NoError(t, pErr)

	pairWrap := pairService.GetPairAndTokens(event.GetPairAddress(), event.GetPossibleProtocolIds())
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
		Program:       program,
	}
	require.True(t, tx.Equal(expectTx))
}

func TestBurnEvent_TokensNotOrdered(t *testing.T) {
	// https://etherscan.io/tx/0x9bb1a9900708f46dea96d6815c57d7b4a5ba47db6bcea5cefe1d75c1fd9e71f5#eventlog#27
	ethLogGetter, pairService := protocol2.PrepareTest()
	receiptLog := ethLogGetter.GetEthLog("0x9bb1a9900708f46dea96d6815c57d7b4a5ba47db6bcea5cefe1d75c1fd9e71f5", 0)

	event, pErr := EventParserBurn.Parse(receiptLog)
	require.NoError(t, pErr)

	pairWrap := pairService.GetPairAndTokens(event.GetPairAddress(), event.GetPossibleProtocolIds())
	event.SetPair(pairWrap.Pair)

	tx := event.GetTx(decimal.NewFromFloat(600.7833261828169461))
	require.Equal(t, tx.Event, "remove")
	expectAmt0, _ := decimal.NewFromString("675887701.296705099")
	expectAmt1, _ := decimal.NewFromString("12.622476492218617898")
	expectTx := &orm.Tx{
		TxHash:        "0x9bb1a9900708f46dea96d6815c57d7b4a5ba47db6bcea5cefe1d75c1fd9e71f5",
		Event:         "remove",
		Token0Amount:  expectAmt0,
		Token1Amount:  expectAmt1,
		Maker:         "0x6e79aF7d51f35fC3E8C266a364566aeD4401586D",
		Token0Address: "0xf816507E690f5Aa4E29d164885EB5fa7a5627860",
		Token1Address: "0xC02aaA39b223FE8D0A0e5C4F27eAD9083C756Cc2",
		Block:         22466385,
		BlockIndex:    3,
		TxIndex:       27,
		PairAddress:   "0x41b5b06Ccf883FC1652E8d4f73A444f6Bb75e384",
		Program:       program,
	}
	require.True(t, tx.Equal(expectTx), "expect: %v, actual: %v", expectTx, tx)
}
