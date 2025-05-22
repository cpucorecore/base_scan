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
	Next() *types.ParseBlockContext
}

type blockGetter struct {
	ctx             context.Context
	ethClient       *ethclient.Client
	wsEthClient     *ethclient.Client
	inputQueue      chan uint64
	outputBuffer    chan *types.ParseBlockContext
	workPool        *ants.Pool
	cache           cache.BlockCache
	stopped         SafeVar[bool]
	blockHeaderChan chan *ethtypes.Header
	blockSequencer  sequencer.BlockSequencer
	headerHeight    SafeVar[uint64]
	retryParams     *config.RetryParams
}

func NewBlockGetter(ethClient *ethclient.Client,
	wsEthClient *ethclient.Client,
	cache cache.BlockCache,
	blockSequencer sequencer.BlockSequencer,
	retryParams *config.RetryParams,
) BlockGetter {
	workPool, err := ants.NewPool(config.G.BlockGetter.PoolSize)
	if err != nil {
		log.Logger.Fatal("ants pool(BlockGetter) init err", zap.Error(err))
	}

	return &blockGetter{
		ctx:             context.Background(),
		ethClient:       ethClient,
		wsEthClient:     wsEthClient,
		inputQueue:      make(chan uint64, config.G.BlockGetter.QueueSize),
		outputBuffer:    make(chan *types.ParseBlockContext, 10),
		workPool:        workPool,
		cache:           cache,
		blockHeaderChan: make(chan *ethtypes.Header, 100),
		blockSequencer:  blockSequencer,
		retryParams:     retryParams,
	}
}

func (bg *blockGetter) getBlock(blockNumber uint64) (*types.ParseBlockContext, error) {
	var (
		block          *ethtypes.Block
		blockReceipts  []*ethtypes.Receipt
		getBlockErr    error
		getReceiptsErr error
		wg             sync.WaitGroup
	)

	wg.Add(2)
	go func() {
		defer wg.Done()
		now := time.Now()
		block, getBlockErr = bg.ethClient.BlockByNumber(bg.ctx, big.NewInt(int64(blockNumber)))
		if getBlockErr == nil {
			duration := time.Since(now)
			metrics.GetBlockDurationMs.Observe(float64(duration.Milliseconds()))
		}
	}()

	go func() {
		defer wg.Done()
		now := time.Now()
		blockReceipts, getReceiptsErr = bg.ethClient.BlockReceipts(bg.ctx, rpc.BlockNumberOrHashWithNumber(rpc.BlockNumber(blockNumber)))
		if getReceiptsErr == nil {
			duration := time.Since(now)
			metrics.GetBlockReceiptsDurationMs.Observe(float64(duration.Milliseconds()))
		}
	}()
	wg.Wait()

	if getBlockErr != nil {
		return nil, getBlockErr
	}
	if getReceiptsErr != nil {
		return nil, getReceiptsErr
	}

	metrics.BlockDelay.Observe(time.Now().Sub(time.Unix((int64)(block.Time()), 0)).Seconds())

	return &types.ParseBlockContext{
		Block:            block,
		BlockReceipts:    blockReceipts,
		HeightTime:       types.GetBlockHeightTime(block.Header()),
		TxIndex2TxSender: make(map[uint]common.Address, block.Transactions().Len()),
	}, nil
}

func (bg *blockGetter) getBlockWithRetry(blockNumber uint64) (*types.ParseBlockContext, error) {
	return retry.DoWithData(func() (*types.ParseBlockContext, error) {
		return bg.getBlock(blockNumber)
	}, bg.retryParams.Attempts, bg.retryParams.Delay)
}

func (bg *blockGetter) GetBlockAsync(blockNumber uint64) {
	bg.inputQueue <- blockNumber
}

func (bg *blockGetter) Next() *types.ParseBlockContext {
	return <-bg.outputBuffer
}

