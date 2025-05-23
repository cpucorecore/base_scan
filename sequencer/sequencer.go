package sequencer

import (
	"base_scan/config"
	"base_scan/log"
	"base_scan/types"
	"go.uber.org/zap"
	"sync"
)

type Sequenceable interface {
	GetSequence() uint64
}

type BlockSequencer interface {
	Init(height uint64)
	Commit(bc *types.ParseBlockContext, output chan *types.ParseBlockContext)
}

type blockSequencer struct {
	active   bool
	mu       sync.Mutex
	cond     *sync.Cond
	sequence uint64
}

func NewBlockSequencer() BlockSequencer {
	s := &blockSequencer{
		active: config.G.EnableSequencer,
	}
	s.cond = sync.NewCond(&s.mu)
	return s
}

func (s *blockSequencer) Init(sequence uint64) {
	log.Logger.Info("init block sequencer", zap.Uint64("sequence", sequence))
	if s.sequence == 0 {
		s.sequence = sequence - 1
	} else {
		log.Logger.Fatal("sequencer init err", zap.Uint64("sequence", sequence), zap.Uint64("old sequence", s.sequence))
	}
}

func (s *blockSequencer) Commit(blockContext *types.ParseBlockContext, outputChan chan *types.ParseBlockContext) {
	if !s.active {
		outputChan <- blockContext
		return
	}

	sequence := blockContext.GetSequence()

	s.mu.Lock()
	for s.sequence+1 != sequence {
		s.cond.Wait()
	}

	outputChan <- blockContext
	s.sequence = sequence

	s.cond.Broadcast()
	s.mu.Unlock()
}
