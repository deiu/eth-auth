# JWT authorization service using Ethereum identities

A goland API that implements a challenge/response mechanism in order to 
issue JWTs using Ethereum identities. It also supports ETH to ENS 
resolution for users that have registered an ENS name.

## Install

```bash
$ go get github.com/deiu/eth-auth
```

## Starting the server

The server needs the following ENV variables to be set before running:
* `INFURA_API_URL` - [optional] the Infura API URL for ENS resolution
* `INFURA_API_KEY` - [optional] the Infura API key for ENS resolution
* `ORIGINS` - list of allowed Origins separated by space (your client app URL)
* `JWT_SECRET` - a passphrase used to encrypt the JWT data
* `LOGGING` - whether to log requests to stdout

Example:

```bash
$ export INFURA_API_URL="https://rinkeby.infura.io/v3/"
$ export INFURA_API_KEY="2702729979....de7a92e689bfff"
$ export ORIGINS="http://localhost:8888 https://example.org"
$ export JWT_SECRET="some-super-secret-passphrase"
$ export LOGGING="true"
$ go run server.go
```

If you want to test the API, check out the `test-client` directory for a quick
client app in JavaScript, which uses Metamask as the Ethereum provider.