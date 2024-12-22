package book

import (
	"epub-tts/internal/str"
	"strings"
)

type Chapter struct {
	ID      string
	Name    string
	Content string
}

func (c Chapter) NameOrID() string {
	if c.Name == "" {
		return c.ID
	}
	return c.Name
}

type TextBook struct {
	Chapters []Chapter
}

func TextBookFromEpub(input Epub) TextBook {
	chapters := []Chapter{}

	for _, v := range input.Sections {
		name := str.SanitizeString(v.Title)
		name = strings.ReplaceAll(name, "\n", "")
		chapter := Chapter{
			ID:      str.SanitizeString(v.ID),
			Name:    name,
			Content: str.SanitizeString(str.RemoveTags(v.HtmlContent)),
		}
		chapters = append(chapters, chapter)
	}

	return TextBook{
		Chapters: chapters,
	}
}
