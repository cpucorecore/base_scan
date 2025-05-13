package uniswap_v3

import (
	"base_scan/parser/protocol"
	"base_scan/types/orm"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestMintEvent_TokenOrdered(t *testing.T) {
	// https://etherscan.io/tx/0x3170ef052778d7ca093f248bcb7bde106e21c051245342e31c437d90da1887d0#eventlog#37
	ethLogGetter, pairService := protocol.PrepareTest()
	receiptLog := ethLogGetter.GetEthLog("0x3170ef052778d7ca093f248bcb7bde106e21c051245342e31c437d90da1887d0", 3)

	event, pErr := EventParserMint.Parse(receiptLog)
	require.NoError(t, pErr)

	pairWrap := pairService.GetPairAndTokens(event.GetPairAddress(), event.GetProtocolId())
	event.SetPair(pairWrap.Pair)

	tx := event.GetTx(decimal.NewFromFloat(614.2276379916430686))
	expectAmt0, _ := decimal.NewFromString("100.137739759714037512")
	expectAmt1, _ := decimal.NewFromString("757.130098870203544310")
	expectTx := &orm.Tx{
		TxHash:        "0x3170ef052778d7ca093f248bcb7bde106e21c051245342e31c437d90da1887d0",
		Event:         "add",
		Token0Amount:  expectAmt0,
		Token1Amount:  expectAmt1,
		Maker:         "0x039e31f5eee7ec10d3E1104514F64c7FEd7A717d",
		Token0Address: "0x97Ad75064b20fb2B2447feD4fa953bF7F007a706",
		Token1Address: "0xC02aaA39b223FE8D0A0e5C4F27eAD9083C756Cc2",
		Block:         22454909,
		BlockIndex:    1,
		TxIndex:       37,
		PairAddress:   "0x6dcba3657EE750A51A13A235B4Ed081317dA3066",
		Program:       program,
	}
	require.True(t, tx.Equal(expectTx))
}

func TestMintEvent_TokenNotOrdered(t *testing.T) {
	// https://etherscan.io/tx/0xe584f7432f7a78c445926a24183c18d5304a63e223027db4c082447b6c48fb1c#eventlog#5
	ethLogGetter, pairService := protocol.PrepareTest()
	receiptLog := ethLogGetter.GetEthLog("0xe584f7432f7a78c445926a24183c18d5304a63e223027db4c082447b6c48fb1c", 5)

	event, pErr := EventParserMint.Parse(receiptLog)
	require.NoError(t, pErr)

	pairWrap := pairService.GetPairAndTokens(event.GetPairAddress(), event.GetProtocolId())
	event.SetPair(pairWrap.Pair)

	tx := event.GetTx(decimal.NewFromFloat(617.9940181222048624))
	expectAmt0, _ := decimal.NewFromString("438146865.549444102")
	expectAmt1, _ := decimal.NewFromString("16.823272345900202140")
	expectTx := &orm.Tx{
		TxHash:        "0xe584f7432f7a78c445926a24183c18d5304a63e223027db4c082447b6c48fb1c",
		Event:         "add",
		Token0Amount:  expectAmt0,
		Token1Amount:  expectAmt1,
		Maker:         "0xdAe033566f063Ccf9631bAc75881c48E0d6fEfDA",
		Token0Address: "0xf816507E690f5Aa4E29d164885EB5fa7a5627860",
		Token1Address: "0xC02aaA39b223FE8D0A0e5C4F27eAD9083C756Cc2",
		Block:         22466385,
		BlockIndex:    0,
		TxIndex:       5,
		PairAddress:   "0x41b5b06Ccf883FC1652E8d4f73A444f6Bb75e384",
		Program:       program,
	}
	require.True(t, tx.Equal(expectTx), "expect: %v, actual: %v", expectTx, tx)
}
