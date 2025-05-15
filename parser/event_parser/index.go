package event_parser

import (
	"base_scan/abi"
	"base_scan/abi/aerodrome"
	uniswapv2 "base_scan/abi/uniswap/v2"
	uniswapv3 "base_scan/abi/uniswap/v3"
	"base_scan/parser/event_parser/event_input_parser"
	"github.com/ethereum/go-ethereum/common"
)

var (
	pairCreatedEventParser = &PairCreatedEventParser{
		FactoryEventParser: FactoryEventParser{
			Topic:                    uniswapv2.PairCreatedTopic0,
			PossibleFactoryAddresses: abi.Topic2FactoryAddresses[uniswapv2.PairCreatedTopic0],
			EventInputParser: event_input_parser.EventInputParser{
				TopicLen:      3,
				DataUnpackLen: 2,
				AbiEvent:      uniswapv2.PairCreatedEvent,
			},
		},
	}

	burnEventParser = &BurnEventParser{
		PoolEventParser: PoolEventParser{
			Topic:               uniswapv2.BurnTopic0,
			PossibleProtocolIds: abi.Topic2ProtocolIds[uniswapv2.BurnTopic0],
			EventInputParser: event_input_parser.EventInputParser{
				TopicLen:      3,
				DataUnpackLen: 2,
				AbiEvent:      uniswapv2.BurnEvent,
			},
		},
	}

	burnEventParserAerodrome = &BurnEventParser{
		PoolEventParser: PoolEventParser{
			Topic:               aerodrome.BurnTopic0,
			PossibleProtocolIds: abi.Topic2ProtocolIds[aerodrome.BurnTopic0],
			EventInputParser: event_input_parser.EventInputParser{
				TopicLen:      3,
				DataUnpackLen: 2,
				AbiEvent:      aerodrome.BurnEvent,
			},
		},
	}

	swapEventParser = &SwapEventParser{
		PoolEventParser: PoolEventParser{
			Topic:               uniswapv2.SwapTopic0,
			PossibleProtocolIds: abi.Topic2ProtocolIds[uniswapv2.SwapTopic0],
			EventInputParser: event_input_parser.EventInputParser{
				TopicLen:      3,
				DataUnpackLen: 5,
				AbiEvent:      uniswapv2.SwapEvent,
			},
		},
	}

	syncEventParser = &SyncEventParser{
		PoolEventParser: PoolEventParser{
			Topic:               uniswapv2.SyncTopic0,
			PossibleProtocolIds: abi.Topic2ProtocolIds[uniswapv2.SyncTopic0],
			EventInputParser: event_input_parser.EventInputParser{
				TopicLen:      1,
				DataUnpackLen: 2,
				AbiEvent:      uniswapv2.SyncEvent,
			},
		},
	}

	mintEventParser = &MintEventParser{
		PoolEventParser: PoolEventParser{
			Topic:               uniswapv2.MintTopic0,
			PossibleProtocolIds: abi.Topic2ProtocolIds[uniswapv2.MintTopic0],
			EventInputParser: event_input_parser.EventInputParser{
				TopicLen:      2,
				DataUnpackLen: 2,
				AbiEvent:      uniswapv2.MintEvent,
			},
		},
	}

	Topic2EventParser = map[common.Hash]EventParser{
		uniswapv2.PairCreatedTopic0: pairCreatedEventParser,
		uniswapv2.MintTopic0:        mintEventParser,
		uniswapv2.BurnTopic0:        burnEventParser,
		uniswapv2.SwapTopic0:        swapEventParser,
		uniswapv2.SyncTopic0:        syncEventParser,

		uniswapv3.PoolCreatedTopic0: &PoolCreatedEventParser{
			FactoryEventParser: FactoryEventParser{
				Topic:                    uniswapv3.PoolCreatedTopic0,
				PossibleFactoryAddresses: abi.Topic2FactoryAddresses[uniswapv3.PoolCreatedTopic0],
				EventInputParser: event_input_parser.EventInputParser{
					TopicLen:      4,
					DataUnpackLen: 2,
					AbiEvent:      uniswapv3.PoolCreatedEvent,
				},
			},
		},
		uniswapv3.MintTopic0: &MintEventParserV3{
			PoolEventParser: PoolEventParser{
				Topic:               uniswapv3.MintTopic0,
				PossibleProtocolIds: abi.Topic2ProtocolIds[uniswapv3.MintTopic0],
				EventInputParser: event_input_parser.EventInputParser{
					TopicLen:      4,
					DataUnpackLen: 4,
					AbiEvent:      uniswapv3.MintEvent,
				},
			},
		},
		uniswapv3.BurnTopic0: &BurnEventParser{
			PoolEventParser: PoolEventParser{
				Topic:               uniswapv3.BurnTopic0,
				PossibleProtocolIds: abi.Topic2ProtocolIds[uniswapv3.BurnTopic0],
				EventInputParser: event_input_parser.EventInputParser{
					TopicLen:      4,
					DataUnpackLen: 3,
					AbiEvent:      uniswapv3.BurnEvent,
				},
			},
		},
		uniswapv3.SwapTopic0: &SwapEventParser{
			PoolEventParser: PoolEventParser{
				Topic:               uniswapv3.SwapTopic0,
				PossibleProtocolIds: abi.Topic2ProtocolIds[uniswapv3.SwapTopic0],
				EventInputParser: event_input_parser.EventInputParser{
					TopicLen:      3,
					DataUnpackLen: 5,
					AbiEvent:      uniswapv3.SwapEvent,
				},
			},
		},

		aerodrome.PairCreatedTopic0: pairCreatedEventParser,
		aerodrome.BurnTopic0:        burnEventParserAerodrome,
		aerodrome.SwapTopic0:        swapEventParser,
		aerodrome.SyncTopic0:        syncEventParser,
	}
)