func (bg *blockGetter) Start() {
	go func() {
		wg := &sync.WaitGroup{}
	tagFor:
		for {
			select {
			case blockNumber, ok := <-bg.inputQueue:
				if !ok {
					log.Logger.Info("block inputQueue is closed")
					break tagFor
				}

				wg.Add(1)
				bg.workPool.Submit(func() {
					defer wg.Done()

					log.Logger.Info("get block start", zap.Uint64("block_number", blockNumber))
					bw, err := bg.getBlockWithRetry(blockNumber)
					if err != nil {
						log.Logger.Error("get block err", zap.Uint64("blockNumber", blockNumber), zap.Error(err))
						return
					}

					log.Logger.Info("get block success", zap.Uint64("blockNumber", blockNumber))
					bg.blockSequencer.Commit(bw, bg.outputBuffer)
				})
			}
		}

		wg.Wait()
		log.Logger.Info("all block getter task finish")
		close(bg.outputBuffer)
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
	if headerHeight > bg.headerHeight.Get() {
		bg.headerHeight.Set(headerHeight)
	}
}

func (bg *blockGetter) getHeaderHeight() uint64 {
	return bg.headerHeight.Get()
}

func (bg *blockGetter) subscribeNewHead() (ethereum.Subscription, <-chan error, error) {
	sub, err := bg.wsEthClient.SubscribeNewHead(bg.ctx, bg.blockHeaderChan)
	if err != nil {
		return nil, nil, err
	}
	return sub, sub.Err(), nil
}

func (bg *blockGetter) reconnectWithBackoff() (ethereum.Subscription, <-chan error) {
	retryDelay := time.Second * 1
	maxRetryDelay := time.Second * 10

	for {
		sub, errChan, err := bg.subscribeNewHead()
		if err == nil {
			log.Logger.Info("WebSocket reconnected successfully")
			return sub, errChan
		}

		log.Logger.Error("WebSocket reconnect failed",
			zap.Error(err),
			zap.Duration("nextRetry", retryDelay),
		)
		time.Sleep(retryDelay)

		retryDelay *= 2
		if retryDelay > maxRetryDelay {
			retryDelay = maxRetryDelay
		}
	}
}

func (bg *blockGetter) startSubscribeNewHead() {
	headerHeight, err := bg.ethClient.BlockNumber(bg.ctx)
	if err != nil {
		log.Logger.Fatal("HeightBigInt() err", zap.Error(err))
	}
	bg.setHeaderHeight(headerHeight)

	sub, errChan, err := bg.subscribeNewHead()
	if err != nil {
		log.Logger.Fatal("subscribeNewHead() err", zap.Error(err))
	}

	go func() {
		noBlockTimeout := time.NewTimer(10 * time.Second)
		defer noBlockTimeout.Stop()

		resetConnection := func() {
			noBlockTimeout.Stop()
			select {
			case <-noBlockTimeout.C:
			default:
			}
			sub.Unsubscribe()

			sub, errChan = bg.reconnectWithBackoff()
			noBlockTimeout.Reset(10 * time.Second)
		}

		for {
			select {
			case err = <-errChan:
				log.Logger.Error("WebSocket error", zap.Error(err))
				resetConnection()
			case blockHeader := <-bg.blockHeaderChan:
				height := blockHeader.Number.Uint64()
				log.Logger.Info("New block", zap.Uint64("height", height))
				bg.setHeaderHeight(height)
				metrics.NewestHeight.Set(float64(height))

				noBlockTimeout.Stop()
				select {
				case <-noBlockTimeout.C:
				default:
				}
				noBlockTimeout.Reset(10 * time.Second)
			case <-noBlockTimeout.C:
				log.Logger.Warn("No new blocks for 10s, reconnect WebSocket")
				resetConnection()
			}
		}
	}()
}

func (bg *blockGetter) startSubscribeNewHead2() {
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
	bg.stopped.Set(true)
}

func (bg *blockGetter) isStopped() bool {
	return bg.stopped.Get()
}

func (bg *blockGetter) doStop() {
	close(bg.inputQueue)
}
