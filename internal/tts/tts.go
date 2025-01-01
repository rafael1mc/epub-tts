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
	ID       int
	BookName string
	Chapter  book.Chapter
}

type jobDone struct {
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
	jobDoneChan := make(chan jobDone, jobCount)

	t.launchWorkers(jobInputChan, jobDoneChan)

	for k, v := range t.textBook.Chapters {
		jobInputChan <- job{ID: k, BookName: t.textBook.Name, Chapter: v}
	}
	close(jobInputChan)

	for range jobCount {
		jobDone := <-jobDoneChan
		if jobDone.Error != nil {
			fmt.Println("Failed to process item", jobDone.Chapter.Name, "with error", jobDone.Error)
		}
	}

	os.RemoveAll(file.TmpDir(t.textBook.Name))
}

func (t TTS) Speak(text string) {
	cmd := fmt.Sprintf(`say "%s"`, text)
	exec.Command("/bin/sh", "-c", cmd).Output()
}

func (t TTS) launchWorkers(jobInputChan <-chan job, jobDoneChan chan<- jobDone) {
	fmt.Println("Launching", t.workerCount, "worker(s)")
	for k := range t.workerCount {
		go t.launchWorker(k, jobInputChan, jobDoneChan)
	}
}

func (t TTS) launchWorker(id int, inputChan <-chan job, doneChan chan<- jobDone) {
	// TODO: use worker id and doneChan with error
	for i := range inputChan {
		if consts.IsDryRun {
			doneChan <- jobDone{job: i}
			continue
		}

		_ = ttsChapter(i.ID, i.BookName, i.Chapter)
		audioConvert(i.ID, i.BookName, i.Chapter)
		// TODO: maybe already delete the aiff file here, to prevent growing then shriking
		// some books generate GBs on aiff
		doneChan <- jobDone{job: i} // not sending errors yet
	}
}

func ttsChapter(pos int, bookName string, chapter book.Chapter) string {
	audioName := file.GetTtsAudioFilename(pos, bookName, chapter)

	fmt.Println("🎤 Narrating chapter: '" + audioName + "' 🎤")
	cmdStr := fmt.Sprintf(`say -f "%s" -o "%s"`, file.GetTextfileName(pos, bookName, chapter), audioName)
	out, _ := exec.Command("/bin/sh", "-c", cmdStr).Output()

	return string(out)
}

func audioConvert(pos int, bookName string, chapter book.Chapter) string {
	ttsAudioName := file.GetTtsAudioFilename(pos, bookName, chapter)
	convertedAudioName := file.GetConvertedAudioFilename(pos, bookName, chapter)

	fmt.Println("🔄 Converting chapter: '" + ttsAudioName + "' 🔄")
	cmdStr := fmt.Sprintf(`ffmpeg -y -i %s %s`, ttsAudioName, convertedAudioName)
	out, _ := exec.Command("/bin/sh", "-c", cmdStr).Output()

	fmt.Println("✅ Chapter '" + convertedAudioName + "' converted ✅")
	return string(out)
}
