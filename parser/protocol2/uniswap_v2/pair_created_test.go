package uniswap_v2

import (
	"base_scan/parser/protocol2"
	"base_scan/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestPairCreatedEvent_TokenOrdered(t *testing.T) {
	// https://etherscan.io/tx/0xe1785ba060973af51be7086ce22f3569e21674f5cb0a43ca736d9ffdb7fcbdf1#eventlog#1
	ethLogGetter, pairService := protocol2.PrepareTest()
	receiptLog := ethLogGetter.GetEthLog("0xe1785ba060973af51be7086ce22f3569e21674f5cb0a43ca736d9ffdb7fcbdf1", 1)
	blockTimestamp := ethLogGetter.GetBlockTimestamp(receiptLog.BlockNumber)

	event, pErr := EventParserPairCreated.Parse(receiptLog)
	require.NoError(t, pErr)

	event.SetBlockTime(time.Unix(int64(blockTimestamp), 0))
	require.True(t, event.CanGetPair())
	pair := event.GetPair()
	pairWrap := pairService.GetTokens(pair)
	event.SetPair(pairWrap.Pair)

	expectPair := &types.Pair{
		Address:        common.HexToAddress("0x8192D5254284a14d85a58dEF7cef5B91Bf247cd9"),
		TokensReversed: false,
		Token0Core: &types.TokenCore{
			Address:  common.HexToAddress("0x962C8A85F500519266269f77DFfBA4CEa0B46Da1"),
			Symbol:   "BERRY",
			Decimals: 9,
		},
		Token1Core: &types.TokenCore{
			Address:  common.HexToAddress("0xC02aaA39b223FE8D0A0e5C4F27eAD9083C756Cc2"),
			Symbol:   "WETH",
			Decimals: 18,
		},
		Block:      22457534,
		BlockAt:    time.Unix(1746935411, 0),
		ProtocolId: types.ProtocolIdUniswapV2,
		Filtered:   false,
		FilterCode: 0,
	}

	require.True(t, pairWrap.Pair.Equal(expectPair))
}

func TestPairCreatedEvent_TokenNotOrdered(t *testing.T) {
	// https://etherscan.io/tx/0xf09ac8ad7e21d15ded627a176ec718903baae5e5a9ce671a611bd852691b24f9#eventlog#87
	ethLogGetter, pairService := protocol2.PrepareTest()
	receiptLog := ethLogGetter.GetEthLog("0xf09ac8ad7e21d15ded627a176ec718903baae5e5a9ce671a611bd852691b24f9", 1)
	blockTimestamp := ethLogGetter.GetBlockTimestamp(receiptLog.BlockNumber)

	event, pErr := EventParserPairCreated.Parse(receiptLog)
	require.NoError(t, pErr)

	event.SetBlockTime(time.Unix(int64(blockTimestamp), 0))
	require.True(t, event.CanGetPair())
	pair := event.GetPair()
	pairWrap := pairService.GetTokens(pair)
	event.SetPair(pairWrap.Pair)

	expectPair := &types.Pair{
		Address:        common.HexToAddress("0x52c77b0CB827aFbAD022E6d6CAF2C44452eDbc39"),
		TokensReversed: true,
		Token0Core: &types.TokenCore{
			Address:  common.HexToAddress("0xE0f63A424a4439cBE457D80E4f4b51aD25b2c56C"),
			Symbol:   "SPX",
			Decimals: 8,
		},
		Token1Core: &types.TokenCore{
			Address:  common.HexToAddress("0xC02aaA39b223FE8D0A0e5C4F27eAD9083C756Cc2"),
			Symbol:   "WETH",
			Decimals: 18,
		},
		Block:      17924533,
		BlockAt:    time.Unix(1692154331, 0),
		ProtocolId: types.ProtocolIdUniswapV2,
		Filtered:   false,
		FilterCode: 0,
	}

	require.True(t, pairWrap.Pair.Equal(expectPair), "expect: %v, actual: %v", expectPair, pairWrap.Pair)
}
