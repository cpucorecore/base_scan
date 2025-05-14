package uniswap_v2

import (
	"base_scan/parser/protocol2"
	"base_scan/parser/protocol2/event_input_parser"
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

func (e *burnEvent) GetProtocolIds() int {
	return protocolId
}

func (e *burnEvent) CanGetTx() bool {
	return true
}

func (e *burnEvent) GetTx(bnbPrice decimal.Decimal) *orm.Tx {
	tx := &orm.Tx{
		TxHash:        e.TxHash.String(),
		Event:         protocol2.EventNameRemove,
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

	tx.Token0Amount, tx.Token1Amount = protocol2.ParseAmountsByPair(e.amount0Wei, e.amount1Wei, e.Pair)
	tx.AmountUsd, tx.PriceUsd = protocol2.CalcAmountAndPrice(bnbPrice, tx.Token0Amount, tx.Token1Amount, e.Pair.Token1Core.Address)
	return tx
}

var _ types.Event = (*burnEvent)(nil)

type BurnEventParser struct {
	EventInputParser event_input_parser.EventInputParser
}

//Parse
/*
https://bscscan.com/tx/0x361179e597dd10ab47cd4e7f49353944f54ad8b33196f8bffc416542d5f092e5#eventlog
Burn (index_topic_1 address sender, uint256 amount0, uint256 amount1, index_topic_2 address to)
*/
func (o *BurnEventParser) Parse(receiptLog *ethtypes.Log) (types.Event, error) {
	input, err := o.EventInputParser.Parse(receiptLog)
	if err != nil {
		return nil, err
	}

	e := &burnEvent{
		EventCommon: types.EventCommonFromEthLog(receiptLog),
		amount0Wei:  input[0].(*big.Int),
		amount1Wei:  input[1].(*big.Int),
	}

	e.Pair = &types.Pair{
		Address:    receiptLog.Address,
		ProtocolId: protocolId,
	}

	return e, nil
}
