# JWT authorization service using Ethereum / ENS identities

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
* `ETH_PRIVKEY` - a private Ethereum key that is used to sign the JWTs (in hexa)
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

**Important!** The Infura API URL should match the network used by the client
app. In the example above, it is set to `Rinkeby`, which means that the client
app should make sure to tell users to set their MetaMask (or similar providers)
current network to `Rinkeby` as well.

## How to use the API

### Obtaining the challenge
Basically, the way the API works is that a client will send a `GET` request to
`/login/{ethAddress}` in order to obtain a unique challenge from the server
that will then be presented to the user in order to be signed with the user's
Ethereum key.

Replace `http://api.example.org` with your own domain.

Request:
```bash
curl 'http://api.example.org/login/0x91ff16a5ffb07e2f58600afc6ff9c1c32ded1f81'
```

Response:
```js
{
  address: "0x91ff16a5ffb07e2f58600afc6ff9c1c32ded1f81",
  challenge: "To prove your identity, please sign this one-time nonce: JiqPLBbLBdCfWZoS"
}
```

### Obtaining the JWT

Next, the client must send a `POST` request to `/login/{ethAddress}`, containing the
signed challenge. The API then will validate the user's signature and issue the JWT
if the signature is good, together with the token's expiration time.

Request:
```bash
curl 'http://localhost:3000/login/0x91ff16a5ffb07e2f58600afc6ff9c1c32ded1f81' \
  -X POST \
  -H 'Content-Type: application/json' \
  --data-binary '{"signature": "0x5114fb7...33f5c031c"}'
```

Response:
```js
{
  expires: "2020-11-06T15:06:38.602022706Z",
  token: "eyJleHAiOjE....fyHd7kPlg",
  user: "0x350F72a69D....67C2EBE98dA"
}
```

**Note:** if you have registered an ENS name for your Ethereum address, the `user`
attribute will return the ENS name instead of a plain Ethereum address.

### Refreshing a JWT before expiration

Clients can obtain a new token as they get closer to the expiration time,
by sending a `GET` request to `/refresh` using the (still valid) JWT as a `Bearer`
token within an `Authorization` header. The response is similar to the one above,
containing a new expiration date and token:

Request:
```bash
curl 'http://localhost:3000/refresh' \
  -H 'Authorization: Bearer eyJleHAiOjE....fyHd7kPlg'
```

Response:
```js
{
  expires: "2020-11-06T15:06:38.602022706Z",
  token: "eyJleHAiOjE2MD....fY1qv8Oxjw",
  user: "0x350F72a69D....67C2EBE98dA"
}
```

If you want to see a full example check out the `test-client` directory from this
repo for a small client app in JavaScript, which uses Metamask as the Ethereum provider.



