package service

import (
	"base_scan/cache"
	"base_scan/config"
	"base_scan/service/contract_caller"
	"base_scan/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/go-redis/redis/v8"
	"github.com/stretchr/testify/require"
	"testing"
)

func TearDown(c cache.Cache, pairAddress, token0Address, token1Address common.Address) {
	c.DelPair(pairAddress)
	c.DelToken(token0Address)
	c.DelToken(token1Address)
}

func TestPairService_GetPairAndTokens_PancakeV2_TokenOrdered(t *testing.T) {
	redisCli := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
	c := cache.NewTwoTierCache(redisCli)

	ethClient, err := ethclient.Dial(config.G.Chain.Endpoint)
	if err != nil {
		t.Fatal(err)
	}

	cc := contract_caller.NewContractCaller(ethClient, config.G.ContractCaller.Retry.GetRetryParams())
	ps := NewPairService(c, cc)

	pairAddress := common.HexToAddress("0xbC42145d5A574EDe9b8860FCa2A49EB7B239Efa5")
	token0Address := common.HexToAddress("0x92aa03137385F18539301349dcfC9EbC923fFb10")
	token1Address := common.HexToAddress("0xbb4CdB9CBd36B01bD1cBaEBF2De08d9173bc095c")
	protocolId := types.ProtocolIdPancakeV2
	expectToken0 := &types.TokenCore{
		Address:  token0Address,
		Symbol:   "SKYAI",
		Decimals: 18,
	}
	expectToken1 := &types.TokenCore{
		Address:  token1Address,
		Symbol:   "WETH",
		Decimals: 18,
	}
	expectPair := &types.Pair{
		Address:        pairAddress,
		TokensReversed: false,
		Token0Core:     expectToken0,
		Token1Core:     expectToken1,
		ProtocolId:     protocolId,
		Filtered:       false,
	}

	defer func() {
		TearDown(c, pairAddress, token0Address, token1Address)
	}()

	pairWrap := ps.GetPairAndTokens(pairAddress, protocolId)
	require.Equal(t, pairWrap.NewPair, true)
	require.Equal(t, pairWrap.NewToken0, true)
	require.Equal(t, pairWrap.NewToken1, true)
	require.True(t, pairWrap.Pair.Equal(expectPair), "pair equal failed", pairWrap.Pair, expectPair)
	require.Equal(t, pairWrap.Pair.TokensReversed, false)

	pairWrapFromCache := ps.GetPairAndTokens(pairAddress, protocolId)
	require.Equal(t, pairWrapFromCache.NewPair, false)
	require.Equal(t, pairWrapFromCache.NewToken0, false)
	require.Equal(t, pairWrapFromCache.NewToken1, false)
	require.True(t, pairWrapFromCache.Pair.Equal(expectPair))
	require.Equal(t, pairWrapFromCache.Pair.TokensReversed, false)
}

func TestPairService_GetTokens_PancakeV2_TokenOrdered(t *testing.T) {
	redisCli := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
	c := cache.NewTwoTierCache(redisCli)

	ethClient, err := ethclient.Dial(config.G.Chain.Endpoint)
	if err != nil {
		t.Fatal(err)
	}

	cc := contract_caller.NewContractCaller(ethClient, config.G.ContractCaller.Retry.GetRetryParams())
	ps := NewPairService(c, cc)

	pairAddress := common.HexToAddress("0x0010B1D2D807182638D19Ec9b6f5beD3E24a5EF9")
	token0Address := common.HexToAddress("0xF251D850898758775958691Df66895d0b5F837AD")
	token1Address := common.HexToAddress("0xbb4CdB9CBd36B01bD1cBaEBF2De08d9173bc095c")
	protocolId := types.ProtocolIdPancakeV2
	expectToken0 := &types.TokenCore{
		Address:  token0Address,
		Symbol:   "SAP",
		Decimals: 9,
	}
	expectToken1 := &types.TokenCore{
		Address:  token1Address,
		Symbol:   "WETH",
		Decimals: 18,
	}
	expectPair := &types.Pair{
		Address:        pairAddress,
		TokensReversed: true,
		Token0Core:     expectToken0,
		Token1Core:     expectToken1,
		ProtocolId:     protocolId,
		Filtered:       false,
	}

	defer func() {
		TearDown(c, pairAddress, token0Address, token1Address)
	}()

	pair := &types.Pair{
		Address:        pairAddress,
		TokensReversed: true,
		Token0Core:     &types.TokenCore{Address: token0Address},
		Token1Core:     &types.TokenCore{Address: token1Address},
		ProtocolId:     protocolId,
		Filtered:       false,
	}
	pairWrap := ps.GetTokens(pair)
	require.Equal(t, pairWrap.NewPair, true)
	require.Equal(t, pairWrap.NewToken0, true)
	require.Equal(t, pairWrap.NewToken1, true)
	require.True(t, pairWrap.Pair.Equal(expectPair))
	require.Equal(t, pairWrap.Pair.TokensReversed, true)

	pairWrapFromCache := ps.GetPairAndTokens(pairAddress, protocolId)
	require.Equal(t, pairWrapFromCache.NewPair, false)
	require.Equal(t, pairWrapFromCache.NewToken0, false)
	require.Equal(t, pairWrapFromCache.NewToken1, false)
	require.True(t, pairWrapFromCache.Pair.Equal(expectPair))
	require.Equal(t, pairWrapFromCache.Pair.TokensReversed, true)
}

