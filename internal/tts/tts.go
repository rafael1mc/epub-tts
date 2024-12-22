package tts

import (
	"epub-tts/internal/book"
	"epub-tts/internal/consts"
	"epub-tts/internal/file"
	"fmt"
	"os"
	"os/exec"
)

type TTS struct {
	workerCount int

	textBook book.TextBook
}

type job struct {
	ID      int
	Chapter book.Chapter
}

type jobError struct {
	job
	Error error
}

func NewTTS(
	workerCount int,
	textBook book.TextBook,
) *TTS {
	return &TTS{
		workerCount: workerCount,
		textBook:    textBook,
	}
}

func (t TTS) Run() {
	fmt.Println("Running text-to-speech")

	jobCount := len(t.textBook.Chapters)
	jobInputChan := make(chan job, jobCount)
	jobDoneChan := make(chan *jobError, jobCount)

	t.launchWorkers(jobInputChan, jobDoneChan)

	for k, v := range t.textBook.Chapters {
		jobInputChan <- job{ID: k, Chapter: v}
	}
	close(jobInputChan)

	for range jobCount {
		jobErr := <-jobDoneChan
		if jobErr != nil {
			fmt.Println("Failed to process item", jobErr.Chapter.Name, "with error", jobErr.Error)
		}
	}

	os.RemoveAll(consts.TmpOutputFolderName)
}

func (t TTS) launchWorkers(jobInputChan <-chan job, jobDoneChan chan<- *jobError) {
	fmt.Println("Launching", t.workerCount, "worker(s)")
	for k := range t.workerCount {
		go t.launchWorker(k, jobInputChan, jobDoneChan)
	}
}

func (t TTS) launchWorker(id int, inputChan <-chan job, doneChan chan<- *jobError) {
	// TODO: use worker id and doneChan with error
	for i := range inputChan {
		if consts.IsDryRun {
			doneChan <- nil
			continue
		}

		_ = ttsChapter(i.ID, i.Chapter)
		audioConvert(i.ID, i.Chapter)
		doneChan <- nil
	}
}

func ttsChapter(pos int, chapter book.Chapter) string {
	audioName := file.GetTtsAudioFilename(pos, chapter)

	fmt.Println("🎤 Narrating chapter: '" + audioName + "' 🎤")
	cmdStr := fmt.Sprintf(`say -f "%s" -o "%s"`, file.GetTextfileName(pos, chapter), audioName)
	out, _ := exec.Command("/bin/sh", "-c", cmdStr).Output()

	return string(out)
}

func audioConvert(pos int, chapter book.Chapter) string {
	ttsAudioName := file.GetTtsAudioFilename(pos, chapter)
	convertedAudioName := file.GetConvertedAudioFilename(pos, chapter)

	fmt.Println("🔄 Converting chapter: '" + ttsAudioName + "' 🔄")
	cmdStr := fmt.Sprintf(`ffmpeg -y -i %s %s`, ttsAudioName, convertedAudioName)
	out, _ := exec.Command("/bin/sh", "-c", cmdStr).Output()

	fmt.Println("✅ Chapter '" + convertedAudioName + "' converted ✅")
	return string(out)
}
