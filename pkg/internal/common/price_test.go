package common

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetCoincapPrices(t *testing.T) {
	coins := []string{"matic", "ethereum", "avalanche"}
	for _, coin := range coins {
		body := make(map[string]interface{})
		err := GetJsonGeneric("https://api.coincap.io/v2/assets?search="+coin, &body)
		assert.NoError(t, err)
		price := body["data"].([]interface{})[0].(map[string]interface{})["priceUsd"].(string)
		assert.NotEqual(t, "", price)
		fmt.Printf("\n"+coin+" PRICE = %+v", price)
	}
}
