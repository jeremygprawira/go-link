package generator

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCalculateLuhnCheckDigit(t *testing.T) {
	t.Run("Calculate Check Digit", func(t *testing.T) {
		testCases := []struct {
			digits      []int
			expectedCD  int
			description string
		}{
			{
				digits:      []int{4, 5, 3, 2, 0, 1, 5, 1, 1, 2, 8, 3, 0, 3, 6},
				expectedCD:  6,
				description: "Credit card example 1",
			},
			{
				digits:      []int{6, 0, 1, 1, 5, 1, 4, 4, 3, 3, 5, 4, 6, 2, 0},
				expectedCD:  1,
				description: "Credit card example 2",
			},
		}

		for _, tc := range testCases {
			checkDigit := calculateLuhnCheckDigit(tc.digits)
			assert.Equal(t, tc.expectedCD, checkDigit, tc.description)
		}
	})
}
