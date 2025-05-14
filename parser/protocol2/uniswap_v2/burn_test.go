package uniswap_v2

import (
	"base_scan/parser/protocol2"
	"base_scan/types"
	"base_scan/types/orm"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestBurnEvent_TokensOrdered(t *testing.T) {
	// https://etherscan.io/tx/0xe665872ff1dd08676bcbfbc4bb8d56548552eb741662f56b4189e66b4d9b7f58#eventlog#213
	ethLogGetter, pairService := protocol2.PrepareTest()
	receiptLog := ethLogGetter.GetEthLog("0xe665872ff1dd08676bcbfbc4bb8d56548552eb741662f56b4189e66b4d9b7f58", 5)

	event, pErr := EventParserBurn.Parse(receiptLog)
	require.NoError(t, pErr)

	pairWrap := pairService.GetPairAndTokens(event.GetPairAddress(), event.GetPossibleProtocolIds())
	event.SetPair(pairWrap.Pair)

	tx := event.GetTx(decimal.NewFromFloat(611.9695865915750796))
	expectAmt0, _ := decimal.NewFromString("24928.900990914")
	expectAmt1, _ := decimal.NewFromString("0.116146453102480079")
	expectTx := &orm.Tx{
		TxHash:        "0xe665872ff1dd08676bcbfbc4bb8d56548552eb741662f56b4189e66b4d9b7f58",
		Event:         "remove",
		Token0Amount:  expectAmt0,
		Token1Amount:  expectAmt1,
		Maker:         "0x7e37a19a0CB97BBc9838EABa34C572887BFe88A0",
		Token0Address: "0x962C8A85F500519266269f77DFfBA4CEa0B46Da1",
		Token1Address: "0xC02aaA39b223FE8D0A0e5C4F27eAD9083C756Cc2",
		Block:         22464243,
		BlockIndex:    65,
		TxIndex:       213,
		PairAddress:   "0x8192D5254284a14d85a58dEF7cef5B91Bf247cd9",
		Program:       types.ProtocolNameUniswapV2,
	}
	require.True(t, tx.Equal(expectTx))
}

func TestBurnEvent_TokensNotOrdered(t *testing.T) {
	// https://etherscan.io/tx/0x7facda9b5c23d6de93a7a84f4245dc8321e1619f0ca39f5c76427114ca79118b#eventlog#387
	ethLogGetter, pairService := protocol2.PrepareTest()
	receiptLog := ethLogGetter.GetEthLog("0x7facda9b5c23d6de93a7a84f4245dc8321e1619f0ca39f5c76427114ca79118b", 5)

	event, pErr := EventParserBurn.Parse(receiptLog)
	require.NoError(t, pErr)

	pairWrap := pairService.GetPairAndTokens(event.GetPairAddress(), event.GetPossibleProtocolIds())
	event.SetPair(pairWrap.Pair)

	tx := event.GetTx(decimal.NewFromFloat(539.0684873312770805))
	expectAmt0, _ := decimal.NewFromString("1016.84116393")
	expectAmt1, _ := decimal.NewFromString("0.353728154862910660")
	expectTx := &orm.Tx{
		TxHash:        "0x7facda9b5c23d6de93a7a84f4245dc8321e1619f0ca39f5c76427114ca79118b",
		Event:         "remove",
		Token0Amount:  expectAmt0,
		Token1Amount:  expectAmt1,
		Maker:         "0xf4c21a1cB819E5F7ABe6dEFde3d118D8F3D61FA7",
		Token0Address: "0xE0f63A424a4439cBE457D80E4f4b51aD25b2c56C",
		Token1Address: "0xC02aaA39b223FE8D0A0e5C4F27eAD9083C756Cc2",
		Block:         22347158,
		BlockIndex:    104,
		TxIndex:       387,
		PairAddress:   "0x52c77b0CB827aFbAD022E6d6CAF2C44452eDbc39",
		Program:       types.ProtocolNameUniswapV2,
	}

	require.True(t, tx.Equal(expectTx), "expect: %v, actual: %v", expectTx, tx)
}
