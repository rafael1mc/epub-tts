package file

import (
	"epub-tts/internal/book"
	"epub-tts/internal/consts"
	"epub-tts/internal/str"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func CreateOutputDirs() error {
	err := os.MkdirAll(consts.TmpOutputFolderName, consts.Perm)
	if err != nil && !errors.Is(err, os.ErrExist) {
		return err
	}

	err = os.MkdirAll(consts.TxtOutputFolderName, consts.Perm)
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
	return GetOutputPath(pos, consts.TxtOutputFolderName, chapter.NameOrID(), "txt")
}

func GetTtsAudioFilename(pos int, chapter book.Chapter) string {
	return GetOutputPath(pos, consts.TmpOutputFolderName, chapter.NameOrID(), "aiff")
}

func GetConvertedAudioFilename(pos int, chapter book.Chapter) string {
	return GetOutputPath(pos, consts.OutputFolderName, chapter.NameOrID(), "mp3")
}

func GetOutputPath(pos int, outputFolder string, name string, extension string) string {
	filename := fmt.Sprintf("%d-%s.%s", pos, name, extension)
	filename = strings.ToLower(filename)
	filename = str.CleanFileName(filename)

	filePath := filepath.Join(outputFolder, filename)

	return filePath
}
