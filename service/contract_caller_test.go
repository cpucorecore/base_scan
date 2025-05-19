package service

import (
	pancakev2 "base_scan/abi/pancake/v2"
	pancakev3 "base_scan/abi/pancake/v3"
	uniswapv2 "base_scan/abi/uniswap/v2"
	uniswapv3 "base_scan/abi/uniswap/v3"
	"base_scan/config"
	"base_scan/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/require"
	"math/big"
	"testing"
)

type TestToken struct {
	address  common.Address
	name     string
	symbol   string
	decimals int
}

type TestPair struct {
	address common.Address
	token0  *TestToken
	token1  *TestToken
	fee     *big.Int
}

var (
	tokenWETH = &TestToken{
		address:  types.WETHAddress,
		name:     "Wrapped Ether",
		symbol:   "WETH",
		decimals: 18,
	}
	tokenAERO = &TestToken{
		address:  common.HexToAddress("0x940181a94A35A4569E4529A3CDfB74e38FD98631"),
		name:     "Aerodrome",
		symbol:   "AERO",
		decimals: 18,
	}
	tokenUSDC = &TestToken{
		address:  common.HexToAddress("0x833589fCD6eDb6E08f4c7C32D4f71b54bdA02913"),
		name:     "USD Coin",
		symbol:   "USDC",
		decimals: 6,
	}
	tokenPEPE = &TestToken{
		address:  common.HexToAddress("0x52b492a33E447Cdb854c7FC19F1e57E8BfA1777D"),
		name:     "BasedPepe",
		symbol:   "PEPE",
		decimals: 18,
	}
	tokenCAKE = &TestToken{
		address:  common.HexToAddress("0x3055913c90Fcc1A6CE9a358911721eEb942013A1"),
		name:     "PancakeSwap Token",
		symbol:   "Cake",
		decimals: 18,
	}
	tokenDEGEN = &TestToken{
		address:  common.HexToAddress("0x4ed4E862860beD51a9570b96d89aF5E1B0Efefed"),
		name:     "Degen",
		symbol:   "DEGEN",
		decimals: 18,
	}
	pairAerodrome = &TestPair{
		address: common.HexToAddress("0x7f670f78B17dEC44d5Ef68a48740b6f8849cc2e6"),
		token0:  tokenWETH,
		token1:  tokenAERO,
	}
	pairUniswapV2 = &TestPair{
		address: common.HexToAddress("0x88A43bbDF9D098eEC7bCEda4e2494615dfD9bB9C"),
		token0:  tokenWETH,
		token1:  tokenUSDC,
	}
	pairUniswapV3 = &TestPair{
		address: common.HexToAddress("0x0FB597D6cFE5bE0d5258A7f017599C2A4Ece34c7"),
		token0:  tokenWETH,
		token1:  tokenPEPE,
		fee:     big.NewInt(10000),
	}
	pairPancakeV2 = &TestPair{
		address: common.HexToAddress("0xc637ab6D3aB0c55a7812B0b23955bA6E40859447"),
		token0:  tokenCAKE,
		token1:  tokenWETH,
	}
	pairPancakeV3 = &TestPair{
		address: common.HexToAddress("0x54D281c7cc029a9Dd71F9ACb7487dd95B1EecF5a"),
		token0:  tokenWETH,
		token1:  tokenDEGEN,
		fee:     big.NewInt(500),
	}
)

func TestContractCaller_CallContract(t *testing.T) {
	cc := GetTestContext().ContractCaller
	address := common.HexToAddress("0x4200000000000000000000000000000000000006")
	req := &types.CallContractReq{
		Address: &address,
	}

	// call erc20 contract with a method not exist, should return non err and empty bytes
	req.Data = types.Name2Data["getReserves"]
	bytes, err := cc.CallContract(req)
	require.Nil(t, err)
	require.Equal(t, 0, len(bytes))

	// call erc20 contract with a method exist, should return non err and non-empty bytes
	req.Data = types.Name2Data["name"]
	bytes, err = cc.CallContract(req)
	require.Nil(t, err)
	require.True(t, len(bytes) > 0)
}

func TestContractCaller_queryValues(t *testing.T) {
	cc := GetTestContext().ContractCaller
	pairAddress := common.HexToAddress("0xc9034c3E7F58003E6ae0C8438e7c8f4598d5ACAA")

	// call pair contract with a method not exist, should return err and empty values
	values, err := cc.queryValues(&pairAddress, "name", 1)
	require.Equal(t, ErrOutputEmpty, err)
	require.Equal(t, 0, len(values))

	// call pair contract with a method exist, should return non err and non-empty values
	values, err = cc.queryValues(&pairAddress, "token0", 1)
	require.Nil(t, err)
	require.True(t, len(values) > 0)
}

