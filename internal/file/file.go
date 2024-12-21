package file

import (
	"epub-tts/internal/book"
	"epub-tts/internal/consts"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func CreateOutputDir(folderName string) error {
	err := os.MkdirAll(folderName, consts.Perm)
	if err != nil && !errors.Is(err, os.ErrExist) {
		return err
	}

	return nil
}

func SaveChapters(textBook book.TextBook) error {
	fmt.Println("Saving chapter text files.")
	for k, v := range textBook.Chapters {
		filename := GetTextfileName(k, v)
		err := os.WriteFile(
			filename,
			[]byte(v.Content),
			consts.Perm,
		)
		if err != nil {
			return err
		}
	}

	return nil
}

func GetTextfileName(pos int, chapter book.Chapter) string {
	return GetOutputPath(pos, consts.OutputFolderName, chapter.Name, "txt")
}

func GetAudioFilename(pos int, chapter book.Chapter) string {
	return GetOutputPath(pos, consts.OutputFolderName, chapter.Name, "aiff")
}

func GetOutputPath(pos int, outputFolder string, name string, extension string) string {
	filename := fmt.Sprintf("%d.%s.%s", pos, name, extension)
	filename = strings.ReplaceAll(filename, " ", "-")
	filename = strings.ReplaceAll(filename, " ", "-")
	filename = strings.ToLower(filename)

	filePath := filepath.Join(outputFolder, filename)

	return filePath
}