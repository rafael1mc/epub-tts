package consts

const (
	Perm                = 0777
	InputFilePath       = "volume/input.epub"
	OutputFolderName    = "output"
	TxtOutputFolderName = "output/txt"
	TmpOutputFolderName = OutputFolderName + "/tmp"
	IsDryRun            = false // if true, will generate text files, but not audio files
	IsDebug             = false // if true, will generate files for section json and html content as well
)