func TestContractCaller_CallXX(t *testing.T) {
	cc := GetTestContext().ContractCaller

	tests := []struct {
		address        string
		expectName     string
		expectSymbol   string
		expectDecimals int
	}{
		{
			address:        "0x4200000000000000000000000000000000000006",
			expectName:     "Wrapped Ether",
			expectSymbol:   "WETH",
			expectDecimals: 18,
		},
		{
			address:        "0x833589fCD6eDb6E08f4c7C32D4f71b54bdA02913",
			expectName:     "USD Coin",
			expectSymbol:   "USDC",
			expectDecimals: 6,
		},
		{
			address:        "0xC2DC84144f625B4feC5a21B888028EBD6c95E38d",
			expectName:     "nightmare.exe",
			expectSymbol:   "nightmare.exe",
			expectDecimals: 18,
		},
		{
			address:        "0xB4bD1fE69dCAA0C64fDb34075a0CC2b332Bd015e",
			expectName:     "PB64",
			expectSymbol:   "PB64",
			expectDecimals: 18,
		},
	}

	for _, test := range tests {
		address := common.HexToAddress(test.address)
		name, callNameErr := cc.CallName(&address)
		require.Nil(t, callNameErr)
		require.Equal(t, test.expectName, name)
		symbol, callSymbolErr := cc.CallSymbol(&address)
		require.Nil(t, callSymbolErr)
		require.Equal(t, test.expectSymbol, symbol)
		decimals, callDecimalsErr := cc.CallDecimals(&address)
		require.Nil(t, callDecimalsErr)
		require.Equal(t, test.expectDecimals, decimals)
		totalSupply, callTotalSupplyErr := cc.CallTotalSupply(&address)
		require.Nil(t, callTotalSupplyErr)
		t.Log(address, totalSupply)
	}
}

func TestContractCaller_CallGetPair_UniswapV2(t *testing.T) {
	cc := GetTestContext().ContractCaller
	tests := []struct {
		exist         bool
		pairAddress   common.Address
		token0Address common.Address
		token1Address common.Address
	}{
		{
			exist:         true,
			pairAddress:   pairUniswapV2.address,
			token0Address: pairUniswapV2.token0.address,
			token1Address: pairUniswapV2.token1.address,
		},
		{
			exist:         false,
			pairAddress:   pairPancakeV2.address,
			token0Address: pairPancakeV2.token0.address,
			token1Address: pairPancakeV2.token1.address,
		},
	}

	for _, test := range tests {
		pairAddress, err := cc.CallGetPair(&uniswapv2.FactoryAddress, &test.token0Address, &test.token1Address)
		require.Nil(t, err)
		require.Equal(t, test.exist, types.IsSameAddress(test.pairAddress, pairAddress))
	}
}

func TestContractCaller_CallGetPair_PancakeV2(t *testing.T) {
	cc := GetTestContext().ContractCaller
	tests := []struct {
		exist         bool
		pairAddress   common.Address
		token0Address common.Address
		token1Address common.Address
	}{
		{
			exist:         false,
			pairAddress:   pairUniswapV2.address,
			token0Address: pairUniswapV2.token0.address,
			token1Address: pairUniswapV2.token1.address,
		},
		{
			exist:         true,
			pairAddress:   pairPancakeV2.address,
			token0Address: pairPancakeV2.token0.address,
			token1Address: pairPancakeV2.token1.address,
		},
	}

	for _, test := range tests {
		pairAddress, err := cc.CallGetPair(&pancakev2.FactoryAddress, &test.token0Address, &test.token1Address)
		require.Nil(t, err)
		require.Equal(t, test.exist, types.IsSameAddress(test.pairAddress, pairAddress))
	}
}

func TestContractCaller_CallGetPool_UniswapV3(t *testing.T) {
	cc := GetTestContext().ContractCaller
	tests := []struct {
		exist         bool
		pairAddress   common.Address
		token0Address common.Address
		token1Address common.Address
		fee           *big.Int
	}{
		{
			exist:         true,
			pairAddress:   pairUniswapV3.address,
			token0Address: pairUniswapV3.token0.address,
			token1Address: pairUniswapV3.token1.address,
			fee:           big.NewInt(10000),
		},
		{
			exist:         false,
			pairAddress:   pairPancakeV3.address,
			token0Address: pairPancakeV3.token0.address,
			token1Address: pairPancakeV3.token1.address,
			fee:           big.NewInt(500),
		},
	}

	for _, test := range tests {
		pairAddress, err := cc.CallGetPool(&uniswapv3.FactoryAddress, &test.token0Address, &test.token1Address, test.fee)
		require.Nil(t, err)
		require.Equal(t, test.exist, types.IsSameAddress(test.pairAddress, pairAddress))
	}
}

