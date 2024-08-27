package main

import (
	"log"
	"os"
)

var OutputDestination *os.File = os.Stdout

func PrintLyric(lyric string) {
	OutputDestination.Truncate(0)
	OutputDestination.Seek(0, 0)
	OutputDestination.WriteString(lyric + "\n")
}

func UpdateOutputDestination(path string) {
	newDest, err := os.Create(path)
	if err != nil {
		log.Println("The output file was set, but I can't create/open it! Permissions issue or wrong path?")
	} else {
		OutputDestination = newDest
	}
}
