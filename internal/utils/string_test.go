package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWrap(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		width    int
		expected string
	}{
		{
			name:     "text within width unchanged",
			input:    "short text",
			width:    20,
			expected: "short text",
		},
		{
			name:     "text exceeding width wraps at word boundary",
			input:    "this is a longer sentence that should wrap",
			width:    20,
			expected: "this is a longer\nsentence that should\nwrap",
		},
		{
			name:     "newlines replaced with spaces before wrapping",
			input:    "hello\nworld",
			width:    50,
			expected: "hello world",
		},
		{
			name:     "carriage returns replaced with spaces",
			input:    "hello\rworld",
			width:    50,
			expected: "hello world",
		},
		{
			name:     "empty string",
			input:    "",
			width:    20,
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, Wrap(tt.input, tt.width))
		})
	}
}

func TestEllipticalTruncate(t *testing.T) {
	tests := []struct {
		name     string
		text     string
		maxLen   int
		expected string
	}{
		{
			name:     "text shorter than maxLen unchanged",
			text:     "short",
			maxLen:   10,
			expected: "short",
		},
		{
			name:     "text longer truncates at last space",
			text:     "this is a long sentence",
			maxLen:   10,
			expected: "this is a...",
		},
		{
			name:     "no spaces truncates at maxLen",
			text:     "abcdefghijklmnop",
			maxLen:   10,
			expected: "abcdefghij...",
		},
		{
			name:     "empty string",
			text:     "",
			maxLen:   10,
			expected: "",
		},
		{
			name:     "exactly at maxLen unchanged",
			text:     "exactly 10",
			maxLen:   10,
			expected: "exactly 10",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, EllipticalTruncate(tt.text, tt.maxLen))
		})
	}
}
