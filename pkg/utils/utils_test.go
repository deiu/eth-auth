package utils

import (
	"encoding/base64"
	"testing"

	"github.com/stretchr/testify/assert"
)

var ADDR = "0x123"

func TestEncodeString(t *testing.T) {
	encoded := EncodeString(ADDR)
	assert.Greater(t, encoded, ADDR)
	assert.NotContains(t, encoded, "0x")
}

func TestDecodeString(t *testing.T) {
	encoded := EncodeString(ADDR)
	decoded := DecodeString(encoded)
	assert.Equal(t, ADDR, decoded)
	assert.Empty(t, DecodeString("foo"))
}

func TestGenerateRandomString(t *testing.T) {
	rand := GenerateRandomString(16)
	_, err := base64.URLEncoding.DecodeString(rand)
	assert.NoError(t, err)
}
