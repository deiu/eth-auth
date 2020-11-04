package resolver

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	ADDR   = "0x350F72a69D5673e9E180Ca04665D367C2EBE98dA"
	ENS    = "sambratest.eth"
	APIKEY = "831a5442dc2e4536a9f8dee4ea1707a6"
	APIURL = "https://rinkeby.infura.io/v3/"
)

func TestInitInfura(t *testing.T) {
	err := InitInfura("", "")
	assert.Error(t, err)
	err = InitInfura(APIKEY, "")
	assert.Error(t, err)
	err = InitInfura(APIURL, APIKEY)
	assert.NoError(t, err)
}

func TestEns2eth(t *testing.T) {
	addr, err := Ens2eth("")
	assert.Error(t, err)
	assert.Empty(t, addr)

	addr, err = Ens2eth(ENS)
	assert.NoError(t, err)
	assert.NotEmpty(t, addr)
	assert.Equal(t, ADDR, addr)
}

func TestEth2ens(t *testing.T) {
	name, err := Eth2ens("")
	assert.Error(t, err)
	assert.Empty(t, name)

	name, err = Eth2ens("foo")
	assert.Error(t, err)
	assert.Empty(t, name)

	name, err = Eth2ens(ADDR)
	assert.NoError(t, err)
	assert.NotEmpty(t, name)
	assert.Equal(t, ENS, name)
}
