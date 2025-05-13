package types

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/shopspring/decimal"
	"math/big"
)

var (
	ZeroBigInt  = big.NewInt(0)
	ZeroDecimal = decimal.NewFromInt(0)
	ZeroAddress = common.Address{}
)
