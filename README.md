# JWT authorization service using Ethereum identities

A goland API that implements a challenge/response mechanism in order to 
issue JWTs using Ethereum identities. It also (optionally) supports ETH to ENS 
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
* `ETH_PRIVKEY` - a private Ethereum key in hexadecimal
* `LOGGING` - whether to log requests to stdout

Example of how to start the server:

```bash
$ export INFURA_API_URL="https://rinkeby.infura.io/v3/"
$ export INFURA_API_KEY="2702729979....de7a92e689bfff"
$ export ORIGINS="http://localhost:8888 https://example.org"
$ export ETH_PRIVKEY="your-exported-eth-key-in-hexa"
$ export LOGGING="true"
$ go run server.go
```

## How to use the API

Basically, the way the API works is that a client will send a request to 
obtain a challenge from the server that will then be presented to the user 
in order to be signed with the user's Ethereum key.

Next, the API will validate the user's signature and issue the JWT if the 
signature is good, together with the token's expiration time.

Clients can obtain a new token as they get closer to the expiration time, 
by sending a GET request to `/refresh`.


If you want to see a full example check out the `test-client` directory 
for a small client app in JavaScript, which uses Metamask as the Ethereum 
provider.



