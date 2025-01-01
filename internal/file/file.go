package file

import (
	"epub-tts/internal/book"
	"epub-tts/internal/consts"
	"epub-tts/internal/str"
	"errors"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"
)

func normalizeBookName(bookName string) string {
	cleanName := strings.ToLower(str.CleanFileName(bookName))
	nameLen := len(cleanName)
	if nameLen > 50 {
		nameLen = 50
	}

	return cleanName[:nameLen]
}

func rootDir(bookName string) string {
	return path.Join(
		consts.OutputRootDir,
		normalizeBookName(bookName),
	)
}

func txtDir(bookName string) string {
	return path.Join(
		rootDir(bookName),
		consts.OutputTxtDir,
	)
}

func TmpDir(bookName string) string {
	return path.Join(
		rootDir(bookName),
		consts.OutputTmpDir,
	)
}

func DebugDir(bookName string) string {
	return path.Join(
		rootDir(bookName),
		consts.OutputDebugDir,
	)
}

func CreateOutputDirs(bookName string) error {
	var err error

	fmt.Println("Creating tmp dir", TmpDir(bookName))
	err = os.MkdirAll(TmpDir(bookName), consts.Perm)
	if err != nil && !errors.Is(err, os.ErrExist) {
		return err
	}

	err = os.MkdirAll(txtDir(bookName), consts.Perm)
	if err != nil && !errors.Is(err, os.ErrExist) {
		return err
	}

	return nil
}

func SaveChapters(textBook book.TextBook) error {
	fmt.Println("Saving chapter text files.")
	for k, v := range textBook.Chapters {
		filename := GetTextfileName(k, textBook.Name, v)
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

func GetTextfileName(pos int, bookName string, chapter book.Chapter) string {
	return GetOutputPath(pos, txtDir(bookName), chapter.NameOrID(), "txt")
}

func GetTtsAudioFilename(pos int, bookName string, chapter book.Chapter) string {
	return GetOutputPath(pos, TmpDir(bookName), chapter.NameOrID(), "aiff")
}

func GetConvertedAudioFilename(pos int, bookName string, chapter book.Chapter) string {
	return GetOutputPath(pos, rootDir(bookName), chapter.NameOrID(), "mp3")
}

func GetOutputPath(pos int, outputFolder string, name string, extension string) string {
	filename := fmt.Sprintf("%d-%s.%s", pos, name, extension)
	filename = strings.ToLower(filename)
	filename = str.CleanFileName(filename)

	filePath := filepath.Join(outputFolder, filename)

	return filePath
}
