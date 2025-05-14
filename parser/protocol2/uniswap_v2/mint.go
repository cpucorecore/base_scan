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

type mintEvent struct {
	*types.EventCommon
	amount0Wei *big.Int
	Amount1Wei *big.Int
}

func (e *mintEvent) GetMintAmount() (decimal.Decimal, decimal.Decimal) {
	return protocol2.ParseAmountsByPair(e.amount0Wei, e.Amount1Wei, e.Pair)
}

func (e *mintEvent) GetProtocolIds() int {
	return protocolId
}

func (e *mintEvent) CanGetTx() bool {
	return true
}

func (e *mintEvent) GetTx(bnbPrice decimal.Decimal) *orm.Tx {
	tx := &orm.Tx{
		TxHash:        e.TxHash.String(),
		Event:         protocol2.EventNameAdd,
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

	tx.Token0Amount, tx.Token1Amount = protocol2.ParseAmountsByPair(e.amount0Wei, e.Amount1Wei, e.Pair)
	tx.AmountUsd, tx.PriceUsd = protocol2.CalcAmountAndPrice(bnbPrice, tx.Token0Amount, tx.Token1Amount, e.Pair.Token1Core.Address)
	return tx
}

func (e *mintEvent) IsMint() bool {
	return true
}

var _ types.Event = (*mintEvent)(nil)

type MintEventParser struct {
	EventInputParser event_input_parser.EventInputParser
}

// Parse
/*
https://bscscan.com/tx/0xabb1ee98091af3ece13ce2b90f6e016fc8a36e17442700ee2c46c49a8d9a8a20#eventlog
Mint (index_topic_1 address sender, uint256 amount0, uint256 amount1)
*/
func (o *MintEventParser) Parse(receiptLog *ethtypes.Log) (types.Event, error) {
	input, err := o.EventInputParser.Parse(receiptLog)
	if err != nil {
		return nil, err
	}

	e := &mintEvent{
		EventCommon: types.EventCommonFromEthLog(receiptLog),
		amount0Wei:  input[0].(*big.Int),
		Amount1Wei:  input[1].(*big.Int),
	}

	e.Pair = &types.Pair{
		Address:    receiptLog.Address,
		ProtocolId: protocolId,
	}

	return e, nil
}
