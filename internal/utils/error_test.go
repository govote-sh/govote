package utils

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestErrMsg_Error(t *testing.T) {
	tests := []struct {
		name     string
		msg      ErrMsg
		expected string
	}{
		{
			name:     "with HTTP status code",
			msg:      ErrMsg{Err: errors.New("not found"), HTTPStatusCode: 404},
			expected: "HTTP 404: not found",
		},
		{
			name:     "with zero status code",
			msg:      ErrMsg{Err: errors.New("something went wrong"), HTTPStatusCode: 0},
			expected: "something went wrong",
		},
		{
			name:     "server error",
			msg:      ErrMsg{Err: errors.New("internal server error"), HTTPStatusCode: 500},
			expected: "HTTP 500: internal server error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.msg.Error())
		})
	}
}

func TestErrMsg_SatisfiesErrorInterface(t *testing.T) {
	var err error = ErrMsg{Err: errors.New("test"), HTTPStatusCode: 0}
	assert.NotNil(t, err)
	assert.Equal(t, "test", err.Error())
}