func TestPairService_GetPairAndTokens_PancakeV2_TokenNotOrdered(t *testing.T) {
	redisCli := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
	c := cache.NewTwoTierCache(redisCli)

	ethClient, err := ethclient.Dial(config.G.Chain.Endpoint)
	if err != nil {
		t.Fatal(err)
	}

	cc := contract_caller.NewContractCaller(ethClient, config.G.ContractCaller.Retry.GetRetryParams())
	ps := NewPairService(c, cc)

	pairAddress := common.HexToAddress("0x0010B1D2D807182638D19Ec9b6f5beD3E24a5EF9")
	token0Address := common.HexToAddress("0xF251D850898758775958691Df66895d0b5F837AD")
	token1Address := common.HexToAddress("0xbb4CdB9CBd36B01bD1cBaEBF2De08d9173bc095c")
	protocolId := types.ProtocolIdPancakeV2
	expectToken0 := &types.TokenCore{
		Address:  token0Address,
		Symbol:   "SAP",
		Decimals: 9,
	}
	expectToken1 := &types.TokenCore{
		Address:  token1Address,
		Symbol:   "WETH",
		Decimals: 18,
	}
	expectPair := &types.Pair{
		Address:        pairAddress,
		TokensReversed: true,
		Token0Core:     expectToken0,
		Token1Core:     expectToken1,
		ProtocolId:     protocolId,
		Filtered:       false,
	}

	defer func() {
		TearDown(c, pairAddress, token0Address, token1Address)
	}()

	pairWrap := ps.GetPairAndTokens(pairAddress, protocolId)
	require.Equal(t, pairWrap.NewPair, true)
	require.Equal(t, pairWrap.NewToken0, true)
	require.Equal(t, pairWrap.NewToken1, true)
	require.True(t, pairWrap.Pair.Equal(expectPair), "pair equal failed", pairWrap.Pair, expectPair)
	require.Equal(t, pairWrap.Pair.TokensReversed, true)

	pairWrapFromCache := ps.GetPairAndTokens(pairAddress, protocolId)
	require.Equal(t, pairWrapFromCache.NewPair, false)
	require.Equal(t, pairWrapFromCache.NewToken0, false)
	require.Equal(t, pairWrapFromCache.NewToken1, false)
	require.True(t, pairWrapFromCache.Pair.Equal(expectPair))
	require.Equal(t, pairWrapFromCache.Pair.TokensReversed, true)
}

func TestPairService_GetPairAndTokens_PancakeV3_TokenNotOrdered(t *testing.T) {
	redisCli := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
	c := cache.NewTwoTierCache(redisCli)

	ethClient, err := ethclient.Dial(config.G.Chain.Endpoint)
	if err != nil {
		t.Fatal(err)
	}

	cc := contract_caller.NewContractCaller(ethClient, config.G.ContractCaller.Retry.GetRetryParams())
	ps := NewPairService(c, cc)

	pairAddress := common.HexToAddress("0x00004Aae5aFF462fF6Ca621624A43d8A212FBFcA")
	token0Address := common.HexToAddress("0xde7EE76A6004157038463C3a316917DEE8863113")
	token1Address := common.HexToAddress("0xbb4CdB9CBd36B01bD1cBaEBF2De08d9173bc095c")
	protocolId := types.ProtocolIdPancakeV3
	expectToken0 := &types.TokenCore{
		Address:  token0Address,
		Symbol:   "RIPDOGE",
		Decimals: 18,
	}
	expectToken1 := &types.TokenCore{
		Address:  token1Address,
		Symbol:   "WETH",
		Decimals: 18,
	}
	expectPair := &types.Pair{
		Address:        pairAddress,
		TokensReversed: true,
		Token0Core:     expectToken0,
		Token1Core:     expectToken1,
		ProtocolId:     protocolId,
		Filtered:       false,
	}

	defer func() {
		TearDown(c, pairAddress, token0Address, token1Address)
	}()

	pairWrap := ps.GetPairAndTokens(pairAddress, protocolId)
	require.Equal(t, pairWrap.NewPair, true)
	require.Equal(t, pairWrap.NewToken0, true)
	require.Equal(t, pairWrap.NewToken1, true)
	require.True(t, pairWrap.Pair.Equal(expectPair))
	require.Equal(t, pairWrap.Pair.TokensReversed, true)

	pairWrapFromCache := ps.GetPairAndTokens(pairAddress, 3)
	require.Equal(t, pairWrapFromCache.NewPair, false)
	require.Equal(t, pairWrapFromCache.NewToken0, false)
	require.Equal(t, pairWrapFromCache.NewToken1, false)
	require.True(t, pairWrapFromCache.Pair.Equal(expectPair))
	require.Equal(t, pairWrapFromCache.Pair.TokensReversed, true)
}

