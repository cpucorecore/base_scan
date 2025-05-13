package types

import (
	chainparams "base_scan/chain/v1_15_11/params"
	"github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/log"
	"github.com/shopspring/decimal"
	"go.uber.org/zap"
	"math/big"
	"time"
)

type BlockHeightTime struct {
	HeightBigInt *big.Int
	Height       uint64
	Timestamp    uint64
	Time         time.Time
}

func GetBlockHeightTime(header *ethtypes.Header) *BlockHeightTime {
	return &BlockHeightTime{
		HeightBigInt: header.Number,
		Height:       header.Number.Uint64(),
		Timestamp:    header.Time,
		Time:         time.Unix((int64)(header.Time), 0).UTC(),
	}
}

type BlockContext struct {
	// input
	Block            *ethtypes.Block
	BlockReceipts    []*ethtypes.Receipt
	HeightTime       *BlockHeightTime
	NativeTokenPrice decimal.Decimal
	TxIndex2TxSender map[uint]common.Address
	// output
	BlockResult *BlockResult
}

func (c *BlockContext) GetBlockNumber() uint64 {
	return c.HeightTime.Height
}

func (c *BlockContext) GetTxSender(txIndex uint) (common.Address, error) {
	if txSender, ok := c.TxIndex2TxSender[txIndex]; ok {
		return txSender, nil
	}

	blockTxs := c.Block.Transactions()
	if txIndex >= uint(len(blockTxs)) {
		log.Warn("receipt txIndex is bigger than block txs length",
			zap.Uint64("blockNumber", c.HeightTime.HeightBigInt.Uint64()),
			zap.Any("blockTxs length", len(blockTxs)),
			zap.Uint("txIndex", txIndex))
		return ZeroAddress, nil
	}

	signer := ethtypes.MakeSigner(chainparams.ChainConfig, c.HeightTime.HeightBigInt, c.HeightTime.Timestamp)
	sender, err := ethtypes.Sender(signer, blockTxs[txIndex])
	if err != nil {
		return ZeroAddress, err
	}
	c.TxIndex2TxSender[txIndex] = sender
	return sender, nil
}
