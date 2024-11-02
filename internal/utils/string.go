package utils

import (
	"strings"
	"unicode"

	"github.com/muesli/reflow/wordwrap"
)

func Wrap(s string, width int) string {
	s = strings.ReplaceAll(s, "\n", " ")
	s = strings.ReplaceAll(s, "\r", " ")
	return wordwrap.String(s, width)
}

// https://stackoverflow.com/questions/59955085/how-can-i-elliptically-truncate-text-in-golang
func EllipticalTruncate(text string, maxLen int) string {
	lastSpaceIx := maxLen
	len := 0
	for i, r := range text {
		if unicode.IsSpace(r) {
			lastSpaceIx = i
		}
		len++
		if len > maxLen {
			return text[:lastSpaceIx] + "..."
		}
	}
	return text
}
