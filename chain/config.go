package chain

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/params"
	"math/big"
)

/*
from geth v1.15.11
TODO: should use op-geth chain config:
https://github.com/ethereum-optimism/op-geth/blob/optimism/params/config.go
*/

func newUint64(v uint64) *uint64 {
	return &v
}

var (
	MainnetTerminalTotalDifficulty, _ = new(big.Int).SetString("58_750_000_000_000_000_000_000", 0)

	DefaultCancunBlobConfig = &params.BlobConfig{
		Target:         3,
		Max:            6,
		UpdateFraction: 3338477,
	}

	DefaultPragueBlobConfig = &params.BlobConfig{
		Target:         6,
		Max:            9,
		UpdateFraction: 5007716,
	}

	ChainConfig = &params.ChainConfig{
		ChainID:                 big.NewInt(Id),
		HomesteadBlock:          big.NewInt(1_150_000),
		DAOForkBlock:            big.NewInt(1_920_000),
		DAOForkSupport:          true,
		EIP150Block:             big.NewInt(2_463_000),
		EIP155Block:             big.NewInt(2_675_000),
		EIP158Block:             big.NewInt(2_675_000),
		ByzantiumBlock:          big.NewInt(4_370_000),
		ConstantinopleBlock:     big.NewInt(7_280_000),
		PetersburgBlock:         big.NewInt(7_280_000),
		IstanbulBlock:           big.NewInt(9_069_000),
		MuirGlacierBlock:        big.NewInt(9_200_000),
		BerlinBlock:             big.NewInt(12_244_000),
		LondonBlock:             big.NewInt(12_965_000),
		ArrowGlacierBlock:       big.NewInt(13_773_000),
		GrayGlacierBlock:        big.NewInt(15_050_000),
		TerminalTotalDifficulty: MainnetTerminalTotalDifficulty, // 58_750_000_000_000_000_000_000
		ShanghaiTime:            newUint64(1681338455),
		CancunTime:              newUint64(1710338135),
		PragueTime:              newUint64(1746612311),
		DepositContractAddress:  common.HexToAddress("0x00000000219ab540356cbb839cbe05303d7705fa"),
		Ethash:                  new(params.EthashConfig),
		BlobScheduleConfig: &params.BlobScheduleConfig{
			Cancun: DefaultCancunBlobConfig,
			Prague: DefaultPragueBlobConfig,
		},
	}
)
