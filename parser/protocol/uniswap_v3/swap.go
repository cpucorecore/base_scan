package uniswap_v3

import (
	"base_scan/log"
	"base_scan/parser/protocol"
	"base_scan/parser/protocol/event_input_parser"
	"base_scan/types"
	"base_scan/types/orm"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/shopspring/decimal"
	"go.uber.org/zap"
	"math/big"
)

type swapEvent struct {
	*types.EventCommon
	Amount0Wei *big.Int
	Amount1Wei *big.Int
}

func (e *swapEvent) GetProtocolId() int {
	return protocolId
}

func (e *swapEvent) CanGetTx() bool {
	return true
}

func (e *swapEvent) GetTx(bnbPrice decimal.Decimal) *orm.Tx {
	tx := &orm.Tx{
		TxHash:        e.TxHash.String(),
		Maker:         e.Maker.String(),
		Token0Address: e.Pair.Token0Core.Address.String(),
		Token1Address: e.Pair.Token1Core.Address.String(),
		Block:         e.BlockNumber,
		BlockAt:       e.BlockTime,
		BlockIndex:    e.TxIndex,
		TxIndex:       e.LogIndex,
		PairAddress:   e.Pair.Address.String(),
		Program:       program,
	}

	tx.Token0Amount, tx.Token1Amount = protocol.ParseAmountsByPair(e.Amount0Wei, e.Amount1Wei, e.Pair)
	if tx.Token0Amount.IsNegative() {
		tx.Event = protocol.EventNameBuy
		tx.Token0Amount = tx.Token0Amount.Neg()
	} else if tx.Token1Amount.IsNegative() {
		tx.Event = protocol.EventNameSell
		tx.Token1Amount = tx.Token1Amount.Neg()
	} else {
		log.Logger.Warn("wrong v3 swap event", zap.Any("event", e))
	}

	tx.AmountUsd, tx.PriceUsd = protocol.CalcAmountAndPrice(bnbPrice, tx.Token0Amount, tx.Token1Amount, e.Pair.Token1Core.Address)
	return tx
}

func (e *swapEvent) CanGetPoolUpdateParameter() bool {
	return true
}

func (e *swapEvent) GetPoolUpdateParameter() *types.PoolUpdateParameter {
	return &types.PoolUpdateParameter{
		BlockNumber:   e.BlockNumber,
		PairAddress:   e.Pair.Address,
		Token0Address: e.Pair.Token0Core.Address,
		Token1Address: e.Pair.Token1Core.Address,
	}
}

var _ types.Event = (*swapEvent)(nil)

type SwapEventParser struct {
	EventInputParser event_input_parser.EventInputParser
}

// Parse
/*
https://bscscan.com/tx/0x840c6c07588af9ed343f993a1d1ade258c11d25f8ece721fc8546295a0d43946#eventlog
Swap (index_topic_1 address sender, index_topic_2 address recipient, int256 amount0, int256 amount1, uint160 sqrtPriceX96, uint128 liquidity, int24 tick, uint128 protocolFeesToken0, uint128 protocolFeesToken1)
*/
func (o *SwapEventParser) Parse(receiptLog *ethtypes.Log) (types.Event, error) {
	input, err := o.EventInputParser.Parse(receiptLog)
	if err != nil {
		return nil, err
	}

	e := &swapEvent{
		EventCommon: types.EventCommonFromEthLog(receiptLog),
		Amount0Wei:  input[0].(*big.Int),
		Amount1Wei:  input[1].(*big.Int),
	}

	e.Pair = &types.Pair{
		Address:    receiptLog.Address,
		ProtocolId: protocolId,
	}

	return e, nil
}
