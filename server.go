package main

import (
	"os"
	"strings"

	"github.com/deiu/eth-auth/internal/server"
)

func main() {
	conf := server.Config{}
	if os.Getenv("INFURA_API_URL") != "" {
		conf.InfuraURL = os.Getenv("INFURA_API_URL")
	}
	if os.Getenv("INFURA_API_KEY") != "" {
		conf.InfuraKey = os.Getenv("INFURA_API_KEY")
	}
	if os.Getenv("ORIGINS") != "" {
		conf.Origins = strings.Fields(os.Getenv("ORIGINS"))
	}
	if os.Getenv("ETH_PRIVKEY") != "" {
		conf.JWTSecret = os.Getenv("ETH_PRIVKEY")
	}
	if os.Getenv("LOGGING") != "" {
		conf.Logging = true
	}
	server.Listen(3000, conf)
}
