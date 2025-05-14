package event_parser

import (
	"base_scan/parser/event_parser/event"
	"base_scan/types"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"math/big"
)

type SwapEventParser struct {
	PoolEventParser
}

func (o *SwapEventParser) Parse(receiptLog *ethtypes.Log) (types.Event, error) {
	eventInput, err := o.EventInputParser.Parse(receiptLog)
	if err != nil {
		return nil, err
	}

	e := &event.SwapEvent{
		EventCommon:   types.EventCommonFromEthLog(receiptLog),
		Amount0InWei:  eventInput[0].(*big.Int),
		Amount1InWei:  eventInput[1].(*big.Int),
		Amount0OutWei: eventInput[2].(*big.Int),
		Amount1OutWei: eventInput[3].(*big.Int),
	}

	e.Pair = &types.Pair{
		Address: receiptLog.Address,
	}

	e.PossibleProtocolIds = o.PossibleProtocolIds

	return e, nil
}
