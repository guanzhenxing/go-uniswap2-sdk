package constants

// ChainID network chain id
type ChainID int

//go:generate stringer -type=ChainID -linecomment
const (
	Mainnet     ChainID = 1
	Ropsten     ChainID = 3
	Rinkeby     ChainID = 4
	Goerli      ChainID = 5
	Kovan       ChainID = 42
	BSC         ChainID = 56
	ArbitrumOne ChainID = 42161
	Polygon     ChainID = 137
	OP          ChainID = 10
	Base        ChainID = 8453
	Linea       ChainID = 59144
	zkSync      ChainID = 324
	Scroll      ChainID = 534352
)
