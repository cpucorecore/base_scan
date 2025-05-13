package parser

import (
	"base_scan/parser/protocol/event_input_parser"
	"base_scan/parser/protocol/uniswap_v2"
	"base_scan/parser/protocol/uniswap_v3"
	"base_scan/types"
	"errors"
	"github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
)

var (
	ErrParserNotFound = errors.New("parser not found")
)

type ReceiptLogParser interface {
	Parse(receiptLog *ethtypes.Log) (types.Event, error)
}

type receiptLogParser struct {
	eventParsers map[common.Hash]event_input_parser.EventParser
}

func NewReceiptLogParser() ReceiptLogParser {
	parser := &receiptLogParser{
		eventParsers: make(map[common.Hash]event_input_parser.EventParser),
	}
	return parser.RegisterParsers()
}

func (p *receiptLogParser) RegisterV2Parsers() {
	p.eventParsers[uniswap_v2.EventParserPairCreated.EventInputParser.Topic0] = uniswap_v2.EventParserPairCreated
	p.eventParsers[uniswap_v2.EventParserMint.EventInputParser.Topic0] = uniswap_v2.EventParserMint
	p.eventParsers[uniswap_v2.EventParserBurn.EventInputParser.Topic0] = uniswap_v2.EventParserBurn
	p.eventParsers[uniswap_v2.EventParserSwap.EventInputParser.Topic0] = uniswap_v2.EventParserSwap
	p.eventParsers[uniswap_v2.EventParserSync.EventInputParser.Topic0] = uniswap_v2.EventParserSync
}

func (p *receiptLogParser) RegisterV3Parsers() {
	p.eventParsers[uniswap_v3.EventParserPoolCreated.EventInputParser.Topic0] = uniswap_v3.EventParserPoolCreated
	p.eventParsers[uniswap_v3.EventParserMint.EventInputParser.Topic0] = uniswap_v3.EventParserMint
	p.eventParsers[uniswap_v3.EventParserBurn.EventInputParser.Topic0] = uniswap_v3.EventParserBurn
	p.eventParsers[uniswap_v3.EventParserSwap.EventInputParser.Topic0] = uniswap_v3.EventParserSwap
}

func (p *receiptLogParser) RegisterParsers() *receiptLogParser {
	p.RegisterV2Parsers()
	p.RegisterV3Parsers()
	return p
}

func (p *receiptLogParser) Parse(receiptLog *ethtypes.Log) (types.Event, error) {
	eventParser, ok := p.eventParsers[receiptLog.Topics[0]]
	if !ok {
		return nil, ErrParserNotFound
	}

	return eventParser.Parse(receiptLog)
}
