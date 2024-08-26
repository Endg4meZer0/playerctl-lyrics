package main

import (
	"fmt"
	"math"
	"time"
)

func SyncLoop() {
	var currentSong SongData
	var currentTimestamps []float64
	var currentLyrics []string

	checkerTicker := time.NewTicker(time.Duration(CurrentConfig.Playerctl.PlayerctlSongCheckInterval*1000) * time.Millisecond)
	positionCheckTicker := time.NewTicker(time.Second)

	position := 0.0
	var timeBeforeGettingLyrics time.Time

	songChanged := make(chan bool, 1)
	fullLyrChan := make(chan bool, 1)

	// Goroutine to check for changes in currently playing song
	go func() {
		for {
			<-checkerTicker.C
			song := GetCurrentSongData()
			if song.Song != currentSong.Song || song.Artist != currentSong.Artist || song.Album != currentSong.Album || song.Duration != currentSong.Duration {
				position = GetCurrentSongPosition()
				timeBeforeGettingLyrics = time.Now()
				currentSong = song
				UpdateData(currentTimestamps, currentLyrics, currentSong)

				songChanged <- true
			}
		}
	}()

	// Goroutine to wait for incoming song metadata (lyrics and instrumental bool)
	go func() {
		for {
			<-songChanged

			// If the duration equals 0s, then there are no supported players out there.
			if currentSong.Duration == 0 {
				currentSong.LyricsType = 4
				fullLyrChan <- false
				continue
			}

			times, lyr := GetSyncedLyrics(&currentSong)
			if currentSong.LyricsType == 5 {
				currentSong.LyricsType = 6
			}

			currentTimestamps = times
			currentLyrics = lyr
			fullLyrChan <- true
		}
	}()

	// Goroutine to watch abnormal changes in player's position
	// For example, seeking on a position bar is counted as abnormal
	go func() {
		for {
			<-positionCheckTicker.C
			initialPosition := GetCurrentSongPosition()
			isPlaying := GetCurrentSongStatus()
			requiredTicks := 10
			for i := 0; i < requiredTicks; i++ {
				time.Sleep(90 * time.Millisecond) // making up for the delay brought by playerctl
				newPosition := GetCurrentSongPosition()
				diff := newPosition - initialPosition
				if diff > -0.1 && diff <= 1.1 && isPlaying { // 0.1 is an okay delta for both sides
					continue
				} else {
					UpdatePosition(newPosition)
					break
				}
			}
		}
	}()

	// Goroutine to update data in the output thread
	go func() {
		for {
			if !<-fullLyrChan {
				currentLyrics = nil
			}

			prevLyric := ""
			count := 1
			for i, lyric := range currentLyrics {
				if CurrentConfig.Output.Romanization.IsEnabled() && IsSupportedAsianLang(lyric) {
					currentLyrics[i] = Romanize(lyric)
				}
				if CurrentConfig.Output.ShowRepeatedLyricsMultiplier {
					if lyric == prevLyric && lyric != "" {
						count++
					} else {
						count = 1
						prevLyric = lyric
					}

					if count != 1 {
						if CurrentConfig.Output.PrintRepeatedLyricsMultiplierToTheRight {
							currentLyrics[i] = fmt.Sprintf(lyric+" "+CurrentConfig.Output.RepeatedLyricsMultiplierFormat, count)
						} else {
							currentLyrics[i] = fmt.Sprintf(CurrentConfig.Output.RepeatedLyricsMultiplierFormat+" "+lyric, count)
						}
					}
				}
			}

			timeAfterGettingLyrics := time.Now()
			position += math.Max(timeAfterGettingLyrics.Sub(timeBeforeGettingLyrics).Seconds(), 0) + 0.1 // tests have shown that additional 0.1 is required to look good

			UpdateData(currentTimestamps, currentLyrics, currentSong)
			UpdatePosition(position)
		}
	}()

	go WriteLyrics()
	go WriteInstrumental()
}
