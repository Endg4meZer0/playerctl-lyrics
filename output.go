package main

import (
	"log"
	"os"
	"time"
)

var outputDestination *os.File = os.Stdout
var outputDestChanged bool = false
var overwrite string = ""

func PrintLyric(lyric string) {
	if outputDestChanged {
		outputDestination.Truncate(0)
		outputDestination.Seek(0, 0)
	} else if CurrentConfig.Output.TerminalOutputInOneLine {
		outputDestination.WriteString("\033[1A\033[K\r")
	}
	if overwrite == "" {
		outputDestination.WriteString(lyric + "\n")
	}
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

func PrintOverwrite(over string) {
	overwrite = over
	if outputDestChanged {
		outputDestination.Truncate(0)
		outputDestination.Seek(0, 0)
	} else if CurrentConfig.Output.TerminalOutputInOneLine {
		outputDestination.WriteString("\033[1A\033[K\r")
	}
	outputDestination.WriteString(overwrite + "\n")
	go func() {
		<-time.NewTimer(5 * time.Second).C
		overwrite = ""
	}()
}
