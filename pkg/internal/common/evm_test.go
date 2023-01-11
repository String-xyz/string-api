package common

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestChecksumValid(t *testing.T) {
	valid := validChecksum("0x44A4b9E2A69d86BA382a511f845CbF2E31286770")
	assert.Equal(t, valid, true)
}

func TestChecksumInvalid(t *testing.T) {
	valid := validChecksum("0x44a4b9E2A69d86BA382a511f845CbF2E31286770")
	assert.Equal(t, valid, false)
}

func TestSanitizeChecksum(t *testing.T) {
	valid := SanitizeChecksum("0x44a4b9E2A69d86BA382a511f845CbF2E31286770")
	assert.Equal(t, valid, "0x44A4b9E2A69d86BA382a511f845CbF2E31286770")
}
