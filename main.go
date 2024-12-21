package main

import (
	"epub-tts/internal/book"
	"epub-tts/internal/consts"
	"epub-tts/internal/debug"
	"epub-tts/internal/file"
	"epub-tts/internal/tts"
	"fmt"
)

func main() {
	fmt.Println(" ---== Execution Started ==--- ")

	err := file.CreateOutputDir(consts.TmpOutputFolderName)
	if err != nil {
		panic(err)
	}

	epub, err := book.ParseEpub(consts.InputFilePath)
	if err != nil {
		panic(err)
	}

	textBook := book.TextBookFromEpub(epub)

	err = file.SaveChapters(textBook)
	if err != nil {
		panic(err)
	}
	debug.GenerateDebugFiles(epub)

	tts := tts.NewTTS(3, textBook)
	tts.Run()

	fmt.Println(" ---== Execution ended ==--- ")
}
