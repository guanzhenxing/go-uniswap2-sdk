package entities

import (
	"fmt"
	"github.com/guanzhenxing/go-uniswap2-sdk/constants"
	"math/big"
)

var (

	// ErrInvalidLiquidity invalid liquidity
	ErrInvalidLiquidity = fmt.Errorf("invalid liquidity")
	// ErrInvalidKLast invalid kLast
	ErrInvalidKLast = fmt.Errorf("invalid kLast")
)

// TokenAmounts warps TokenAmount array
type TokenAmounts [2]*TokenAmount

// Tokens warps Token array
type Tokens [2]*Token

// NewTokenAmounts creates a TokenAmount
func NewTokenAmounts(tokenAmountA, tokenAmountB *TokenAmount) (TokenAmounts, error) {
	ok, err := tokenAmountA.Token.SortsBefore(tokenAmountB.Token)
	if err != nil {
		return TokenAmounts{}, err
	}
	if ok {
		return TokenAmounts{tokenAmountA, tokenAmountB}, nil
	}
	return TokenAmounts{tokenAmountB, tokenAmountA}, nil
}

// Pair warps uniswap pair
type Pair struct {
	// sorted tokens
	TokenAmounts
}

// NewPair creates Pair
func NewPair(tokenAmountA, tokenAmountB *TokenAmount) (*Pair, error) {
	tokenAmounts, err := NewTokenAmounts(tokenAmountA, tokenAmountB)
	if err != nil {
		return nil, err
	}

	pair := &Pair{
		TokenAmounts: tokenAmounts,
	}

	return pair, err
}

// InvolvesToken Returns true if the token is either token0 or token1
// @param token to check
func (p *Pair) InvolvesToken(token *Token) bool {
	return token.Equals(p.TokenAmounts[0].Token) || token.Equals(p.TokenAmounts[1].Token)
}

// Token0Price Returns the current mid price of the pair in terms of token0, i.e. the ratio of reserve1 to reserve0
func (p *Pair) Token0Price() *Price {
	return NewPrice(p.Token0().Currency, p.Token1().Currency, p.TokenAmounts[0].Raw(), p.TokenAmounts[1].Raw())
}

// Token1Price Returns the current mid price of the pair in terms of token1, i.e. the ratio of reserve0 to reserve1
func (p *Pair) Token1Price() *Price {
	return NewPrice(p.Token1().Currency, p.Token0().Currency, p.TokenAmounts[1].Raw(), p.TokenAmounts[0].Raw())
}

// PriceOf Returns the price of the given token in terms of the other token in the pair.
// @param token token to return price of
func (p *Pair) PriceOf(token *Token) (*Price, error) {
	if !p.InvolvesToken(token) {
		return nil, ErrDiffToken
	}

	if token.Equals(p.Token0()) {
		return p.Token0Price(), nil
	}
	return p.Token1Price(), nil
}

// Token0 returns the first token in the pair
func (p *Pair) Token0() *Token {
	return p.TokenAmounts[0].Token
}

// Token1 returns the last token in the pair
func (p *Pair) Token1() *Token {
	return p.TokenAmounts[1].Token
}

// Reserve0 returns the first TokenAmount in the pair
func (p *Pair) Reserve0() *TokenAmount {
	return p.TokenAmounts[0]
}

// Reserve1 returns the last TokenAmount in the pair
func (p *Pair) Reserve1() *TokenAmount {
	return p.TokenAmounts[1]
}

// ReserveOf returns the TokenAmount that equals to the token
func (p *Pair) ReserveOf(token *Token) (*TokenAmount, error) {
	if !p.InvolvesToken(token) {
		return nil, ErrDiffToken
	}

	if token.Equals(p.Token0()) {
		return p.Reserve0(), nil
	}
	return p.Reserve1(), nil
}

// GetOutputAmount returns OutputAmount and a Pair for the InputAmout
func (p *Pair) GetOutputAmount(inputAmount *TokenAmount) (*TokenAmount, *Pair, error) {
	if !p.InvolvesToken(inputAmount.Token) {
		return nil, nil, ErrDiffToken
	}

	if p.Reserve0().Raw().Cmp(constants.Zero) == 0 ||
		p.Reserve1().Raw().Cmp(constants.Zero) == 0 {
		return nil, nil, ErrInsufficientReserves
	}

	inputReserve, err := p.ReserveOf(inputAmount.Token)
	if err != nil {
		return nil, nil, err
	}
	token := p.Token0()
	if inputAmount.Token.Equals(p.Token0()) {
		token = p.Token1()
	}
	outputReserve, err := p.ReserveOf(token)
	if err != nil {
		return nil, nil, err
	}

	inputAmountWithFee := big.NewInt(0).Mul(inputAmount.Raw(), constants.B997)
	numerator := big.NewInt(0).Mul(inputAmountWithFee, outputReserve.Raw())
	denominator := big.NewInt(0).Add(big.NewInt(0).Mul(inputReserve.Raw(), constants.B1000), inputAmountWithFee)
	outputAmount, err := NewTokenAmount(token, big.NewInt(0).Div(numerator, denominator))
	if err != nil {
		return nil, nil, err
	}
	if outputAmount.Raw().Cmp(constants.Zero) == 0 {
		return nil, nil, ErrInsufficientInputAmount
	}

	tokenAmountA, err := inputAmount.Add(inputReserve)
	if err != nil {
		return nil, nil, err
	}
	tokenAmountB, err := outputReserve.Subtract(outputAmount)
	if err != nil {
		return nil, nil, err
	}
	pair, err := NewPair(tokenAmountA, tokenAmountB)
	if err != nil {
		return nil, nil, err
	}
	return outputAmount, pair, nil
}

// GetInputAmount returns InputAmout and a Pair for the OutputAmount
func (p *Pair) GetInputAmount(outputAmount *TokenAmount) (*TokenAmount, *Pair, error) {
	if !p.InvolvesToken(outputAmount.Token) {
		return nil, nil, ErrDiffToken
	}

	outputReserve, err := p.ReserveOf(outputAmount.Token)
	if err != nil {
		return nil, nil, err
	}
	if p.Reserve0().Raw().Cmp(constants.Zero) == 0 ||
		p.Reserve1().Raw().Cmp(constants.Zero) == 0 ||
		outputAmount.Raw().Cmp(outputReserve.Raw()) >= 0 {
		return nil, nil, ErrInsufficientReserves
	}

	token := p.Token0()
	if outputAmount.Token.Equals(p.Token0()) {
		token = p.Token1()
	}
	inputReserve, err := p.ReserveOf(token)
	if err != nil {
		return nil, nil, err
	}

	numerator := big.NewInt(0).Mul(inputReserve.Raw(), outputAmount.Raw())
	numerator.Mul(numerator, constants.B1000)
	denominator := big.NewInt(0).Sub(outputReserve.Raw(), outputAmount.Raw())
	denominator.Mul(denominator, constants.B997)
	amount := big.NewInt(0).Div(numerator, denominator)
	amount.Add(amount, constants.One)
	inputAmount, err := NewTokenAmount(token, amount)
	if err != nil {
		return nil, nil, err
	}

	tokenAmountA, err := inputAmount.Add(inputReserve)
	if err != nil {
		return nil, nil, err
	}
	tokenAmountB, err := outputReserve.Subtract(outputAmount)
	if err != nil {
		return nil, nil, err
	}
	pair, err := NewPair(tokenAmountA, tokenAmountB)
	if err != nil {
		return nil, nil, err
	}
	return inputAmount, pair, nil
}
