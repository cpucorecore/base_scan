package uniswap_v2

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
	Amount0InWei  *big.Int
	Amount1InWei  *big.Int
	Amount0OutWei *big.Int
	Amount1OutWei *big.Int
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

	if e.Amount0InWei.Cmp(types.ZeroBigInt) > 0 {
		tx.Token0Amount, tx.Token1Amount = protocol.ParseAmountsByPair(e.Amount0InWei, e.Amount1OutWei, e.Pair)
		if !e.Pair.TokensReversed {
			tx.Event = protocol.EventNameSell
		} else {
			tx.Event = protocol.EventNameBuy
		}
	} else if e.Amount1InWei.Cmp(types.ZeroBigInt) > 0 {
		tx.Token0Amount, tx.Token1Amount = protocol.ParseAmountsByPair(e.Amount0OutWei, e.Amount1InWei, e.Pair)
		if !e.Pair.TokensReversed {
			tx.Event = protocol.EventNameBuy
		} else {
			tx.Event = protocol.EventNameSell
		}
	} else {
		log.Logger.Warn("wrong v2 swap event", zap.Any("event", e))
	}

	tx.AmountUsd, tx.PriceUsd = protocol.CalcAmountAndPrice(bnbPrice, tx.Token0Amount, tx.Token1Amount, e.Pair.Token1Core.Address)
	return tx
}

var _ types.Event = (*swapEvent)(nil)

type SwapEventParser struct {
	EventInputParser event_input_parser.EventInputParser
}

// Parse
/*
https://bscscan.com/tx/0xb306341fa6d7f7ff6b59be7be12881eb1d3ce199f25a1ad20a5e70dd2048e586#eventlog
Swap (index_topic_1 address sender, uint256 amount0In, uint256 amount1In, uint256 amount0Out, uint256 amount1Out, index_topic_2 address to)
*/
func (o *SwapEventParser) Parse(receiptLog *ethtypes.Log) (types.Event, error) {
	eventInput, err := o.EventInputParser.Parse(receiptLog)
	if err != nil {
		return nil, err
	}

	e := &swapEvent{
		EventCommon:   types.EventCommonFromEthLog(receiptLog),
		Amount0InWei:  eventInput[0].(*big.Int),
		Amount1InWei:  eventInput[1].(*big.Int),
		Amount0OutWei: eventInput[2].(*big.Int),
		Amount1OutWei: eventInput[3].(*big.Int),
	}

	e.Pair = &types.Pair{
		Address:    receiptLog.Address,
		ProtocolId: protocolId,
	}

	return e, nil
}
