package utils

import (
	"encoding/base64"
	"math/rand"
	"strings"
	"time"
)

const (
	// MessageForbidden is returned for HTTP 403 errors
	MessageForbidden = "You are not allowed to access this page"
)

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
const (
	letterIdxBits = 6                    // 6 bits to represent a letter index
	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
)

var src = rand.NewSource(time.Now().UnixNano())

// GenerateRandomString returns a URL-safe, base64 encoded
// securely generated random string.
func GenerateRandomString(n int) string {
	sb := strings.Builder{}
	sb.Grow(n)
	// A src.Int63() generates 63 random bits, enough for letterIdxMax characters!
	for i, cache, remain := n-1, src.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			sb.WriteByte(letterBytes[idx])
			i--
		}
		cache >>= letterIdxBits
		remain--
	}

	return sb.String()
}

// EncodeString makes sure the ETH address is safe for storage
func EncodeString(str string) string {
	return base64.URLEncoding.EncodeToString([]byte(str))
}

// DecodeString takes a safely encoded string and decodes it into hex string
func DecodeString(str string) string {
	dec, err := base64.URLEncoding.DecodeString(str)
	if err != nil {
		return ""
	}
	return string(dec)
}

// IndexOfString returns the index (position) of value within slice
func IndexOfString(slice []string, value string) int {
	for p, v := range slice {
		if v == value {
			return p
		}
	}
	return -1
}
