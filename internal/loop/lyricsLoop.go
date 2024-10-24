package loop

import (
	"math"
	"time"

	"lrcsnc/internal/output"
	"lrcsnc/internal/pkg/global"
	"lrcsnc/internal/player"
)

var lyricsTimer = time.NewTimer(5 * time.Second)
var currentPosition = 0.0
var writtenTimestamp = 0.0

func UpdatePosition(newPosition float64) {
	currentPosition = newPosition
	lyricsTimer.Reset(1)
}

func SyncLyrics() {
	go func() {
		for {
			<-lyricsTimer.C
			if global.CurrentSong.LyricsData.LyricsType >= 2 {
				output.OutputControllers[global.CurrentConfig.Global.Output].DisplayCurrentLyric(-1)
			} else {
				playerData := player.PlayerInfoProviders[global.CurrentConfig.Player.PlayerProvider].GetPlayerInfo()

				if math.Abs(currentPosition-playerData.Position) > 1 {
					currentPosition = playerData.Position
				}

				// 5999.99s is basically the maximum limit of .lrc files' timestamps AFAIK, so 6000s is unreachable
				currentLyricTimestamp := -1.0
				nextLyricTimestamp := 6000.0
				timestampIndex := -1

				for i, timestamp := range global.CurrentSong.LyricsData.LyricTimestamps {
					if timestamp <= currentPosition && currentLyricTimestamp <= timestamp {
						currentLyricTimestamp = timestamp
						timestampIndex = i
					}
				}

				if timestampIndex != len(global.CurrentSong.LyricsData.LyricTimestamps)-1 {
					nextLyricTimestamp = global.CurrentSong.LyricsData.LyricTimestamps[timestampIndex+1]
				}

				lyricsTimerDuration := time.Duration(int64(math.Abs(nextLyricTimestamp-currentPosition-0.01)*1000)) * time.Millisecond // tests have shown that it slows down and starts to mismatch without additional 0.01 offset

				if currentLyricTimestamp == -1 || (global.CurrentPlayer.IsPlaying && writtenTimestamp != currentLyricTimestamp) {
					output.OutputControllers[global.CurrentConfig.Global.Output].DisplayCurrentLyric(timestampIndex)
				}

				writtenTimestamp = currentLyricTimestamp
				currentPosition = nextLyricTimestamp
				lyricsTimer.Reset(lyricsTimerDuration)
			}
		}
	}()
}
