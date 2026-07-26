package utils

import (
	"testing"
	"unicode/utf8"
)

func TestWrap(t *testing.T) {
	tests := []struct {
		name  string
		input string
		width int
		want  string
	}{
		{
			name:  "text within width unchanged",
			input: "short text",
			width: 20,
			want:  "short text",
		},
		{
			name:  "text exceeding width wraps at word boundary",
			input: "this is a longer sentence that should wrap",
			width: 20,
			want:  "this is a longer\nsentence that should\nwrap",
		},
		{
			name:  "newlines replaced with spaces before wrapping",
			input: "hello\nworld",
			width: 50,
			want:  "hello world",
		},
		{
			name:  "carriage returns replaced with spaces",
			input: "hello\rworld",
			width: 50,
			want:  "hello world",
		},
		{
			name:  "empty string",
			input: "",
			width: 20,
			want:  "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Wrap(tt.input, tt.width); got != tt.want {
				t.Errorf("Wrap() = %q, want %q", got, tt.want)
			}
		})
	}
}

// The multi-byte cases are the point: maxLen used to be applied as a byte
// offset, slicing through a rune and producing invalid UTF-8.
func TestEllipticalTruncate(t *testing.T) {
	tests := []struct {
		name   string
		text   string
		maxLen int
		want   string
	}{
		{
			name:   "text shorter than maxLen unchanged",
			text:   "short",
			maxLen: 10,
			want:   "short",
		},
		{
			name:   "text longer truncates at last space",
			text:   "this is a long sentence",
			maxLen: 10,
			want:   "this is a...",
		},
		{
			name:   "no spaces truncates at maxLen",
			text:   "abcdefghijklmnop",
			maxLen: 10,
			want:   "abcdefghij...",
		},
		{
			name:   "empty string",
			text:   "",
			maxLen: 10,
			want:   "",
		},
		{
			name:   "exactly at maxLen unchanged",
			text:   "exactly 10",
			maxLen: 10,
			want:   "exactly 10",
		},
		{
			name:   "multi-byte with no space cuts on a rune boundary",
			text:   "ñññññ",
			maxLen: 3,
			want:   "ñññ...",
		},
		{
			name:   "multi-byte with no space counts runes, not bytes",
			text:   "Proposición-Ñ-Presupuesto",
			maxLen: 15,
			want:   "Proposición-Ñ-P...",
		},
		{
			name:   "multi-byte still breaks at the last space",
			text:   "Ñañá es bueno",
			maxLen: 8,
			want:   "Ñañá es...",
		},
		{
			name:   "multi-byte exactly at maxLen unchanged",
			text:   "ñññ",
			maxLen: 3,
			want:   "ñññ",
		},
		{
			name:   "zero maxLen",
			text:   "hello",
			maxLen: 0,
			want:   "...",
		},
		{
			name:   "negative maxLen does not panic",
			text:   "hello",
			maxLen: -1,
			want:   "...",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := EllipticalTruncate(tt.text, tt.maxLen)
			if got != tt.want {
				t.Errorf("EllipticalTruncate() = %q, want %q", got, tt.want)
			}
			if !utf8.ValidString(got) {
				t.Errorf("EllipticalTruncate() = %q, which is not valid UTF-8", got)
			}
		})
	}
}
