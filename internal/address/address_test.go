package address

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInputAddress_IsEmpty(t *testing.T) {
	tests := []struct {
		name string
		addr InputAddress
		want bool
	}{
		{
			name: "all fields empty",
			addr: InputAddress{},
			want: true,
		},
		{
			name: "only street populated",
			addr: InputAddress{Street: "123 Main St"},
			want: false,
		},
		{
			name: "only city populated",
			addr: InputAddress{City: "Richmond"},
			want: false,
		},
		{
			name: "only state populated",
			addr: InputAddress{State: "VA"},
			want: false,
		},
		{
			name: "only postal code populated",
			addr: InputAddress{PostalCode: "23220"},
			want: false,
		},
		{
			name: "all fields populated",
			addr: InputAddress{
				Street:     "123 Main St",
				City:       "Richmond",
				State:      "VA",
				PostalCode: "23220",
			},
			want: false,
		},
		{
			name: "partial address - street and city",
			addr: InputAddress{
				Street: "123 Main St",
				City:   "Richmond",
			},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.addr.IsEmpty()
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestInputAddress_String(t *testing.T) {
	tests := []struct {
		name string
		addr InputAddress
		want string
	}{
		{
			name: "all fields populated",
			addr: InputAddress{
				Street:     "123 Main St",
				City:       "Richmond",
				State:      "VA",
				PostalCode: "23220",
			},
			want: "123 Main St, Richmond, VA, 23220",
		},
		{
			name: "missing street",
			addr: InputAddress{
				City:       "Richmond",
				State:      "VA",
				PostalCode: "23220",
			},
			want: "Richmond, VA, 23220",
		},
		{
			name: "missing city",
			addr: InputAddress{
				Street:     "123 Main St",
				State:      "VA",
				PostalCode: "23220",
			},
			want: "123 Main St, VA, 23220",
		},
		{
			name: "missing state",
			addr: InputAddress{
				Street:     "123 Main St",
				City:       "Richmond",
				PostalCode: "23220",
			},
			want: "123 Main St, Richmond, 23220",
		},
		{
			name: "missing postal code",
			addr: InputAddress{
				Street: "123 Main St",
				City:   "Richmond",
				State:  "VA",
			},
			want: "123 Main St, Richmond, VA",
		},
		{
			name: "only street",
			addr: InputAddress{
				Street: "123 Main St",
			},
			want: "123 Main St",
		},
		{
			name: "only city",
			addr: InputAddress{
				City: "Richmond",
			},
			want: "Richmond",
		},
		{
			name: "only state",
			addr: InputAddress{
				State: "VA",
			},
			want: "VA",
		},
		{
			name: "only postal code",
			addr: InputAddress{
				PostalCode: "23220",
			},
			want: "23220",
		},
		{
			name: "all empty",
			addr: InputAddress{},
			want: "",
		},
		{
			name: "street and postal code only",
			addr: InputAddress{
				Street:     "123 Main St",
				PostalCode: "23220",
			},
			want: "123 Main St, 23220",
		},
		{
			name: "city and state only",
			addr: InputAddress{
				City:  "Richmond",
				State: "VA",
			},
			want: "Richmond, VA",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.addr.String()
			assert.Equal(t, tt.want, got)
		})
	}
}
