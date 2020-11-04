package auth

import (
	"testing"

	"github.com/dgrijalva/jwt-go"
	"github.com/stretchr/testify/assert"
)

var (
	secret = "foobar"
	user   = "test.eth"
)

func TestNewJWT(t *testing.T) {
	token, err := NewJWT(user, secret)
	assert.NoError(t, err)
	assert.NotEmpty(t, token)
}

func TestGetUserFromJWT(t *testing.T) {
	tokenStr, err := NewJWT(user, secret)
	assert.NoError(t, err)

	parsed, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		return token, nil
	})
	u := GetUserFromJWT(parsed)
	assert.Equal(t, user, u)
}
