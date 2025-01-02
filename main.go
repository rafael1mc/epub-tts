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

	epub, err := book.ParseEpub(consts.InputFilePath)
	if err != nil {
		panic(err)
	}

	textBook := book.TextBookFromEpub(epub)

	err = file.CreateOutputDirs(textBook.Name)
	if err != nil {
		panic(err)
	}

	err = file.SaveChapters(textBook)
	if err != nil {
		panic(err)
	}
	debug.GenerateDebugFiles(epub)

	tts := tts.NewTTS(textBook)
	tts.Run()

	if consts.SpeakProcessCompletion {
		tts.Speak(consts.SpeakCompletionMessage)
	}
	fmt.Println(" ---== Execution ended ==--- ")
}
