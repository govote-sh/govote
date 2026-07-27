package api

import (
	"net/url"
	"testing"
)

// An address whose only street component is Line2/Line3 used to run together
// with everything after it ("Apt 4Richmond").
func TestAddressString(t *testing.T) {
	tests := []struct {
		name string
		addr Address
		want string
	}{
		{
			name: "line2 only, then city",
			addr: Address{Line2: "Apt 4", City: "Richmond"},
			want: "Apt 4, Richmond",
		},
		{
			name: "line3 only, then city and state",
			addr: Address{Line3: "Rear Entrance", City: "Richmond", State: "VA"},
			want: "Rear Entrance, Richmond, VA",
		},
		{
			name: "line2 and line3 with no location name or line1",
			addr: Address{Line2: "Apt 4", Line3: "Rear Entrance", City: "Richmond", State: "VA", Zip: "23220"},
			want: "Apt 4, Rear Entrance, Richmond, VA 23220",
		},
		{
			name: "all fields",
			addr: Address{
				LocationName: "Main St Community Center",
				Line1:        "100 Main St",
				Line2:        "Apt 4",
				Line3:        "Rear Entrance",
				City:         "Richmond",
				State:        "VA",
				Zip:          "23220",
			},
			want: "Main St Community Center, 100 Main St, Apt 4, Rear Entrance, Richmond, VA 23220",
		},
		{
			name: "typical polling place",
			addr: Address{
				LocationName: "Main St Community Center",
				Line1:        "100 Main St",
				City:         "Richmond",
				State:        "VA",
				Zip:          "23220",
			},
			want: "Main St Community Center, 100 Main St, Richmond, VA 23220",
		},
		{
			name: "state and zip join with a space, not a comma",
			addr: Address{City: "Richmond", State: "VA", Zip: "23220"},
			want: "Richmond, VA 23220",
		},
		{
			name: "zip without state",
			addr: Address{City: "Richmond", Zip: "23220"},
			want: "Richmond, 23220",
		},
		{
			name: "state without zip",
			addr: Address{City: "Richmond", State: "VA"},
			want: "Richmond, VA",
		},
		{
			name: "location name only",
			addr: Address{LocationName: "Main St Community Center"},
			want: "Main St Community Center",
		},
		{
			name: "zip only",
			addr: Address{Zip: "23220"},
			want: "23220",
		},
		{
			name: "empty address",
			addr: Address{},
			want: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.addr.String(); got != tt.want {
				t.Errorf("Address.String() = %q, want %q", got, tt.want)
			}
		})
	}
}

// The maps link is where a bad address string actually costs a voter.
func TestPollingPlaceGetMapsUrl(t *testing.T) {
	const mapsPrefix = "https://www.google.com/maps/search/?api=1&query="

	t.Run("address with only line2 is separated in the query", func(t *testing.T) {
		p := PollingPlace{Address: Address{Line2: "Apt 4", City: "Richmond", State: "VA", Zip: "23220"}}

		got, err := p.GetMapsUrl()
		if err != nil {
			t.Fatalf("GetMapsUrl() error = %v", err)
		}
		want := mapsPrefix + url.QueryEscape("Apt 4, Richmond, VA 23220")
		if got != want {
			t.Errorf("GetMapsUrl() = %q, want %q", got, want)
		}
	})

	t.Run("falls back to coordinates when the address is empty", func(t *testing.T) {
		p := PollingPlace{Latitude: 37.5407, Longitude: -77.4360}

		got, err := p.GetMapsUrl()
		if err != nil {
			t.Fatalf("GetMapsUrl() error = %v", err)
		}
		want := mapsPrefix + url.QueryEscape("37.540700,-77.436000")
		if got != want {
			t.Errorf("GetMapsUrl() = %q, want %q", got, want)
		}
	})

	t.Run("errors when there is neither an address nor coordinates", func(t *testing.T) {
		p := PollingPlace{}

		if _, err := p.GetMapsUrl(); err == nil {
			t.Error("GetMapsUrl() error = nil, want an error")
		}
	})
}
