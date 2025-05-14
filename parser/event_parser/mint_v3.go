package event_parser

import (
	"base_scan/parser/event_parser/event"
	"base_scan/types"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"math/big"
)

type MintEventParserV3 struct {
	PoolEventParser
}

func (o *MintEventParserV3) Parse(receiptLog *ethtypes.Log) (types.Event, error) {
	input, err := o.EventInputParser.Parse(receiptLog)
	if err != nil {
		return nil, err
	}

	e := &event.MintEvent{
		EventCommon: types.EventCommonFromEthLog(receiptLog),
		Amount0Wei:  input[2].(*big.Int),
		Amount1Wei:  input[3].(*big.Int),
	}

	e.Pair = &types.Pair{
		Address: receiptLog.Address,
	}

	e.PossibleProtocolIds = o.PossibleProtocolIds

	return e, nil
}
