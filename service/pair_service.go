package service

import (
	pancakev2 "base_scan/abi/pancake/v2"
	pancakev3 "base_scan/abi/pancake/v3"
	uniswapv2 "base_scan/abi/uniswap/v2"
	uniswapv3 "base_scan/abi/uniswap/v3"
	"base_scan/cache"
	"base_scan/log"
	"base_scan/metrics"
	"base_scan/service/contract_caller"
	"base_scan/types"
	"context"
	"github.com/ethereum/go-ethereum/common"
	"github.com/shopspring/decimal"
	"go.uber.org/zap"
	"golang.org/x/sync/singleflight"
	"time"
)

type PairService interface {
	SetPair(pair *types.Pair)
	GetTokens(pair *types.Pair) *types.PairWrap
	GetPairAndTokens(address common.Address, protocolIds []int) *types.PairWrap
}

type pairService struct {
	ctx            context.Context
	cache          cache.Cache
	contractCaller *contract_caller.ContractCaller
	group          singleflight.Group
}

func NewPairService(
	cache cache.Cache,
	contractCaller *contract_caller.ContractCaller,
) PairService {
	return &pairService{
		ctx:            context.Background(),
		cache:          cache,
		contractCaller: contractCaller,
	}
}

func (s *pairService) SetPair(pair *types.Pair) {
	s.cache.SetPair(pair)
}

func (s *pairService) getToken(tokenAddress common.Address) (*types.Token, error) {
	// TODO parallelize contract call
	token := &types.Token{
		Address: tokenAddress,
	}

	name, callNameErr := s.contractCaller.CallName(&tokenAddress)
	if callNameErr == nil {
		token.Name = name
	}

	symbol, callSymbolErr := s.contractCaller.CallSymbol(&tokenAddress)
	if callSymbolErr == nil {
		token.Symbol = symbol
	}

	decimals, callDecimalsErr := s.contractCaller.CallDecimals(&tokenAddress)
	if callDecimalsErr != nil {
		token.Filtered = true
		return token, callDecimalsErr
	}
	token.Decimals = (int8)(decimals)

	totalSupply, callTotalSupplyErr := s.contractCaller.CallTotalSupply(&tokenAddress)
	if callTotalSupplyErr == nil {
		token.TotalSupply = decimal.NewFromBigInt(totalSupply, -(int32)(token.Decimals))
	}

	return token, nil
}

func (s *pairService) GetToken(tokenAddress common.Address) (*types.Token, error, bool) {
	cacheToken, ok := s.cache.GetToken(tokenAddress)
	if ok {
		return cacheToken, nil, true
	}

	now := time.Now()
	token, err := s.getToken(tokenAddress)
	if err != nil {
		s.cache.SetToken(token)
		return nil, err, false
	}
	metrics.GetTokenDuration.Observe(time.Since(now).Seconds())

	s.cache.SetToken(token)
	return token, nil, false
}

func (s *pairService) getTokens(pair *types.Pair) *types.PairWrap {
	pairWrap := &types.PairWrap{
		Pair:    pair,
		NewPair: true,
	}

	token0, getToken0Err, token0FromCache := s.GetToken(pair.Token0Core.Address)
	if getToken0Err != nil {
		pair.Filtered = true
		pair.FilterCode = types.FilterCodeGetToken0
		return pairWrap
	}

	token1, getToken1Err, token1FromCache := s.GetToken(pair.Token1Core.Address)
	if getToken1Err != nil {
		pair.Filtered = true
		pair.FilterCode = types.FilterCodeGetToken1
		return pairWrap
	}

	pair.Token0 = token0
	pair.Token1 = token1

	pair.Token0Core.Symbol = token0.Symbol
	pair.Token0Core.Decimals = token0.Decimals
	pair.Token1Core.Symbol = token1.Symbol
	pair.Token1Core.Decimals = token1.Decimals

	tokensReversed := pair.OrderToken0Token1()
	if tokensReversed {
		pairWrap.NewToken0 = !token1FromCache
		pairWrap.NewToken1 = !token0FromCache
	} else {
		pairWrap.NewToken0 = !token0FromCache
		pairWrap.NewToken1 = !token1FromCache
	}

	return pairWrap
}

func (s *pairService) GetTokens(pair *types.Pair) *types.PairWrap {
	pairWrap := s.getTokens(pair)
	s.SetPair(pair)
	return pairWrap
}

func (s *pairService) getPairAndTokens(address common.Address, protocolIds []int) *types.PairWrap {
	doResult, _, _ := s.group.Do(address.String(), func() (interface{}, error) {
		pair := s.getPair(address)
		if pair.Filtered {
			s.SetPair(pair)
			return &types.PairWrap{
				Pair:      pair,
				NewPair:   false,
				NewToken0: false,
				NewToken1: false,
			}, nil
		}

		if !s.verifyPair(pair, protocolIds) {
			s.SetPair(pair)
			return &types.PairWrap{
				Pair:      pair,
				NewPair:   false,
				NewToken0: false,
				NewToken1: false,
			}, nil
		}

		return s.GetTokens(pair), nil
	})

	return doResult.(*types.PairWrap)
}

func (s *pairService) GetPairAndTokens(address common.Address, protocolIds []int) *types.PairWrap {
	cachePair, ok := s.cache.GetPair(address)
	if ok {
		return &types.PairWrap{
			Pair: cachePair,
		}
	}

	pairWrap := s.getPairAndTokens(address, protocolIds)
	return pairWrap
}

