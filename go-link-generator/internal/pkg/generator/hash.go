package generator

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/binary"

	"github.com/jeremygprawira/go-link-generator/internal/pkg/encoder"
)

// HMAC generates a Base62 encoded HMAC signature of the specified digit length.
func HMAC(secret string, data string, digit int) (string, error) {
	h := hmac.New(sha256.New, []byte(secret))
	h.Write([]byte(data))
	hashBytes := h.Sum(nil)

	// Take the first 8 bytes of the hash to create a numeric value
	hashNum := binary.BigEndian.Uint64(hashBytes[:8])

	// Encode the numeric hash into base62 to get the HMAC signature characters
	return encoder.Encode(hashNum, digit).Base62()
}
