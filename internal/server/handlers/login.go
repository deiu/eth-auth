package handlers

import (
	"fmt"
	"time"

	"github.com/deiu/eth-auth/pkg/auth"
	"github.com/deiu/eth-auth/pkg/resolver"
	"github.com/gofiber/fiber"
)

// ResolveENS is used as a flag to decide where to use ENS or not
var ResolveENS bool

// EthSignKey contains the hexadecimal representation of the Eth private
// key that is used to sign JWTs
var EthSignKey string

const (
	// ErrMessageForbidden is returned for HTTP 403 errors
	ErrMessageForbidden = "You are not allowed to access this page"
)

// Challenge handles HTTP requests for issuing a login challenge
func Challenge(c *fiber.Ctx) {
	// Bad request, wrong address length
	if len(c.Params("addr")) != 42 {
		c.SendStatus(fiber.StatusBadRequest) // HTTP 400
		c.Send("The address provided is not a valid Ethereum address")
		return
	}

	challenge, err := auth.NewChallenge(c.Params("addr"))
	if err != nil {
		c.SendStatus(fiber.StatusInternalServerError) // HTTP 500
		c.Send(err.Error())
	}
	// Respond
	response := fiber.Map{}
	response["address"] = c.Params("addr")
	response["challenge"] = challenge
	c.JSON(response)
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
		status := fiber.StatusInternalServerError
		if err.Error() == auth.ErrChallengeBadRequest {
			status = fiber.StatusBadRequest
		} else if err.Error() == auth.ErrChallengeNotFound {
			status = fiber.StatusNotFound
		}
		c.SendStatus(status) // HTTP 403
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
	expTime := time.Now().Add(auth.JWTExpiration)
	token, err := auth.NewJWT(user, expTime.Unix(), EthSignKey)
	if err != nil {
		c.SendStatus(fiber.StatusInternalServerError) // HTTP 500
		c.Send(err.Error())
		return
	}
	// Respond
	response := fiber.Map{}
	response["token"] = token
	response["expires"] = expTime.UTC()
	c.JSON(response)
}

// Refresh handles HTTP requests for refreshing a JWT
func Refresh(c *fiber.Ctx) {
	fmt.Println("BAD JWT")
	// Check if we got a valid JTW and extract the user
	user := auth.GetUserFromJWT(c.Locals("user"))
	if len(user) == 0 {
		c.SendStatus(fiber.StatusForbidden) // HTTP 403
		c.Send(ErrMessageForbidden)
		return
	}
	expTime := time.Now().Add(auth.JWTExpiration)
	token, err := auth.NewJWT(user, expTime.Unix(), EthSignKey)
	if err != nil {
		c.SendStatus(fiber.StatusInternalServerError) // HTTP 500
		c.Send(err.Error())
		return
	}
	// Respond
	response := fiber.Map{}
	response["token"] = token
	response["expires"] = expTime.UTC()
	c.JSON(response)
}
