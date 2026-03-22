package address

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInputAddress_IsEmpty(t *testing.T) {
	tests := []struct {
		name     string
		addr     InputAddress
		expected bool
	}{
		{
			name:     "all fields empty",
			addr:     InputAddress{},
			expected: true,
		},
		{
			name:     "street set",
			addr:     InputAddress{Street: "123 Main St"},
			expected: false,
		},
		{
			name:     "city set",
			addr:     InputAddress{City: "Richmond"},
			expected: false,
		},
		{
			name:     "state set",
			addr:     InputAddress{State: "VA"},
			expected: false,
		},
		{
			name:     "postal code set",
			addr:     InputAddress{PostalCode: "23220"},
			expected: false,
		},
		{
			name:     "all fields set",
			addr:     InputAddress{Street: "123 Main St", City: "Richmond", State: "VA", PostalCode: "23220"},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.addr.IsEmpty())
		})
	}
}

func TestInputAddress_String(t *testing.T) {
	tests := []struct {
		name     string
		addr     InputAddress
		expected string
	}{
		{
			name:     "full address",
			addr:     InputAddress{Street: "123 Main St", City: "Richmond", State: "VA", PostalCode: "23220"},
			expected: "123 Main St, Richmond, VA, 23220",
		},
		{
			name:     "missing street",
			addr:     InputAddress{City: "Richmond", State: "VA", PostalCode: "23220"},
			expected: "Richmond, VA, 23220",
		},
		{
			name:     "single field",
			addr:     InputAddress{State: "VA"},
			expected: "VA",
		},
		{
			name:     "all empty",
			addr:     InputAddress{},
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.addr.String())
		})
	}
}
