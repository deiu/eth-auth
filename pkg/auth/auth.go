package auth

import (
	"time"

	"github.com/dgrijalva/jwt-go"
)

var (
	// JWTExpiration describes how long the token should be valid, in hours
	JWTExpiration = time.Duration(time.Hour * 48)
)

// NewJWT creates a new JWT for the specfied user, based on the
// provided secret
func NewJWT(user string, expTime int64, hexKey string) (string, error) {
	// Create token
	token := jwt.New(jwt.SigningMethodES256)

	// Set claims
	claims := token.Claims.(jwt.MapClaims)
	claims["name"] = user
	claims["exp"] = expTime

	// Generate encoded token and send it as response.
	key, err := ParseHexKey(hexKey)
	if err != nil {
		return "", err
	}
	t, err := token.SignedString(key)
	if err != nil {
		return "", err
	}
	return t, nil
}

// GetUserFromJWT checks if the auth token is valid
func GetUserFromJWT(token interface{}) string {
	user := token.(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	return claims["name"].(string)
}
