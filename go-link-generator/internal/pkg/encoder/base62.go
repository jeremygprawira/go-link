package encoder

// charset is the base62 alphabet — order matters, index = encoded value
// digits first, then uppercase, then lowercase
const charset = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
const base = uint64(len(charset)) // 62

type Builder struct {
	number uint64
	length int
	err    error
}

// Encode is the standalone entry point to format ANY number into a string.
// Use this if you already have a number (like a database ID) and just want to encode it.
// Example: code, err := encoder.Encode(12345, 5).Base62()
func Encode(number uint64, length int) *Builder {
	return &Builder{
		number: number,
		length: length,
	}
}

// NewBuilder creates a builder that knows how to format a number into a string.
func NewBuilder(number uint64, length int, err error) *Builder {
	return &Builder{
		number: number,
		length: length,
		err:    err,
	}
}

// Base62 generates a short ID string from the random number using the Base-62 math logic.
// This is the final step in the generator.Random(length).Encode().Base62() chain.
func (b *Builder) Base62() (string, error) {
	if b.err != nil {
		return "", b.err
	}

	// Create a placeholder (buffer) for the short string we're about to build.
	buf := make([]byte, b.length)
	number := b.number

	// The Core Logic - Translate the giant number into Base62 characters.
	// We work from right to left (filling the last character of the ID first).
	for i := b.length - 1; i >= 0; i-- {
		// Use the "remainder" of dividing the number by 62 to pick a character from our list.
		// If the remainder is 0, we pick '0'; if it's 10, we pick 'A'; if it's 61, we pick 'z'.
		buf[i] = charset[number%base]

		// Divide the number by 62 to "shift" it over, and repeat to find the next character.
		number /= base
	}

	// If the random number was too small, the ID is padded with '0's at the start.
	// If it was too large, the leftmost bits are dropped.
	// This ensures the final ID is always exactly the length we asked for.
	return string(buf), nil
}
