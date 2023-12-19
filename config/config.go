package config

import "github.com/ethereum/go-ethereum/common"

type UniSwap2Config struct {
	TokenDecimals  int
	Univ2Symbol    string
	Univ2Name      string
	FactoryAddress common.Address
	InitCodeHash   []byte
}

var (
	Conf *UniSwap2Config
)

func Init(tokenDecimals int, univ2Symbol, univ2Name string,
	factoryAddress common.Address, initCodeHash []byte) {
	Conf = &UniSwap2Config{
		TokenDecimals:  tokenDecimals,
		Univ2Symbol:    univ2Symbol,
		Univ2Name:      univ2Name,
		FactoryAddress: factoryAddress,
		InitCodeHash:   initCodeHash,
	}
}
