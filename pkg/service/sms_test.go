// TODO: Make this service instantiable

package service

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSendSMS(t *testing.T) {
	err := MessageTeam("This is a test of the String Messaging Service!")
	assert.NoError(t, err)
}