func TestPairService_GetTokens_PancakeV3_TokenNotOrdered(t *testing.T) {
	redisCli := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
	c := cache.NewTwoTierCache(redisCli)

	ethClient, err := ethclient.Dial(config.G.Chain.Endpoint)
	if err != nil {
		t.Fatal(err)
	}

	cc := contract_caller.NewContractCaller(ethClient, config.G.ContractCaller.Retry.GetRetryParams())
	ps := NewPairService(c, cc)

	pairAddress := common.HexToAddress("0x00004Aae5aFF462fF6Ca621624A43d8A212FBFcA")
	token0Address := common.HexToAddress("0xde7EE76A6004157038463C3a316917DEE8863113")
	token1Address := common.HexToAddress("0xbb4CdB9CBd36B01bD1cBaEBF2De08d9173bc095c")
	protocolId := types.ProtocolIdPancakeV3
	expectToken0 := &types.TokenCore{
		Address:  token0Address,
		Symbol:   "RIPDOGE",
		Decimals: 18,
	}
	expectToken1 := &types.TokenCore{
		Address:  token1Address,
		Symbol:   "WETH",
		Decimals: 18,
	}
	expectPair := &types.Pair{
		Address:        pairAddress,
		TokensReversed: true,
		Token0Core:     expectToken0,
		Token1Core:     expectToken1,
		ProtocolId:     protocolId,
		Filtered:       false,
	}

	defer func() {
		TearDown(c, pairAddress, token0Address, token1Address)
	}()

	pair := &types.Pair{
		Address:        pairAddress,
		TokensReversed: true,
		Token0Core:     &types.TokenCore{Address: token0Address},
		Token1Core:     &types.TokenCore{Address: token1Address},
		ProtocolId:     protocolId,
		Filtered:       false,
	}
	pairWrap := ps.GetTokens(pair)
	require.Equal(t, pairWrap.NewPair, true)
	require.Equal(t, pairWrap.NewToken0, true)
	require.Equal(t, pairWrap.NewToken1, true)
	require.True(t, pairWrap.Pair.Equal(expectPair))
	require.True(t, pair.Equal(expectPair))
	require.Equal(t, pairWrap.Pair.TokensReversed, true)

	pairWrapFromCache := ps.GetPairAndTokens(pairAddress, 3)
	require.Equal(t, pairWrapFromCache.NewPair, false)
	require.Equal(t, pairWrapFromCache.NewToken0, false)
	require.Equal(t, pairWrapFromCache.NewToken1, false)
	require.True(t, pairWrapFromCache.Pair.Equal(expectPair))
	require.True(t, pair.Equal(expectPair))
	require.Equal(t, pairWrapFromCache.Pair.TokensReversed, true)
}

func TestPairService_GetTokens_FourMeme(t *testing.T) {
	redisCli := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
	c := cache.NewTwoTierCache(redisCli)

	ethClient, err := ethclient.Dial(config.G.Chain.Endpoint)
	if err != nil {
		t.Fatal(err)
	}

	cc := contract_caller.NewContractCaller(ethClient, config.G.ContractCaller.Retry.GetRetryParams())
	ps := NewPairService(c, cc)

	pairAddress := common.HexToAddress("0x00000471e70672d9Be3C277A3258AeC13E6Da7fd")
	token0Address := common.HexToAddress("0x00000471e70672d9Be3C277A3258AeC13E6Da7fd")
	token1Address := types.WETHAddress
	protocolId := types.ProtocolIdFourMeme
	expectToken0 := &types.TokenCore{
		Address:  token0Address,
		Symbol:   "1000000000",
		Decimals: 18,
	}
	expectToken1 := &types.TokenCore{
		Address:  token1Address,
		Symbol:   "WETH",
		Decimals: 18,
	}

	expectPair := &types.Pair{
		Address:        pairAddress,
		TokensReversed: false,
		Token0Core:     expectToken0,
		Token1Core:     expectToken1,
		ProtocolId:     protocolId,
		Filtered:       false,
	}

	defer func() {
		TearDown(c, pairAddress, token0Address, token1Address)
	}()

	pairWrap := ps.GetTokens(&types.Pair{
		Address:        pairAddress,
		TokensReversed: false,
		Token0Core:     &types.TokenCore{Address: token0Address},
		Token1Core:     &types.TokenCore{Address: token1Address},
		ProtocolId:     protocolId,
	})
	require.Equal(t, pairWrap.NewPair, true)
	require.Equal(t, pairWrap.NewToken0, true)
	require.Equal(t, pairWrap.NewToken1, true)
	require.True(t, pairWrap.Pair.Equal(expectPair))
	require.Equal(t, pairWrap.Pair.TokensReversed, false)

	pairWrapFromCache := ps.GetPairAndTokens(pairAddress, 3)
	require.Equal(t, pairWrapFromCache.NewPair, false)
	require.Equal(t, pairWrapFromCache.NewToken0, false)
	require.Equal(t, pairWrapFromCache.NewToken1, false)
	require.True(t, pairWrapFromCache.Pair.Equal(expectPair))
	require.Equal(t, pairWrapFromCache.Pair.TokensReversed, false)
}
