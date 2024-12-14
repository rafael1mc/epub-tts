package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
	"unicode"
)

type EpubSection struct {
	ID          string `json:"id"`
	HtmlContent string `json:"htmlString"`
}

type Chapter struct {
	Name    string
	Content string
}

const (
	perm             = 0777
	outputFolderName = "output"
	isDryRun         = false // if true, will generate text files, but not audio files
	isDebug          = false // if true, will generate files for json and html content as well
)

func main() {
	fmt.Println(" ---== Execution started ==--- ")

	createOutputDir(outputFolderName)
	epubSections := parseEpub()
	generateDebugFiles(epubSections)
	chapters := epubSectionsToCharters(epubSections)
	saveChaptersText(chapters)
	ttsChapters(chapters)

	fmt.Println(" ---== Execution ended ==--- ")
}

func parseEpub() []EpubSection {
	fmt.Println("Parsing epub")
	cmdStr := "docker run --rm -v ./volume:/app/input epub-parser input.epub"
	out, _ := exec.Command("/bin/sh", "-c", cmdStr).Output()

	result := []EpubSection{}
	err := json.Unmarshal(out, &result)
	if err != nil {
		panic(err)
	}

	if isDebug {
		writeFile(getOutputPath(0, "___debug_secions", "json"), string(out))
	}

	return result
}

func epubSectionsToCharters(sections []EpubSection) []Chapter {
	chapters := []Chapter{}

	for _, v := range sections {
		chapters = append(chapters, epubSectionToCharter(v))
	}

	return chapters
}
func epubSectionToCharter(s EpubSection) Chapter {
	return Chapter{
		Name:    sanitizeString(s.ID),
		Content: sanitizeString(removeTags(s.HtmlContent)),
	}
}

func sanitizeString(str string) string {
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

func removeTags(input string) string {
	// Define the regex pattern to match HTML tags
	tagRegex := regexp.MustCompile(`<[^>]+>`)
	// Replace all occurrences of the tag pattern with an empty string
	return tagRegex.ReplaceAllString(input, "")
}

func saveChaptersText(chapters []Chapter) {
	fmt.Println("Saving chapters as text")

	createOutputDir(outputFolderName)

	for k, v := range chapters {
		filename := txtFilename(k, v)
		writeFile(filename, v.Content)
	}
}

func ttsChapters(chapters []Chapter) {
	if isDryRun {
		return
	}
	fmt.Println("Narrating chapters")

	wg := sync.WaitGroup{}
	wg.Add(len(chapters))

	// TODO maybe use channels, but for now it's enough
	for k, v := range chapters {
		go func(pos int, chapter Chapter) {
			ttsChapter(pos, chapter)
			wg.Done()
		}(k, v)
	}

	fmt.Println("--")

	wg.Wait()
}

func ttsChapter(pos int, chapter Chapter) string {
	audioName := audioFilename(pos, chapter)
	fmt.Println("Narrating chapter:" + audioName)
	cmdStr := fmt.Sprintf(`say -f "%s" -o "%s"`, txtFilename(pos, chapter), audioName)
	out, _ := exec.Command("/bin/sh", "-c", cmdStr).Output()

	fmt.Println("Chapter", audioName, "narrated")
	return string(out)
}

func createOutputDir(folderName string) {
	err := os.MkdirAll(folderName, perm)
	if err != nil && !errors.Is(err, os.ErrExist) {
		panic(err)
	}
}

func txtFilename(pos int, chapter Chapter) string {
	return getOutputPath(pos, chapter.Name, "txt")
}

func audioFilename(pos int, chapter Chapter) string {
	return getOutputPath(pos, chapter.Name, "aiff")
}

func getOutputPath(pos int, name string, extension string) string {
	filename := fmt.Sprintf("%d.%s.%s", pos, name, extension)
	filename = strings.ReplaceAll(filename, " ", "-")
	filename = strings.ReplaceAll(filename, " ", "-")
	filename = strings.ToLower(filename)

	filePath := filepath.Join(outputFolderName, filename)

	return filePath
}

func writeFile(name string, content string) {
	err := os.WriteFile(
		name,
		[]byte(content),
		perm,
	)
	if err != nil {
		panic(err)
	}
}

func generateDebugFiles(items []EpubSection) {
	if !isDebug {
		return
	}
	fmt.Println("Saving debug files")

	for k, v := range items {
		//
		// JSON
		//
		jsonContent, err := json.Marshal(v)
		if err != nil {
			panic(err)
		}

		writeFile(getOutputPath(k, v.ID, "json"), string(jsonContent))

		//
		// HTML
		//
		writeFile(getOutputPath(k, v.ID, "html"), v.HtmlContent)
	}
}
