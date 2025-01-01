package str

import (
	"regexp"
	"strings"
	"unicode"

	"golang.org/x/text/runes"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
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

// CleanFileName removes invalid characters for filenames
// and also removes accents and special characters.
func CleanFileName(input string) string {
	// Normalize the input string to remove accents
	normalized, err := normalize(input)
	if err != nil {
		// TODO add log
		normalized = input
	}

	normalized = strings.ReplaceAll(normalized, "—", "_")
	normalized = strings.ReplaceAll(normalized, ":", "_")

	// Define a regular expression that allows only alphanumeric characters, dashes, and underscores
	re := regexp.MustCompile(`[^a-zA-Z0-9\s\-_\.]`)

	// Remove any character that is not a word character, whitespace, dash, or period
	cleaned := re.ReplaceAllString(normalized, "")

	// Optionally replace spaces with underscores or dashes
	cleaned = strings.ReplaceAll(cleaned, " ", "_")
	cleaned = strings.ReplaceAll(cleaned, `\n`, "")
	cleaned = strings.ReplaceAll(cleaned, "\n", "")

	return cleaned
}

// https://stackoverflow.com/a/65981868
func normalize(s string) (string, error) {
	t := transform.Chain(norm.NFD, runes.Remove(runes.In(unicode.Mn)), norm.NFC)
	result, _, err := transform.String(t, s)
	if err != nil {
		return "", err
	}

	return result, nil
}

func RemoveTags(input string) string {
	// Define the regex pattern to match HTML tags
	tagRegex := regexp.MustCompile(`<[^>]+>`)
	// Replace all occurrences of the tag pattern with an empty string
	return tagRegex.ReplaceAllString(input, "")
}
