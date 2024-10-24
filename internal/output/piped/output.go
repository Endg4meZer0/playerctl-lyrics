package piped

import (
	"log"
	"lrcsnc/internal/pkg/global"
	"os"
	"strings"
	"time"
)

var outputDestination *os.File = os.Stdout
var outputDestChanged = false
var overwrite = ""
var instrTimer = time.NewTimer(5 * time.Minute)
var writeInstrumental bool = false

func Init() {
	i := 1
	instrTimer.Reset(time.Duration(global.CurrentConfig.Output.Instrumental.Interval*1000) * time.Millisecond)
	for {
		<-instrTimer.C
		note := global.CurrentConfig.Output.Instrumental.Symbol
		j := int(global.CurrentConfig.Output.Instrumental.MaxCount + 1)

		// Only update instrumental stuff if the song is playing
		if global.CurrentPlayer.IsPlaying && writeInstrumental {
			stringToPrint := ""
			switch global.CurrentSong.LyricsData.LyricsType {
			case 1:
				if global.CurrentConfig.Output.ShowNotSyncedLyricsWarning {
					stringToPrint += "This song's lyrics are not synced on LrcLib! "
				}
			case 3:
				if global.CurrentConfig.Output.ShowSongNotFoundWarning {
					stringToPrint += "This song was not found on LrcLib! "
				}
			case 5:
				if global.CurrentConfig.Output.ShowGettingLyricsMessage {
					stringToPrint += "Getting lyrics... "
				}
			case 6:
				stringToPrint += "Failed to get lyrics! "
			}

			stringToPrint += strings.Repeat(note, i%j)

			outputDestination.WriteString(stringToPrint + "\n")

			i++
			// Don't want to cause any overflow here
			if i >= j {
				i = 1
			}
		}
		instrTimer.Reset(time.Duration(global.CurrentConfig.Output.Instrumental.Interval*1000) * time.Millisecond)
	}
}

func PrintLyric(lyric string) {
	if outputDestChanged {
		outputDestination.Truncate(0)
		outputDestination.Seek(0, 0)
	}
	if overwrite == "" {
		if lyric == "" {
			writeInstrumental = true
			instrTimer.Reset(1)
		} else {
			writeInstrumental = false
			instrTimer.Stop()
			outputDestination.WriteString(lyric + "\n")
		}
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
	}
	outputDestination.WriteString(overwrite + "\n")
	go func() {
		<-time.NewTimer(5 * time.Second).C
		overwrite = ""
	}()
}
