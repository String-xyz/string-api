package model

import "errors"

type Chain struct {
	ChainID       uint32  `json:"chainID" db:"chainID"`
	RPC           string  `json:"RPC" db:"RPC"`
	CoingeckoName string  `json:"coingeckoName" db:"coingeckoName"`
	OwlracleName  string  `json:"owlracleName" db:"owlracleName"`
	StringFee     float32 `json:"stringFee" db:"stringFee"`
}

var chains = []Chain{
	{
		ChainID:       1,
		RPC:           "https://mainnet.infura.io/v3",
		CoingeckoName: "ethereum",
		OwlracleName:  "eth",
		StringFee:     0.035,
	},
	{
		ChainID:       5,
		RPC:           "https://goerli.infura.io/v3/9aa3d95b3bc440fa88ea12eaa4456161",
		CoingeckoName: "ethereum",
		OwlracleName:  "eth",
		StringFee:     0.035,
	},
	{
		ChainID:       137,
		RPC:           "https://rpc-mainnet.matic.quiknode.pro",
		CoingeckoName: "matic-network",
		OwlracleName:  "poly",
		StringFee:     0.05,
	},
	{
		ChainID:       43113,
		RPC:           "https://api.avax-test.network/ext/bc/C/rpc",
		CoingeckoName: "avalanche-2",
		OwlracleName:  "avax",
		StringFee:     0.05,
	},
	{
		ChainID:       43114,
		RPC:           "https://api.avax.network/ext/bc/C/rpc",
		CoingeckoName: "avalanche-2",
		OwlracleName:  "avax",
		StringFee:     0.05,
	},
	{
		ChainID:       80001,
		RPC:           "https://matic-mumbai.chainstacklabs.com",
		CoingeckoName: "matic-network",
		OwlracleName:  "poly",
		StringFee:     0.05,
	},
}

func ChainInfo(chainID uint32) (Chain, error) {
	for _, c := range chains {
		if c.ChainID == chainID {
			return c, nil
		}
	}
	return Chain{}, errors.New("ChainInfo: chainID not supported")
}
