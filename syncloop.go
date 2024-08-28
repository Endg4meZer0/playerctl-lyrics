package main

import (
	"fmt"
	"math"
	"time"
)

func SyncLoop() {
	var currentSong SongData

	checkerTicker := time.NewTicker(time.Duration(CurrentConfig.Playerctl.PlayerctlSongCheckInterval*1000) * time.Millisecond)
	positionCheckTicker := time.NewTimer(time.Second)
	positionInnerCheckTicker := time.NewTicker(time.Second)
	positionInnerCheckTicker.Stop()

	position := 0.0
	var timeBeforeGettingLyrics time.Time

	songChanged := make(chan bool, 1)
	fullLyrChan := make(chan bool, 1)

	// Goroutine to check for changes in currently playing song
	go func() {
		for {
			<-checkerTicker.C
			song := GetSongData()
			if song.Song != currentSong.Song || song.Artist != currentSong.Artist || song.Album != currentSong.Album || song.Duration != currentSong.Duration {
				_, position = GetPlayerData()
				timeBeforeGettingLyrics = time.Now()
				currentSong = song
				UpdateData(currentSong)

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

			GetSyncedLyrics(&currentSong)
			if currentSong.LyricsType == 5 {
				currentSong.LyricsType = 6
			}

			fullLyrChan <- true
		}
	}()

	// Goroutine to watch abnormal changes in player's position
	// For example, seeking on a position bar is counted as abnormal
	go func() {
		for {
			<-positionCheckTicker.C
			_, initialPosition := GetPlayerData()
			requiredTicks := 10
			positionInnerCheckTicker.Reset(100 * time.Millisecond)
			for i := 0; i < requiredTicks; i++ {
				<-positionInnerCheckTicker.C
				isStillPlaying, newPosition := GetPlayerData()
				expectedPosition := (initialPosition + 0.1*(float64(i)+1))
				diff := newPosition - expectedPosition
				if !(((diff <= 0.21 && diff >= -1.11) || (diff >= 0.89 && diff <= 1.01)) || !isStillPlaying) {
					UpdatePosition(newPosition)
					break
				}
			}
			positionInnerCheckTicker.Stop()
			positionCheckTicker.Reset(1)
		}
	}()

	// Goroutine to update data in the output thread
	go func() {
		for {
			<-fullLyrChan

			prevLyric := ""
			count := 1
			for i, lyric := range currentSong.Lyrics {
				if CurrentConfig.Output.Romanization.IsEnabled() && IsSupportedAsianLang(lyric) {
					currentSong.Lyrics[i] = Romanize(lyric)
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
							currentSong.Lyrics[i] = fmt.Sprintf(lyric+" "+CurrentConfig.Output.RepeatedLyricsMultiplierFormat, count)
						} else {
							currentSong.Lyrics[i] = fmt.Sprintf(CurrentConfig.Output.RepeatedLyricsMultiplierFormat+" "+lyric, count)
						}
					}
				}
			}

			timeAfterGettingLyrics := time.Now()
			position += math.Max(timeAfterGettingLyrics.Sub(timeBeforeGettingLyrics).Seconds(), 0) + 0.1 // tests have shown that additional 0.1 is required to look good

			UpdateData(currentSong)
			UpdatePosition(position)
		}
	}()

	go WriteLyrics()
	go WriteInstrumental()
}
