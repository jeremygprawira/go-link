package generator

import (
	"fmt"
	"sync"
	"time"

	"github.com/jeremygprawira/go-link-generator/internal/pkg/encoder"
)

type Snowflake struct {
	mu        sync.Mutex
	lastStamp int64
	sequence  int64
	machineID int64
}

const (
	epoch        = int64(1700000000) // custom epoch (SECONDS) - approx Nov 2023
	machineBits  = 2
	sequenceBits = 4

	maxSequence = (1 << sequenceBits) - 1
)

var defaultSnowflake Snowflake

// SnowflakeID generates a Snowflake ID based on the provided machine ID.
// It returns a SnowflakeBuilder to support fluent method chaining.
func SnowflakeID(machineId int64) *SnowflakeBuilder {
	defaultSnowflake.mu.Lock()
	defer defaultSnowflake.mu.Unlock()

	defaultSnowflake.machineID = machineId & ((1 << machineBits) - 1)

	timestamp := time.Now().Unix()
	if timestamp < defaultSnowflake.lastStamp {
		return &SnowflakeBuilder{err: fmt.Errorf("clock moved backward")}
	}

	if timestamp == defaultSnowflake.lastStamp {
		defaultSnowflake.sequence = (defaultSnowflake.sequence + 1) & maxSequence
		if defaultSnowflake.sequence == 0 {
			timestamp = tilNextSecond(defaultSnowflake.lastStamp)
		}
	} else {
		defaultSnowflake.sequence = 0
	}

	defaultSnowflake.lastStamp = timestamp

	id := ((timestamp - epoch) << (machineBits + sequenceBits)) | (defaultSnowflake.machineID << sequenceBits) | defaultSnowflake.sequence

	return &SnowflakeBuilder{id: id}
}

func tilNextSecond(lastStamp int64) int64 {
	timestamp := time.Now().Unix()
	for timestamp <= lastStamp {
		time.Sleep(10 * time.Millisecond)
		timestamp = time.Now().Unix()
	}
	return timestamp
}

// ----------------------------------------------------------------------------
// Fluent API Builders
// ----------------------------------------------------------------------------

// SnowflakeBuilder represents an intermediate step holding the generated ID.
type SnowflakeBuilder struct {
	id  int64
	err error
}

// Base62 encodes the generated ID into a base62 string.
// Returns a SnowflakeEncoded which continues the fluent chain.
func (b *SnowflakeBuilder) Base62() *SnowflakeEncoded {
	if b.err != nil {
		return &SnowflakeEncoded{err: b.err}
	}

	// Our custom 35-bit Mini-Snowflake comfortably fits into 6 Base62 characters
	idStr, err := encoder.Encode(uint64(b.id), 6).Base62()
	return &SnowflakeEncoded{
		idStr: idStr,
		err:   err,
	}
}

// SnowflakeEncoded represents the intermediate step holding the Base62 encoded ID.
type SnowflakeEncoded struct {
	idStr string
	err   error
}

// AddHMAC appends an HMAC signature of 'digit' length using the provided 'secret'.
func (e *SnowflakeEncoded) AddHMAC(digit int, secret string) (string, error) {
	if e.err != nil {
		return "", e.err
	}

	signature, err := HMAC(secret, e.idStr, digit)
	if err != nil {
		return "", fmt.Errorf("failed to encode hmac: %w", err)
	}

	return e.idStr + signature, nil
}

// GetString allows retrieving the standalone Base62 generated string if HMAC is not needed.
func (e *SnowflakeEncoded) GetString() (string, error) {
	return e.idStr, e.err
}
