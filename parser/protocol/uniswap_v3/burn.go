package uniswap_v3

import (
	"base_scan/parser/protocol"
	"base_scan/parser/protocol/event_input_parser"
	"base_scan/types"
	"base_scan/types/orm"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/shopspring/decimal"
	"math/big"
)

type burnEvent struct {
	*types.EventCommon
	amount0Wei *big.Int
	amount1Wei *big.Int
}

func (e *burnEvent) GetProtocolId() int {
	return protocolId
}

func (e *burnEvent) CanGetTx() bool {
	return true
}

func (e *burnEvent) GetTx(bnbPrice decimal.Decimal) *orm.Tx {
	tx := &orm.Tx{
		TxHash:        e.TxHash.String(),
		Event:         protocol.EventNameRemove,
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

	tx.Token0Amount, tx.Token1Amount = protocol.ParseAmountsByPair(e.amount0Wei, e.amount1Wei, e.Pair)
	tx.AmountUsd, tx.PriceUsd = protocol.CalcAmountAndPrice(bnbPrice, tx.Token0Amount, tx.Token1Amount, e.Pair.Token1Core.Address)
	return tx
}

var _ types.Event = (*burnEvent)(nil)

type BurnEventParser struct {
	EventInputParser event_input_parser.EventInputParser
}

// Parse
/*
https://bscscan.com/tx/0xc08489e828a7032ea8ae0cdb4795fa3e1cf04e19eb97981ba2ef1556087be1d9#eventlog
Burn(index_topic_1 address owner, index_topic_2 int24 tickLower, index_topic_3 int24 tickUpper, uint128 amount, uint256 amount0, uint256 amount1)
*/
func (o *BurnEventParser) Parse(receiptLog *ethtypes.Log) (types.Event, error) {
	input, err := o.EventInputParser.Parse(receiptLog)
	if err != nil {
		return nil, err
	}

	e := &burnEvent{
		EventCommon: types.EventCommonFromEthLog(receiptLog),
		amount0Wei:  input[1].(*big.Int),
		amount1Wei:  input[2].(*big.Int),
	}

	e.Pair = &types.Pair{
		Address:    receiptLog.Address,
		ProtocolId: protocolId,
	}

	return e, nil
}
