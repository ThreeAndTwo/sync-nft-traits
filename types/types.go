package types

import (
	"math/big"

	"github.com/alethio/web3-multicall-go/multicall"
)

type NftInfo struct {
	URI         string
	TotalSupply *big.Int
}

type ChainInfo struct {
	McClient *multicall.Multicall
}

type NftData struct {
	Key         string
	Val         string
	Count       int64
	RarityRate  float64
	TotalSupply int64
}

type Attribute struct {
	TraitType string `json:"trait_type"`
	Value     string `json:"value"`
}
