package consts

const (
	Perm          = 0777
	InputFilePath = "volume/input.epub"

	IsDryRun = false // if true, will generate text files, but not audio files
	IsDebug  = false // if true, will generate files for section json and html content as well

	SpeakProcessCompletion = true // if true, will say something at the end of the process
)

const (
	OutputRootDir  = "output"
	OutputTxtDir   = "txt"
	OutputTmpDir   = "tmp"
	OutputDebugDir = "debug"
)
