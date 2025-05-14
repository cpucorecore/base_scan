package event_parser

import (
	"base_scan/parser/event_parser/event_input_parser"
	"github.com/ethereum/go-ethereum/common"
)

type PoolEventParser struct {
	Topic               common.Hash
	PossibleProtocolIds []int
	EventInputParser    event_input_parser.EventInputParser
}
