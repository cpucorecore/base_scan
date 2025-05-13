package uniswap_v2

import (
	"base_scan/parser/protocol"
	"base_scan/parser/protocol/event_input_parser"
	"base_scan/types"
	"github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"math/big"
)

type syncEvent struct {
	*types.EventCommon
	amount0Wei *big.Int
	amount1Wei *big.Int
}

func (e *syncEvent) GetProtocolId() int {
	return protocolId
}

func (e *syncEvent) CanGetPoolUpdate() bool {
	return true
}

func (e *syncEvent) GetPoolUpdate() *types.PoolUpdate {
	pu := &types.PoolUpdate{
		Program:       program,
		LogIndex:      e.LogIndex,
		Address:       e.ContractAddress,
		Token0Address: e.Pair.Token0Core.Address,
		Token1Address: e.Pair.Token1Core.Address,
	}

	pu.Token0Amount, pu.Token1Amount = protocol.ParseAmountsByPair(e.amount0Wei, e.amount1Wei, e.Pair)

	return pu
}

var _ types.Event = (*syncEvent)(nil)

type SyncEventParser struct {
	Topic0           common.Hash
	EventInputParser event_input_parser.EventInputParser
}

// Parse
/*
https://bscscan.com/tx/0xabb1ee98091af3ece13ce2b90f6e016fc8a36e17442700ee2c46c49a8d9a8a20#eventlog#7
Sync (uint112 reserve0, uint112 reserve1)
*/
func (o *SyncEventParser) Parse(receiptLog *ethtypes.Log) (types.Event, error) {
	input, err := o.EventInputParser.Parse(receiptLog)
	if err != nil {
		return nil, err
	}

	e := &syncEvent{
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
