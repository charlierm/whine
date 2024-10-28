package govee

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestConvertFahrenheitToCelsius(t *testing.T) {
	type TestCase struct {
		description string
		fahrenheit  float64
		expected    float64
	}
	testCases := []TestCase{
		{
			description: "NegativeTemperaturesShouldBeConvertedCorrectly",
			fahrenheit:  10.0,
			expected:    -12.2222,
		},
		{
			description: "PostiveTemperaturesShouldBeConvertedCorrectly",
			fahrenheit:  100.0,
			expected:    37.7778,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			actual := convertFahrenheitToCelsius(tc.fahrenheit)
			assert.InDelta(t, tc.expected, actual, 0.001)
		})
	}
}
