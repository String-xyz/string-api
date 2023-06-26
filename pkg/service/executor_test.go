package service

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func setupTest() (Executor, error) {
	// Stub out the RPC
	chain := Chain{
		ChainId:       43113,
		RPC:           "https://api.avax-test.network/ext/bc/C/rpc",
		Explorer:      "https://testnet.snowtrace.io",
		CoincapName:   "avalanche",
		CoingeckoName: "avalanche-2",
		OwlracleName:  "avax",
		StringFee:     0.03,
		UUID:          "N/A",
		GasTokenId:    "N/A",
	}
	// Dial the executor
	executor := NewExecutor()
	err := executor.Initialize(chain)
	return executor, err
}

func TestGetEventData(t *testing.T) {
	e, err := setupTest()
	assert.NoError(t, err)

	eventData, err := e.GetEventData("0xe27a6a4b4ee6cbd51242faf21044941de70f5ba65ea86673d7abde75eb6c2f56",
		"Transfer(address,address,uint256)")
	assert.NoError(t, err)
	assert.Equal(t, "0x00000000000000000000000044a4b9e2a69d86ba382a511f845cbf2e31286770", eventData[0].Topics[2].Hex())
}

func TestGetTokensTransferred(t *testing.T) {
	e, err := setupTest()
	assert.NoError(t, err)

	tokens, err := e.GetTokenIds("0xe27a6a4b4ee6cbd51242faf21044941de70f5ba65ea86673d7abde75eb6c2f56")
	assert.NoError(t, err)

	assert.Equal(t, []string{"167"}, tokens)
}

func TestFilterEventData(t *testing.T) {
	e, err := setupTest()
	assert.NoError(t, err)

	eventData, err := e.GetEventData("0xe27a6a4b4ee6cbd51242faf21044941de70f5ba65ea86673d7abde75eb6c2f56",
		"Transfer(address,address,uint256)")
	assert.NoError(t, err)

	filteredEvents := FilterEventData(eventData, []int{2}, []string{"44a4b9e2a69d86ba382a511f845cbf2e31286770"})

	assert.Equal(t, "0x00000000000000000000000044a4b9e2a69d86ba382a511f845cbf2e31286770", filteredEvents[0].Topics[2].Hex())

	filteredEvents = FilterEventData(eventData, []int{1}, []string{"44a4b9e2a69d86ba382a511f845cbf2e31286770"})

	assert.Equal(t, 0, len(filteredEvents))
}

// Can't do this without loading the env :(
// func TestForwardToken(t *testing.T) {
// 	e, err := setupTest()
// 	assert.NoError(t, err)

// 	txIds, gas, err := e.ForwardTokens("0xe27a6a4b4ee6cbd51242faf21044941de70f5ba65ea86673d7abde75eb6c2f56",
// 		"0x44A4b9E2A69d86BA382a511f845CbF2E31286770")
// 	assert.NoError(t, err)

// 	fmt.Printf("TxIds: %v\n", txIds)
// 	fmt.Printf("Gas: %v\n", gas)
// }
