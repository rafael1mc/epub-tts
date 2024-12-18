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
 - Docker
 - Golang

# How to use
 1. Clone this repo
 2. Build the image for the dependency:
 ```
 docker build -t epub-parser ./parser
 ```
 3. Replace the file inside `volume/input.epub` with the book you want to convert to audio (keep file name)
 4. Execute the program (note that it will take quite some time, but you should see _some_ output during execution):
```
go run .
```
 5. You should see a new `output` folder with each text and audio file.

# TODO
 - [ ] Add support for Ubuntu TTS
 - [ ] Reduce output audio size
 - [ ] Remove duplicated information in each chapter (eg book title)
 - [ ] Enhance information retrieval
 - [ ] Find a Golang epub parser
 - [ ] Make it fully run inside the container
 - [ ] Add Web UI to Drag and Drop epub files
 - [ ] Add TTS progress (look into `say` progress)
 - [ ] Support other languages
 - [ ] Organize code a little
 - [ ] ?

### Dependencies
 - Uses [gaoxiaoliangz/epub-parser](https://github.com/gaoxiaoliangz/epub-parser) (for now, simplest parser that worked how I wanted)
 - MacOS `say` command
 - Note: The example book in this repo is taken from [Project Guttenber](https://www.gutenberg.org/about/), with Copyright Status as "Public domain in the USA"
<hr>

# License
Check [LICENSE](https://github.com/rafael1mc/epub-tts/blob/main/LICENSE) file.