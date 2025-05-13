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

type mintEvent struct {
	*types.EventCommon
	amount0Wei *big.Int
	amount1Wei *big.Int
}

func (e *mintEvent) GetProtocolId() int {
	return protocolId
}

func (e *mintEvent) CanGetTx() bool {
	return true
}

func (e *mintEvent) GetTx(bnbPrice decimal.Decimal) *orm.Tx {
	tx := &orm.Tx{
		TxHash:        e.TxHash.String(),
		Event:         protocol.EventNameAdd,
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

func (e *mintEvent) IsMint() bool {
	return true
}

func (e *mintEvent) GetMintAmount() (decimal.Decimal, decimal.Decimal) {
	return protocol.ParseAmountsByPair(e.amount0Wei, e.amount1Wei, e.Pair)
}

var _ types.Event = (*mintEvent)(nil)

type MintEventParser struct {
	event_input_parser.EventInputParser
}

// Parse
/*
https://bscscan.com/tx/0x47e4cedd1a33fceab488df4fb48691b740fcb736413d68e6dcd1b1b53cf1e46f#eventlog#522
Mint (address sender, index_topic_1 address owner, index_topic_2 int24 tickLower, index_topic_3 int24 tickUpper, uint128 amount, uint256 amount0, uint256 amount1)
*/
func (o *MintEventParser) Parse(receiptLog *ethtypes.Log) (types.Event, error) {
	input, err := o.EventInputParser.Parse(receiptLog)
	if err != nil {
		return nil, err
	}

	e := &mintEvent{
		EventCommon: types.EventCommonFromEthLog(receiptLog),
		amount0Wei:  input[2].(*big.Int),
		amount1Wei:  input[3].(*big.Int),
	}

	e.Pair = &types.Pair{
		Address:    receiptLog.Address,
		ProtocolId: protocolId,
	}

	return e, nil
}
