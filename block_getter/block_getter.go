package block_getter

import (
	"base_scan/cache"
	"base_scan/config"
	"base_scan/log"
	"base_scan/metrics"
	"base_scan/sequencer"
	"base_scan/types"
	"context"
	"github.com/avast/retry-go/v4"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/panjf2000/ants/v2"
	"go.uber.org/zap"
	"math/big"
	"sync"
	"time"
)

type BlockGetter interface {
	Start()
	GetStartBlockNumber(startBlockNumber uint64) uint64
	StartDispatch(startBlockNumber uint64)
	Stop()
	GetBlockAsync(blockNumber uint64)
	Next() *types.BlockContext
}

type blockGetter struct {
	ethClient            *ethclient.Client
	wsEthClient          *ethclient.Client
	queue                chan uint64
	buffer               chan *types.BlockContext
	workPool             *ants.Pool
	cache                cache.BlockCache
	stopped              bool
	stoppedLock          sync.RWMutex
	blockHeaderChan      chan *ethtypes.Header
	blockGetterSequencer sequencer.BlockSequencer
	headerHeightLock     sync.RWMutex
	headerHeight         uint64
	retryParams          *config.RetryParams
}

func NewBlockGetter(ethClient *ethclient.Client,
	wsEthClient *ethclient.Client,
	cache cache.BlockCache,
	blockGetterSequencer sequencer.BlockSequencer,
	retryParams *config.RetryParams,
) BlockGetter {
	workPool, err := ants.NewPool(config.G.BlockGetter.PoolSize)
	if err != nil {
		log.Logger.Fatal("ants pool(BlockGetter) init err", zap.Error(err))
	}

	return &blockGetter{
		ethClient:            ethClient,
		wsEthClient:          wsEthClient,
		queue:                make(chan uint64, config.G.BlockGetter.QueueSize),
		buffer:               make(chan *types.BlockContext, 10),
		workPool:             workPool,
		cache:                cache,
		blockHeaderChan:      make(chan *ethtypes.Header, 100),
		blockGetterSequencer: blockGetterSequencer,
		retryParams:          retryParams,
	}
}

func (bg *blockGetter) getBlockSync(blockNumber uint64) (*types.BlockContext, error) {
	now := time.Now()
	block, getBlockErr := bg.ethClient.BlockByNumber(context.Background(), big.NewInt(int64(blockNumber)))
	if getBlockErr != nil {
		return nil, getBlockErr
	}
	duration := time.Since(now)
	metrics.GetBlockDuration.Observe(duration.Seconds())

	now = time.Now()
	blockReceipts, getBlockReceiptErr := bg.ethClient.BlockReceipts(context.Background(), rpc.BlockNumberOrHashWithNumber(rpc.BlockNumber(blockNumber)))
	if getBlockReceiptErr != nil {
		return nil, getBlockReceiptErr
	}
	duration = time.Since(now)
	metrics.GetBlockReceiptsDuration.Observe(duration.Seconds())
	metrics.BlockDelay.Observe(time.Now().Sub(time.Unix((int64)(block.Time()), 0)).Seconds())
	return &types.BlockContext{
		Block:            block,
		BlockReceipts:    blockReceipts,
		HeightTime:       types.GetBlockHeightTime(block.Header()),
		TxIndex2TxSender: make(map[uint]common.Address, 200),
	}, nil
}

func (bg *blockGetter) getBlockSyncWithRetry(blockNumber uint64) (*types.BlockContext, error) {
	return retry.DoWithData(func() (*types.BlockContext, error) {
		return bg.getBlockSync(blockNumber)
	}, bg.retryParams.Attempts, bg.retryParams.Delay)
}

func (bg *blockGetter) GetBlockAsync(blockNumber uint64) {
	bg.queue <- blockNumber
}

func (bg *blockGetter) Next() *types.BlockContext {
	return <-bg.buffer
}

func (bg *blockGetter) latestBlockNumber() (uint64, error) {
	blockNumber, err := bg.ethClient.BlockNumber(context.Background())
	if err != nil {
		log.Logger.Error("block number err", zap.Error(err))
		return 0, err
	}
	return blockNumber, nil
}

func (bg *blockGetter) latestBlockNumberWithRetry() (uint64, error) {
	return retry.DoWithData(func() (uint64, error) {
		return bg.latestBlockNumber()
	}, bg.retryParams.Attempts, bg.retryParams.Delay)
}

