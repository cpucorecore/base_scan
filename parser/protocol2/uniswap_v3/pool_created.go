package uniswap_v3

import (
	"base_scan/parser/protocol2"
	"base_scan/parser/protocol2/event_input_parser"
	"base_scan/types"
	"github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
)

const (
	ErrCodeWrongFactory = iota + types.ProtocolIdUniswapV3*10000 + 1
	ErrCodeParseEventInput
)

type poolCreatedEvent struct {
	*types.EventCommon
	MintEvent types.Event
}

func (e *poolCreatedEvent) GetProtocolIds() int {
	return protocolId
}

func (e *poolCreatedEvent) CanGetPair() bool {
	return true
}

func (e *poolCreatedEvent) GetPair() *types.Pair {
	if e.MintEvent != nil {
		e.Pair.Token0InitAmount, e.Pair.Token1InitAmount = e.MintEvent.GetMintAmount()
	}
	e.Pair.BlockAt = e.BlockTime
	return e.Pair
}

func (e *poolCreatedEvent) IsCreatePair() bool {
	return true
}

func (e *poolCreatedEvent) LinkEvent(event types.Event) {
	e.MintEvent = event
}

var _ types.Event = &poolCreatedEvent{}

type PoolCreatedEventParser struct {
	EventInputParser event_input_parser.EventInputParser
	FactoryAddress   common.Address
}

// Parse
/*
https://bscscan.com/tx/0x66ca21086d70078b736c8a72d266808d2102d01e0277f885082bcfcdcf926e37#eventlog#18
PoolCreated (index_topic_1 address token0, index_topic_2 address token1, index_topic_3 uint24 fee, int24 tickSpacing, address pool)
*/
func (o *PoolCreatedEventParser) Parse(receiptLog *ethtypes.Log) (types.Event, error) {
	pair := &types.Pair{}

	if receiptLog.Address.Cmp(o.FactoryAddress) != 0 {
		pair.Filtered = true
		pair.FilterCode = ErrCodeWrongFactory
		return nil, protocol2.ErrWrongFactoryAddress
	}

	input, err := o.EventInputParser.Parse(receiptLog)
	if err != nil {
		pair.Filtered = true
		pair.FilterCode = ErrCodeParseEventInput
		return nil, err
	}

	e := &poolCreatedEvent{
		EventCommon: types.EventCommonFromEthLog(receiptLog),
	}

	pair.Address = input[1].(common.Address)
	pair.Token0Core = &types.TokenCore{
		Address: common.BytesToAddress(receiptLog.Topics[1].Bytes()[12:]),
	}
	pair.Token1Core = &types.TokenCore{
		Address: common.BytesToAddress(receiptLog.Topics[2].Bytes()[12:]),
	}
	pair.Block = receiptLog.BlockNumber
	pair.ProtocolId = protocolId

	pair.FilterByToken0AndToken1()

	e.Pair = pair

	return e, nil
}
