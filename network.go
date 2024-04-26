package main

type Chain int

const (
	Regtest Chain = iota
	Testnet3
)

func (c Chain) MagicValue() uint32 {
	switch c {
	case Regtest:
		return 0xDAB5BFFA
	case Testnet3:
		return 0x0709110B
	default:
		panic("invalid chain")
	}
}