func (bg *blockGetter) Start() {
	go func() {
		wg := &sync.WaitGroup{}
	tagFor:
		for {
			select {
			case blockNumber, ok := <-bg.queue:
				if !ok {
					log.Logger.Info("block queue is closed")
					break tagFor
				}
				wg.Add(1)
				bg.workPool.Submit(func() {
					defer wg.Done()
					log.Logger.Info("get block start", zap.Uint64("block_number", blockNumber))
					bw, err := bg.getBlockSyncWithRetry(blockNumber)
					if err != nil {
						log.Logger.Error("get block err", zap.Uint64("blockNumber", blockNumber), zap.Error(err))
						return
					}
					log.Logger.Info("get block success", zap.Uint64("blockNumber", blockNumber))
					bg.blockGetterSequencer.Commit(bw, bg.buffer)
				})
			}
		}

		taskNumber := bg.workPool.Waiting()
		log.Logger.Debug("wait block getter task finish", zap.Int("taskNumber", taskNumber))
		wg.Wait()
		log.Logger.Debug("all block getter task finish")
		close(bg.buffer)
	}()
}

func (bg *blockGetter) GetStartBlockNumber(startBlockNumber uint64) uint64 {
	newestBlockNumber, err := bg.ethClient.BlockNumber(context.Background())
	if err != nil {
		log.Logger.Fatal("ethClient.HeightBigInt() err", zap.Error(err))
	}

	if startBlockNumber == 0 {
		startBlockNumber = bg.cache.GetFinishedBlock()
	}

	if startBlockNumber == 0 {
		startBlockNumber = newestBlockNumber
	}

	return startBlockNumber
}

func (bg *blockGetter) setHeaderHeight(headerHeight uint64) {
	bg.headerHeightLock.Lock()
	defer bg.headerHeightLock.Unlock()
	if headerHeight > bg.headerHeight {
		bg.headerHeight = headerHeight
	}
}

func (bg *blockGetter) getHeaderHeight() uint64 {
	bg.headerHeightLock.RLock()
	defer bg.headerHeightLock.RUnlock()
	return bg.headerHeight
}

func (bg *blockGetter) subscribeNewHead() (ethereum.Subscription, <-chan error, error) {
	sub, err := bg.wsEthClient.SubscribeNewHead(context.Background(), bg.blockHeaderChan)
	if err != nil {
		return nil, nil, err
	}
	return sub, sub.Err(), nil
}

func (bg *blockGetter) startSubscribeNewHead() {
	headerHeight, err := bg.ethClient.BlockNumber(context.Background())
	if err != nil {
		log.Logger.Fatal("HeightBigInt() err", zap.Error(err))
	}
	bg.setHeaderHeight(headerHeight)

	sub, subErrChan, subErr := bg.subscribeNewHead()
	if subErr != nil {
		log.Logger.Fatal("subscribeNewHead() err", zap.Error(subErr))
	}

	go func() {
		for {
			select {
			case err = <-subErrChan:
				log.Logger.Error("receive block err", zap.Error(err))
				sub.Unsubscribe()
				for {
					sub, subErrChan, subErr = bg.subscribeNewHead()
					if subErr != nil {
						log.Logger.Error("subscribeNewHead() err", zap.Error(subErr))
						time.Sleep(time.Second * 1)
						continue
					}
					log.Logger.Info("subscribeNewHead() success")
					break
				}

			case blockHeader := <-bg.blockHeaderChan:
				log.Logger.Info("receive block header", zap.Any("height", blockHeader.Number))
				headerHeight = blockHeader.Number.Uint64()
				metrics.NewestHeight.Set(float64(headerHeight))
				bg.setHeaderHeight(headerHeight)
			}
		}
	}()
}

func (bg *blockGetter) dispatchRange(from, to uint64) (stopped bool, nextBlock uint64) {
	for i := from; i <= to; i++ {
		if bg.isStopped() {
			return true, i
		}
		bg.GetBlockAsync(i)
	}
	return false, 0
}

func (bg *blockGetter) StartDispatch(startBlockNumber uint64) {
	bg.startSubscribeNewHead()

	go func() {
		cur := startBlockNumber
		for {
			headerHeight := bg.getHeaderHeight()
			if headerHeight < cur {
				time.Sleep(100 * time.Millisecond)
				continue
			}

			stopped, nextBlockHeight := bg.dispatchRange(cur, headerHeight)
			if stopped {
				log.Logger.Info("dispatch interrupted", zap.Uint64("nextBlockHeight", nextBlockHeight))
				bg.doStop()
				return
			}

			cur = headerHeight + 1
		}
	}()
}

func (bg *blockGetter) Stop() {
	bg.stoppedLock.Lock()
	defer bg.stoppedLock.Unlock()
	bg.stopped = true
}

func (bg *blockGetter) isStopped() bool {
	bg.stoppedLock.RLock()
	defer bg.stoppedLock.RUnlock()
	return bg.stopped
}

func (bg *blockGetter) doStop() {
	close(bg.queue)
}
