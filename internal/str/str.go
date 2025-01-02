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

// SplitAfterN will find the next 'divider' after 'n' characeters in 'text'
// Ex: "Oh well. This is a small sentence. This is a bigger sentence. This is the final sentence."
// SplitAfterN(text, '.', 10)
// will return return
// ["Oh well. This is a smaller sentence.", " This is a bigger sentence.", " This is the final sentence"]
// Because that's where '.' appears after 10 characters, cyclically
// If 'divider' happens exactly at 'n', the split will also happen
func SplitAfterN(text string, divider rune, n int) []string {
	items := []string{}
	lookingForDivider := false
	substr := ""
	for i, v := range text {
		substr += string(v)
		if (i+1)%n == 0 {
			lookingForDivider = true
		}

		if lookingForDivider && v == divider { // waiting for at least n letters to pass
			lookingForDivider = false
			items = append(items, substr)
			substr = ""
		}
	}
	items = append(items, substr)

	// remove last item if string ends exactly at Nth divider
	if len(items) > 0 && items[len(items)-1] == "" {
		items = append([]string{}, items[:len(items)-1]...)
	}
	return items
}
