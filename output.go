package main

import (
	"log"
	"os"
)

var outputDestination *os.File = os.Stdout
var outputDestChanged bool = false

func PrintLyric(lyric string) {
	if outputDestChanged {
		outputDestination.Truncate(0)
		outputDestination.Seek(0, 0)
	} else if CurrentConfig.Output.TerminalOutputInOneLine {
		outputDestination.WriteString("\033[1A\033[K\r")
	}
	outputDestination.WriteString(lyric + "\n")
}

func UpdateOutputDestination(path string) {
	newDest, err := os.Create(path)
	if err != nil {
		log.Println("The output file was set, but I can't create/open it! Permissions issue or wrong path?")
	} else {
		outputDestination = newDest
		outputDestChanged = true
	}
}

func CloseOutput() {
	outputDestination.Close()
}
