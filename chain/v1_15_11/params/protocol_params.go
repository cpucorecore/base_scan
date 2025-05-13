package params

const (
	DefaultBaseFeeChangeDenominator = 8 // Bounds the amount the base fee can change between blocks.
	DefaultElasticityMultiplier     = 2 // Bounds the maximum gas limit an EIP-1559 block may have.
)
