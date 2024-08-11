package main

import (
	"fmt"
	"math"
	"strings"
	"time"
)

func WriteLyrics(lyricsTimer *time.Timer, instrTimer *time.Timer, currentLyrics *map[float64]string, isPlaying *bool, currentSong *SongData) {
	if currentSong.LyricsType == 4 {
		instrTimer.Stop()
		fmt.Println()
	} else if currentSong.LyricsType >= 2 {
		instrTimer.Reset(1)
	} else {
		wasPaused := !*isPlaying
		*isPlaying = GetCurrentSongStatus()
		currentTimestamp := GetCurrentSongPosition()

		// if a floating part is less than that value down there (tested on cmus, may differ between players)
		// then make an assumption that the player uses integers as position markers
		// and we'll allow the timestamps to be equal or greater only by seconds, not milliseconds
		// 99% sure it can be done better but since it works as of now...
		_, currentTimestampFloatPart := math.Modf(currentTimestamp)
		playerUsesIntegerPosition := currentTimestampFloatPart < 0.000100

		// 5999.99s is basically the maximum limit of .lrc files' timestamps, so 6000s is unreachable
		firstLyricTimestamp := 6000.0
		currentLyricTimestamp := 0.0
		nextLyricTimestamp := 6000.0
		lyric := ""

		for lyricTimestamp, l := range *currentLyrics {
			if firstLyricTimestamp > lyricTimestamp {
				firstLyricTimestamp = lyricTimestamp
			}
			if lyricTimestamp < currentTimestamp && currentLyricTimestamp <= lyricTimestamp {
				currentLyricTimestamp = lyricTimestamp
				lyric = l
			}
			if lyricTimestamp > currentTimestamp && nextLyricTimestamp > lyricTimestamp {
				nextLyricTimestamp = lyricTimestamp
			}
		}

		fmt.Println(currentTimestamp, currentLyricTimestamp, nextLyricTimestamp, firstLyricTimestamp, time.Now())
		// If the nextLyricTimestamp remained at 6000s, then there are no more lyrics.
		// If that's the case, we'll need to account that the same song may be put on repeat
		// So the idea would be to change the value nextLyricTimestamp to the playing song's duration
		// and maybe add a bit like 0.25s to be 100% sure
		if nextLyricTimestamp == 6000.0 {
			nextLyricTimestamp = math.Abs(currentSong.Duration) + 0.25
		}
		// If the currentTimestamp is less than even the first timestamp of the lyrics
		// then reset an instrumental ticker until the first lyric shows up
		if currentTimestamp < firstLyricTimestamp {
			instrTimer.Reset(1)
		} else if *isPlaying { // If paused then don't print the lyric and instead try once more time later
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
				if !wasPaused { // if the playback was paused, that usually causes lyric to print itself twice, so here's a little fuse
					fmt.Println(lyric)
				}
			}
		}
		var lyricsTimerDuration time.Duration
		if playerUsesIntegerPosition {
			currentTimestampIntPart := math.Floor(currentTimestamp)
			nextLyricTimestampIntPart := math.Floor(nextLyricTimestamp)

			lyricsTimerDuration = time.Duration(int64(math.Abs(nextLyricTimestampIntPart-currentTimestampIntPart)*1000)) * time.Millisecond
		} else {
			lyricsTimerDuration = time.Duration(int64(math.Abs(nextLyricTimestamp-currentTimestamp)*1000)) * time.Millisecond
		}
		lyricsTimer.Reset(lyricsTimerDuration)
		if lyricsTimerDuration/time.Millisecond > 1500 {
			positionCheckTicker := time.NewTicker(1.5 * 1000 * time.Millisecond)
			expectedTicks := int(math.Floor(float64(lyricsTimerDuration/time.Millisecond/1000) / 1.5))
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
			go WriteLyrics(lyricsTimer, instrTimer, currentLyrics, isPlaying, currentSong)
		}()
	}
}

// instrTimer.Stop to stop writing instrumental
// instrTimer.Reset to continue again
func WriteInstrumental(instrTimer *time.Timer, isPlaying *bool, currentSong *SongData) {
	note := "♪"
	i := 1
	for {
		<-instrTimer.C
		*isPlaying = GetCurrentSongStatus()
		// Not playing? Don't change anything, or it will look kinda strange
		if !*isPlaying {
			continue
		} else {
			if currentSong.LyricsType == 3 {
				fmt.Println("Current song was not found on LrcLib! " + strings.Repeat(note, i%4))
			} else if currentSong.LyricsType == 5 {
				fmt.Println("Failed to get lyrics! " + strings.Repeat(note, i%4))
			} else {
				fmt.Println(strings.Repeat(note, i%4))
			}
			i++
			// Don't want to cause any overflow here
			if i > 3 {
				i = 1
			}
			instrTimer.Reset(500 * time.Millisecond)
		}
	}
}
