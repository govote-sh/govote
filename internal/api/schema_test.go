package api

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAddress_String(t *testing.T) {
	tests := []struct {
		name string
		addr Address
		want string
	}{
		{
			name: "all fields populated",
			addr: Address{
				LocationName: "City Hall",
				Line1:        "123 Main St",
				Line2:        "Suite 200",
				City:         "Richmond",
				State:        "VA",
				Zip:          "23220",
			},
			want: "City Hall, 123 Main St, Suite 200, Richmond, VA 23220",
		},
		{
			name: "without location name",
			addr: Address{
				Line1: "123 Main St",
				City:  "Richmond",
				State: "VA",
				Zip:   "23220",
			},
			want: "123 Main St, Richmond, VA 23220",
		},
		{
			name: "with line3",
			addr: Address{
				Line1: "123 Main St",
				Line2: "Suite 200",
				Line3: "Building A",
				City:  "Richmond",
				State: "VA",
				Zip:   "23220",
			},
			want: "123 Main St, Suite 200, Building A, Richmond, VA 23220",
		},
		{
			name: "minimal address",
			addr: Address{
				City:  "Richmond",
				State: "VA",
			},
			want: "Richmond, VA",
		},
		{
			name: "empty address",
			addr: Address{},
			want: "",
		},
		{
			name: "only zip",
			addr: Address{
				Zip: "23220",
			},
			want: "23220",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.addr.String()
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestPollingPlace_GetMapsUrl(t *testing.T) {
	tests := []struct {
		name    string
		place   PollingPlace
		want    string
		wantErr bool
	}{
		{
			name: "with valid address",
			place: PollingPlace{
				Address: Address{
					Line1: "123 Main St",
					City:  "Richmond",
					State: "VA",
					Zip:   "23220",
				},
			},
			want:    "https://www.google.com/maps/search/?api=1&query=123+Main+St%2C+Richmond%2C+VA+23220",
			wantErr: false,
		},
		{
			name: "with coordinates and no address",
			place: PollingPlace{
				Address:   Address{},
				Latitude:  37.5407,
				Longitude: -77.4360,
			},
			want:    "https://www.google.com/maps/search/?api=1&query=37.540700%2C-77.436000",
			wantErr: false,
		},
		{
			name: "with zero coordinates and no address",
			place: PollingPlace{
				Address:   Address{},
				Latitude:  0,
				Longitude: 0,
			},
			want:    "",
			wantErr: true,
		},
		{
			name: "with address and coordinates (address takes precedence)",
			place: PollingPlace{
				Address: Address{
					Line1: "123 Main St",
					City:  "Richmond",
					State: "VA",
				},
				Latitude:  37.5407,
				Longitude: -77.4360,
			},
			want:    "https://www.google.com/maps/search/?api=1&query=123+Main+St%2C+Richmond%2C+VA",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.place.GetMapsUrl()
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}

func TestPollingPlace_FilterValue(t *testing.T) {
	tests := []struct {
		name  string
		place PollingPlace
		want  string
	}{
		{
			name: "with name",
			place: PollingPlace{
				Name: "Central Voting Location",
			},
			want: "Central Voting Location",
		},
		{
			name: "without name, with location",
			place: PollingPlace{
				Address: Address{
					LocationName: "City Hall",
				},
			},
			want: "City Hall",
		},
		{
			name: "without name or location",
			place: PollingPlace{
				Address: Address{
					Line1: "123 Main St",
					City:  "Richmond",
				},
			},
			want: "123 Main St, Richmond",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.place.FilterValue()
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestContest_Title(t *testing.T) {
	tests := []struct {
		name    string
		contest Contest
		want    string
	}{
		{
			name: "short title",
			contest: Contest{
				BallotTitle: "Governor",
			},
			want: "Governor",
		},
		{
			name: "title exactly 80 chars",
			contest: Contest{
				BallotTitle: "1234567890123456789012345678901234567890123456789012345678901234567890123456789",
			},
			want: "1234567890123456789012345678901234567890123456789012345678901234567890123456789",
		},
		{
			name: "title over 80 chars",
			contest: Contest{
				BallotTitle: "This is a very long ballot title that exceeds eighty characters and should be truncated with ellipsis",
			},
			want: "This is a very long ballot title that exceeds eighty characters and should be...",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.contest.Title()
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestCandidate_Title(t *testing.T) {
	candidate := Candidate{
		Name:  "John Doe",
		Party: "Independent",
	}

	assert.Equal(t, "John Doe", candidate.Title())
}

func TestCandidate_Description(t *testing.T) {
	candidate := Candidate{
		Name:  "John Doe",
		Party: "Independent",
	}

	assert.Equal(t, "Independent", candidate.Description())
}