func (s *pairService) getPair(pairAddress common.Address) *types.Pair {
	pair := &types.Pair{
		Address: pairAddress,
	}

	token0Address, err0 := s.contractCaller.CallToken0(&pairAddress)
	if err0 != nil {
		log.Logger.Error("CallToken0 err, this pair will filtered", zap.Error(err0), zap.String("address", pairAddress.String()))
		pair.Filtered = true
		pair.FilterCode = types.FilterCodeGetToken0
		return pair
	}

	pair.Token0Core = &types.TokenCore{
		Address: token0Address,
	}

	token1Address, err1 := s.contractCaller.CallToken1(&pairAddress)
	if err1 != nil {
		log.Logger.Error("CallToken1 err, this pair will filtered", zap.Error(err1), zap.String("address", pairAddress.String()))
		pair.Filtered = true
		pair.FilterCode = types.FilterCodeGetToken1
		return pair
	}

	pair.Token1Core = &types.TokenCore{
		Address: token1Address,
	}

	pair.FilterByToken0AndToken1()

	return pair
}

func (s *pairService) verifyPair(pair *types.Pair, protocolIds []int) bool {
	for _, protocolId := range protocolIds {
		switch protocolId {
		case types.ProtocolIdUniswapV2:
			pairAddressQueried, getPairErr := s.contractCaller.CallGetPair(&uniswapv2.FactoryAddress, &pair.Token0Core.Address, &pair.Token1Core.Address)
			if getPairErr != nil {
				continue
			}

			if types.IsSameAddress(pairAddressQueried, pair.Address) {
				pair.ProtocolId = protocolId
				return true
			}

		case types.ProtocolIdPancakeV2:
			pairAddressQueried, getPairErr := s.contractCaller.CallGetPair(&pancakev2.FactoryAddress, &pair.Token0Core.Address, &pair.Token1Core.Address)
			if getPairErr != nil {
				continue
			}

			if types.IsSameAddress(pairAddressQueried, pair.Address) {
				pair.ProtocolId = protocolId
				return true
			}

		case types.ProtocolIdUniswapV3:
			fee, callFeeErr := s.contractCaller.CallFee(&pair.Address)
			if callFeeErr != nil {
				continue
			}

			pairAddressQueried, getPairErr := s.contractCaller.CallGetPool(&uniswapv3.FactoryAddress, &pair.Token0Core.Address, &pair.Token1Core.Address, fee)
			if getPairErr != nil {
				continue
			}

			if types.IsSameAddress(pairAddressQueried, pair.Address) {
				pair.ProtocolId = protocolId
				return true
			}

		case types.ProtocolIdPancakeV3:
			fee, callFeeErr := s.contractCaller.CallFee(&pair.Address)
			if callFeeErr != nil {
				continue
			}

			pairAddressQueried, getPairErr := s.contractCaller.CallGetPool(&pancakev3.FactoryAddress, &pair.Token0Core.Address, &pair.Token1Core.Address, fee)
			if getPairErr != nil {
				continue
			}

			if types.IsSameAddress(pairAddressQueried, pair.Address) {
				pair.ProtocolId = protocolId
				return true
			}

		case types.ProtocolIdAerodrome:
			isPool, isPoolErr := s.contractCaller.CallIsPool(&pair.Address)
			if isPoolErr != nil {
				continue
			}

			if isPool {
				pair.ProtocolId = protocolId
				return true
			}
		}
	}

	pair.Filtered = true
	pair.FilterCode = types.FilterCodeVerifyFailed
	return false
}

func (s *pairService) getPairV2(protocolId int, pairAddress common.Address) *types.Pair {
	now := time.Now()

	pair := &types.Pair{
		ProtocolId: protocolId,
		Address:    pairAddress,
	}

	token0Address, err0 := s.contractCaller.CallToken0(&pairAddress)
	if err0 != nil {
		log.Logger.Error("CallToken0 err, this pair will filtered", zap.Error(err0), zap.String("address", pairAddress.String()))
		pair.Filtered = true
		pair.FilterCode = types.FilterCodeGetToken0
		return pair
	}

	pair.Token0Core = &types.TokenCore{
		Address: token0Address,
	}

	token1Address, err1 := s.contractCaller.CallToken1(&pairAddress)
	if err1 != nil {
		log.Logger.Error("CallToken1 err, this pair will filtered", zap.Error(err1), zap.String("address", pairAddress.String()))
		pair.Filtered = true
		pair.FilterCode = types.FilterCodeGetToken1
		return pair
	}

	pair.Token1Core = &types.TokenCore{
		Address: token1Address,
	}

	if pair.FilterByToken0AndToken1() {
		return pair
	}

	isPool, getPairErr := s.contractCaller.CallIsPool(&pairAddress)
	if getPairErr != nil {
		log.Logger.Error("CallGetPair err, this pair will filtered", zap.Error(getPairErr), zap.String("address", pairAddress.String()))
		pair.Filtered = true
		pair.FilterCode = types.FilterCodeGetPair
		return pair
	}

	if !isPool {
		pair.Filtered = true
		pair.FilterCode = types.FilterCodeVerifyFailed
		return pair
	}

	metrics.GetV2PairDuration.Observe(time.Since(now).Seconds())
	return pair
}
