package utils

import "fmt"

// Updated ErrMsg to include HTTPStatusCode for better error context
type ErrMsg struct {
	Err            error
	HTTPStatusCode int // New field to capture the HTTP status code
}

func (e ErrMsg) Error() string {
	if e.HTTPStatusCode > 0 {
		return fmt.Sprintf("HTTP %d: %v", e.HTTPStatusCode, e.Err)
	}
	return e.Err.Error()
}
