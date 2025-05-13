package uniswap_v2

import (
	"base_scan/abi/uniswap/v2"
	"base_scan/parser/protocol/event_input_parser"
)

var (
	EventParserPairCreated *PairCreatedEventParser
	EventParserMint        *MintEventParser
	EventParserBurn        *BurnEventParser
	EventParserSwap        *SwapEventParser
	EventParserSync        *SyncEventParser
)

func init() {
	EventParserPairCreated = &PairCreatedEventParser{
		EventInputParser: event_input_parser.EventInputParser{
			Topic0:        v2.PairCreatedTopic0,
			TopicLen:      3,
			DataUnpackLen: 2,
			AbiEvent:      v2.PairCreatedEvent,
		},
		FactoryAddress: v2.FactoryAddress,
	}

	EventParserMint = &MintEventParser{
		EventInputParser: event_input_parser.EventInputParser{
			Topic0:        v2.MintTopic0,
			TopicLen:      2,
			DataUnpackLen: 2,
			AbiEvent:      v2.MintEvent,
		},
	}

	EventParserBurn = &BurnEventParser{
		EventInputParser: event_input_parser.EventInputParser{
			Topic0:        v2.BurnTopic0,
			TopicLen:      3,
			DataUnpackLen: 2,
			AbiEvent:      v2.BurnEvent,
		},
	}

	EventParserSwap = &SwapEventParser{
		EventInputParser: event_input_parser.EventInputParser{
			Topic0:        v2.SwapTopic0,
			TopicLen:      3,
			DataUnpackLen: 4,
			AbiEvent:      v2.SwapEvent,
		},
	}

	EventParserSync = &SyncEventParser{
		EventInputParser: event_input_parser.EventInputParser{
			Topic0:        v2.SyncTopic0,
			TopicLen:      1,
			DataUnpackLen: 2,
			AbiEvent:      v2.SyncEvent,
		},
	}
}
