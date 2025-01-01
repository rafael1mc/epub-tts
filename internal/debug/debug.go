package debug

import (
	"encoding/json"
	"epub-tts/internal/book"
	"epub-tts/internal/consts"
	"epub-tts/internal/file"
	"errors"
	"fmt"
	"os"
)

func GenerateDebugFiles(epub book.Epub) {
	if !consts.IsDebug {
		return
	}
	fmt.Println("Saving debug files")

	err := os.MkdirAll(file.DebugDir(epub.Name), consts.Perm)
	if err != nil && !errors.Is(err, os.ErrExist) {
		panic(err)
	}

	for k, v := range epub.Sections {
		//
		// JSON
		//
		jsonContent, err := json.Marshal(v)
		if err != nil {
			panic(err)
		}

		err = os.WriteFile(
			file.GetOutputPath(k, file.DebugDir(epub.Name), v.ID, "json"),
			jsonContent,
			consts.Perm,
		)
		if err != nil {
			fmt.Println("Failed to save json debug file")
		}

		//
		// HTML
		//
		err = os.WriteFile(
			file.GetOutputPath(k, file.DebugDir(epub.Name), v.ID, "html"),
			[]byte(v.HtmlContent),
			consts.Perm,
		)
		if err != nil {
			fmt.Println("Failed to save html debug file")
		}
	}
}
