package service

import (
	"testing"

	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
)

func TestSendSMS(t *testing.T) {
	err := godotenv.Load("../../.env")
	assert.NoError(t, err)
	err = MessageStaff("This is a test of the String Messaging Service!")
	assert.NoError(t, err)
}
