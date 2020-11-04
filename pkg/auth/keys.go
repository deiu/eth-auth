package auth

import (
	"crypto/ecdsa"

	"github.com/ethereum/go-ethereum/crypto"
)

// ParseHexKey parses a key in hexa (string) into an ECDSA private key
func ParseHexKey(hexKey string) (*ecdsa.PrivateKey, error) {
	key, err := crypto.HexToECDSA(hexKey)
	if err != nil {
		return key, err
	}

	return key, nil
}
