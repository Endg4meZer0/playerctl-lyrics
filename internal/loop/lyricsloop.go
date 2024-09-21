package loop

import (
	"fmt"
	"lrcsnc/internal/output"
	"lrcsnc/internal/player"
	"lrcsnc/pkg"
	"math"
	"strings"
	"time"
)

var lyricsTimer = time.NewTimer(5 * time.Minute)
var instrTimer = time.NewTimer(5 * time.Minute)
var currentPosition = 0.0
var writtenTimestamp = 0.0
var instrumentalLyric = false

func UpdatePosition(newPosition float64) {
	currentPosition = newPosition
	lyricsTimer.Reset(1)
}

func WriteLyrics() {
	go func() {
		for {
			<-lyricsTimer.C
			if pkg.CurrentSong.LyricsType == 4 {
				instrTimer.Stop()
				fmt.Println()
			} else if pkg.CurrentSong.LyricsType >= 2 {
				instrumentalLyric = true
				instrTimer.Reset(1)
			} else {
				isPlaying, currentPlayerPosition := player.GetPlayerData()
				if math.Abs(currentPosition-currentPlayerPosition) > 1 {
					currentPosition = currentPlayerPosition
				}

				// 5999.99s is basically the maximum limit of .lrc files' timestamps AFAIK, so 6000s is unreachable
				currentLyricTimestamp := -1.0
				nextLyricTimestamp := 6000.0
				lyric := ""
				timestampIndex := -1

				for i, timestamp := range pkg.CurrentSong.LyricTimestamps {
					if timestamp <= currentPosition && currentLyricTimestamp <= timestamp {
						currentLyricTimestamp = timestamp
						lyric = pkg.CurrentSong.Lyrics[i]
						timestampIndex = i
					}
				}

				if timestampIndex != len(pkg.CurrentSong.LyricTimestamps)-1 {
					nextLyricTimestamp = pkg.CurrentSong.LyricTimestamps[timestampIndex+1]
				}

				lyricsTimerDuration := time.Duration(int64(math.Abs(nextLyricTimestamp-currentPosition-0.01)*1000)) * time.Millisecond // tests have shown that it slows down and mismatches without additional 0.01 offset

				writtenTimestamp = currentLyricTimestamp
				// If the currentLyricTimestamp remained at -1.0
				// then reset an instrumental ticker until the first lyric shows up
				if currentLyricTimestamp == -1 {
					instrumentalLyric = true
					instrTimer.Reset(1)
				} else if isPlaying && writtenTimestamp != currentLyricTimestamp { // If paused then don't print the lyric and instead try once more time later
					if lyric == "" {
						// An empty lyric basically means instrumental part,
						// so we reset the instrumental ticker and moving on
						instrumentalLyric = true
						instrTimer.Reset(1)
					} else {
						// An actual lyric when all the conditions are met needs to
						// 1) stop instrumental ticker
						// 2) print itself
						// 3) call the next writing goroutine
						instrumentalLyric = false
						instrTimer.Stop()
						output.PrintLyric(lyric)
					}
				}
				currentPosition = nextLyricTimestamp
				lyricsTimer.Reset(lyricsTimerDuration)
			}
		}
	}()
}

// instrTimer.Stop to stop writing instrumental
// instrTimer.Reset to continue again
func WriteInstrumental() {
	i := 1
	instrTimer.Reset(time.Duration(pkg.CurrentConfig.Output.Instrumental.Interval*1000) * time.Millisecond)
	for {
		<-instrTimer.C
		note := pkg.CurrentConfig.Output.Instrumental.Symbol
		j := int(pkg.CurrentConfig.Output.Instrumental.MaxCount + 1)
		// Not playing? Don't change anything, or it will look kinda strange
		if isPlaying, _ := player.GetPlayerData(); isPlaying {
			if !instrumentalLyric {
				continue
			}
			stringToPrint := ""
			switch pkg.CurrentSong.LyricsType {
			case 1:
				if pkg.CurrentConfig.Output.ShowNotSyncedLyricsWarning {
					stringToPrint += "This song's lyrics are not synced on LrcLib! "
				}
			case 3:
				if pkg.CurrentConfig.Output.ShowSongNotFoundWarning {
					stringToPrint += "Current song was not found on LrcLib! "
				}
			case 5:
				if pkg.CurrentConfig.Output.ShowGettingLyricsMessage {
					stringToPrint += "Getting lyrics... "
				}
			case 6:
				stringToPrint += "Failed to get lyrics! "
			}

			stringToPrint += strings.Repeat(note, i%j)

			output.PrintLyric(stringToPrint)

			i++
			// Don't want to cause any overflow here
			if i > j-1 {
				i = 1
			}
		}
		instrTimer.Reset(time.Duration(pkg.CurrentConfig.Output.Instrumental.Interval*1000) * time.Millisecond)
	}
}
