package main

import (
	"database/sql"
	"fmt"
	"os"
	"time"

	"github.com/String-xyz/string-api/api"
	"github.com/String-xyz/string-api/pkg/model"
	"github.com/String-xyz/string-api/pkg/store"
	"github.com/joho/godotenv"
	"github.com/rs/zerolog"
)

// Set this!
const stringPublicAddress = "0x44A4b9E2A69d86BA382a511f845CbF2E31286771"

func main() {
	var env string
	if len(os.Args) > 1 {
		env = os.Args[1]
	}
	if env == "local" {
		fmt.Printf("\n\nSeeding Mock Data")
		mockSeeding()
	} else {
		fmt.Printf("\n\nSeeding Production Data")
		dataSeeding()
	}
}

func dataSeeding() {
	// Initialize repos
	godotenv.Load(".env") // removed the err since in cloud this wont be loaded
	port := os.Getenv("PORT")
	if port == "" {
		panic("no port!")
	}
	lg := zerolog.New(os.Stdout)

	config := api.APIConfig{
		DB:     store.MustNewPG(),
		Redis:  store.NewRedisStore(),
		Port:   port,
		Logger: &lg,
	}

	repos := api.NewRepos(config)
	// api.Start(config)

	// Write to repos

	// Networks without GasTokenID
	networkPolygon, err := repos.Network.Create(model.Network{Name: "Polygon Mainnet", NetworkID: 137, ChainID: 137, GasOracle: "poly", RPCUrl: "https://rpc-mainnet.matic.quiknode.pro", ExplorerUrl: "https://polygonscan.com"})
	if err != nil {
		fmt.Printf("%+v", err)
		return
	}
	networkMumbai, err := repos.Network.Create(model.Network{Name: "Mumbai Testnet", NetworkID: 80001, ChainID: 80001, GasOracle: "poly", RPCUrl: "https://matic-mumbai.chainstacklabs.com", ExplorerUrl: "https://mumbai.polygonscan.com/"})
	if err != nil {
		panic(err)
	}
	networkGoerli, err := repos.Network.Create(model.Network{Name: "Goerli Testnet", NetworkID: 5, ChainID: 5, GasOracle: "eth", RPCUrl: "https://goerli.infura.io/v3/9aa3d95b3bc440fa88ea12eaa4456161", ExplorerUrl: "https://goerli.etherscan.io"})
	if err != nil {
		panic(err)
	}
	networkEthereum, err := repos.Network.Create(model.Network{Name: "Ethereum Mainnet", NetworkID: 1, ChainID: 1, GasOracle: "eth", RPCUrl: "https://rpc.ankr.com/eth", ExplorerUrl: "https://etherscan.io/"})
	if err != nil {
		panic(err)
	}
	networkFuji, err := repos.Network.Create(model.Network{Name: "Fuji Testnet", NetworkID: 43113, ChainID: 43113, GasOracle: "avax", RPCUrl: "https://api.avax-test.network/ext/bc/C/rpc", ExplorerUrl: "https://testnet.snowtrace.io"})
	if err != nil {
		panic(err)
	}
	networkAvalanche, err := repos.Network.Create(model.Network{Name: "Avalanche Mainnet", NetworkID: 43114, ChainID: 43114, GasOracle: "avax", RPCUrl: "https://api.avax.network/ext/bc/C/rpc", ExplorerUrl: "https://snowtrace.io"})
	if err != nil {
		panic(err)
	}

	// Assets
	assetAvalanche, err := repos.Asset.Create(model.Asset{Name: "AVAX", Description: "Avalanche", Decimals: 18, IsCrypto: true, NetworkID: nullString(networkAvalanche.ID), ValueOracle: nullString("avalanche-2")})
	if err != nil {
		panic(err)
	}
	assetEthereum, err := repos.Asset.Create(model.Asset{Name: "ETH", Description: "Ethereum", Decimals: 18, IsCrypto: true, NetworkID: nullString(networkEthereum.ID), ValueOracle: nullString("ethereum")})
	if err != nil {
		panic(err)
	}
	assetMatic, err := repos.Asset.Create(model.Asset{Name: "MATIC", Description: "Matic", Decimals: 18, IsCrypto: true, NetworkID: nullString(networkPolygon.ID), ValueOracle: nullString("matic-network")})
	if err != nil {
		panic(err)
	}
	/*assetUSD*/
	_, err = repos.Asset.Create(model.Asset{Name: "USD", Description: "United States Dollar", Decimals: 6, IsCrypto: false})
	if err != nil {
		panic(err)
	}

	// Update Networks with GasTokenIDs
	err = repos.Network.Update(networkPolygon.ID, model.NetworkUpdates{GasTokenID: &assetMatic.ID})
	if err != nil {
		panic(err)
	}
	err = repos.Network.Update(networkMumbai.ID, model.NetworkUpdates{GasTokenID: &assetMatic.ID})
	if err != nil {
		panic(err)
	}
	err = repos.Network.Update(networkGoerli.ID, model.NetworkUpdates{GasTokenID: &assetEthereum.ID})
	if err != nil {
		panic(err)
	}
	err = repos.Network.Update(networkEthereum.ID, model.NetworkUpdates{GasTokenID: &assetEthereum.ID})
	if err != nil {
		panic(err)
	}
	err = repos.Network.Update(networkFuji.ID, model.NetworkUpdates{GasTokenID: &assetAvalanche.ID})
	if err != nil {
		panic(err)
	}
	err = repos.Network.Update(networkAvalanche.ID, model.NetworkUpdates{GasTokenID: &assetAvalanche.ID})
	if err != nil {
		panic(err)
	}

	// String User
	userString, err := repos.User.Create(model.User{Type: "Internal", Status: "Internal"})
	if err != nil {
		panic(err)
	}

	// Instruments, used in TX Legs
	/*instrumentDeveloperCard*/
	_, err = repos.Instrument.Create(model.Instrument{Type: "Bank Account", Status: "Live", Network: "bankprov", PublicKey: "420481286", UserID: userString.ID})
	if err != nil {
		panic(err)
	}
	/*instrumentDeveloperWallet*/
	_, err = repos.Instrument.Create(model.Instrument{Type: "Crypto Wallet", Status: "Internal", Network: "EVM", PublicKey: stringPublicAddress, UserID: userString.ID})
	if err != nil {
		panic(err)
	}
}

