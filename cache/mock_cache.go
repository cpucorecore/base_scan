package cache

import (
	"base_scan/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/shopspring/decimal"
	"math/big"
)

type MockCache struct{}

func (m MockCache) DelToken(address common.Address) {
}

func (m MockCache) DelPair(address common.Address) {
}

func (m MockCache) SetPrice(blockNumber *big.Int, price decimal.Decimal) {
}

func (m MockCache) GetPrice(blockNumber *big.Int) (decimal.Decimal, bool) {
	return decimal.Decimal{}, false
}

func (m MockCache) SetToken(token *types.Token) {
}

func (m MockCache) GetToken(address common.Address) (*types.Token, bool) {
	return nil, false
}

func (m MockCache) SetPair(pair *types.Pair) {
}

func (m MockCache) GetPair(address common.Address) (*types.Pair, bool) {
	return nil, false
}

func (m MockCache) PairExist(address common.Address) bool {
	return false
}

func (m MockCache) SetFinishedBlock(blockNumber uint64) {
}

func (m MockCache) GetFinishedBlock() uint64 {
	return 0
}

var _ Cache = MockCache{}
