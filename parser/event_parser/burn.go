package event_parser

import (
	"base_scan/parser/event_parser/event"
	"base_scan/types"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"math/big"
)

type BurnEventParser struct {
	PoolEventParser
}

func (o *BurnEventParser) Parse(receiptLog *ethtypes.Log) (types.Event, error) {
	input, err := o.EventInputParser.Parse(receiptLog)
	if err != nil {
		return nil, err
	}

	e := &event.BurnEvent{
		EventCommon: types.EventCommonFromEthLog(receiptLog),
		Amount0Wei:  input[0].(*big.Int),
		Amount1Wei:  input[1].(*big.Int),
	}

	e.Pair = &types.Pair{
		Address: receiptLog.Address,
	}

	e.PossibleProtocolIds = o.PossibleProtocolIds

	return e, nil
}
