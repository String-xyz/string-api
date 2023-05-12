package common

import (
	"testing"

	libcommon "github.com/String-xyz/go-lib/v2/common"
	"github.com/String-xyz/string-api/pkg/model"
	"github.com/stretchr/testify/assert"
)

func TestRecoverSignature(t *testing.T) {
	addr, err := RecoverAddress("I_yuW7V0SrAWPdYyMqEsF",
		"0x6838f8b71e879e48cfbe62db6a510bb7600a58e3ff389784a94bcd85bd4cebdc432ba6ab9af2f5676d0fd84c27372f2ece60f486ec8c5b23b04bccc14a3c3a4d1c")
	assert.NoError(t, err)
	assert.Equal(t, "0x8db97C7cEcE249c2b98bDC0226Cc4C2A57BF52FC", addr.Hex())
}

// TODO: This test should be moved to the go-lib repo
func TestKeysAndValues(t *testing.T) {
	mType := "type"
	m := model.ContactUpdates{Type: &mType}
	names, vals := libcommon.KeysAndValues(m)
	assert.Len(t, names, 1)
	assert.Len(t, vals, 1)
}
