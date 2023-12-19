package entities

import (
	"fmt"
	"strings"

	"github.com/ethereum/go-ethereum/common"
)

var (
	ErrDiffChainID = fmt.Errorf("diff chain id")
	ErrDiffToken   = fmt.Errorf("diff token")
	ErrSameAddr    = fmt.Errorf("same address")
)

/**
 * Represents an ERC20 token with a unique address and some metadata.
 */
type Token struct {
	*Currency

	common.Address
}

func NewToken(address common.Address, decimals int, symbol, name string) (*Token, error) {
	currency, err := newCurrency(decimals, symbol, name)
	if err != nil {
		return nil, err
	}

	return &Token{
		Currency: currency,
		Address:  address,
	}, nil
}

/**
 * Returns true if the two tokens are equivalent, i.e. have the same chainId and address.
 * @param other other token to compare
 */
func (t *Token) Equals(other *Token) bool {
	if t == other {
		return true
	}

	return t.Address == other.Address
}

/**
 * Returns true if the address of this token sorts before the address of the other token
 * @param other other token to compare
 * @throws if the tokens have the same address
 * @throws if the tokens are on different chains
 */
func (t *Token) SortsBefore(other *Token) (bool, error) {

	if t.Address == other.Address {
		return false, ErrSameAddr
	}

	return strings.ToLower(t.Address.String()) < strings.ToLower(other.Address.String()), nil
}
