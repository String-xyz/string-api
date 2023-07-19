package service

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetGasPrices(t *testing.T) {
	gas, tip, err := GetGasRate(Chain{RPC: "https://api.avax.network/ext/bc/C/rpc"}, 10, 10)
	assert.NoError(t, err)
	assert.Greater(t, gas.Int64(), int64(0))
	assert.Greater(t, tip.Int64(), int64(0))
	fmt.Printf("gas: %d, tip: %d", gas.Int64(), tip.Int64())

}
