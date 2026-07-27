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

// EllipticalTruncate truncates text to at most maxLen runes of the original text,
// breaking at the last space before the limit when possible, and appending "...".
// maxLen counts runes but slicing counts bytes, so truncation uses range's byte offset.
// https://stackoverflow.com/questions/59955085/how-can-i-elliptically-truncate-text-in-golang
func EllipticalTruncate(text string, maxLen int) string {
	lastSpaceIx := -1 // byte offset, -1 if no space yet
	count := 0
	for i, r := range text {
		if count >= maxLen {
			if lastSpaceIx >= 0 {
				return text[:lastSpaceIx] + "..."
			}
			return text[:i] + "..."
		}
		if unicode.IsSpace(r) {
			lastSpaceIx = i
		}
		count++
	}
	return text
}
