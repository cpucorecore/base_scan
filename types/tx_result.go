package types

import (
	"base_scan/log"
	"github.com/ethereum/go-ethereum/common"
	"go.uber.org/zap"
)

type TxPairEvent struct {
	V2 []Event
	V3 []Event
}

func (tpe *TxPairEvent) AddEvent(event Event) {
	switch event.GetProtocolId() {
	case ProtocolIdUniswapV2:
		if tpe.V2 == nil {
			tpe.V2 = make([]Event, 0, 10)
		}
		tpe.V2 = append(tpe.V2, event)
	case ProtocolIdUniswapV3:
		if tpe.V3 == nil {
			tpe.V3 = make([]Event, 0, 10)
		}
		tpe.V3 = append(tpe.V3, event)
	}
}

func (tpe *TxPairEvent) LinkEvents() {
	tpe.LinkEventsPancakeV2()
	tpe.LinkEventsV3()
}

func LinkPairCreatedEventAndMintEvent(pairCreatedEvents, mintEvents []Event) {
	mintEventsLen := len(mintEvents)
	for i, pairCreatedEvent := range pairCreatedEvents {
		if i < mintEventsLen {
			pairCreatedEvent.LinkEvent(mintEvents[i])
		} else {
			log.Logger.Warn("pair have no related mint event", zap.Any("pairCreatedEvent", pairCreatedEvent))
		}
	}
}

func (tpe *TxPairEvent) LinkEventsPancakeV2() {
	mintEvents := make([]Event, 0)
	pairCreatedEvents := make([]Event, 0)
	for _, v2Event := range tpe.V2 {
		if v2Event.IsMint() {
			mintEvents = append(mintEvents, v2Event)
		} else if v2Event.IsCreatePair() {
			pairCreatedEvents = append(pairCreatedEvents, v2Event)
		}
	}
	LinkPairCreatedEventAndMintEvent(pairCreatedEvents, mintEvents)
}

func (tpe *TxPairEvent) LinkEventsV3() {
	mintEvents := make([]Event, 0)
	pairCreatedEvents := make([]Event, 0)
	for _, v3Event := range tpe.V3 {
		if v3Event.IsMint() {
			mintEvents = append(mintEvents, v3Event)
		} else if v3Event.IsCreatePair() {
			pairCreatedEvents = append(pairCreatedEvents, v3Event)
		}
	}
	LinkPairCreatedEventAndMintEvent(pairCreatedEvents, mintEvents)
}

type TxResult struct {
	Maker                   common.Address
	PairCreatedEvents       []Event
	PairAddress2TxPairEvent map[common.Address]*TxPairEvent
}

func NewTxResult(maker common.Address) *TxResult {
	return &TxResult{
		Maker:                   maker,
		PairCreatedEvents:       make([]Event, 0, 10),
		PairAddress2TxPairEvent: make(map[common.Address]*TxPairEvent),
	}
}

func (tr *TxResult) AddEvent(event Event) {
	event.SetMaker(tr.Maker)
	if event.IsCreatePair() {
		tr.PairCreatedEvents = append(tr.PairCreatedEvents, event)
	}

	pairAddress := event.GetPairAddress()
	txPairEvent, ok := tr.PairAddress2TxPairEvent[pairAddress]
	if ok {
		txPairEvent.AddEvent(event)
		return
	}

	txPairEvent = &TxPairEvent{}
	txPairEvent.AddEvent(event)
	tr.PairAddress2TxPairEvent[pairAddress] = txPairEvent
}

func (tr *TxResult) LinkEvents() {
	for _, pairEvent := range tr.PairAddress2TxPairEvent {
		pairEvent.LinkEvents()
	}
}
