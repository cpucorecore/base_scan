package event_parser

import (
	"base_scan/abi"
	common2 "base_scan/parser/event_parser/common"
	"base_scan/parser/event_parser/event"
	"base_scan/types"
	"github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
)

type PoolCreatedEventParser struct {
	FactoryEventParser
}

func (o *PoolCreatedEventParser) Parse(receiptLog *ethtypes.Log) (types.Event, error) {
	pair := &types.Pair{}

	_, ok := o.PossibleFactoryAddresses[receiptLog.Address]
	if !ok {
		pair.Filtered = true
		pair.FilterCode = types.FilterCodeWrongFactory
		return nil, common2.ErrWrongFactoryAddress
	}

	input, err := o.EventInputParser.Parse(receiptLog)
	if err != nil {
		pair.Filtered = true
		pair.FilterCode = types.ErrCodeParseEventInput
		return nil, err
	}

	e := &event.PairCreatedEvent{
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
	pair.ProtocolId = abi.FactoryAddress2ProtocolId[receiptLog.Address]

	pair.FilterByToken0AndToken1()

	e.Pair = pair

	return e, nil
}
