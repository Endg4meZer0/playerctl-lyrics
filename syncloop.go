package main

import (
	"time"
)

func SyncLoop() {
	var currentSong SongData
	var isPlaying bool

	lyricsTimer := time.NewTimer(time.Second)
	lyricsTimer.Stop()
	instrTimer := time.NewTimer(500 * time.Millisecond)        // using timer instead of ticker allows to use different durations when necessary without much thoughts
	instrTimer.Stop()                                          // stopping because the timer should be reset when it's needed
	go WriteInstrumental(instrTimer, &isPlaying, &currentSong) // also starting the instrumental thread at the same time to not create additional instances and only work with the ticker

	checkerTicker := time.NewTicker(time.Second)

	songChanged := make(chan bool, 1)
	fullLyrChan := make(chan map[float64]string, 1)

	// Goroutine to wait for incoming song metadata (lyrics and instrumental bool)
	go func() {
		for {
			<-songChanged

			lyricsTimer.Stop()
			instrTimer.Stop()

			// If the duration equals 0s, then there are no supported players out there.
			if currentSong.Duration == 0 {
				currentSong.LyricsType = 4
				fullLyrChan <- nil
				continue
			}

			lyr := GetSyncedLyrics(&currentSong)
			if currentSong.LyricsType == 5 {
				currentSong.LyricsType = 6
			}
			fullLyrChan <- lyr
		}
	}()

	// Goroutine to check for changes in currently playing song
	go func() {
		for {
			<-checkerTicker.C
			song := GetCurrentSongData()
			if song.Song != currentSong.Song || song.Artist != currentSong.Artist || song.Album != currentSong.Album || song.Duration != currentSong.Duration {
				currentSong = song
				songChanged <- true
			}
		}
	}()

	go func() {
		for {
			currentLyrics := <-fullLyrChan

			for i, lyric := range currentLyrics {
				if IsSupportedAsianLang(lyric) {
					currentLyrics[i] = Romanize(lyric)
				}
			}

			go WriteLyrics(lyricsTimer, instrTimer, &currentLyrics, &isPlaying, &currentSong, "", 0)
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
