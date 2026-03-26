package validator

import (
	"github.com/jeremygprawira/go-link-generator/internal/pkg/generator"
)

// SnowflakeSystemCode validates whether a given system code has a valid HMAC signature.
// This is used to verify the integrity of codes generated via generator.SnowflakeID(...).Base62().AddHMAC(...)
func SnowflakeSystemCode(code string, secret string, hmacLen int) bool {
	if len(code) <= hmacLen {
		return false
	}

	idStr := code[:len(code)-hmacLen]
	expectedSignature := code[len(code)-hmacLen:]

	calculatedSignature, err := generator.HMAC(secret, idStr, hmacLen)
	if err != nil {
		return false
	}

	return expectedSignature == calculatedSignature
}
