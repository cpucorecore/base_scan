package uniswap_v3

import (
	"base_scan/parser/protocol"
	"base_scan/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestPoolCreatedEvent_TokenOrdered(t *testing.T) {
	// https://etherscan.io/tx/0xef1ebf77c3be94747f9c5bc68622027102bd25c0b3eb391bc69a7d8ba4b5aa79#eventlog#151
	ethLogGetter, pairService := protocol.PrepareTest()
	receiptLog := ethLogGetter.GetEthLog("0xef1ebf77c3be94747f9c5bc68622027102bd25c0b3eb391bc69a7d8ba4b5aa79", 0)
	blockTimestamp := ethLogGetter.GetBlockTimestamp(receiptLog.BlockNumber)

	event, pErr := EventParserPoolCreated.Parse(receiptLog)
	require.NoError(t, pErr)

	event.SetBlockTime(time.Unix(int64(blockTimestamp), 0))
	require.True(t, event.CanGetPair())
	pair := event.GetPair()
	pairWrap := pairService.GetTokens(pair)
	event.SetPair(pairWrap.Pair)

	expectPair := &types.Pair{
		Address:        common.HexToAddress("0x6dcba3657EE750A51A13A235B4Ed081317dA3066"),
		TokensReversed: false,
		Token0Core: &types.TokenCore{
			Address:  common.HexToAddress("0x97Ad75064b20fb2B2447feD4fa953bF7F007a706"),
			Symbol:   "beraSTONE",
			Decimals: 18,
		},
		Token1Core: &types.TokenCore{
			Address:  common.HexToAddress("0xC02aaA39b223FE8D0A0e5C4F27eAD9083C756Cc2"),
			Symbol:   "WETH",
			Decimals: 18,
		},
		Block:      21465486,
		BlockAt:    time.Unix(1734960539, 0),
		ProtocolId: types.ProtocolIdUniswapV3,
		Filtered:   false,
		FilterCode: 0,
	}

	require.True(t, pairWrap.Pair.Equal(expectPair))
}

func TestPoolCreatedEvent_TokenNotOrdered(t *testing.T) {
	// https://etherscan.io/tx/0x80d0577afed311fb803bc93dacefe68656a282bb44b5b6e50e9d49413337672e#eventlog#198
	ethLogGetter, pairService := protocol.PrepareTest()
	receiptLog := ethLogGetter.GetEthLog("0x80d0577afed311fb803bc93dacefe68656a282bb44b5b6e50e9d49413337672e", 0)
	blockTimestamp := ethLogGetter.GetBlockTimestamp(receiptLog.BlockNumber)

	event, pErr := EventParserPoolCreated.Parse(receiptLog)
	require.NoError(t, pErr)

	event.SetBlockTime(time.Unix(int64(blockTimestamp), 0))
	require.True(t, event.CanGetPair())
	pair := event.GetPair()
	pairWrap := pairService.GetTokens(pair)
	event.SetPair(pairWrap.Pair)

	expectPair := &types.Pair{
		Address:        common.HexToAddress("0x41b5b06Ccf883FC1652E8d4f73A444f6Bb75e384"),
		TokensReversed: true,
		Token0Core: &types.TokenCore{
			Address:  common.HexToAddress("0xf816507E690f5Aa4E29d164885EB5fa7a5627860"),
			Symbol:   "RATO",
			Decimals: 9,
		},
		Token1Core: &types.TokenCore{
			Address:  common.HexToAddress("0xC02aaA39b223FE8D0A0e5C4F27eAD9083C756Cc2"),
			Symbol:   "WETH",
			Decimals: 18,
		},
		Block:      22450480,
		BlockAt:    time.Unix(1746849611, 0),
		ProtocolId: protocolId,
		Filtered:   false,
		FilterCode: 0,
	}

	require.True(t, pairWrap.Pair.Equal(expectPair), "expect: %v, actual: %v", expectPair, pairWrap.Pair)
}
