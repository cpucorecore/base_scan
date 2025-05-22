package sequencer

import (
	"base_scan/config"
	"base_scan/log"
	"base_scan/types"
	"go.uber.org/zap"
	"sync"
)

type BlockSequencer interface {
	Init(height uint64)
	Commit(bc *types.ParseBlockContext, output chan *types.ParseBlockContext)
}

type blockSequencer struct {
	active bool
	mu     sync.Mutex
	cond   *sync.Cond
	height uint64
}

func NewBlockSequencer() BlockSequencer {
	s := &blockSequencer{
		active: config.G.EnableSequencer,
	}
	s.cond = sync.NewCond(&s.mu)
	return s
}

func (s *blockSequencer) Init(startHeight uint64) {
	log.Logger.Info("init block sequencer", zap.Uint64("startHeight", startHeight))
	if s.height == 0 {
		s.height = startHeight - 1
	} else {
		log.Logger.Fatal("sequencer init err", zap.Uint64("startHeight", startHeight), zap.Uint64("old height", s.height))
	}
}

func (s *blockSequencer) Commit(blockContext *types.ParseBlockContext, outputChan chan *types.ParseBlockContext) {
	if !s.active {
		outputChan <- blockContext
		return
	}

	s.mu.Lock()
	for s.height+1 != blockContext.GetBlockNumber() {
		s.cond.Wait()
	}

	outputChan <- blockContext
	s.height = blockContext.GetBlockNumber()
	s.cond.Broadcast()
	s.mu.Unlock()
}
