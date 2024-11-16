package lrclib

import (
	"lrcsnc/internal/pkg/structs"
	"strconv"
	"strings"
)

type LrcLibDTO struct {
	Title        string  `json:"trackName"`
	Artist       string  `json:"artistName"`
	Album        string  `json:"albumName"`
	Duration     float64 `json:"duration"`
	Instrumental bool    `json:"instrumental"`
	PlainLyrics  string  `json:"plainLyrics"`
	SyncedLyrics string  `json:"syncedLyrics"`
}

func (dto LrcLibDTO) ToLyricsData() (out structs.LyricsData) {
	if dto.Instrumental {
		out.LyricsType = 2
		return
	}

	if dto.PlainLyrics != "" && dto.SyncedLyrics == "" {
		out.Lyrics = strings.Split(dto.PlainLyrics, "\n")
		out.LyricsType = 1
		return
	}

	out.LyricsType = 0
	resultLyrics, resultTimestamps := parseSyncedLyrics(dto.SyncedLyrics)

	out.Lyrics = resultLyrics
	out.LyricTimestamps = resultTimestamps

	return
}

func parseSyncedLyrics(lyrics string) (out []string, timestamps []float64) {
	unsorted := false
	syncedLyrics := strings.Split(lyrics, "\n")

	resultLyrics := make([]string, 0, len(syncedLyrics))
	resultTimestamps := make([]float64, 0, len(syncedLyrics))

	for _, lyric := range syncedLyrics {
		lyricParts := strings.Split(lyric, "]")
		if len(lyricParts) < 2 {
			continue
		}

		lyricStr := strings.TrimRight(strings.TrimSpace(lyricParts[len(lyricParts)-1]), "\r")
		unsorted = unsorted || len(lyricParts) > 2

		for _, timecodeStr := range lyricParts[:len(lyricParts)-1] {
			timecode := timecodeToFloat(timecodeStr)
			if timecode == -1 {
				continue
			}
			resultLyrics = append(resultLyrics, lyricStr)
			resultTimestamps = append(resultTimestamps, timecode)
		}
	}

	if unsorted {
		// Sort the timestamps and lyrics
		// Not really a lot of data, and the "repeating" format of lyrics appears very rarely, so using bubble sort is not going to impact general performance
		for i := 0; i < len(resultTimestamps); i++ {
			for j := i + 1; j < len(resultTimestamps); j++ {
				if resultTimestamps[i] > resultTimestamps[j] {
					resultTimestamps[i], resultTimestamps[j] = resultTimestamps[j], resultTimestamps[i]
					resultLyrics[i], resultLyrics[j] = resultLyrics[j], resultLyrics[i]
				}
			}
		}
	}

	return resultLyrics, resultTimestamps
}

func timecodeToFloat(timecode string) float64 {
	// [01:23.45
	if len(timecode) != 9 {
		return -1
	}
	minutes, err := strconv.ParseFloat(timecode[1:3], 64)
	if err != nil {
		return -1
	}
	seconds, err := strconv.ParseFloat(timecode[4:9], 64)
	if err != nil {
		return -1
	}
	return minutes*60.0 + seconds
}