func mockSeeding() {
	// Initialize repos
	godotenv.Load(".env") // removed the err since in cloud this wont be loaded
	port := os.Getenv("PORT")
	if port == "" {
		panic("no port!")
	}
	lg := zerolog.New(os.Stdout)

	config := api.APIConfig{
		DB:     store.MustNewPG(),
		Redis:  store.NewRedisStore(),
		Port:   port,
		Logger: &lg,
	}

	repos := api.NewRepos(config)
	// api.Start(config)

	// Write to repos

	// Networks without GasTokenID
	networkPolygon, err := repos.Network.Create(model.Network{Name: "Polygon Mainnet", NetworkID: 137, ChainID: 137, GasOracle: "poly", RPCUrl: "https://rpc-mainnet.matic.quiknode.pro", ExplorerUrl: "https://polygonscan.com"})
	if err != nil {
		fmt.Printf("%+v", err)
		return
	}
	networkMumbai, err := repos.Network.Create(model.Network{Name: "Mumbai Testnet", NetworkID: 80001, ChainID: 80001, GasOracle: "poly", RPCUrl: "https://matic-mumbai.chainstacklabs.com", ExplorerUrl: "https://mumbai.polygonscan.com/"})
	if err != nil {
		panic(err)
	}
	networkGoerli, err := repos.Network.Create(model.Network{Name: "Goerli Testnet", NetworkID: 5, ChainID: 5, GasOracle: "eth", RPCUrl: "https://goerli.infura.io/v3/9aa3d95b3bc440fa88ea12eaa4456161", ExplorerUrl: "https://goerli.etherscan.io"})
	if err != nil {
		panic(err)
	}
	networkEthereum, err := repos.Network.Create(model.Network{Name: "Ethereum Mainnet", NetworkID: 1, ChainID: 1, GasOracle: "eth", RPCUrl: "https://rpc.ankr.com/eth", ExplorerUrl: "https://etherscan.io/"})
	if err != nil {
		panic(err)
	}
	networkFuji, err := repos.Network.Create(model.Network{Name: "Fuji Testnet", NetworkID: 43113, ChainID: 43113, GasOracle: "avax", RPCUrl: "https://api.avax-test.network/ext/bc/C/rpc", ExplorerUrl: "https://testnet.snowtrace.io"})
	if err != nil {
		panic(err)
	}
	networkAvalanche, err := repos.Network.Create(model.Network{Name: "Avalanche Mainnet", NetworkID: 43114, ChainID: 43114, GasOracle: "avax", RPCUrl: "https://api.avax.network/ext/bc/C/rpc", ExplorerUrl: "https://snowtrace.io"})
	if err != nil {
		panic(err)
	}

	// Assets
	assetAvalanche, err := repos.Asset.Create(model.Asset{Name: "AVAX", Description: "Avalanche", Decimals: 18, IsCrypto: true, NetworkID: nullString(networkAvalanche.ID), ValueOracle: nullString("avalanche-2")})
	if err != nil {
		panic(err)
	}
	assetEthereum, err := repos.Asset.Create(model.Asset{Name: "ETH", Description: "Ethereum", Decimals: 18, IsCrypto: true, NetworkID: nullString(networkEthereum.ID), ValueOracle: nullString("ethereum")})
	if err != nil {
		panic(err)
	}
	assetMatic, err := repos.Asset.Create(model.Asset{Name: "MATIC", Description: "Matic", Decimals: 18, IsCrypto: true, NetworkID: nullString(networkPolygon.ID), ValueOracle: nullString("matic-network")})
	if err != nil {
		panic(err)
	}
	/*assetUSD*/
	_, err = repos.Asset.Create(model.Asset{Name: "USD", Description: "United States Dollar", Decimals: 6, IsCrypto: false})
	if err != nil {
		panic(err)
	}

	// Update Networks with GasTokenIDs
	err = repos.Network.Update(networkPolygon.ID, model.NetworkUpdates{GasTokenID: &assetMatic.ID})
	if err != nil {
		panic(err)
	}
	err = repos.Network.Update(networkMumbai.ID, model.NetworkUpdates{GasTokenID: &assetMatic.ID})
	if err != nil {
		panic(err)
	}
	err = repos.Network.Update(networkGoerli.ID, model.NetworkUpdates{GasTokenID: &assetEthereum.ID})
	if err != nil {
		panic(err)
	}
	err = repos.Network.Update(networkEthereum.ID, model.NetworkUpdates{GasTokenID: &assetEthereum.ID})
	if err != nil {
		panic(err)
	}
	err = repos.Network.Update(networkFuji.ID, model.NetworkUpdates{GasTokenID: &assetAvalanche.ID})
	if err != nil {
		panic(err)
	}
	err = repos.Network.Update(networkAvalanche.ID, model.NetworkUpdates{GasTokenID: &assetAvalanche.ID})
	if err != nil {
		panic(err)
	}

	// String User
	userString, err := repos.User.Create(model.User{Type: "Internal", Status: "Internal"})
	if err != nil {
		panic(err)
	}

	// Devices, this is used in TX LEG
	/*deviceDeveloper*/
	_, err = repos.Device.Create(model.Device{LastUsedAt: time.Now(), Type: "Admin", Description: "Developer Laptop", UserID: userString.ID})
	if err != nil {
		panic(err)
	}

	// Instruments, used in TX Legs
	/*instrumentDeveloperCard*/
	_, err = repos.Instrument.Create(model.Instrument{Type: "Bank Account", Status: "Live", Network: "bankprov", PublicKey: "420481286", UserID: userString.ID})
	if err != nil {
		panic(err)
	}
	/*instrumentDeveloperWallet*/
	_, err = repos.Instrument.Create(model.Instrument{Type: "Crypto Wallet", Status: "Internal", Network: "EVM", PublicKey: stringPublicAddress, UserID: userString.ID})
	if err != nil {
		panic(err)
	}

	// Platforms, pointless
	/*platformDeveloper*/
	// _, err = repos.Platform.Create(model.Platform{Type: "Game", Status: "Verified", Name: "Nintendo", ApiKey: "Developer", Authentication: "Email"})
	// if err != nil {
	// 	panic(err)
	// }

}

func nullString(str string) sql.NullString {
	return sql.NullString{String: str, Valid: true}
}
