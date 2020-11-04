package handlers

import (
	"os"

	"github.com/deiu/eth-auth/pkg/auth"
	"github.com/deiu/eth-auth/pkg/resolver"
	"github.com/gofiber/fiber"
)

// ResolveENS is used as a flag to decide where to use ENS or not
var ResolveENS bool

const (
	// ErrMessageForbidden is returned for HTTP 403 errors
	ErrMessageForbidden = "You are not allowed to access this page"
)

// Challenge handles HTTP requests for issuing a login challenge
func Challenge(c *fiber.Ctx) {
	// Bad request, wrong address length
	if len(c.Params("addr")) != 42 {
		c.SendStatus(fiber.StatusBadRequest) // HTTP 400
		c.Send(c.Params("addr") + " is not a valid Ethereum address")
		return
	}

	challenge, err := auth.NewChallenge(c.Params("addr"))
	if err != nil {
		c.SendStatus(fiber.StatusInternalServerError) // HTTP 500
		c.Send(err.Error())
	}

	c.Send(challenge)
}

// Validate handles HTTP requests for validating a login challenge
func Validate(c *fiber.Ctx) {
	// Bad request, wrong address length
	if len(c.Params("addr")) != 42 {
		c.SendStatus(fiber.StatusBadRequest) // HTTP 400
		return
	}
	// Parse request body
	params := new(struct {
		Signature string `json:"signature"`
	})
	err := c.BodyParser(&params)
	if err != nil || len(params.Signature) == 0 {
		// bad request
		c.SendStatus(fiber.StatusBadRequest) // HTTP 400
		c.Send(err.Error())
		return
	}
	// Verify challenge
	address, err := auth.VerifyChallenge(c.Params("addr"), params.Signature)
	if err != nil {
		c.SendStatus(fiber.StatusForbidden) // HTTP 403
		c.Send(err.Error())
		return
	}
	user := address
	// Resolve (reverse) addr into ENS name if we're using ENS
	if ResolveENS == true {
		ensName, err := resolver.Eth2ens(address)
		if err == nil {
			// Fallback on address if we can't resolve the name
			user = ensName
		}
	}
	// Load secret from the env
	jwtSecret := ""
	if os.Getenv("JWT_SECRET") != "" {
		jwtSecret = os.Getenv("JWT_SECRET")
	}
	token, err := auth.NewJWT(user, jwtSecret)
	if err != nil {
		c.SendStatus(fiber.StatusInternalServerError) // HTTP 500
		c.Send(err.Error())
		return
	}
	// Respond
	response := fiber.Map{}
	response["token"] = token
	c.JSON(response)
}

// Refresh handles HTTP requests for refreshing a JWT
func Refresh(c *fiber.Ctx) {
	// Check if we got a valid JTW and extract the user
	user := auth.GetUserFromJWT(c.Locals("user"))
	if len(user) == 0 {
		c.SendStatus(fiber.StatusForbidden) // HTTP 403
		c.Send(ErrMessageForbidden)
		return
	}
	// Load secret from the env
	jwtSecret := ""
	if os.Getenv("JWT_SECRET") != "" {
		jwtSecret = os.Getenv("JWT_SECRET")
	}
	token, err := auth.NewJWT(user, jwtSecret)
	if err != nil {
		c.SendStatus(fiber.StatusInternalServerError) // HTTP 500
		c.Send(err.Error())
		return
	}
	// Respond
	response := fiber.Map{}
	response["token"] = token
	c.JSON(response)
}
