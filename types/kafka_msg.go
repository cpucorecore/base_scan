package types

import (
	"base_scan/repository/orm"
)

type BlockInfo struct {
	BlockNumber          uint64
	BlockUnixTimestamp   uint64
	BnbPrice             string
	Txs                  []*orm.Tx
	NewTokens            []*orm.Token
	NewPairs             []*orm.Pair
	PoolUpdates          []*PoolUpdate
	PoolUpdateParameters []*PoolUpdateParameter
}

type BlockInfoOld struct {
	BlockNumber            uint64
	BlockAt                uint64
	BnbPrice               string
	Txs                    []*orm.Tx
	PoolUpdatesV2          []*PoolUpdate
	PoolUpdateParametersV3 []*PoolUpdateParameter
}
