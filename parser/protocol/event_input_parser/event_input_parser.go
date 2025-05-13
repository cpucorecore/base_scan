package event_input_parser

import (
	"errors"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
)

var (
	ErrWrongTopicLen      = errors.New("wrong topic length")
	ErrWrongDataUnpackLen = errors.New("wrong data unpack length")
)

type EventInputParser struct {
	Topic0        common.Hash
	TopicLen      int
	DataUnpackLen int
	AbiEvent      *abi.Event
}

func (p *EventInputParser) Parse(receiptLog *ethtypes.Log) ([]interface{}, error) {
	if len(receiptLog.Topics) != p.TopicLen {
		return nil, ErrWrongTopicLen
	}

	eventInput, err := p.AbiEvent.Inputs.Unpack(receiptLog.Data)
	if err != nil {
		return nil, err
	}

	if len(eventInput) != p.DataUnpackLen {
		return nil, ErrWrongDataUnpackLen
	}

	return eventInput, nil
}
