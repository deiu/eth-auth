#!/bin/sh

export INFURA_API_URL="https://rinkeby.infura.io/v3/"; \
 export INFURA_API_KEY="27027299799b4c2ca2de7a92e689b017"; \
 export ORIGINS="http://localhost:8888 http://localhost:8080"; \
 export JWT_SECRET="supersecret"; \
 export LOGGING="true"; \
 go run server.go


 # db3a5e23ea6074098a1e8b112b9e8e5996a080a809a1befd50df03a15a182e2e