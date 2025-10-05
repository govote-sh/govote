package address

import "strings"

// InputAddress represents an address entered by the user
type InputAddress struct {
	Street     string
	City       string
	State      string
	PostalCode string
}

// IsEmpty returns true if all fields are empty
func (a InputAddress) IsEmpty() bool {
	return a.Street == "" && a.City == "" && a.State == "" && a.PostalCode == ""
}

// String formats the address for display, omitting empty fields
func (a InputAddress) String() string {
	var parts []string

	if a.Street != "" {
		parts = append(parts, a.Street)
	}
	if a.City != "" {
		parts = append(parts, a.City)
	}
	if a.State != "" {
		parts = append(parts, a.State)
	}
	if a.PostalCode != "" {
		parts = append(parts, a.PostalCode)
	}

	return strings.Join(parts, ", ")
}
