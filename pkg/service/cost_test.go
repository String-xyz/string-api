package service

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetTokenPrices(t *testing.T) {
	keyMap, err := GetCoingeckoCoinMapping()
	assert.NoError(t, err)
	name, ok := keyMap[CoinKey{ChainId: 43114, Address: "0xb97ef9ef8734c71904d8002f8b6bc66dd9c48a6e"}.String()]
	assert.True(t, ok)
	assert.Equal(t, "usd-coin", name)
}

func TestGetTokenData(t *testing.T) {
	coin, err := GetCoingeckoCoinData("defi-kingdoms")
	assert.NoError(t, err)
	fmt.Printf("%+v", coin)
}
