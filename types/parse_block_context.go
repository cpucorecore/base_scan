package types

import (
	chainparams "base_scan/chain"
	"base_scan/log"
	"errors"
	"github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/shopspring/decimal"
	"go.uber.org/zap"
	"math/big"
	"time"
)

var (
	errx = errors.New("xx")
)

type BlockHeightTime struct {
	Height       uint64
	Timestamp    uint64
	HeightBigInt *big.Int
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

type ParseBlockContext struct {
	// input
	Block            *ethtypes.Block
	BlockReceipts    []*ethtypes.Receipt
	HeightTime       *BlockHeightTime
	NativeTokenPrice decimal.Decimal
	TxIndex2TxSender map[uint]common.Address
	// output
	BlockResult *BlockResult
}

func (c *ParseBlockContext) GetBlockNumber() uint64 {
	return c.HeightTime.Height
}

func (c *ParseBlockContext) GetTxSender(txIndex uint) (common.Address, error) {
	if txSender, ok := c.TxIndex2TxSender[txIndex]; ok {
		return txSender, nil
	}

	transactions := c.Block.Transactions()
	transactionsLen := uint(transactions.Len())
	if txIndex >= transactionsLen {
		log.Logger.Info("Waring: txIndex out of range",
			zap.Uint64("height", c.HeightTime.Height),
			zap.Any("transactions length", transactionsLen),
			zap.Uint("txIndex", txIndex),
		)
		return ZeroAddress, errx
	}

	signer := ethtypes.MakeSigner(chainparams.ChainConfig, c.HeightTime.HeightBigInt, c.HeightTime.Timestamp)
	sender, err := ethtypes.Sender(signer, transactions[txIndex])
	if err != nil {
		return ZeroAddress, err
	}

	c.TxIndex2TxSender[txIndex] = sender
	return sender, nil
}
