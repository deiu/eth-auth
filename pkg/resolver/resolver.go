package resolver

import (
	"errors"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	ens "github.com/wealdtech/go-ens/v3"
)

var (
	client *ethclient.Client
)

// InitInfura is used to initialize the API for ENS resolution
func InitInfura(infuraURL string, infuraKey string) error {
	if len(infuraURL) == 0 || len(infuraKey) == 0 {
		return errors.New("Infura API URL or key was not provided")
	}
	// Init client connection to Infura
	var err error
	client, err = ethclient.Dial(infuraURL + infuraKey)
	if err != nil {
		return err
	}
	return nil
}

// Ens2eth resolves an ENS name to an ETH address
func Ens2eth(domain string) (string, error) {
	if len(domain) == 0 {
		return "", errors.New("No domain was provided for ENS lookup")
	}
	address, err := ens.Resolve(client, domain)
	if err != nil {
		panic(err)
	}
	return address.Hex(), nil
}

// Eth2ens reverse resolves an ETH address to an ENS name
func Eth2ens(address string) (string, error) {
	if len(address) == 0 {
		return "", errors.New("No address was provided for ENS reverse lookup")
	}
	addr := common.HexToAddress(address)
	domain, err := ens.ReverseResolve(client, addr)
	if err != nil {
		return "", errors.New(err.Error() + " Addr: " + address)
	}

	return domain, nil
}
