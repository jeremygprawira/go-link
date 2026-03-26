package generator

import (
	"crypto/rand"
	"encoding/binary"

	"github.com/jeremygprawira/go-link-generator/internal/pkg/encoder"
)

type RandomResult struct {
	number uint64
	length int
	err    error
}

// Random generates a cryptographically secure random number and returns a result that can be encoded.
// length: The desired output length when encoded.
func Random(length int) *RandomResult {
	// Grab 8 bytes of "digital noise" (cryptographically secure random bytes).
	var raw [8]byte
	_, err := rand.Read(raw[:])

	return &RandomResult{
		// Combine those 8 random bytes into one giant 64-bit number.
		number: binary.BigEndian.Uint64(raw[:]),
		length: length,
		err:    err,
	}
}

// Encode starts the encoding process for the random number, passing along any previous error.
func (r *RandomResult) Encode() *encoder.Builder {
	return encoder.NewBuilder(r.number, r.length, r.err)
}
