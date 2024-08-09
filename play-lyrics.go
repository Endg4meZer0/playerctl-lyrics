package main

import (
	"fmt"
	"math"
	"strings"
	"time"
)

func PlayLyrics() {
	var currentSong SongData
	var currentLyrics map[float64]string
	var currentlyInstrumental bool
	var currentSongIsNotFound bool
	var noPlayersFound bool
	var isPlaying bool

	lyricsTimer := time.NewTimer(time.Second)
	lyricsTimer.Stop()
	instrTicker := time.NewTicker(time.Second)
	instrTicker.Stop()                                                      // stopping because the ticker should be reset when it's needed, and Go doesn't close ticker's channel after stop
	go WriteInstrumental(instrTicker.C, &isPlaying, &currentSongIsNotFound) // also starting the instrumental thread at the same time to not create additional instances and only work with the ticker

	checkerTicker := time.NewTicker(time.Second)

	songChanged := make(chan bool, 1)
	fullLyrChan := make(chan map[float64]string, 1)

	// Goroutine to wait for incoming song metadata (lyrics and instrumental bool)
	go func() {
		for {
			<-songChanged
			if currentSong.Song == "" && currentSong.Artist == "" && currentSong.Album == "" {
				fullLyrChan <- nil
				noPlayersFound = true
				continue
			} else {
				noPlayersFound = false
			}

			lyr, instr := GetSyncedLyrics(&currentSong)
			if lyr == nil {
				currentSongIsNotFound = !instr
				currentlyInstrumental = true
			} else {
				currentSongIsNotFound = false
				currentlyInstrumental = false
			}
			fullLyrChan <- lyr
		}
	}()

	// Goroutine to check for changes in currently playing song
	go func() {
		for {
			<-checkerTicker.C
			song := GetCurrentSongData()
			if song != currentSong {
				currentSong = song
				songChanged <- true
				//currentLyrics, currentlyInstrumental := GetSyncedLyrics(song)
			}
		}
	}()

	go func() {
		for {
			currentLyrics = <-fullLyrChan
			lyricsTimer.Stop()
			instrTicker.Stop()
			go WriteLyrics(lyricsTimer, instrTicker, &currentLyrics, &isPlaying, &currentlyInstrumental, &noPlayersFound)
			/*
				Timer is made like this:
				1) get the lyric from the map based on timestamp (we need the next lyric AFTER that timestamp)
				2) make a goroutine that writes the lyric to stdout when the timestamp comes
				3) should be recursive, so after the timer ticks, everything should begin from 1 again.
				Probably will make an extra function for that
			*/
		}
	}()
}

func WriteLyrics(lyricsTimer *time.Timer, instrTicker *time.Ticker, currentLyrics *map[float64]string, isPlaying *bool, currentlyInstrumental *bool, noPlayersFound *bool) {
	if *noPlayersFound {
		instrTicker.Stop()
		fmt.Println()
	} else if *currentlyInstrumental {
		instrTicker.Reset(time.Second)
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
			instrTicker.Reset(time.Second)
		} else if *isPlaying { // If paused then don't print the lyric and instead try once more time later
			if lyric == "" {
				// An empty lyric basically means instrumental part,
				// so we reset the instrumental ticker and moving on
				instrTicker.Reset(time.Second)
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
