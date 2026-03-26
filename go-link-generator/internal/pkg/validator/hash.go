package validator

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
)

// HMAC validates whether a given string has a valid HMAC signature.
// This is used to verify the integrity of codes generated via generator.SnowflakeID(...).Base62().AddHMAC(...)
func HMAC(secret, code string, length int) bool {
	if len(code) <= length {
		return false
	}

	idStr := code[:len(code)-length]
	expectedSignature := code[len(code)-length:]

	hmac := hmac.New(sha256.New, []byte(secret))
	hmac.Write([]byte(idStr))
	calculatedSignature := base64.RawURLEncoding.EncodeToString(hmac.Sum(nil))[:length]

	return expectedSignature == calculatedSignature
}
