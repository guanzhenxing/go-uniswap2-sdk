package utils

import (
	"fmt"
	"github.com/guanzhenxing/go-uniswap2-sdk/constants"

	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

// ValidateSolidityTypeInstance determines if a value is a legal SolidityType
func ValidateSolidityTypeInstance(value *big.Int, t constants.SolidityType) error {
	if value.Cmp(constants.Zero) < 0 || value.Cmp(constants.SolidityTypeMaxima[t]) > 0 {
		return fmt.Errorf(`%v is not a %s`, value, t)
	}
	return nil
}

// ValidateAndParseAddress warns if addresses are not checksummed
func ValidateAndParseAddress(address string) common.Address {
	return common.HexToAddress(address)
}
