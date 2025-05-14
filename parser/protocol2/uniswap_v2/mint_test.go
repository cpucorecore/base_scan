package uniswap_v2

import (
	"base_scan/parser/protocol2"
	"base_scan/types"
	"base_scan/types/orm"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestMintEvent_TokensOrdered(t *testing.T) {
	// https://etherscan.io/tx/0x2fc8cdd95d90c1fcd52b3373c759015ffef6a8f650e0d128d348c3bcad94d162#eventlog#314
	ethLogGetter, pairService := protocol2.PrepareTest()
	receiptLog := ethLogGetter.GetEthLog("0x2fc8cdd95d90c1fcd52b3373c759015ffef6a8f650e0d128d348c3bcad94d162", 6)

	event, pErr := EventParserMint.Parse(receiptLog)
	require.NoError(t, pErr)
	require.False(t, event.CanGetPair())

	pairWrap := pairService.GetPairAndTokens(event.GetPairAddress(), event.GetPossibleProtocolIds())
	event.SetPair(pairWrap.Pair)

	tx := event.GetTx(decimal.NewFromFloat(597.3087764529163277))
	expectAmt0, _ := decimal.NewFromString("4875.639626653")
	expectAmt1, _ := decimal.NewFromString("0.034922643770150808")
	expectTx := &orm.Tx{
		TxHash:        "0x2fc8cdd95d90c1fcd52b3373c759015ffef6a8f650e0d128d348c3bcad94d162",
		Event:         "add",
		Token0Amount:  expectAmt0,
		Token1Amount:  expectAmt1,
		Maker:         "0xf4c21a1cB819E5F7ABe6dEFde3d118D8F3D61FA7",
		Token0Address: "0x962C8A85F500519266269f77DFfBA4CEa0B46Da1",
		Token1Address: "0xC02aaA39b223FE8D0A0e5C4F27eAD9083C756Cc2",
		Block:         22462116,
		BlockIndex:    89,
		TxIndex:       314,
		PairAddress:   "0x8192D5254284a14d85a58dEF7cef5B91Bf247cd9",
		Program:       types.ProtocolNameUniswapV2,
	}
	require.True(t, expectTx.Equal(tx), "expect: %v, actual: %v", expectTx, tx)
}

func TestMintEvent_TokensNotOrdered(t *testing.T) {
	// https://etherscan.io/tx/0x5e26d9ed01582ca467e567385ca6d7be32476f2a069d226b3cc8a364fc244c5b#eventlog#143
	ethLogGetter, pairService := protocol2.PrepareTest()
	receiptLog := ethLogGetter.GetEthLog("0x5e26d9ed01582ca467e567385ca6d7be32476f2a069d226b3cc8a364fc244c5b", 5)

	event, pErr := EventParserMint.Parse(receiptLog)
	require.NoError(t, pErr)
	require.False(t, event.CanGetPair())

	pairWrap := pairService.GetPairAndTokens(event.GetPairAddress(), event.GetPossibleProtocolIds())
	event.SetPair(pairWrap.Pair)

	tx := event.GetTx(decimal.NewFromFloat(597.3087764529163277))
	expectAmt0, _ := decimal.NewFromString("18.09668611")
	expectAmt1, _ := decimal.NewFromString("0.005006873309713942")
	expectTx := &orm.Tx{
		TxHash:        "0x5e26d9ed01582ca467e567385ca6d7be32476f2a069d226b3cc8a364fc244c5b",
		Event:         "add",
		Token0Amount:  expectAmt0,
		Token1Amount:  expectAmt1,
		Maker:         "0xf4c21a1cB819E5F7ABe6dEFde3d118D8F3D61FA7",
		Token0Address: "0xE0f63A424a4439cBE457D80E4f4b51aD25b2c56C",
		Token1Address: "0xC02aaA39b223FE8D0A0e5C4F27eAD9083C756Cc2",
		Block:         22425947,
		BlockIndex:    25,
		TxIndex:       143,
		PairAddress:   "0x52c77b0CB827aFbAD022E6d6CAF2C44452eDbc39",
		Program:       types.ProtocolNameUniswapV2,
	}

	require.True(t, tx.Equal(expectTx), "expect: %v, actual: %v", expectTx, tx)
}
