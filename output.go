package main

import (
	"fmt"
	"math"
	"strings"
	"time"
)

var lyricsTimer = time.NewTimer(100 * time.Minute)
var instrTimer = time.NewTimer(500 * time.Millisecond)
var currentTimestamps []float64
var currentLyrics []string
var currentSong = SongData{LyricsType: 5}
var isPlaying = false
var currentPosition = 0.0
var firstInstanceForSong = false

func UpdateData(newTimes []float64, newLyrics []string, newSong SongData, position float64) {
	currentTimestamps = newTimes
	currentLyrics = newLyrics
	currentSong = newSong
	currentPosition = position
	lyricsTimer.Reset(100)
	instrTimer.Stop()
}

func UpdatePosition(newPosition float64) {
	currentPosition = newPosition
	lyricsTimer.Reset(100)
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
				firstLyricTimestamp := 6000.0
				currentLyricTimestamp := -1.0
				nextLyricTimestamp := 6000.0
				lyric := ""

				for i, timestamp := range currentTimestamps {
					if firstLyricTimestamp > timestamp {
						firstLyricTimestamp = timestamp
					}
					if timestamp <= currentPosition && currentLyricTimestamp <= timestamp {
						currentLyricTimestamp = timestamp
						lyric = currentLyrics[i]
					}
					if timestamp > currentPosition && nextLyricTimestamp > timestamp {
						nextLyricTimestamp = timestamp
					}
				}

				if nextLyricTimestamp == 6000.0 {
					// If the nextLyricTimestamp remained at 6000s, then there are no more lyrics.
					// If that's the case, we'll need to account that the same song may be put on repeat
					// So the idea would be to change the value nextLyricTimestamp to the playing song's duration
					// and maybe add a bit like 0.1s to be 100% sure we're on a new song iteration
					nextLyricTimestamp = math.Abs(currentSong.Duration) + 0.1
				}

				lyricsTimerDuration := time.Duration(int64(math.Abs(nextLyricTimestamp-currentPosition-0.01)*1000)) * time.Millisecond // tests have shown that it slows down and mismatches without additional 0.01 offset

				// If the currentPosition is less than even the first timestamp of the lyrics
				// then reset an instrumental ticker until the first lyric shows up
				if currentPosition < firstLyricTimestamp {
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
		isPlaying = GetCurrentSongStatus()
		// Not playing? Don't change anything, or it will look kinda strange
		if !isPlaying {
			instrTimer.Reset(time.Duration(CurrentConfig.Output.Instrumental.Interval*1000) * time.Millisecond)
			continue
		} else {
			messagePrinted := false
			switch currentSong.LyricsType {
			case 1:
				if CurrentConfig.Output.ShowNotSyncedLyricsWarning {
					fmt.Println("This song's lyrics are not synced on LrcLib! " + strings.Repeat(note, i%j))
					messagePrinted = true
				}
			case 3:
				if CurrentConfig.Output.ShowSongNotFoundWarning {
					fmt.Println("Current song was not found on LrcLib! " + strings.Repeat(note, i%j))
					messagePrinted = true
				}
			case 5:
				if CurrentConfig.Output.ShowGettingLyricsMessage {
					fmt.Println("Getting lyrics... " + strings.Repeat(note, i%j))
					messagePrinted = true
				}
			case 6:
				fmt.Println("Failed to get lyrics! " + strings.Repeat(note, i%j))
				messagePrinted = true
			}

			if !messagePrinted {
				fmt.Println(strings.Repeat(note, i%j))
			}

			i++
			// Don't want to cause any overflow here
			if i > j-1 {
				i = 1
			}
			instrTimer.Reset(time.Duration(CurrentConfig.Output.Instrumental.Interval*1000) * time.Millisecond)
		}
	}
}
