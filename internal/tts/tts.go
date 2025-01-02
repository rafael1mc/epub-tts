package tts

import (
	"epub-tts/internal/book"
	"epub-tts/internal/consts"
	"epub-tts/internal/file"
	"epub-tts/internal/pool"
	"epub-tts/internal/str"
	"fmt"
	"os"
	"os/exec"
	"path"
)

type TTS struct {
	textBook book.TextBook
}

type job struct {
	ID       int
	BookName string
	Chapter  book.Chapter
}

func NewTTS(
	textBook book.TextBook,
) *TTS {
	return &TTS{
		textBook: textBook,
	}
}

func (t TTS) Run() {
	if consts.IsDryRun {
		fmt.Println("DryRun -> Skipping text-to-speech")
		return
	}
	fmt.Println("Running text-to-speech")

	workerPool := pool.NewPool[job](consts.SegmentWorkoutCount)
	for k, chapter := range t.textBook.Chapters {
		// create chapter folder inside tmp
		chapterDirName := file.GetChapterDirName(k, t.textBook.Name, chapter)
		err := os.MkdirAll(chapterDirName, consts.Perm)
		if err != nil {
			fmt.Println("Failed to create chapter directory with error:", err)
			continue
		}

		// add segment file to be used by ffmpeg when merging
		segments := str.SplitAfterN(chapter.Content, '.', consts.MinSegmentLength)
		orderFileContent := ""
		for j := range segments {
			segmentMp3Name := path.Base(file.GetSegmentMp3Filename(k, j, t.textBook.Name, chapter))
			orderFileContent += "file " + segmentMp3Name + "\n"
		}
		segmentOrderFile := file.GetSegmentOrderfilename(k, t.textBook.Name, chapter)
		err = os.WriteFile(
			segmentOrderFile,
			[]byte(orderFileContent),
			consts.Perm,
		)
		if err != nil {
			fmt.Println("Failed to create segment order file with error:", err)
			continue
		}

		for j, segment := range segments {
			workerPool.AddWork(func() {
				// create txt for each segment
				segmentTextName := file.GetSegmentTextFilename(k, j, t.textBook.Name, chapter)
				err := os.WriteFile(segmentTextName, []byte(segment), consts.Perm)
				if err != nil {
					fmt.Println("Failed to create segment file with error:", err)
					return
				}

				// generate audio for each segment
				segmentAiffName := file.GetSegmentAiffFilename(k, j, t.textBook.Name, chapter)
				segmentMp3Name := file.GetSegmentMp3Filename(k, j, t.textBook.Name, chapter)
				ttsFile(segmentTextName, segmentAiffName)
				convertAudio(segmentAiffName, segmentMp3Name)
			})
		}
		workerPool.Start() // will wait
		chapterFilename := file.GetFullChapterFilename(k, t.textBook.Name, chapter)
		mergeChapter(segmentOrderFile, chapterFilename)

		err = os.RemoveAll(chapterDirName)
		if err != nil {
			fmt.Println("failed to remove chapter dir with error:", err)
		}
		fmt.Println("✅ Chapter '" + path.Base(chapterDirName) + "' completed ✅")
	}
}

func (t TTS) Speak(text string) {
	cmd := fmt.Sprintf(`say "%s"`, text)
	exec.Command("/bin/sh", "-c", cmd).Output()
}

func ttsFile(inputPath string, outputPath string) {
	fmt.Println("🎤 Narrating: '" + path.Base(inputPath) + "' 🎤")
	cmdStr := fmt.Sprintf(`say -f "%s" -o "%s"`, inputPath, outputPath)
	exec.Command("/bin/sh", "-c", cmdStr).Output()
}

func convertAudio(inputPath, outputPath string) {
	fmt.Println("🔄 Converting: '" + path.Base(inputPath) + "' 🔄")
	cmdStr := fmt.Sprintf(`ffmpeg -y -i %s %s`, inputPath, outputPath)
	exec.Command("/bin/sh", "-c", cmdStr).Output()
	fmt.Println("✓ Converted: '" + path.Base(outputPath) + "' ✓")
}

func mergeChapter(orderPath string, outputPath string) {
	fmt.Println("📦 Merging '" + path.Base(outputPath) + "' 📦")
	cmdStr := fmt.Sprintf(`ffmpeg -y -f concat -safe 0 -i %s -c copy %s`, orderPath, outputPath)
	exec.Command("/bin/sh", "-c", cmdStr).Output()
}
