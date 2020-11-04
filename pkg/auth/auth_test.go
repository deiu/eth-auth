package auth

import (
	"testing"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/stretchr/testify/assert"
)

var (
	hexKey  = "f1568c42e5f46532b07e09b4f53c6780b7dc5ee4f07266baec1bb99d912b9c32"
	user    = "test.eth"
	expTime = time.Now().Add(JWTExpiration).Unix()
)

func TestNewJWT(t *testing.T) {
	token, err := NewJWT(user, expTime, "foo")
	assert.Error(t, err)
	assert.Empty(t, token)
	token, err = NewJWT(user, expTime, hexKey)
	assert.NoError(t, err)
	assert.NotEmpty(t, token)
}

func TestToken(t *testing.T) {
	tokenString, err := NewJWT(user, expTime, hexKey)
	assert.NoError(t, err)
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		privKey, _ := ParseHexKey(hexKey)
		return &privKey.PublicKey, nil
	})
	assert.NoError(t, err)
	claims, _ := token.Claims.(jwt.MapClaims)
	assert.True(t, token.Valid)
	assert.Equal(t, user, claims["name"])
}

func TestGetUserFromJWT(t *testing.T) {
	tokenStr, err := NewJWT(user, expTime, hexKey)
	assert.NoError(t, err)

	parsed, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		return token, nil
	})
	u := GetUserFromJWT(parsed)
	assert.Equal(t, user, u)
}
