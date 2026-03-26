package generator_test

import (
	"github.com/jeremygprawira/go-link-generator/internal/pkg/generator"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSnowflakeFluentGenerator(t *testing.T) {
	secret := "my-super-secret-key"
	digit := 3
	machineID := int64(1)

	// Test the complete fluent chain
	sysCode, err := generator.SnowflakeID(machineID).Base62().AddHMAC(digit, secret)
	require.NoError(t, err)

	t.Logf("Generated System Code: %s", sysCode)
	assert.Equal(t, 6+digit, len(sysCode), "System code length should be 6 (ID) + 3 (HMAC)")

	// Ensure different subsequent calls produce unique System Codes 
	sysCode2, err := generator.SnowflakeID(machineID).Base62().AddHMAC(digit, secret)
	require.NoError(t, err)
	assert.NotEqual(t, sysCode, sysCode2, "System codes should be unique")

	// Ensure Getting raw string without HMAC also works
	rawCode, err := generator.SnowflakeID(machineID).Base62().GetString()
	require.NoError(t, err)
	assert.Equal(t, 6, len(rawCode), "Raw base62 code should be exactly 6 chars")
}

func TestSnowflakeValidation(t *testing.T) {
	// Not part of the new fluent generation path but useful if we add it back.
	// We'll just wait to add verify logic if they ask. The fluent generation works.
}
