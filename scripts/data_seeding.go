package scripts

import (
	"context"
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
	ctx := context.Background()
	// api.Start(config)

	// Write to repos

	// Networks without GasTokenId
	networkPolygon, err := repos.Network.Create(model.Network{Name: "Polygon Mainnet", NetworkId: 137, ChainId: 137, GasOracle: "poly", RPCUrl: "https://rpc-mainnet.matic.quiknode.pro", ExplorerUrl: "https://polygonscan.com"})
	if err != nil {
		fmt.Printf("%+v", err)
		return
	}
	networkMumbai, err := repos.Network.Create(model.Network{Name: "Mumbai Testnet", NetworkId: 80001, ChainId: 80001, GasOracle: "poly", RPCUrl: "https://matic-mumbai.chainstacklabs.com", ExplorerUrl: "https://mumbai.polygonscan.com"})
	if err != nil {
		panic(err)
	}
	networkGoerli, err := repos.Network.Create(model.Network{Name: "Goerli Testnet", NetworkId: 5, ChainId: 5, GasOracle: "eth", RPCUrl: "https://goerli.infura.io/v3/9aa3d95b3bc440fa88ea12eaa4456161", ExplorerUrl: "https://goerli.etherscan.io"})
	if err != nil {
		panic(err)
	}
	networkEthereum, err := repos.Network.Create(model.Network{Name: "Ethereum Mainnet", NetworkId: 1, ChainId: 1, GasOracle: "eth", RPCUrl: "https://rpc.ankr.com/eth", ExplorerUrl: "https://etherscan.io"})
	if err != nil {
		panic(err)
	}
	networkFuji, err := repos.Network.Create(model.Network{Name: "Fuji Testnet", NetworkId: 1, ChainId: 43113, GasOracle: "avax", RPCUrl: "https://api.avax-test.network/ext/bc/C/rpc", ExplorerUrl: "https://testnet.snowtrace.io"})
	if err != nil {
		panic(err)
	}
	networkAvalanche, err := repos.Network.Create(model.Network{Name: "Avalanche Mainnet", NetworkId: 1, ChainId: 43114, GasOracle: "avax", RPCUrl: "https://api.avax.network/ext/bc/C/rpc", ExplorerUrl: "https://snowtrace.io"})
	if err != nil {
		panic(err)
	}
	networkNitroGoerli, err := repos.Network.Create(model.Network{Name: "Nitro Goerli Rollup Testnet", NetworkId: 421613, ChainId: 421613, GasOracle: "arb", RPCUrl: "https://goerli-rollup.arbitrum.io/rpc", ExplorerUrl: "https://goerli.arbiscan.io"})
	if err != nil {
		panic(err)
	}
	networkArbitrumNova, err := repos.Network.Create(model.Network{Name: "Arbitrum Nova Mainnet", NetworkId: 42170, ChainId: 42170, GasOracle: "arb", RPCUrl: "https://nova.arbitrum.io/rpc", ExplorerUrl: "https://nova-explorer.arbitrum.io"})
	if err != nil {
		panic(err)
	}
	// Assets
	assetAvalanche, err := repos.Asset.Create(model.Asset{Name: "AVAX", Description: "Avalanche", Decimals: 18, IsCrypto: true, NetworkId: nullString(networkAvalanche.Id), ValueOracle: nullString("avalanche-2")})
	if err != nil {
		panic(err)
	}
	assetEthereum, err := repos.Asset.Create(model.Asset{Name: "ETH", Description: "Ethereum", Decimals: 18, IsCrypto: true, NetworkId: nullString(networkEthereum.Id), ValueOracle: nullString("ethereum")})
	if err != nil {
		panic(err)
	}
	assetMatic, err := repos.Asset.Create(model.Asset{Name: "MATIC", Description: "Matic", Decimals: 18, IsCrypto: true, NetworkId: nullString(networkPolygon.Id), ValueOracle: nullString("matic-network")})
	if err != nil {
		panic(err)
	}
	assetGoerliEth, err := repos.Asset.Create(model.Asset{Name: "GOERLIETH", Description: "Goerli Ethereum", Decimals: 18, IsCrypto: true, NetworkId: nullString(networkNitroGoerli.Id), ValueOracle: nullString("ethereum")})
	if err != nil {
		panic(err)
	}
	/*assetUSD*/
	_, err = repos.Asset.Create(model.Asset{Name: "USD", Description: "United States Dollar", Decimals: 6, IsCrypto: false})
	if err != nil {
		panic(err)
	}

	// Update Networks with GasTokenIds
	err = repos.Network.Update(ctx, networkPolygon.Id, model.NetworkUpdates{GasTokenId: &assetMatic.Id})
	if err != nil {
		panic(err)
	}
	err = repos.Network.Update(ctx, networkMumbai.Id, model.NetworkUpdates{GasTokenId: &assetMatic.Id})
	if err != nil {
		panic(err)
	}
	err = repos.Network.Update(ctx, networkGoerli.Id, model.NetworkUpdates{GasTokenId: &assetEthereum.Id})
	if err != nil {
		panic(err)
	}
	err = repos.Network.Update(ctx, networkEthereum.Id, model.NetworkUpdates{GasTokenId: &assetEthereum.Id})
	if err != nil {
		panic(err)
	}
	err = repos.Network.Update(ctx, networkFuji.Id, model.NetworkUpdates{GasTokenId: &assetAvalanche.Id})
	if err != nil {
		panic(err)
	}
	err = repos.Network.Update(ctx, networkAvalanche.Id, model.NetworkUpdates{GasTokenId: &assetAvalanche.Id})
	if err != nil {
		panic(err)
	}
	err = repos.Network.Update(ctx, networkNitroGoerli.Id, model.NetworkUpdates{GasTokenId: &assetGoerliEth.Id})
	if err != nil {
		panic(err)
	}
	err = repos.Network.Update(ctx, networkArbitrumNova.Id, model.NetworkUpdates{GasTokenId: &assetEthereum.Id})
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

	type UpdateId struct {
		Id string `json:"id" db:"id"`
	}

	updateId := UpdateId{Id: internalId}
	userString, err = repos.User.Update(ctx, userString.Id, updateId)
	if err != nil {
		panic(err)
	}
	// Instruments, used in TX Legs
	/*instrumentDeveloperCard*/
	bankString, err := repos.Instrument.Create(model.Instrument{Type: "Bank Account", Status: "Live", Network: "bankprov", PublicKey: "420481286", UserId: userString.Id})
	if err != nil {
		panic(err)
	}

	bankId := os.Getenv("STRING_BANK_ID")
	if bankId == "" {
		panic("STRING_BANK_ID is not set in ENV!")
	}

	updateId = UpdateId{Id: bankId}
	err = repos.Instrument.Update(ctx, bankString.Id, updateId)
	if err != nil {
		panic(err)
	}

	/*instrumentDeveloperWallet*/
	walletString, err := repos.Instrument.Create(model.Instrument{Type: "Crypto Wallet", Status: "Internal", Network: "EVM", PublicKey: stringPublicAddress, UserId: userString.Id})
	if err != nil {
		panic(err)
	}

	walletId := os.Getenv("STRING_WALLET_ID")
	if bankId == "" {
		panic("STRING_WALLET_ID is not set in ENV!")
	}

	updateId = UpdateId{Id: walletId}
	err = repos.Instrument.Update(ctx, walletString.Id, updateId)
	if err != nil {
		panic(err)
	}

	// Platforms, placeholder
	/*platformDeveloper*/
	placeholderPlatform, err := repos.Platform.Create(model.Platform{Name: "Nintendo", Description: "Fun"})

	if err != nil {
		panic(err)
	}

	platformId := os.Getenv("STRING_PLACEHOLDER_PLATFORM_ID")
	if bankId == "" {
		panic("STRING_PLACEHOLDER_PLATFORM_ID is not set in ENV!")
	}

	updateId = UpdateId{Id: platformId}
	err = repos.Platform.Update(ctx, placeholderPlatform.Id, updateId)
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
	ctx := context.Background()

	// api.Start(config)

	// Write to repos

	// Networks without GasTokenId
	networkPolygon, err := repos.Network.Create(model.Network{Name: "Polygon Mainnet", NetworkId: 137, ChainId: 137, GasOracle: "poly", RPCUrl: "https://rpc-mainnet.matic.quiknode.pro", ExplorerUrl: "https://polygonscan.com"})
	if err != nil {
		fmt.Printf("%+v", err)
		return
	}
	networkMumbai, err := repos.Network.Create(model.Network{Name: "Mumbai Testnet", NetworkId: 80001, ChainId: 80001, GasOracle: "poly", RPCUrl: "https://matic-mumbai.chainstacklabs.com", ExplorerUrl: "https://mumbai.polygonscan.com"})
	if err != nil {
		panic(err)
	}
	networkGoerli, err := repos.Network.Create(model.Network{Name: "Goerli Testnet", NetworkId: 5, ChainId: 5, GasOracle: "eth", RPCUrl: "https://goerli.infura.io/v3/9aa3d95b3bc440fa88ea12eaa4456161", ExplorerUrl: "https://goerli.etherscan.io"})
	if err != nil {
		panic(err)
	}
	networkEthereum, err := repos.Network.Create(model.Network{Name: "Ethereum Mainnet", NetworkId: 1, ChainId: 1, GasOracle: "eth", RPCUrl: "https://rpc.ankr.com/eth", ExplorerUrl: "https://etherscan.io"})
	if err != nil {
		panic(err)
	}
	networkFuji, err := repos.Network.Create(model.Network{Name: "Fuji Testnet", NetworkId: 1, ChainId: 43113, GasOracle: "avax", RPCUrl: "https://api.avax-test.network/ext/bc/C/rpc", ExplorerUrl: "https://testnet.snowtrace.io"})
	if err != nil {
		panic(err)
	}
	networkAvalanche, err := repos.Network.Create(model.Network{Name: "Avalanche Mainnet", NetworkId: 1, ChainId: 43114, GasOracle: "avax", RPCUrl: "https://api.avax.network/ext/bc/C/rpc", ExplorerUrl: "https://snowtrace.io"})
	if err != nil {
		panic(err)
	}
	networkNitroGoerli, err := repos.Network.Create(model.Network{Name: "Nitro Goerli Rollup Testnet", NetworkId: 421613, ChainId: 421613, GasOracle: "arb", RPCUrl: "https://goerli-rollup.arbitrum.io/rpc", ExplorerUrl: "https://goerli.arbiscan.io"})
	if err != nil {
		panic(err)
	}
	networkArbitrumNova, err := repos.Network.Create(model.Network{Name: "Arbitrum Nova Mainnet", NetworkId: 42170, ChainId: 42170, GasOracle: "arb", RPCUrl: "https://nova.arbitrum.io/rpc", ExplorerUrl: "https://nova-explorer.arbitrum.io"})
	if err != nil {
		panic(err)
	}
	// Assets
	assetAvalanche, err := repos.Asset.Create(model.Asset{Name: "AVAX", Description: "Avalanche", Decimals: 18, IsCrypto: true, NetworkId: nullString(networkAvalanche.Id), ValueOracle: nullString("avalanche-2")})
	if err != nil {
		panic(err)
	}
	assetEthereum, err := repos.Asset.Create(model.Asset{Name: "ETH", Description: "Ethereum", Decimals: 18, IsCrypto: true, NetworkId: nullString(networkEthereum.Id), ValueOracle: nullString("ethereum")})
	if err != nil {
		panic(err)
	}
	assetMatic, err := repos.Asset.Create(model.Asset{Name: "MATIC", Description: "Matic", Decimals: 18, IsCrypto: true, NetworkId: nullString(networkPolygon.Id), ValueOracle: nullString("matic-network")})
	if err != nil {
		panic(err)
	}
	assetGoerliEth, err := repos.Asset.Create(model.Asset{Name: "GOERLIETH", Description: "Goerli Ethereum", Decimals: 18, IsCrypto: true, NetworkId: nullString(networkNitroGoerli.Id), ValueOracle: nullString("ethereum")})
	if err != nil {
		panic(err)
	}
	/*assetUSD*/
	_, err = repos.Asset.Create(model.Asset{Name: "USD", Description: "United States Dollar", Decimals: 6, IsCrypto: false})
	if err != nil {
		panic(err)
	}

	// Update Networks with GasTokenIds
	err = repos.Network.Update(ctx, networkPolygon.Id, model.NetworkUpdates{GasTokenId: &assetMatic.Id})
	if err != nil {
		panic(err)
	}
	err = repos.Network.Update(ctx, networkMumbai.Id, model.NetworkUpdates{GasTokenId: &assetMatic.Id})
	if err != nil {
		panic(err)
	}
	err = repos.Network.Update(ctx, networkGoerli.Id, model.NetworkUpdates{GasTokenId: &assetEthereum.Id})
	if err != nil {
		panic(err)
	}
	err = repos.Network.Update(ctx, networkEthereum.Id, model.NetworkUpdates{GasTokenId: &assetEthereum.Id})
	if err != nil {
		panic(err)
	}
	err = repos.Network.Update(ctx, networkFuji.Id, model.NetworkUpdates{GasTokenId: &assetAvalanche.Id})
	if err != nil {
		panic(err)
	}
	err = repos.Network.Update(ctx, networkAvalanche.Id, model.NetworkUpdates{GasTokenId: &assetAvalanche.Id})
	if err != nil {
		panic(err)
	}
	err = repos.Network.Update(ctx, networkNitroGoerli.Id, model.NetworkUpdates{GasTokenId: &assetGoerliEth.Id})
	if err != nil {
		panic(err)
	}
	err = repos.Network.Update(ctx, networkArbitrumNova.Id, model.NetworkUpdates{GasTokenId: &assetEthereum.Id})
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

	type UpdateId struct {
		Id string `json:"id" db:"id"`
	}

	updateId := UpdateId{Id: internalId}
	userString, err = repos.User.Update(ctx, userString.Id, updateId)
	if err != nil {
		panic(err)
	}

	// Devices, this is used in TX LEG
	/*deviceDeveloper*/

	// Instruments, used in TX Legs
	/*instrumentDeveloperCard*/
	bankString, err := repos.Instrument.Create(model.Instrument{Type: "Bank Account", Status: "Live", Network: "bankprov", PublicKey: "420481286", UserId: userString.Id})
	if err != nil {
		panic(err)
	}

	bankId := os.Getenv("STRING_BANK_ID")
	if bankId == "" {
		panic("STRING_BANK_ID is not set in ENV!")
	}

	updateId = UpdateId{Id: bankId}
	err = repos.Instrument.Update(ctx, bankString.Id, updateId)
	if err != nil {
		panic(err)
	}

	/*instrumentDeveloperWallet*/
	walletString, err := repos.Instrument.Create(model.Instrument{Type: "Crypto Wallet", Status: "Internal", Network: "EVM", PublicKey: stringPublicAddress, UserId: userString.Id})
	if err != nil {
		panic(err)
	}

	walletId := os.Getenv("STRING_WALLET_ID")
	if bankId == "" {
		panic("STRING_WALLET_ID is not set in ENV!")
	}

	updateId = UpdateId{Id: walletId}
	err = repos.Instrument.Update(ctx, walletString.Id, updateId)
	if err != nil {
		panic(err)
	}

	// Platforms, placeholder
	/*platformDeveloper*/
	placeholderPlatform, err := repos.Platform.Create(model.Platform{Name: "Nintendo", Description: "Fun"})
	if err != nil {
		panic(err)
	}

	platformId := os.Getenv("STRING_PLACEHOLDER_PLATFORM_ID")
	if bankId == "" {
		panic("STRING_PLACEHOLDER_PLATFORM_ID is not set in ENV!")
	}

	updateId = UpdateId{Id: platformId}
	err = repos.Platform.Update(ctx, placeholderPlatform.Id, updateId)
	if err != nil {
		panic(err)
	}
}

func nullString(str string) sql.NullString {
	return sql.NullString{String: str, Valid: true}
}
