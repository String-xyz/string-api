package service

import "github.com/String-xyz/string-api/pkg/repository"

type Chain struct {
	ChainID       uint64
	RPC           string
	CoingeckoName string
	OwlracleName  string
	StringFee     float64
}

// TODO: should we store this in a DB or determine it dynamically???  Previously this was defined in the preprocessor in the Chain array
func stringFee(chainId uint64) (float64, error) {
	return 0.05, nil
}

func ChainInfo(chainId uint64, networkRepo repository.Network, assetRepo repository.Asset) (Chain, error) {
	network, err := networkRepo.GetChainID(chainId)
	if err != nil {
		return Chain{}, err
	}
	asset, err := assetRepo.GetID(network.GasTokenID)
	if err != nil {
		return Chain{}, err
	}
	fee, err := stringFee(chainId)
	if err != nil {
		return Chain{}, err
	}
	return Chain{ChainID: chainId, RPC: network.RPCUrl, CoingeckoName: asset.ValueOracle, OwlracleName: network.GasOracle, StringFee: fee}, nil
}
