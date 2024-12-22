package main

import (
	"epub-tts/internal/book"
	"epub-tts/internal/consts"
	"epub-tts/internal/debug"
	"epub-tts/internal/file"
	"epub-tts/internal/tts"
	"fmt"
	"time"
)

func main() {
	fmt.Println(" ---== Execution Started ==--- ")

	time.Sleep(2 * time.Second)
	foo()
	time.Sleep(2 * time.Second)
	bar()
	time.Sleep(2 * time.Second)
	foobar()

	return

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

	tts := tts.NewTTS(8, textBook)
	tts.Run()

	fmt.Println(" ---== Execution ended ==--- ")
}

func foo() {
	fmt.Printf("\033[0;0H")
	fmt.Println("Item1: Done                                                  	")
	fmt.Println("Item2: Converting												")
	fmt.Println("Item3: Text-to-speech											")
	fmt.Println("Item4: Waiting													")
	fmt.Println("Item5: Waiting													")
	fmt.Println("Item6: Waiting													")
	fmt.Println("Item7: Waiting													")
}

func bar() {
	fmt.Printf("\033[0;0H")
	fmt.Println("Item1: Done                                                  	")
	fmt.Println("Item2: Done													")
	fmt.Println("Item3: Converting												")
	fmt.Println("Item4: Text-to-speech											")
	fmt.Println("Item5: Text-to-speech											")
	fmt.Println("Item6: Waiting													")
	fmt.Println("Item7: Waiting													")
}

func foobar() {
	fmt.Printf("\033[0;0H")
	fmt.Println("Item1: Done                                                  	")
	fmt.Println("Item2: Done													")
	fmt.Println("Item3: Done													")
	fmt.Println("Item4: Done													")
	fmt.Println("Item5: Done													")
	fmt.Println("Item6: Waiting													")
	fmt.Println("Item7: Waiting													")
}
