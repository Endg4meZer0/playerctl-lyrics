package loop

import (
	"math"
	"time"

	"lrcsnc/internal/output"
	"lrcsnc/internal/pkg/global"
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
			if global.Player.Song.LyricsData.LyricsType >= 2 {
				output.OutputControllers[global.Config.Global.Output].DisplayCurrentLyric(-1)
			} else {
				// 5999.99s is basically the maximum limit of .lrc files' timestamps AFAIK, so 6000s is unreachable
				currentLyricTimestamp := -1.0
				nextLyricTimestamp := 6000.0
				timestampIndex := -1

				for i, timestamp := range global.Player.Song.LyricsData.LyricTimestamps {
					if timestamp <= currentPosition && currentLyricTimestamp <= timestamp {
						currentLyricTimestamp = timestamp
						timestampIndex = i
					}
				}

				if timestampIndex != len(global.Player.Song.LyricsData.LyricTimestamps)-1 {
					nextLyricTimestamp = global.Player.Song.LyricsData.LyricTimestamps[timestampIndex+1]
				}

				lyricsTimerDuration := time.Duration(int64(math.Abs(nextLyricTimestamp-currentPosition-0.01)*1000)) * time.Millisecond // tests have shown that it slows down and starts to mismatch without additional 0.01 offset

				if currentLyricTimestamp == -1 || (global.Player.IsPlaying && writtenTimestamp != currentLyricTimestamp) {
					output.OutputControllers[global.Config.Global.Output].DisplayCurrentLyric(timestampIndex)
				}

				writtenTimestamp = currentLyricTimestamp
				currentPosition = nextLyricTimestamp
				lyricsTimer.Reset(lyricsTimerDuration)
			}
		}
	}()
}
