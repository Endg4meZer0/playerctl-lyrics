package loop

import (
	"fmt"
	"math"
	"time"

	"lrcsnc/internal/lrclib"
	"lrcsnc/internal/player"
	"lrcsnc/internal/romanization"
	"lrcsnc/pkg"
)

func SyncLoop() {
	checkerTicker := time.NewTicker(time.Duration(pkg.CurrentConfig.Playerctl.PlayerctlSongCheckInterval*1000) * time.Millisecond)
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
			song := player.GetSongData()
			if song.Song != pkg.CurrentSong.Song || song.Artist != pkg.CurrentSong.Artist || song.Album != pkg.CurrentSong.Album || song.Duration != pkg.CurrentSong.Duration {
				_, position = player.GetPlayerData()
				timeBeforeGettingLyrics = time.Now()
				pkg.CurrentSong = song

				songChanged <- true
			}
		}
	}()

	// Goroutine to wait for incoming song metadata (lyrics and instrumental bool)
	go func() {
		for {
			<-songChanged

			// If the duration equals 0s, then there are no supported players out there.
			if pkg.CurrentSong.Duration == 0 {
				pkg.CurrentSong.LyricsType = 4
				fullLyrChan <- false
				continue
			}

			lrclib.GetSyncedLyrics(&pkg.CurrentSong)
			if pkg.CurrentSong.LyricsType == 5 {
				pkg.CurrentSong.LyricsType = 6
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
			pkg.IsPlaying, initialPosition = player.GetPlayerData()
			if pkg.IsPlaying {
				if pkg.CurrentConfig.Global.DisableActiveSync {
					positionInnerCheckTicker.Reset(time.Second)
					<-positionInnerCheckTicker.C
					var newPosition float64
					pkg.IsPlaying, newPosition = player.GetPlayerData()
					diff := newPosition - initialPosition
					if !(diff >= 0.9 && diff <= 1.1) {
						UpdatePosition(newPosition)
					}
				} else {
					requiredTicks := 10
					positionInnerCheckTicker.Reset(100 * time.Millisecond)
					for i := 0; i < requiredTicks; i++ {
						<-positionInnerCheckTicker.C
						var newPosition float64
						pkg.IsPlaying, newPosition = player.GetPlayerData()
						expectedPosition := (initialPosition + 0.1*(float64(i)+1))
						diff := newPosition - expectedPosition
						if !(((diff <= 0.21 && diff >= -1.11) || (diff >= 0.89 && diff <= 1.01)) || !pkg.IsPlaying) {
							UpdatePosition(newPosition)
							break
						}
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
			for i, lyric := range pkg.CurrentSong.Lyrics {
				if pkg.CurrentConfig.Output.Romanization.IsEnabled() && romanization.IsSupportedAsianLang(lyric) {
					pkg.CurrentSong.Lyrics[i] = romanization.Romanize(lyric)
				}
				if pkg.CurrentConfig.Output.ShowRepeatedLyricsMultiplier {
					if lyric == prevLyric && lyric != "" {
						count++
					} else {
						count = 1
						prevLyric = lyric
					}

					if count != 1 {
						if pkg.CurrentConfig.Output.PrintRepeatedLyricsMultiplierToTheRight {
							pkg.CurrentSong.Lyrics[i] = fmt.Sprintf(lyric+" "+pkg.CurrentConfig.Output.RepeatedLyricsMultiplierFormat, count)
						} else {
							pkg.CurrentSong.Lyrics[i] = fmt.Sprintf(pkg.CurrentConfig.Output.RepeatedLyricsMultiplierFormat+" "+lyric, count)
						}
					}
				}
			}

			if pkg.CurrentConfig.Output.TimestampOffset != 0 {
				for i, timestamp := range pkg.CurrentSong.LyricTimestamps {
					pkg.CurrentSong.LyricTimestamps[i] = timestamp + (float64(pkg.CurrentConfig.Output.TimestampOffset) / 1000)
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
