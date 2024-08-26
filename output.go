package main

import (
	"fmt"
	"math"
	"strings"
	"time"
)

var lyricsTimer = time.NewTimer(5 * time.Minute)
var instrTimer = time.NewTimer(5 * time.Minute)
var currentTimestamps []float64
var currentLyrics []string
var currentSong = SongData{LyricsType: 5}
var isPlaying = false
var currentPosition = 0.0

func UpdateData(newTimes []float64, newLyrics []string, newSong SongData) {
	currentTimestamps = newTimes
	currentLyrics = newLyrics
	currentSong = newSong
	lyricsTimer.Reset(1)
	instrTimer.Stop()
}

func UpdatePosition(newPosition float64) {
	currentPosition = newPosition
	lyricsTimer.Reset(1)
	instrTimer.Stop()
}

func WriteLyrics() {
	go func() {
		for {
			<-lyricsTimer.C
			if currentSong.LyricsType == 4 {
				instrTimer.Stop()
				fmt.Println()
			} else if currentSong.LyricsType >= 2 {
				instrTimer.Reset(1)
			} else {
				isPlaying = GetCurrentSongStatus()

				// 5999.99s is basically the maximum limit of .lrc files' timestamps AFAIK, so 6000s is unreachable
				currentLyricTimestamp := -1.0
				nextLyricTimestamp := 6000.0
				lyric := ""
				timestampIndex := -1

				for i, timestamp := range currentTimestamps {
					if timestamp <= currentPosition && currentLyricTimestamp <= timestamp {
						currentLyricTimestamp = timestamp
						lyric = currentLyrics[i]
						timestampIndex = i
					}
				}

				if timestampIndex != len(currentTimestamps) {
					nextLyricTimestamp = currentTimestamps[timestampIndex+1]
				}

				lyricsTimerDuration := time.Duration(int64(math.Abs(nextLyricTimestamp-currentPosition-0.01)*1000)) * time.Millisecond // tests have shown that it slows down and mismatches without additional 0.01 offset

				// If the currentLyricTimestamp remained at -1.0
				// then reset an instrumental ticker until the first lyric shows up
				if currentLyricTimestamp == -1 {
					instrTimer.Reset(1)
				} else if isPlaying { // If paused then don't print the lyric and instead try once more time later
					if lyric == "" {
						// An empty lyric basically means instrumental part,
						// so we reset the instrumental ticker and moving on
						instrTimer.Reset(1)
					} else {
						// An actual lyric when all the conditions are met needs to
						// 1) stop instrumental ticker
						// 2) print itself
						// 3) call the next writing goroutine
						instrTimer.Stop()
						fmt.Println(lyric)
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
	note := CurrentConfig.Output.Instrumental.Symbol
	i := 1
	j := int(CurrentConfig.Output.Instrumental.MaxCount + 1)
	instrTimer.Reset(time.Duration(CurrentConfig.Output.Instrumental.Interval*1000) * time.Millisecond)
	for {
		<-instrTimer.C
		// Not playing? Don't change anything, or it will look kinda strange
		if GetCurrentSongStatus() {
			switch currentSong.LyricsType {
			case 1:
				if CurrentConfig.Output.ShowNotSyncedLyricsWarning {
					fmt.Println("This song's lyrics are not synced on LrcLib! " + strings.Repeat(note, i%j))
				}
			case 3:
				if CurrentConfig.Output.ShowSongNotFoundWarning {
					fmt.Println("Current song was not found on LrcLib! " + strings.Repeat(note, i%j))
				}
			case 5:
				if CurrentConfig.Output.ShowGettingLyricsMessage {
					fmt.Println("Getting lyrics... " + strings.Repeat(note, i%j))
				}
			case 6:
				fmt.Println("Failed to get lyrics! " + strings.Repeat(note, i%j))
			default:
				fmt.Println(strings.Repeat(note, i%j))
			}

			i++
			// Don't want to cause any overflow here
			if i > j-1 {
				i = 1
			}
		}
		instrTimer.Reset(time.Duration(CurrentConfig.Output.Instrumental.Interval*1000) * time.Millisecond)
	}
}
