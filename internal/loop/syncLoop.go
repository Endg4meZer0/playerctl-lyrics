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
	songChanged := make(chan bool, 1)
	fullLyrChan := make(chan bool, 1)
	var err error

	// Initial player data
	global.CurrentPlayer, err = player.PlayerProviders[global.CurrentConfig.Player.PlayerProvider].GetPlayerInfo()
	if err != nil {
		// TODO: logger :)
	}

	global.CurrentSong, err = player.PlayerProviders[global.CurrentConfig.Player.PlayerProvider].GetSongInfo()
	if err != nil {
		// TODO: logger :)
	}

	songChanged <- true

	// Set up timers and tickers
	checkerTicker := time.NewTicker(time.Duration(global.CurrentConfig.Player.SongCheckInterval*1000) * time.Millisecond)
	positionCheckTimer := time.NewTimer(time.Second)
	positionInnerCheckTicker := time.NewTicker(time.Second)
	positionInnerCheckTicker.Stop()

	// Some additional syncing variables to make up for the time spent on getting the data itself
	position := 0.0
	var timeBeforeGettingLyrics time.Time

	// Goroutine to check for changes in currently playing song
	go func() {
		for {
			<-checkerTicker.C
			playerData, err := player.PlayerProviders[global.CurrentConfig.Player.PlayerProvider].GetPlayerInfo()
			if err != nil {
				// TODO: logger :)
				continue
			}
			song, err := player.PlayerProviders[global.CurrentConfig.Player.PlayerProvider].GetSongInfo()
			if err != nil {
				// TODO: logger :)
				continue
			}
			if song.Title != global.CurrentSong.Title || song.Artist != global.CurrentSong.Artist || song.Album != global.CurrentSong.Album || song.Duration != global.CurrentSong.Duration {
				position = playerData.Position
				timeBeforeGettingLyrics = time.Now()
				global.CurrentSong = song

				songChanged <- true
			}
		}
	}()

	// Goroutine to wait for incoming song metadata (lyrics and instrumental bool)
	go func() {
		for {
			<-songChanged

			output.OutputControllers[global.CurrentConfig.Global.Output].OnSongInfoChange()

			// If the duration equals 0s, then there are no supported players out there.
			if global.CurrentSong.Duration == 0 {
				global.CurrentSong.LyricsData.LyricsType = 4
				fullLyrChan <- false
				continue
			}

			global.CurrentSong.LyricsData = lyrics.LyricsDataProviders[global.CurrentConfig.Global.LyricsProvider].GetLyricsData(global.CurrentSong)
			if global.CurrentSong.LyricsData.LyricsType == 5 {
				global.CurrentSong.LyricsData.LyricsType = 6
			}

			fullLyrChan <- true
		}
	}()

	// Goroutine to watch abnormal changes in player's position
	// For example, seeking on a position bar is counted as abnormal
	// TODO: replace with signal-based communication from D-Bus
	go func() {
		for {
			<-positionCheckTimer.C
			var initialPosition float64
			var err error
			global.CurrentPlayer, err = player.PlayerProviders[global.CurrentConfig.Player.PlayerProvider].GetPlayerInfo()
			if err != nil {
				// TODO: logger :)
				continue
			}

			initialPosition = global.CurrentPlayer.Position
			if global.CurrentPlayer.IsPlaying {
				if global.CurrentConfig.Global.EnableActiveSync {
					requiredTicks := 10
					positionInnerCheckTicker.Reset(100 * time.Millisecond)
					for i := 0; i < requiredTicks; i++ {
						<-positionInnerCheckTicker.C
						global.CurrentPlayer, err = player.PlayerProviders[global.CurrentConfig.Player.PlayerProvider].GetPlayerInfo()
						if err != nil {
							// TODO: logger :)
							continue
						}

						newPosition := global.CurrentPlayer.Position
						expectedPosition := (initialPosition + 0.1*(float64(i)+1))
						diff := newPosition - expectedPosition
						if !(((diff <= 0.21 && diff >= -1.11) || (diff >= 0.89 && diff <= 1.01)) || !global.CurrentPlayer.IsPlaying) {
							output.OutputControllers[global.CurrentConfig.Global.Output].OnPlayerInfoChange()
							UpdatePosition(newPosition)
							break
						}
					}
				} else {
					positionInnerCheckTicker.Reset(time.Second)
					<-positionInnerCheckTicker.C
					global.CurrentPlayer, err = player.PlayerProviders[global.CurrentConfig.Player.PlayerProvider].GetPlayerInfo()
					if err != nil {
						// TODO: logger :)
						continue
					}
					newPosition := global.CurrentPlayer.Position
					diff := newPosition - initialPosition
					if !(diff >= 0.9 && diff <= 1.1) {
						output.OutputControllers[global.CurrentConfig.Global.Output].OnPlayerInfoChange()
						UpdatePosition(newPosition)
					}
				}
				positionInnerCheckTicker.Stop()
				positionCheckTimer.Reset(1)
			} else {
				positionCheckTimer.Reset(time.Second)
			}
		}
	}()

	// Goroutine to update data in the output thread
	go func() {
		for {
			<-fullLyrChan

			prevLyric := ""
			count := 1
			for i, lyric := range global.CurrentSong.LyricsData.Lyrics {
				lyric = strings.TrimSpace(lyric)

				if global.CurrentConfig.Lyrics.Romanization.IsEnabled() && romanization.IsSupportedAsianLang(lyric) {
					lyric = romanization.Romanize(lyric)
				}
				if global.CurrentConfig.Global.Output == "piped" && global.CurrentConfig.Output.Piped.ShowRepeatedLyricsMultiplier {
					if lyric == prevLyric && lyric != "" {
						count++
					} else {
						count = 1
						prevLyric = lyric
					}

					if count != 1 {
						if global.CurrentConfig.Output.Piped.PrintRepeatedLyricsMultiplierToTheRight {
							lyric = fmt.Sprintf(lyric+" "+global.CurrentConfig.Output.Piped.RepeatedLyricsMultiplierFormat, count)
						} else {
							lyric = fmt.Sprintf(global.CurrentConfig.Output.Piped.RepeatedLyricsMultiplierFormat+" "+lyric, count)
						}
					}
				}

				global.CurrentSong.LyricsData.Lyrics[i] = lyric
			}

			if global.CurrentConfig.Lyrics.TimestampOffset != 0 {
				for i, timestamp := range global.CurrentSong.LyricsData.LyricTimestamps {
					global.CurrentSong.LyricsData.LyricTimestamps[i] = timestamp + global.CurrentConfig.Lyrics.TimestampOffset
				}
			}

			output.OutputControllers[global.CurrentConfig.Global.Output].OnSongInfoChange()

			timeAfterGettingLyrics := time.Now()
			position += math.Max(timeAfterGettingLyrics.Sub(timeBeforeGettingLyrics).Seconds(), 0) + 0.1 // tests have shown that additional 0.1 is required to look good

			UpdatePosition(position)
		}
	}()

	// Launch the goroutine for lyrics syncing
	go SyncLyrics()
}
