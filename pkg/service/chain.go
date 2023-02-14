// TODO: Make this service instantiable

package service

import (
	"github.com/String-xyz/string-api/pkg/internal/common"
	"github.com/String-xyz/string-api/pkg/repository"
)

type Chain struct {
	ChainID       uint64
	RPC           string
	Explorer      string
	CoingeckoName string
	OwlracleName  string
	StringFee     float64
	UUID          string
	GasTokenID    string
}

// TODO: should we store this in a DB or determine it dynamically???  Previously this was defined in the preprocessor in the Chain array
func stringFee(chainId uint64) (float64, error) {
	return 0.03, nil
}

func ChainInfo(chainId uint64, networkRepo repository.Network, assetRepo repository.Asset) (Chain, error) {
	network, err := networkRepo.GetByChainId(chainId)
	if err != nil {
		return Chain{}, common.StringError(err)
	}
	asset, err := assetRepo.GetById(network.GasTokenID)
	if err != nil {
		return Chain{}, common.StringError(err)
	}
	fee, err := stringFee(chainId)
	if err != nil {
		return Chain{}, common.StringError(err)
	}
	return Chain{ChainID: chainId, RPC: network.RPCUrl, Explorer: network.ExplorerUrl, CoingeckoName: asset.ValueOracle.String, OwlracleName: network.GasOracle, StringFee: fee, UUID: network.ID, GasTokenID: network.GasTokenID}, nil
}
