package event_parser

import (
	"base_scan/abi"
	"base_scan/parser/event_parser/common"
	"base_scan/parser/event_parser/event"
	"base_scan/types"
	ethcommon "github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
)

type PairCreatedEventParser struct {
	FactoryEventParser
}

func (o *PairCreatedEventParser) Parse(receiptLog *ethtypes.Log) (types.Event, error) {
	pair := &types.Pair{}

	_, ok := o.PossibleFactoryAddresses[receiptLog.Address]
	if !ok {
		pair.Filtered = true
		pair.FilterCode = types.FilterCodeWrongFactory
		return nil, common.ErrWrongFactoryAddress
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

	pair.Address = input[0].(ethcommon.Address)
	pair.Token0Core = &types.TokenCore{
		Address: ethcommon.BytesToAddress(receiptLog.Topics[1].Bytes()[12:]),
	}
	pair.Token1Core = &types.TokenCore{
		Address: ethcommon.BytesToAddress(receiptLog.Topics[2].Bytes()[12:]),
	}
	pair.Block = receiptLog.BlockNumber
	pair.ProtocolId = abi.FactoryAddress2ProtocolId[receiptLog.Address]

	pair.FilterByToken0AndToken1()

	e.EventCommon.Pair = pair

	return e, nil
}
