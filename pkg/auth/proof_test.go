package auth

import (
	"strings"
	"testing"
	"time"

	"github.com/deiu/eth-auth/pkg/utils"
	"github.com/patrickmn/go-cache"

	"github.com/stretchr/testify/assert"
)

const (
	ADDR      = "0x91ff16a5ffb07e2f58600afc6ff9c1c32ded1f81"
	BADADDR   = "0x91ff16a5ffb07e2f58600afc6ff9c1c32ded1f82"
	NONCE     = "iZMCIAnQFggYXNgG"
	SIGNATURE = "0x0ffd5de2c0cc7daf88e3c824ea40ed91c1b001c38cacdaed8b54c6c5a95c2f17108ff663599a28ae53303550d625fc0a360925a9bad207ebc203166cac74bb2d1c"
	DIFSIG    = "0x6f1c94f3d7a10b326483fa8460f1f37044b37fa5679d367896140be71ff9bac434449fa54dbbccbc5989a8c3f60219a105b15a8605a8745a5074b94d423a788d1b"
	BADSIG    = "0x910a25f93937048962410890ed991ece1b8463c2c02b26ea97c9b8caf5ee183646f90d129a64f03d4dc4d9d381d37b5a230c7e262829a15cd72585bdb9668a510b"
)

func TestNewChallenge(t *testing.T) {
	_, err := NewChallenge("")
	assert.Error(t, err)
	challenge, err := NewChallenge(ADDR)
	assert.NoError(t, err)
	assert.Contains(t, challenge, challengeText)
	data, found := db.Get(utils.EncodeString(ADDR))
	assert.True(t, found)
	ch := data.(*Challenge)
	assert.NotNil(t, ch.Nonce)
	assert.NotNil(t, ch.Date)
}

func TestVerifyChallenge(t *testing.T) {
	_, err := VerifyChallenge("", SIGNATURE)
	assert.Error(t, err)
	_, err = VerifyChallenge(ADDR, "")
	assert.Error(t, err)
	_, err = VerifyChallenge("0xabc", SIGNATURE)
	assert.Error(t, err)
	_, err = VerifyChallenge(ADDR, BADSIG)
	assert.Error(t, err)

	challenge := Challenge{
		Nonce: NONCE,
		Date:  time.Now().UTC(),
	}
	// add to pending list of challenges
	db.Set(utils.EncodeString(ADDR), &challenge, cache.DefaultExpiration)

	_, err = VerifyChallenge(ADDR, DIFSIG)
	assert.Error(t, err)

	// add to pending list of challenges
	db.Set(utils.EncodeString(ADDR), &challenge, cache.DefaultExpiration)

	addr, err := VerifyChallenge(ADDR, SIGNATURE)
	assert.NoError(t, err)
	assert.Equal(t, ADDR, strings.ToLower(addr))
}

func TestSignHash(t *testing.T) {
	text := []byte("hello")
	assert.Greater(t, len(signHash(text)), len(text))
}

func TestDB(t *testing.T) {
	assert.NotNil(t, DB())
}
