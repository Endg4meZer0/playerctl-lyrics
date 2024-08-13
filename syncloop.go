package main

import (
	"time"
)

func SyncLoop() {
	var currentSong SongData
	var currentLyrics map[float64]string

	checkerTicker := time.NewTicker(time.Second)

	songChanged := make(chan bool, 1)
	fullLyrChan := make(chan bool, 1)

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

			lyr := GetSyncedLyrics(&currentSong)
			if currentSong.LyricsType == 5 {
				currentSong.LyricsType = 6
			}

			currentLyrics = lyr
			fullLyrChan <- true
		}
	}()

	// Goroutine to check for changes in currently playing song
	go func() {
		for {
			<-checkerTicker.C
			song := GetCurrentSongData()
			if song.Song != currentSong.Song || song.Artist != currentSong.Artist || song.Album != currentSong.Album || song.Duration != currentSong.Duration {
				currentSong = song
				UpdateData(currentLyrics, currentSong)

				songChanged <- true
				//currentLyrics, currentlyInstrumental := GetSyncedLyrics(song)
			}
		}
	}()

	go func() {
		for {
			if !<-fullLyrChan {
				currentLyrics = nil
			}
			UpdateData(currentLyrics, currentSong)
			/*
				Timer is made like this:
				1) get the lyric from the map based on timestamp (we need the next lyric AFTER that timestamp)
				2) make a goroutine that writes the lyric to stdout when the timestamp comes
				3) should be recursive, so after the timer ticks, everything should begin from 1 again.
				Probably will make an extra function for that
			*/
		}
	}()

	go WriteLyrics()
	go WriteInstrumental()
}
