package loop

import (
	"fmt"
	"math"
	"strings"
	"time"

	"lrcsnc/internal/lyrics"
	"lrcsnc/internal/output"
	"lrcsnc/internal/pkg/global"
	player "lrcsnc/internal/player/providers"
	"lrcsnc/internal/romanization"
)

func SyncLoop() {
	// Channels to communicate between goroutines
	playerChanges := player.PlayerProviders[global.Config.Player.PlayerProvider].Subscribe()
	songChanged := make(chan bool, 1)
	fullLyrChan := make(chan bool, 1)
	var err error

	// Initial player data
	global.Player, err = player.PlayerProviders[global.Config.Player.PlayerProvider].GetInfo()
	if err != nil {
		// TODO: logger :)
	}

	// Set up timers and tickers
	positionFollowTick := make(chan bool, 3)

	// Some additional syncing variables to make up for the time spent on getting the data itself
	position := 0.0
	var timeBeforeGettingLyrics time.Time

	// Goroutine that reacts to the changes in player's state
	go func() {
		for {
			newState := <-playerChanges
			positionFollowTick <- true

			if newState.Song.Title != global.Player.Song.Title || newState.Song.Artist != global.Player.Song.Artist || newState.Song.Album != global.Player.Song.Album || newState.Song.Duration != global.Player.Song.Duration {
				position = newState.Position
				timeBeforeGettingLyrics = time.Now()
				global.Player.Song = newState.Song

				songChanged <- true
			}
		}
	}()

	// Goroutine that handles the process of getting lyrics data
	go func() {
		for {
			<-songChanged

			output.OutputControllers[global.Config.Global.Output].OnSongInfoChange()

			// If the duration equals 0s, then there are no supported players out there.
			if global.Player.Song.Duration == 0 {
				global.Player.Song.LyricsData.LyricsType = 4
				fullLyrChan <- false
				continue
			}

			lyricsData, err := lyrics.GetLyricsData(global.Player.Song)
			if err != nil {
				// TODO: logger :)
			}

			global.Player.Song.LyricsData = lyricsData

			fullLyrChan <- true
		}
	}()

	// Goroutine to watch abnormal changes in player's position
	// For example, seeking on a position bar is counted as abnormal
	// TODO: replace with signal-based communication from D-Bus
	go func() {
		for {
			<-positionFollowTick
			UpdatePosition(global.Player.Position)
		}
	}()

	// Goroutine to update data in the output thread
	go func() {
		for {
			<-fullLyrChan

			prevLyric := ""
			count := 1
			for i, lyric := range global.Player.Song.LyricsData.Lyrics {
				lyric = strings.TrimSpace(strings.ReplaceAll(lyric, "\r", ""))

				// Apply lyrics multiplier
				if global.Config.Global.Output == "piped" && global.Config.Output.Piped.ShowRepeatedLyricsMultiplier {
					if lyric == prevLyric && lyric != "" {
						count++
					} else {
						count = 1
						prevLyric = lyric
					}

					if count != 1 {
						if global.Config.Output.Piped.PrintRepeatedLyricsMultiplierToTheRight {
							lyric = fmt.Sprintf(lyric+" "+global.Config.Output.Piped.RepeatedLyricsMultiplierFormat, count)
						} else {
							lyric = fmt.Sprintf(global.Config.Output.Piped.RepeatedLyricsMultiplierFormat+" "+lyric, count)
						}
					}
				}

				global.Player.Song.LyricsData.Lyrics[i] = lyric
			}

			// Romanization
			if lang := romanization.GetLang(global.Player.Song.LyricsData.Lyrics); global.Config.Lyrics.Romanization.IsEnabled() && lang != 0 {
				global.Player.Song.LyricsData.Lyrics = romanization.Romanize(global.Player.Song.LyricsData.Lyrics, lang)
			}

			if global.Config.Lyrics.TimestampOffset != 0 {
				for i, timestamp := range global.Player.Song.LyricsData.LyricTimestamps {
					global.Player.Song.LyricsData.LyricTimestamps[i] = timestamp + global.Config.Lyrics.TimestampOffset
				}
			}

			output.OutputControllers[global.Config.Global.Output].OnSongInfoChange()

			timeAfterGettingLyrics := time.Now()
			position += math.Max(timeAfterGettingLyrics.Sub(timeBeforeGettingLyrics).Seconds(), 0) + 0.1 // tests have shown that additional 0.1 is required to look good

			UpdatePosition(position)
		}
	}()

	// Launch the goroutine for lyrics syncing
	go SyncLyrics()
}
