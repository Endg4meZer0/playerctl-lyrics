package main

import (
	"fmt"
	"math"
	"strings"
	"time"
)

var lyricsTimer = time.NewTimer(100 * time.Minute)
var instrTimer = time.NewTimer(500 * time.Millisecond)
var currentLyrics map[float64]string
var currentSong = SongData{LyricsType: 5}
var prevLyric = ""
var lyricsRepeated uint = 0
var isPlaying = false

func UpdateData(newLyrics map[float64]string, newSong SongData) {
	currentLyrics = newLyrics
	currentSong = newSong
	lyricsTimer.Reset(100)
	instrTimer.Stop()
}

func WriteLyrics() {
	go func() {
		for {
			<-lyricsTimer.C
			if currentSong.LyricsType == 4 {
				instrTimer.Stop()
				fmt.Println()
			} else if currentSong.LyricsType >= 2 {
				instrTimer.Reset(1)
			} else {
				wasPaused := !isPlaying
				isPlaying = GetCurrentSongStatus()
				currentTimestamp := GetCurrentSongPosition()

				// if a floating part is less than that value down there (tested on cmus, may differ between players)
				// then make an assumption that the player uses integers as position markers
				// and we'll allow the timestamps to be equal or greater only by seconds, not milliseconds
				// 99% sure it can be done better but since it works as of now...
				_, currentTimestampFloatPart := math.Modf(currentTimestamp)
				playerUsesIntegerPosition := currentTimestampFloatPart < 0.000100

				// 5999.99s is basically the maximum limit of .lrc files' timestamps, so 6000s is unreachable
				firstLyricTimestamp := 6000.0
				currentLyricTimestamp := -1.0
				nextLyricTimestamp := 6000.0
				lyric := ""

				for lyricTimestamp, l := range currentLyrics {
					if firstLyricTimestamp > lyricTimestamp {
						firstLyricTimestamp = lyricTimestamp
					}
					if lyricTimestamp < currentTimestamp && currentLyricTimestamp <= lyricTimestamp {
						currentLyricTimestamp = lyricTimestamp
						lyric = l
					}
					if lyricTimestamp > currentTimestamp && nextLyricTimestamp > lyricTimestamp {
						nextLyricTimestamp = lyricTimestamp
					}
				}

				if lyric == prevLyric {
					lyricsRepeated++
				} else {
					lyricsRepeated = 1
				}

				if nextLyricTimestamp == 6000.0 {
					// If the nextLyricTimestamp remained at 6000s, then there are no more lyrics.
					// If that's the case, we'll need to account that the same song may be put on repeat
					// So the idea would be to change the value nextLyricTimestamp to the playing song's duration
					// and maybe add a bit like 0.1s to be 100% sure we're on a new song iteration
					nextLyricTimestamp = math.Abs(currentSong.Duration) + 0.1
				}

				prevLyric = lyric

				// If the currentTimestamp is less than even the first timestamp of the lyrics
				// then reset an instrumental ticker until the first lyric shows up
				if currentTimestamp < firstLyricTimestamp {
					instrTimer.Reset(1)
				} else if isPlaying { // If paused then don't print the lyric and instead try once more time later
					if lyric == "" {
						// An empty lyric basically means instrumental part,
						// so we reset the instrumental ticker and moving on
						instrTimer.Reset(1)
					} else {
						// An actual lyric when all the conditions are met needs to
						// 1) stop instrumental ticker
						// 2) print itself
						// 3) call the next writing goroutine
						instrTimer.Stop()
						if !wasPaused && (!playerUsesIntegerPosition || math.Abs(nextLyricTimestamp-currentTimestamp) >= 1.0) { // if the playback was paused, that usually causes lyric to print itself twice, so here's a little fuse
							if CurrentConfig.Output.ShowRepeatedLyricsMultiplier && lyricsRepeated > 1 {
								if CurrentConfig.Output.PrintRepeatedLyricsMultiplierToTheRight {
									fmt.Print(lyric + " ")
									fmt.Println(fmt.Sprintf(CurrentConfig.Output.RepeatedLyricsMultiplierFormat, lyricsRepeated))
								} else {
									fmt.Print(fmt.Sprintf(CurrentConfig.Output.RepeatedLyricsMultiplierFormat, lyricsRepeated) + " ")
									fmt.Println(lyric)
								}
							} else {
								fmt.Println(lyric)
							}
						}
					}
				}
				lyricsTimerDuration := time.Duration(int64(math.Abs(nextLyricTimestamp-currentTimestamp)*1000)) * time.Millisecond
				lyricsTimer.Reset(lyricsTimerDuration)
				if lyricsTimerDuration/time.Millisecond > 1500 {
					positionCheckTicker := time.NewTicker(1.5 * 1000 * time.Millisecond)
					expectedTicks := int(math.Floor(float64(lyricsTimerDuration/time.Millisecond/1000) / 1.5))
					currentTick := 0
					// Resets the lyric timer if it sees an unusual position change
					go func() {
						for {
							<-positionCheckTicker.C
							currentTick++
							receivedPosition := GetCurrentSongPosition()
							if receivedPosition < math.Floor(currentLyricTimestamp) || receivedPosition > math.Ceil(nextLyricTimestamp) || currentTick >= expectedTicks {
								positionCheckTicker.Stop()
								if currentTick < expectedTicks {
									lyricsTimer.Reset(1)
								}
								return
							}
						}
					}()
				}
			}
		}
	}()
}

// instrTimer.Stop to stop writing instrumental
// instrTimer.Reset to continue again
func WriteInstrumental() {
	note := CurrentConfig.Output.Instrumental.Symbol
	i := 1
	j := int(CurrentConfig.Output.Instrumental.MaxCount + 1)
	instrTimer.Reset(time.Duration(CurrentConfig.Output.Instrumental.Interval*1000) * time.Millisecond)
	for {
		<-instrTimer.C
		isPlaying = GetCurrentSongStatus()
		// Not playing? Don't change anything, or it will look kinda strange
		if !isPlaying {
			instrTimer.Reset(time.Duration(CurrentConfig.Output.Instrumental.Interval*1000) * time.Millisecond)
			continue
		} else {
			messagePrinted := false
			switch currentSong.LyricsType {
			case 1:
				if CurrentConfig.Output.ShowNotSyncedLyricsWarning {
					fmt.Println("This song's lyrics are not synced on LrcLib! " + strings.Repeat(note, i%j))
					messagePrinted = true
				}
			case 3:
				if CurrentConfig.Output.ShowSongNotFoundWarning {
					fmt.Println("Current song was not found on LrcLib! " + strings.Repeat(note, i%j))
					messagePrinted = true
				}
			case 5:
				if CurrentConfig.Output.ShowGettingLyricsMessage {
					fmt.Println("Getting lyrics... " + strings.Repeat(note, i%j))
					messagePrinted = true
				}
			case 6:
				fmt.Println("Failed to get lyrics! " + strings.Repeat(note, i%j))
				messagePrinted = true
			}

			if !messagePrinted {
				fmt.Println(strings.Repeat(note, i%j))
			}

			i++
			// Don't want to cause any overflow here
			if i > j-1 {
				i = 1
			}
			instrTimer.Reset(time.Duration(CurrentConfig.Output.Instrumental.Interval*1000) * time.Millisecond)
		}
	}
}
