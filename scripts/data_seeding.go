package scripts

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/String-xyz/string-api/api"
	"github.com/String-xyz/string-api/pkg/model"
	"github.com/String-xyz/string-api/pkg/store"
	"github.com/joho/godotenv"
	"github.com/rs/zerolog"
)

// Set this!
const stringPublicAddress = "0x44A4b9E2A69d86BA382a511f845CbF2E31286771"

func DataSeeding() {
	// Initialize repos
	godotenv.Load(".env") // removed the err since in cloud this wont be loaded
	port := os.Getenv("PORT")
	if port == "" {
		panic("no port!")
	}

	lg := zerolog.New(os.Stdout)

	// Note: This will panic if the env is set to use docker and you run this script from the command line
	config := api.APIConfig{
		DB:     store.MustNewPG(),
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
	networkFuji, err := repos.Network.Create(model.Network{Name: "Fuji Testnet", NetworkID: 1, ChainID: 43113, GasOracle: "avax", RPCUrl: "https://api.avax-test.network/ext/bc/C/rpc", ExplorerUrl: "https://testnet.snowtrace.io"})
	if err != nil {
		panic(err)
	}
	networkAvalanche, err := repos.Network.Create(model.Network{Name: "Avalanche Mainnet", NetworkID: 1, ChainID: 43114, GasOracle: "avax", RPCUrl: "https://api.avax.network/ext/bc/C/rpc", ExplorerUrl: "https://snowtrace.io"})
	if err != nil {
		panic(err)
	}
	networkNitroGoerli, err := repos.Network.Create(model.Network{Name: "Arbitrum Nova Testnet", NetworkID: 421613, ChainID: 421613, GasOracle: "arb", RPCUrl: "https://goerli-rollup.arbitrum.io/rpc", ExplorerUrl: "https://goerli.arbiscan.io/"})
	if err != nil {
		panic(err)
	}
	networkArbitrumNova, err := repos.Network.Create(model.Network{Name: "Arbitrum Nova Mainnet", NetworkID: 42170, ChainID: 42170, GasOracle: "arb", RPCUrl: "https://nova.arbitrum.io/rpc", ExplorerUrl: "https://nova-explorer.arbitrum.io/"})
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
	assetGoerliEth, err := repos.Asset.Create(model.Asset{Name: "GOERLIETH", Description: "Goerli Ethereum", Decimals: 18, IsCrypto: true, NetworkID: nullString(networkNitroGoerli.ID), ValueOracle: nullString("ethereum")})
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
	err = repos.Network.Update(networkNitroGoerli.ID, model.NetworkUpdates{GasTokenID: &assetGoerliEth.ID})
	if err != nil {
		panic(err)
	}
	err = repos.Network.Update(networkArbitrumNova.ID, model.NetworkUpdates{GasTokenID: &assetEthereum.ID})
	if err != nil {
		panic(err)
	}

	// String User
	userString, err := repos.User.Create(model.User{Type: "Internal", Status: "Internal"})
	if err != nil {
		panic(err)
	}

	// Set String User ID to what's defined in the ENV
	internalId := os.Getenv("STRING_INTERNAL_ID")
	if internalId == "" {
		panic("STRING_INTERNAL_ID is not set in ENV!")
	}

	type UpdateID struct {
		ID string `json:"id" db:"id"`
	}

	updateId := UpdateID{ID: internalId}
	userString, err = repos.User.Update(userString.ID, updateId)
	if err != nil {
		panic(err)
	}
	// Instruments, used in TX Legs
	/*instrumentDeveloperCard*/
	bankString, err := repos.Instrument.Create(model.Instrument{Type: "Bank Account", Status: "Live", Network: "bankprov", PublicKey: "420481286", UserID: userString.ID})
	if err != nil {
		panic(err)
	}

	bankId := os.Getenv("STRING_BANK_ID")
	if bankId == "" {
		panic("STRING_BANK_ID is not set in ENV!")
	}

	updateId = UpdateID{ID: bankId}
	err = repos.Instrument.Update(bankString.ID, updateId)
	if err != nil {
		panic(err)
	}

	/*instrumentDeveloperWallet*/
	walletString, err := repos.Instrument.Create(model.Instrument{Type: "Crypto Wallet", Status: "Internal", Network: "EVM", PublicKey: stringPublicAddress, UserID: userString.ID})
	if err != nil {
		panic(err)
	}

	walletId := os.Getenv("STRING_WALLET_ID")
	if bankId == "" {
		panic("STRING_WALLET_ID is not set in ENV!")
	}

	updateId = UpdateID{ID: walletId}
	err = repos.Instrument.Update(walletString.ID, updateId)
	if err != nil {
		panic(err)
	}

	// Platforms, placeholder
	/*platformDeveloper*/
	placeholderPlatform, err := repos.Platform.Create(model.Platform{Type: "Game", Status: "Verified", Name: "Nintendo", ApiKey: "Internal", Authentication: "Email"})
	if err != nil {
		panic(err)
	}

	platformId := os.Getenv("STRING_PLACEHOLDER_PLATFORM_ID")
	if bankId == "" {
		panic("STRING_PLACEHOLDER_PLATFORM_ID is not set in ENV!")
	}

	updateId = UpdateID{ID: platformId}
	err = repos.Platform.Update(placeholderPlatform.ID, updateId)
	if err != nil {
		panic(err)
	}
}

func MockSeeding() {
	// Initialize repos
	godotenv.Load(".env") // removed the err since in cloud this wont be loaded
	port := os.Getenv("PORT")
	if port == "" {
		panic("no port!")
	}
	lg := zerolog.New(os.Stdout)

	// Note: This will panic if the env is set to use docker and you run this script from the command line
	config := api.APIConfig{
		DB:     store.MustNewPG(),
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
	networkFuji, err := repos.Network.Create(model.Network{Name: "Fuji Testnet", NetworkID: 1, ChainID: 43113, GasOracle: "avax", RPCUrl: "https://api.avax-test.network/ext/bc/C/rpc", ExplorerUrl: "https://testnet.snowtrace.io"})
	if err != nil {
		panic(err)
	}
	networkAvalanche, err := repos.Network.Create(model.Network{Name: "Avalanche Mainnet", NetworkID: 1, ChainID: 43114, GasOracle: "avax", RPCUrl: "https://api.avax.network/ext/bc/C/rpc", ExplorerUrl: "https://snowtrace.io"})
	if err != nil {
		panic(err)
	}
	networkNitroGoerli, err := repos.Network.Create(model.Network{Name: "Arbitrum Nova Testnet", NetworkID: 421613, ChainID: 421613, GasOracle: "arb", RPCUrl: "https://goerli-rollup.arbitrum.io/rpc", ExplorerUrl: "https://goerli.arbiscan.io/"})
	if err != nil {
		panic(err)
	}
	networkArbitrumNova, err := repos.Network.Create(model.Network{Name: "Arbitrum Nova Mainnet", NetworkID: 42170, ChainID: 42170, GasOracle: "arb", RPCUrl: "https://nova.arbitrum.io/rpc", ExplorerUrl: "https://nova-explorer.arbitrum.io/"})
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
	assetGoerliEth, err := repos.Asset.Create(model.Asset{Name: "GOERLIETH", Description: "Goerli Ethereum", Decimals: 18, IsCrypto: true, NetworkID: nullString(networkNitroGoerli.ID), ValueOracle: nullString("ethereum")})
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
	err = repos.Network.Update(networkNitroGoerli.ID, model.NetworkUpdates{GasTokenID: &assetGoerliEth.ID})
	if err != nil {
		panic(err)
	}
	err = repos.Network.Update(networkArbitrumNova.ID, model.NetworkUpdates{GasTokenID: &assetEthereum.ID})
	if err != nil {
		panic(err)
	}

	// String User
	userString, err := repos.User.Create(model.User{Type: "Internal", Status: "Internal"})
	if err != nil {
		panic(err)
	}

	// Set String User ID to what's defined in the ENV
	internalId := os.Getenv("STRING_INTERNAL_ID")
	if internalId == "" {
		panic("STRING_INTERNAL_ID is not set in ENV!")
	}

	type UpdateID struct {
		ID string `json:"id" db:"id"`
	}

	updateId := UpdateID{ID: internalId}
	userString, err = repos.User.Update(userString.ID, updateId)
	if err != nil {
		panic(err)
	}

	// Devices, this is used in TX LEG
	/*deviceDeveloper*/

	// Instruments, used in TX Legs
	/*instrumentDeveloperCard*/
	bankString, err := repos.Instrument.Create(model.Instrument{Type: "Bank Account", Status: "Live", Network: "bankprov", PublicKey: "420481286", UserID: userString.ID})
	if err != nil {
		panic(err)
	}

	bankId := os.Getenv("STRING_BANK_ID")
	if bankId == "" {
		panic("STRING_BANK_ID is not set in ENV!")
	}

	updateId = UpdateID{ID: bankId}
	err = repos.Instrument.Update(bankString.ID, updateId)
	if err != nil {
		panic(err)
	}

	/*instrumentDeveloperWallet*/
	walletString, err := repos.Instrument.Create(model.Instrument{Type: "Crypto Wallet", Status: "Internal", Network: "EVM", PublicKey: stringPublicAddress, UserID: userString.ID})
	if err != nil {
		panic(err)
	}

	walletId := os.Getenv("STRING_WALLET_ID")
	if bankId == "" {
		panic("STRING_WALLET_ID is not set in ENV!")
	}

	updateId = UpdateID{ID: walletId}
	err = repos.Instrument.Update(walletString.ID, updateId)
	if err != nil {
		panic(err)
	}

	// Platforms, placeholder
	/*platformDeveloper*/
	placeholderPlatform, err := repos.Platform.Create(model.Platform{Type: "Game", Status: "Verified", Name: "Nintendo", ApiKey: "Internal", Authentication: "Email"})
	if err != nil {
		panic(err)
	}

	platformId := os.Getenv("STRING_PLACEHOLDER_PLATFORM_ID")
	if bankId == "" {
		panic("STRING_PLACEHOLDER_PLATFORM_ID is not set in ENV!")
	}

	updateId = UpdateID{ID: platformId}
	err = repos.Platform.Update(placeholderPlatform.ID, updateId)
	if err != nil {
		panic(err)
	}
}

func nullString(str string) sql.NullString {
	return sql.NullString{String: str, Valid: true}
}
