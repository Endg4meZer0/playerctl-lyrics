package main

import (
	"fmt"
	"math"
	"strings"
	"time"
)

func WriteLyrics(lyricsTimer *time.Timer, instrTicker *time.Ticker, currentLyrics *map[float64]string, isPlaying *bool, currentlyInstrumental *bool, noPlayersFound *bool) {
	if *noPlayersFound {
		instrTicker.Stop()
		fmt.Println()
	} else if *currentlyInstrumental {
		instrTicker.Reset(500 * time.Millisecond)
	} else {
		*isPlaying = GetCurrentSongStatus()
		currentTimestamp := GetCurrentSongPosition()
		playerUsesIntegerPosition := false
		if _, d := math.Modf(currentTimestamp); d < 0.000100 {
			// if a floating part is less than this value (tested on cmus, may differ between players)
			// then make an assumption that the player uses integers as position markers
			// 99% sure it can be done better but since it works as of now...
			playerUsesIntegerPosition = true
		}
		firstTimestamp := 6000.0
		currentLyricTimestamp := -1.0
		nextLyricTimestamp := 6000.0
		lyric := ""
		for lyricTimestamp, l := range *currentLyrics {
			if firstTimestamp > lyricTimestamp {
				firstTimestamp = lyricTimestamp
			}
			if lyricTimestamp < currentTimestamp && currentLyricTimestamp < lyricTimestamp {
				currentLyricTimestamp = lyricTimestamp
				lyric = l
			}
			if lyricTimestamp > currentTimestamp && nextLyricTimestamp > lyricTimestamp {
				nextLyricTimestamp = lyricTimestamp
			}
		}
		// If the currentTimestamp is less than even the first timestamp of the lyrics
		// then reset an instrumental ticker until the first lyric shows up
		if currentTimestamp < firstTimestamp {
			instrTicker.Reset(500 * time.Millisecond)
		} else if *isPlaying { // If paused then don't print the lyric and instead try once more time later
			if lyric == "" {
				// An empty lyric basically means instrumental part,
				// so we reset the instrumental ticker and moving on
				instrTicker.Reset(500 * time.Millisecond)
			} else {
				// An actual lyric when all the conditions are met needs to
				// 1) stop instrumental ticker
				// 2) print itself
				// 3) call the next writing goroutine
				instrTicker.Stop()
				if !playerUsesIntegerPosition || math.Abs(nextLyricTimestamp-currentTimestamp) >= 1.0 {
					fmt.Println(lyric)
				}
			}
		}
		lyricsTimerDuration := time.Duration(int64(math.Abs(nextLyricTimestamp-currentTimestamp)*1000)) * time.Millisecond
		lyricsTimer.Reset(lyricsTimerDuration)
		if lyricsTimerDuration/time.Millisecond > 2500 {
			positionCheckTicker := time.NewTicker(2.5 * 1000 * time.Millisecond)
			expectedTicks := int(math.Floor(float64(lyricsTimerDuration/time.Millisecond/1000) / 2.5))
			currentTick := 0
			// Resets the lyric timer if it sees an unusual position change
			go func() {
				for {
					<-positionCheckTicker.C
					currentTick++
					receivedPosition := GetCurrentSongPosition()
					if receivedPosition < math.Floor(currentLyricTimestamp) || receivedPosition > math.Ceil(nextLyricTimestamp) || currentTick >= expectedTicks {
						positionCheckTicker.Stop()
						if currentTick < expectedTicks {
							lyricsTimer.Reset(1)
						}
						return
					}
				}
			}()
		}
		go func() {
			<-lyricsTimer.C
			go WriteLyrics(lyricsTimer, instrTicker, currentLyrics, isPlaying, currentlyInstrumental, noPlayersFound)
		}()
	}
}

// ticker.Stop to stop writing instrumental
// ticker.Reset to continue again
// Should be the same instance (probably? i hope?)
func WriteInstrumental(channel <-chan time.Time, isPlaying *bool, currentSongIsNotFound *bool) {
	note := "♪"
	i := 1
	for {
		<-channel
		// Not playing? Don't change anything, or it will look kinda strange
		if !*isPlaying {
			continue
		} else {
			if *currentSongIsNotFound {
				fmt.Println("Current song was not found on LrcLib! " + strings.Repeat(note, i%4))
			} else {
				fmt.Println(strings.Repeat(note, i%4))
			}
			i++
			// Don't want to cause any overflow here
			if i > 3 {
				i = 1
			}
		}
	}
}
