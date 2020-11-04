package auth

import (
	"encoding/hex"

	"testing"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/assert"
)

func TestLoadHexKey(t *testing.T) {
	key := "f1568c42e5f46532b07e09b4f53c6780b7dc5ee4f07266baec1bb99d912b9c32"
	privateKey := LoadHexKey(key)
	assert.Equal(t, privateKey.PrivateHexKey, hex.EncodeToString(crypto.FromECDSA(privateKey.PrivateRawKey)))
}
