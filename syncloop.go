package main

import (
	"time"
)

func SyncLoop() {
	var currentSong SongData
	var currentLyrics map[float64]string
	var currentlyInstrumental bool
	var currentSongIsNotFound bool
	var noPlayersFound bool
	var isPlaying bool

	lyricsTimer := time.NewTimer(time.Second)
	lyricsTimer.Stop()
	instrTicker := time.NewTicker(500 * time.Millisecond)
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
