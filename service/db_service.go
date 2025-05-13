package service

import (
	"base_scan/repository"
	"base_scan/types/orm"
)

type DBService interface {
	AddTokens(tokens []*orm.Token) error
	AddPairs(pairs []*orm.Pair) error
	UpdateToken(tokenAddress, mainPairAddress string) error
	AddTxs(txs []*orm.Tx) error
}

type dbService struct {
	tokenRepository *repository.TokenRepository
	pairRepository  *repository.PairRepository
	txRepository    *repository.TxRepository
}

func (s *dbService) AddTokens(tokens []*orm.Token) error {
	return s.tokenRepository.CreateBatch(tokens, "address", "chain_id")
}

func (s *dbService) AddPairs(pairs []*orm.Pair) error {
	return s.pairRepository.CreateBatch(pairs, "address", "chain_id")
}

func (s *dbService) UpdateToken(tokenAddress, mainPairAddress string) error {
	return s.tokenRepository.UpdateMainPair(tokenAddress, mainPairAddress)
}

func (s *dbService) AddTxs(txs []*orm.Tx) error {
	return s.txRepository.CreateBatch(txs, "token0_address", "block", "block_index", "tx_index")
}

func NewDBService(
	tokenRepository *repository.TokenRepository,
	pairRepository *repository.PairRepository,
	txRepository *repository.TxRepository,
) DBService {
	return &dbService{
		tokenRepository: tokenRepository,
		pairRepository:  pairRepository,
		txRepository:    txRepository,
	}
}
