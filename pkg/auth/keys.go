package auth

import (
	"crypto/ecdsa"
	"encoding/hex"
	"log"

	"github.com/ethereum/go-ethereum/crypto"
)

// EthKey represents the Ethereum private key in both ECDSA and hexa
type EthKey struct {
	PrivateRawKey *ecdsa.PrivateKey
	PrivateHexKey string
}

// NewHexKey creates a random eth address for signing JWT
func NewHexKey() (*EthKey, error) {
	key := &EthKey{}
	privateKey, err := crypto.GenerateKey()
	if err != nil {
		return key, err
	}

	key.PrivateRawKey = privateKey
	key.PrivateHexKey = hex.EncodeToString(crypto.FromECDSA(privateKey))
	return key, nil
}

// LoadHexKey parses a key in hexa (string) into an ECDSA struct
func LoadHexKey(hexKey string) *EthKey {
	privateKey, err := crypto.HexToECDSA(hexKey)
	if err != nil {
		log.Fatal(err)
	}

	return &EthKey{
		PrivateRawKey: privateKey,
		PrivateHexKey: hexKey,
	}
}
