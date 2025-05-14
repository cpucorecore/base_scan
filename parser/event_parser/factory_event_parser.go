package event_parser

import (
	"base_scan/parser/event_parser/event_input_parser"
	"github.com/ethereum/go-ethereum/common"
)

type FactoryEventParser struct {
	Topic                    common.Hash
	PossibleFactoryAddresses map[common.Address]struct{}
	EventInputParser         event_input_parser.EventInputParser
}
