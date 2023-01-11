package handler

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSanitizeChecksumArray(t *testing.T) {
	arrayByValue := []string{"0x44a4b9E2A69d86BA382a511f845CbF2E31286770"}
	var cxParams []*string
	for _, p := range arrayByValue {
		cxParams = append(cxParams, &p)
	}
	SanitizeChecksums(cxParams...)
	assert.Equal(t, *cxParams[0], "0x44A4b9E2A69d86BA382a511f845CbF2E31286770")
}
