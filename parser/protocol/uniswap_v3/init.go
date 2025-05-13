package uniswap_v3

import (
	"base_scan/abi/uniswap/v3"
	"base_scan/parser/protocol/event_input_parser"
)

var (
	EventParserPoolCreated *PoolCreatedEventParser
	EventParserMint        *MintEventParser
	EventParserBurn        *BurnEventParser
	EventParserSwap        *SwapEventParser
)

func init() {
	EventParserPoolCreated = &PoolCreatedEventParser{
		EventInputParser: event_input_parser.EventInputParser{
			Topic0:        v3.PoolCreatedTopic0,
			TopicLen:      4,
			DataUnpackLen: 2,
			AbiEvent:      v3.PoolCreatedEvent,
		},
		FactoryAddress: v3.FactoryAddress,
	}

	EventParserMint = &MintEventParser{
		EventInputParser: event_input_parser.EventInputParser{
			Topic0:        v3.MintTopic0,
			TopicLen:      4,
			DataUnpackLen: 4,
			AbiEvent:      v3.MintEvent,
		},
	}

	EventParserBurn = &BurnEventParser{
		EventInputParser: event_input_parser.EventInputParser{
			Topic0:        v3.BurnTopic0,
			TopicLen:      4,
			DataUnpackLen: 3,
			AbiEvent:      v3.BurnEvent,
		},
	}

	EventParserSwap = &SwapEventParser{
		EventInputParser: event_input_parser.EventInputParser{
			Topic0:        v3.SwapTopic0,
			TopicLen:      3,
			DataUnpackLen: 5,
			AbiEvent:      v3.SwapEvent,
		},
	}
}
