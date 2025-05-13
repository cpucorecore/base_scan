package uniswap_v2

import (
	"base_scan/parser/protocol"
	"base_scan/parser/protocol/event_input_parser"
	"base_scan/types"
	"github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
)

const (
	ErrCodeWrongFactory = iota + types.ProtocolIdUniswapV2*10000 + 1
	ErrCodeParseEventInput
)

type pairCreatedEvent struct {
	*types.EventCommon
	MintEvent types.Event
}

func (e *pairCreatedEvent) GetProtocolId() int {
	return protocolId
}

func (e *pairCreatedEvent) CanGetPair() bool {
	return true
}

func (e *pairCreatedEvent) GetPair() *types.Pair {
	if e.MintEvent != nil {
		e.Pair.Token0InitAmount, e.Pair.Token1InitAmount = e.MintEvent.GetMintAmount()
	}
	e.Pair.BlockAt = e.BlockTime
	return e.Pair
}

func (e *pairCreatedEvent) IsCreatePair() bool {
	return true
}

func (e *pairCreatedEvent) LinkEvent(event types.Event) { // for pair initial token0/token1 amount
	e.MintEvent = event
}

var _ types.Event = (*pairCreatedEvent)(nil)

type PairCreatedEventParser struct {
	EventInputParser event_input_parser.EventInputParser
	FactoryAddress   common.Address
}

// Parse
/*
https://bscscan.com/tx/0xb566ecc18c7b854eaf5c868dcf1c6de1742e5c63e1739bc9aa161114d4bd9628#eventlog#137
PairCreated (index_topic_1 address token0, index_topic_2 address token1, address pair, uint256)
*/
func (o *PairCreatedEventParser) Parse(receiptLog *ethtypes.Log) (types.Event, error) {
	pair := &types.Pair{}

	if receiptLog.Address.Cmp(o.FactoryAddress) != 0 {
		pair.Filtered = true
		pair.FilterCode = ErrCodeWrongFactory
		return nil, protocol.ErrWrongFactoryAddress
	}

	input, err := o.EventInputParser.Parse(receiptLog)
	if err != nil {
		pair.Filtered = true
		pair.FilterCode = ErrCodeParseEventInput
		return nil, err
	}

	e := &pairCreatedEvent{
		EventCommon: types.EventCommonFromEthLog(receiptLog),
	}

	pair.Address = input[0].(common.Address)
	pair.Token0Core = &types.TokenCore{
		Address: common.BytesToAddress(receiptLog.Topics[1].Bytes()[12:]),
	}
	pair.Token1Core = &types.TokenCore{
		Address: common.BytesToAddress(receiptLog.Topics[2].Bytes()[12:]),
	}
	pair.Block = receiptLog.BlockNumber
	pair.ProtocolId = protocolId

	pair.FilterByToken0AndToken1()

	e.EventCommon.Pair = pair

	return e, nil
}
