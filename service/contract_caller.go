package service

import (
	"base_scan/abi/aerodrome"
	v2 "base_scan/abi/uniswap/v2"
	v3 "base_scan/abi/uniswap/v3"
	"base_scan/config"
	"base_scan/metrics"
	"base_scan/types"
	"context"
	"errors"
	"github.com/avast/retry-go/v4"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"math/big"
	"strings"
	"time"
)

var (
	ErrUnpackerNotFound  = errors.New("unpacker not found")
	ErrOutputEmpty       = errors.New("output is empty")
	ErrWrongOutputLength = errors.New("wrong output length")
	ErrReserve0NotBigInt = errors.New("reverse0 is not *big.Int")
	ErrReserve1NotBigInt = errors.New("reverse1 is not *big.Int")
)

type ContractCaller struct {
	ctx         context.Context
	ethClient   *ethclient.Client
	retryParams *config.RetryParams
}

func NewContractCaller(ethClient *ethclient.Client, retryParams *config.RetryParams) *ContractCaller {
	return &ContractCaller{
		ctx:         context.Background(),
		ethClient:   ethClient,
		retryParams: retryParams,
	}
}

func IsRetryableErr(err error) bool {
	errMsg := err.Error()
	if strings.Contains(errMsg, "execution reverted") ||
		strings.Contains(errMsg, "out of gas") ||
		strings.Contains(errMsg, "abi: cannot marshal in to go slice") {
		return false
	}
	return true
}

func (c *ContractCaller) callContract(req *CallContractReq) ([]byte, error) {
	now := time.Now()
	bytes, err := c.ethClient.CallContract(
		c.ctx,
		ethereum.CallMsg{
			To:   req.Address,
			Data: req.Data,
		},
		req.BlockNumber,
	)
	metrics.CallContractDurationMs.Observe(float64(time.Since(now).Milliseconds()))

	if err != nil {
		if IsRetryableErr(err) {
			return nil, err
		}
		return nil, nil
	}

	return bytes, nil
}

func (c *ContractCaller) CallContract(req *CallContractReq) ([]byte, error) {
	ctxWithTimeout, _ := context.WithTimeout(c.ctx, c.retryParams.Timeout)
	return retry.DoWithData(func() ([]byte, error) {
		return c.callContract(req)
	}, c.retryParams.Attempts, c.retryParams.Delay, retry.Context(ctxWithTimeout))
}

func (c *ContractCaller) queryValues(address *common.Address, name string, outputLength int) ([]interface{}, error) {
	req := &CallContractReq{
		Address: address,
		Data:    Name2Data[name], // TODO check
	}

	bytes, err := c.CallContract(req)
	if err != nil {
		return nil, err
	}

	if len(bytes) == 0 {
		return nil, ErrOutputEmpty
	}

	upk, ok := Name2Unpacker[name]
	if !ok {
		return nil, ErrUnpackerNotFound
	}

	values, unpackErr := upk.Unpack(name, bytes, outputLength)
	if unpackErr != nil {
		return nil, unpackErr
	}

	return values, nil
}

func (c *ContractCaller) queryString(address *common.Address, name string) (string, error) {
	values, err := c.queryValues(address, name, 1)
	if err != nil {
		return "", err
	}
	return ParseString(values[0])
}

func (c *ContractCaller) CallName(address *common.Address) (string, error) {
	return c.queryString(address, "name")
}

func (c *ContractCaller) CallSymbol(address *common.Address) (string, error) {
	return c.queryString(address, "symbol")
}

func (c *ContractCaller) queryInt(address *common.Address, name string) (int, error) {
	values, err := c.queryValues(address, name, 1)
	if err != nil {
		return 0, err
	}
	return ParseInt(values[0])
}

func (c *ContractCaller) CallDecimals(address *common.Address) (int, error) {
	return c.queryInt(address, "decimals")
}

func (c *ContractCaller) queryBigInt(address *common.Address, name string) (*big.Int, error) {
	values, err := c.queryValues(address, name, 1)
	if err != nil {
		return nil, err
	}
	return ParseBigInt(values[0])
}

func (c *ContractCaller) CallTotalSupply(address *common.Address) (*big.Int, error) {
	return c.queryBigInt(address, "totalSupply")
}

func (c *ContractCaller) queryAddress(address *common.Address, name string) (common.Address, error) {
	values, err := c.queryValues(address, name, 1)
	if err != nil {
		return types.ZeroAddress, err
	}
	return ParseAddress(values[0])
}

