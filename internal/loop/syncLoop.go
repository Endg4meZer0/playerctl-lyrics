package loop

import (
	"fmt"
	"math"
	"time"

	"lrcsnc/internal/lyrics"
	"lrcsnc/internal/player"
	"lrcsnc/internal/romanization"
	"lrcsnc/pkg/global"
)

func SyncLoop() {
	checkerTicker := time.NewTicker(time.Duration(global.CurrentConfig.Player.SongCheckInterval*1000) * time.Millisecond)
	positionCheckTimer := time.NewTimer(time.Second)
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
			playerData := player.PlayerDataProviders[global.CurrentConfig.Player.PlayerProvider].GetPlayerData()
			song := player.PlayerDataProviders[global.CurrentConfig.Player.PlayerProvider].GetSongData()
			if song.Song != global.CurrentSong.Song || song.Artist != global.CurrentSong.Artist || song.Album != global.CurrentSong.Album || song.Duration != global.CurrentSong.Duration {
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
	go func() {
		for {
			<-positionCheckTimer.C
			var initialPosition float64
			global.CurrentPlayer = player.PlayerDataProviders[global.CurrentConfig.Player.PlayerProvider].GetPlayerData()
			initialPosition = global.CurrentPlayer.Position
			if global.CurrentPlayer.IsPlaying {
				if global.CurrentConfig.Global.EnableActiveSync {
					requiredTicks := 10
					positionInnerCheckTicker.Reset(100 * time.Millisecond)
					for i := 0; i < requiredTicks; i++ {
						<-positionInnerCheckTicker.C
						global.CurrentPlayer = player.PlayerDataProviders[global.CurrentConfig.Player.PlayerProvider].GetPlayerData()
						newPosition := global.CurrentPlayer.Position
						expectedPosition := (initialPosition + 0.1*(float64(i)+1))
						diff := newPosition - expectedPosition
						if !(((diff <= 0.21 && diff >= -1.11) || (diff >= 0.89 && diff <= 1.01)) || !global.CurrentPlayer.IsPlaying) {
							UpdatePosition(newPosition)
							break
						}
					}
				} else {
					positionInnerCheckTicker.Reset(time.Second)
					<-positionInnerCheckTicker.C
					global.CurrentPlayer = player.PlayerDataProviders[global.CurrentConfig.Player.PlayerProvider].GetPlayerData()
					newPosition := global.CurrentPlayer.Position
					diff := newPosition - initialPosition
					if !(diff >= 0.9 && diff <= 1.1) {
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
				if global.CurrentConfig.Output.Romanization.IsEnabled() && romanization.IsSupportedAsianLang(lyric) {
					global.CurrentSong.LyricsData.Lyrics[i] = romanization.Romanize(lyric)
				}
				if global.CurrentConfig.Output.ShowRepeatedLyricsMultiplier {
					if lyric == prevLyric && lyric != "" {
						count++
					} else {
						count = 1
						prevLyric = lyric
					}

					if count != 1 {
						if global.CurrentConfig.Output.PrintRepeatedLyricsMultiplierToTheRight {
							global.CurrentSong.LyricsData.Lyrics[i] = fmt.Sprintf(lyric+" "+global.CurrentConfig.Output.RepeatedLyricsMultiplierFormat, count)
						} else {
							global.CurrentSong.LyricsData.Lyrics[i] = fmt.Sprintf(global.CurrentConfig.Output.RepeatedLyricsMultiplierFormat+" "+lyric, count)
						}
					}
				}
			}

			if global.CurrentConfig.Output.TimestampOffset != 0 {
				for i, timestamp := range global.CurrentSong.LyricsData.LyricTimestamps {
					global.CurrentSong.LyricsData.LyricTimestamps[i] = timestamp + (float64(global.CurrentConfig.Output.TimestampOffset) / 1000)
				}
			}

			timeAfterGettingLyrics := time.Now()
			position += math.Max(timeAfterGettingLyrics.Sub(timeBeforeGettingLyrics).Seconds(), 0) + 0.1 // tests have shown that additional 0.1 is required to look good

			UpdatePosition(position)
		}
	}()

	go WriteLyrics()
	go WriteInstrumental()
}
