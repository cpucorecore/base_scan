package parser

import (
	"base_scan/abi"
	uniswapv2 "base_scan/abi/uniswap/v2"
	uniswapv3 "base_scan/abi/uniswap/v3"
	"base_scan/parser/event_parser"
	"base_scan/parser/event_parser/event_input_parser"
	"github.com/ethereum/go-ethereum/common"
)

var Topic2EventParser = map[common.Hash]event_parser.EventParser{
	uniswapv3.BurnTopic0: &event_parser.BurnEventParser{
		PoolEventParser: event_parser.PoolEventParser{
			Topic:               uniswapv3.BurnTopic0,
			PossibleProtocolIds: abi.Topic2ProtocolIds[uniswapv3.BurnTopic0],
			EventInputParser: event_input_parser.EventInputParser{
				TopicLen:      4,
				DataUnpackLen: 3,
				AbiEvent:      uniswapv3.BurnEvent,
			},
		},
	},
	uniswapv2.PairCreatedTopic0: &event_parser.PairCreatedEventParser{
		FactoryEventParser: event_parser.FactoryEventParser{
			Topic:                    uniswapv2.PairCreatedTopic0,
			PossibleFactoryAddresses: abi.Topic2FactoryAddresses[uniswapv2.PairCreatedTopic0],
			EventInputParser: event_input_parser.EventInputParser{
				TopicLen:      3,
				DataUnpackLen: 2,
				AbiEvent:      uniswapv2.PairCreatedEvent,
			},
		},
	},
	uniswapv3.SwapTopic0: &event_parser.SwapEventParser{
		PoolEventParser: event_parser.PoolEventParser{
			Topic:               uniswapv3.SwapTopic0,
			PossibleProtocolIds: abi.Topic2ProtocolIds[uniswapv3.SwapTopic0],
			EventInputParser: event_input_parser.EventInputParser{
				TopicLen:      3,
				DataUnpackLen: 5,
				AbiEvent:      uniswapv3.SwapEvent,
			},
		},
	},
	uniswapv2.SyncTopic0: &event_parser.SyncEventParser{
		PoolEventParser: event_parser.PoolEventParser{
			Topic:               uniswapv2.SyncTopic0,
			PossibleProtocolIds: abi.Topic2ProtocolIds[uniswapv2.SyncTopic0],
			EventInputParser: event_input_parser.EventInputParser{
				TopicLen:      1,
				DataUnpackLen: 2,
				AbiEvent:      uniswapv2.SyncEvent,
			},
		},
	},
}
