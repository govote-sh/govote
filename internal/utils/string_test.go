package utils

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWrap(t *testing.T) {
	tests := []struct {
		name  string
		input string
		width int
		want  string
	}{
		{
			name:  "text within width",
			input: "Hello World",
			width: 50,
			want:  "Hello World",
		},
		{
			name:  "text exceeding width",
			input: "This is a long sentence that should be wrapped at the specified width",
			width: 20,
			want:  "This is a long\nsentence that should\nbe wrapped at the\nspecified width",
		},
		{
			name:  "text with existing newlines",
			input: "Line one\nLine two\nLine three",
			width: 50,
			want:  "Line one Line two Line three",
		},
		{
			name:  "text with carriage returns",
			input: "Line one\rLine two\rLine three",
			width: 50,
			want:  "Line one Line two Line three",
		},
		{
			name:  "text with both newlines and carriage returns",
			input: "Line\n\rMixed\r\nContent",
			width: 50,
			want:  "Line  Mixed  Content",
		},
		{
			name:  "empty string",
			input: "",
			width: 20,
			want:  "",
		},
		{
			name:  "single character",
			input: "A",
			width: 10,
			want:  "A",
		},
		{
			name:  "very narrow width",
			input: "Hello World",
			width: 5,
			want:  "Hello\nWorld",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Wrap(tt.input, tt.width)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestEllipticalTruncate(t *testing.T) {
	tests := []struct {
		name   string
		input  string
		maxLen int
		want   string
	}{
		{
			name:   "text shorter than maxLen",
			input:  "Hello",
			maxLen: 10,
			want:   "Hello",
		},
		{
			name:   "text exactly at maxLen",
			input:  "Hello World",
			maxLen: 11,
			want:   "Hello World",
		},
		{
			name:   "text longer than maxLen with spaces",
			input:  "This is a long sentence that needs truncation",
			maxLen: 20,
			want:   "This is a long...",
		},
		{
			name:   "text longer than maxLen without spaces",
			input:  "Thisisaverylongwordwithoutanyspaces",
			maxLen: 10,
			want:   "Thisisaver...",
		},
		{
			name:   "text with space at truncation point",
			input:  "Hello World and more text",
			maxLen: 11,
			want:   "Hello World...",
		},
		{
			name:   "empty string",
			input:  "",
			maxLen: 10,
			want:   "",
		},
		{
			name:   "maxLen zero",
			input:  "Hello",
			maxLen: 0,
			want:   "...",
		},
		{
			name:   "maxLen less than ellipsis length",
			input:  "Hello World",
			maxLen: 2,
			want:   "He...", // Actual behavior: truncates to maxLen chars
		},
		{
			name:   "unicode text with emoji",
			input:  "Hello 👋 World 🌍 This is a test",
			maxLen: 20,
			want:   "Hello 👋 World 🌍 This...", // Counts runes, not bytes
		},
		{
			name:   "text with accented characters",
			input:  "Héllo Wörld this is a test with áccents",
			maxLen: 20,
			want:   "Héllo Wörld this is...",
		},
		{
			name:   "text with multiple spaces",
			input:  "Hello    World    with    spaces",
			maxLen: 15,
			want:   "Hello    World ...", // Preserves consecutive spaces
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := EllipticalTruncate(tt.input, tt.maxLen)
			assert.Equal(t, tt.want, got)
		})
	}
}

// TestEllipticalTruncate_PreservesWords verifies that truncation happens at word boundaries
func TestEllipticalTruncate_PreservesWords(t *testing.T) {
	input := "The quick brown fox jumps over the lazy dog"
	result := EllipticalTruncate(input, 25)

	// Should truncate at a word boundary, not in the middle of a word
	assert.True(t, strings.HasSuffix(result, "..."))
	assert.False(t, strings.Contains(result, "jump...")) // Should not cut mid-word
}
