# epub-tts

Convert ePUB into audio files.

Code will parse the ePUB into sections (which roughly correlates to book chapters) and 'text-to-speech' each section into its own audio file.<br>Output will be prefixed with a number to maintain order.

<sub>
This is an alpha, proof of concept version.
To me, it's supposed to be a simple alternative for when eyes are tired but the mind is not :)
</sub>

<br>

# Requirements
 - Run on MacOS
 - ffmpeg installed and available in $PATH
 - Golang

# How to use
 1. Clone this repo
 2. Replace the file inside `volume/input.epub` with the book you want to convert to audio (keep file name)
 3. Execute the program (note that it will take quite some time, but you should see _some_ output during execution):
```
go run .
```
 5. You should see a new `output` folder with each text and audio file.

# TODO
 - [x] Parse ePUB from Golang
 - [x] Organize code
 - [x] Add worker pools for batch conversion and less CPU strain
 - [x] Reduce output audio size
 - [x] Extract chapter info
 - [ ] Add more sample ePUBs
 - [ ] Add automated tests
 - [ ] Add multiple input (separate output by folder same name as input file)
 - [ ] Organize the code some more
 - [ ] Support other languages beyond english
 - [ ] Display progress
 - [ ] Add support for Ubuntu TTS
 - [ ] Add Web UI to Drag and Drop epub files
 - [ ] ?

### Dependencies
 - MacOS `say` command
 - Note: The example book in this repo is taken from [Project Guttenber](https://www.gutenberg.org/about/), with Copyright Status as "Public domain in the USA"
<hr>

# License
Check [LICENSE](https://github.com/rafael1mc/epub-tts/blob/main/LICENSE) file.