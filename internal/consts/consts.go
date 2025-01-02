package consts

const (
	Perm = 0777
	// InputFilePath = "input/radical-acceptance.epub"
	InputFilePath = "input/Dracula.epub"

	IsDryRun = false // if true, will generate text files, but not audio files
	IsDebug  = false // if true, will generate files for section json and html content as well

	SpeakCompletionMessage = "TTS completed"
	SpeakProcessCompletion = true // if true, will say something at the end of the process

	SegmentWorkoutCount  = 10
	MinSegmentLength     = 500 // used to tts parts of chapter simultaneously
	SegmentOrderFilename = "-segment-order.txt"
)

const (
	OutputRootDir  = "output"
	OutputTxtDir   = "txt"
	OutputTmpDir   = "tmp"
	OutputDebugDir = "debug"
)
