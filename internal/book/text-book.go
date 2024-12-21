package book

import "epub-tts/internal/str"

type Chapter struct {
	Name    string
	Content string
}

type TextBook struct {
	Chapters []Chapter
}

func TextBookFromEpub(input Epub) TextBook {
	chapters := []Chapter{}

	for _, v := range input.Sections {
		chapter := Chapter{
			Name:    str.SanitizeString(v.ID),
			Content: str.SanitizeString(str.RemoveTags(v.HtmlContent)),
		}
		chapters = append(chapters, chapter)
	}

	return TextBook{
		Chapters: chapters,
	}
}