func TestContractCaller_CallGetPool_PancakeV3(t *testing.T) {
	cc := GetTestContext().ContractCaller
	tests := []struct {
		exist         bool
		pairAddress   common.Address
		token0Address common.Address
		token1Address common.Address
		fee           *big.Int
	}{
		{
			exist:         false,
			pairAddress:   pairUniswapV3.address,
			token0Address: pairUniswapV3.token0.address,
			token1Address: pairUniswapV3.token1.address,
			fee:           big.NewInt(10000),
		},
		{
			exist:         true,
			pairAddress:   pairPancakeV3.address,
			token0Address: pairPancakeV3.token0.address,
			token1Address: pairPancakeV3.token1.address,
			fee:           big.NewInt(500),
		},
	}

	for _, test := range tests {
		pairAddress, err := cc.CallGetPool(&pancakev3.FactoryAddress, &test.token0Address, &test.token1Address, test.fee)
		require.Nil(t, err)
		require.Equal(t, test.exist, types.IsSameAddress(test.pairAddress, pairAddress))
	}
}

func TestContractCaller_CallIsPool(t *testing.T) {
	cc := GetTestContext().ContractCaller

	tests := []struct {
		isPool      bool
		pairAddress common.Address
	}{
		{
			isPool:      true,
			pairAddress: pairAerodrome.address,
		},
		{
			isPool:      false,
			pairAddress: pairUniswapV2.address,
		},
		{
			isPool:      false,
			pairAddress: pairUniswapV3.address,
		},
		{
			isPool:      false,
			pairAddress: pairPancakeV2.address,
		},
		{
			isPool:      false,
			pairAddress: pairPancakeV3.address,
		},
	}

	for _, test := range tests {
		isPool, err := cc.CallIsPool(&test.pairAddress)
		require.Nil(t, err)
		require.Equal(t, test.isPool, isPool)
	}
}

func TestContractCaller_CallToken0AndCallToken1(t *testing.T) {
	cc := GetTestContext().ContractCaller

	tests := []struct {
		pairAddress   common.Address
		token0Address common.Address
		token1Address common.Address
	}{
		{
			pairAddress:   pairAerodrome.address,
			token0Address: pairAerodrome.token0.address,
			token1Address: pairAerodrome.token1.address,
		},
		{
			pairAddress:   pairUniswapV2.address,
			token0Address: pairUniswapV2.token0.address,
			token1Address: pairUniswapV2.token1.address,
		},
		{
			pairAddress:   pairUniswapV3.address,
			token0Address: pairUniswapV3.token0.address,
			token1Address: pairUniswapV3.token1.address,
		},
		{
			pairAddress:   pairPancakeV2.address,
			token0Address: pairPancakeV2.token0.address,
			token1Address: pairPancakeV2.token1.address,
		},
		{
			pairAddress:   pairPancakeV3.address,
			token0Address: pairPancakeV3.token0.address,
			token1Address: pairPancakeV3.token1.address,
		},
	}

	for _, test := range tests {
		token0Address, err0 := cc.CallToken0(&test.pairAddress)
		require.Nil(t, err0)
		require.Equal(t, test.token0Address, token0Address)
		token1Address, err1 := cc.CallToken1(&test.pairAddress)
		require.Nil(t, err1)
		require.Equal(t, test.token1Address, token1Address)
	}
}

func TestContractCaller_CallFee(t *testing.T) {
	cc := GetTestContext().ContractCaller

	tests := []struct {
		callErr     bool
		pairAddress common.Address
		expectFee   *big.Int
	}{
		{
			callErr:     true,
			pairAddress: pairAerodrome.address,
		},
		{
			callErr:     true,
			pairAddress: pairUniswapV2.address,
		},
		{
			callErr:     false,
			pairAddress: pairUniswapV3.address,
			expectFee:   pairUniswapV3.fee,
		},
		{
			callErr:     true,
			pairAddress: pairPancakeV2.address,
		},
		{
			callErr:     false,
			pairAddress: pairPancakeV3.address,
			expectFee:   pairPancakeV3.fee,
		},
	}

	for _, test := range tests {
		fee, err := cc.CallFee(&test.pairAddress)
		if test.callErr {
			require.NotNil(t, err, test.pairAddress)
		} else {
			require.Nil(t, err, test.pairAddress)
			require.Equal(t, test.expectFee.String(), fee.String(), test.pairAddress)
		}
	}
}

func TestContractCaller_GetReservesByBlockNumber(t *testing.T) {
	t.Skip()
	ethClient, err := ethclient.Dial(config.G.Chain.Endpoint)
	if err != nil {
		t.Fatal(err)
	}
	cc := NewContractCaller(ethClient, config.G.ContractCaller.Retry.GetRetryParams())

	r0, r1, err := cc.GetReservesByBlockNumber(big.NewInt(30423400))
	if err != nil {
		t.Fatal(err)
	}
	USDCAmountDivWETHAmount := decimal.NewFromBigInt(r1, -6).Div(decimal.NewFromBigInt(r0, -18))
	t.Log(USDCAmountDivWETHAmount)
}
