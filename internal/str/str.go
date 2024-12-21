package str

import (
	"regexp"
	"strings"
	"unicode"
)

func SanitizeString(str string) string {
	str = strings.Trim(str, "\r\n\t ")
	str = strings.ReplaceAll(str, "\r\n", "\n")

	// make lines with only spaces to be just lines so they can be grouped below
	blankLineRegex := regexp.MustCompile(`(?m)^\s*$`)
	str = blankLineRegex.ReplaceAllString(str, "\n")

	str = strings.Map(func(r rune) rune {
		if unicode.IsPrint(r) || r == '\n' {
			return r
		}
		return -1
	}, str)

	// remove excess line breaks
	for strings.Contains(str, "\n\n\n") {
		str = strings.ReplaceAll(str, "\n\n\n", "\n\n")
	}

	return str
}

func RemoveTags(input string) string {
	// Define the regex pattern to match HTML tags
	tagRegex := regexp.MustCompile(`<[^>]+>`)
	// Replace all occurrences of the tag pattern with an empty string
	return tagRegex.ReplaceAllString(input, "")
}
