package constants

import (
	"math/big"
)

type TradeType int

const (
	ExactInput TradeType = iota
	ExactOutput
)

type Rounding int

const (
	RoundDown Rounding = iota
	RoundHalfUp
	RoundUp
)

// Valid check this rounding mode is valid
func (r Rounding) Valid() bool {
	return r == RoundDown ||
		r == RoundHalfUp ||
		r == RoundUp
}

var (
	MinimumLiquidity = big.NewInt(1000)

	Zero  = big.NewInt(0)
	One   = big.NewInt(1)
	Two   = big.NewInt(2)
	Three = big.NewInt(3)
	Five  = big.NewInt(5)
	Ten   = big.NewInt(10)

	B100  = big.NewInt(100)
	B997  = big.NewInt(997)
	B1000 = big.NewInt(1000)
)
