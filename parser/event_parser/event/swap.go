package event

import (
	"base_scan/parser/protocol2"
	"base_scan/types"
	"base_scan/types/orm"
	"github.com/shopspring/decimal"
	"math/big"
)

type SwapEvent struct {
	*types.EventCommon
	Amount0InWei  *big.Int
	Amount1InWei  *big.Int
	Amount0OutWei *big.Int
	Amount1OutWei *big.Int
}

func (e *SwapEvent) CanGetTx() bool {
	return true
}

func (e *SwapEvent) GetTx(bnbPrice decimal.Decimal) *orm.Tx {
	tx := &orm.Tx{
		TxHash:        e.TxHash.String(),
		Maker:         e.Maker.String(),
		Token0Address: e.Pair.Token0Core.Address.String(),
		Token1Address: e.Pair.Token1Core.Address.String(),
		Block:         e.BlockNumber,
		BlockAt:       e.BlockTime,
		BlockIndex:    e.TxIndex,
		TxIndex:       e.LogIndex,
		PairAddress:   e.Pair.Address.String(),
		Program:       types.GetProtocolName(e.Pair.ProtocolId),
	}

	if e.Amount0InWei.Cmp(types.ZeroBigInt) > 0 {
		tx.Token0Amount, tx.Token1Amount = protocol2.ParseAmountsByPair(e.Amount0InWei, e.Amount1OutWei, e.Pair)
		if !e.Pair.TokensReversed {
			tx.Event = protocol2.EventNameSell
		} else {
			tx.Event = protocol2.EventNameBuy
		}
	} else if e.Amount1InWei.Cmp(types.ZeroBigInt) > 0 {
		tx.Token0Amount, tx.Token1Amount = protocol2.ParseAmountsByPair(e.Amount0OutWei, e.Amount1InWei, e.Pair)
		if !e.Pair.TokensReversed {
			tx.Event = protocol2.EventNameBuy
		} else {
			tx.Event = protocol2.EventNameSell
		}
	} else {
	}

	tx.AmountUsd, tx.PriceUsd = protocol2.CalcAmountAndPrice(bnbPrice, tx.Token0Amount, tx.Token1Amount, e.Pair.Token1Core.Address)
	return tx
}

var _ types.Event = (*SwapEvent)(nil)