func (c *ContractCaller) CallToken0(address *common.Address) (common.Address, error) {
	return c.queryAddress(address, "token0")
}

func (c *ContractCaller) CallToken1(address *common.Address) (common.Address, error) {
	return c.queryAddress(address, "token1")
}

/*
CallGetPair
for uniswap/pancake v2
*/
func (c *ContractCaller) CallGetPair(factoryAddress, token0Address, token1Address *common.Address) (common.Address, error) {
	req := BuildCallContractReqDynamic(nil, factoryAddress, v2.FactoryAbi, "getPair", token0Address, token1Address)

	bytes, err := c.CallContract(req)
	if err != nil {
		return types.ZeroAddress, err
	}

	if len(bytes) == 0 {
		return types.ZeroAddress, ErrOutputEmpty
	}

	values, unpackErr := PancakeV2FactoryUnpacker.Unpack("getPair", bytes, 1)
	if unpackErr != nil {
		return types.ZeroAddress, unpackErr
	}

	if len(values) != 1 {
		return types.ZeroAddress, ErrWrongOutputLength
	}

	return ParseAddress(values[0])
}

/*
CallFee
for uniswap/pancake v3
*/
func (c *ContractCaller) CallFee(address *common.Address) (*big.Int, error) {
	return c.queryBigInt(address, "fee")
}

/*
CallGetPool
for uniswap/pancake v3
*/
func (c *ContractCaller) CallGetPool(factoryAddress, token0Address, token1Address *common.Address, fee *big.Int) (common.Address, error) {
	req := BuildCallContractReqDynamic(nil, factoryAddress, v3.FactoryAbi, "getPool", token0Address, token1Address, fee)

	bytes, err := c.CallContract(req)
	if err != nil {
		return types.ZeroAddress, err
	}

	if len(bytes) == 0 {
		return types.ZeroAddress, ErrOutputEmpty
	}

	values, unpackErr := PancakeV3FactoryUnpacker.Unpack("getPool", bytes, 1)
	if unpackErr != nil {
		return types.ZeroAddress, unpackErr
	}

	if len(values) != 1 {
		return types.ZeroAddress, ErrWrongOutputLength
	}

	return ParseAddress(values[0])
}

/*
CallIsPool
for aerodrome
*/
func (c *ContractCaller) CallIsPool(poolAddress *common.Address) (bool, error) {
	req := BuildCallContractReqDynamic(nil, &aerodrome.FactoryAddress, aerodrome.FactoryAbi, "isPool", poolAddress)

	bytes, err := c.CallContract(req)
	if err != nil {
		return false, err
	}

	if len(bytes) == 0 {
		return false, ErrOutputEmpty
	}

	values, unpackErr := AerodromeV2FactoryUnpacker.Unpack("isPool", bytes, 1)
	if unpackErr != nil {
		return false, unpackErr
	}

	if len(values) != 1 {
		return false, ErrWrongOutputLength
	}

	return ParseBool(values[0])
}

/*
callGetReserves
for uniswap/pancake v2
*/
func (c *ContractCaller) callGetReserves(blockNumber *big.Int) ([]interface{}, error) {
	req := BuildCallContractReqDynamic(blockNumber, &types.WETHUSDCPairAddressUniswapV2, v2.PairAbi, "getReserves")

	bytes, err := c.CallContract(req)
	if err != nil {
		return nil, err
	}

	if len(bytes) == 0 {
		return nil, ErrOutputEmpty
	}

	values, unpackErr := PancakeV2Unpacker.Unpack("getReserves", bytes, 3)
	if unpackErr != nil {
		return nil, unpackErr
	}

	if len(values) != 3 {
		return nil, ErrWrongOutputLength
	}

	return values, nil
}

func (c *ContractCaller) GetReservesByBlockNumber(blockNumber *big.Int) (*big.Int, *big.Int, error) {
	values, err := c.callGetReserves(blockNumber)
	if err != nil {
		return nil, nil, err
	}

	reserve0, ok0 := values[0].(*big.Int)
	if !ok0 {
		return nil, nil, ErrReserve0NotBigInt
	}

	reserve1, ok1 := values[1].(*big.Int)
	if !ok1 {
		return nil, nil, ErrReserve1NotBigInt
	}

	return reserve0, reserve1, nil
}
