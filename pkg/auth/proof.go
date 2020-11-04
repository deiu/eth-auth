package auth

import (
	"errors"
	"fmt"
	"time"

	"github.com/deiu/eth-auth/pkg/utils"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/patrickmn/go-cache"
)

// Challenge defines the object that makes up the challenge
type Challenge struct {
	Nonce string
	Date  time.Time
}

var (
	db *cache.Cache
)

const (
	// This should be replaced by EIP 712 in the future
	challengeText = `To prove your identity, please sign this one-time nonce: `
)

func init() {
	// Init in-memory store of pending challenges
	db = cache.New(0*time.Minute, 1*time.Minute)
}

// NewChallenge prepares a new challenge message to be signed by the user
// client-side, using Web3.
func NewChallenge(address string) (string, error) {
	if len(address) == 0 {
		return "", errors.New("No address provided when generating the challenge")
	}
	nonce := utils.GenerateRandomString(16)
	challenge := Challenge{
		Nonce: nonce,
		Date:  time.Now().UTC(),
	}
	// add to pending list of challenges
	db.Set(utils.EncodeString(address), &challenge, cache.DefaultExpiration)

	message := challengeText + nonce
	return message, nil
}

// VerifyChallenge is used to check the signature in the response and return the
// address of the user who signed the message.
func VerifyChallenge(address string, signature string) (string, error) {
	if len(address) == 0 || len(signature) == 0 {
		return "", errors.New("No address or signature provided when attempting to prove the challenge")
	}

	// Fetch pending challenge from db
	data, found := db.Get(utils.EncodeString(address))
	if found == false {
		return "", errors.New("Could not find a matching challenge for address: " +
			address)
	}
	defer db.Delete(utils.EncodeString(address))
	challenge := data.(*Challenge)

	// Verify signature
	fromAddr := common.HexToAddress(address)
	// We need to transform the original signature first
	sig := hexutil.MustDecode(signature)
	// https://github.com/ethereum/go-ethereum/blob/55599ee95d4151a2502465e0afc7c47bd1acba77/internal/ethapi/api.go#L442
	if sig[64] != 27 && sig[64] != 28 {
		return "", errors.New("invalid Ethereum signature")
	}
	sig[64] -= 27 // Transform yellow paper V from 27/28 to 0/1

	msg := []byte(challengeText + challenge.Nonce)
	// Recover public key from signature
	pubKey, err := crypto.SigToPub(signHash(msg), sig)
	if err != nil {
		return "", err
	}
	// Recover address from public key
	recoveredAddr := crypto.PubkeyToAddress(*pubKey)

	if fromAddr != recoveredAddr {
		return "", errors.New("Signature address does not match. " + fromAddr.String() +
			" (expected) != " + recoveredAddr.String() + "(recovered)")
	}

	return recoveredAddr.String(), nil
}

// This gives context to the signed message and prevents\
// signing of transactions.
func signHash(data []byte) []byte {
	msg := fmt.Sprintf("\x19Ethereum Signed Message:\n%d%s", len(data), data)
	return crypto.Keccak256([]byte(msg))
}

// DB returns the cache DB object (only used for testing)
func DB() *cache.Cache {
	return db
}
