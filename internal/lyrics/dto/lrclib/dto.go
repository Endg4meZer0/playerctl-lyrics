package lrclib

import (
	"lrcsnc/internal/pkg/structs"
	"lrcsnc/internal/pkg/util"
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
	syncedLyrics := strings.Split(dto.SyncedLyrics, "\n")

	resultLyrics := make([]string, len(syncedLyrics))
	resultTimestamps := make([]float64, len(syncedLyrics))

	for i, lyric := range syncedLyrics {
		lyricParts := strings.SplitN(lyric, " ", 2)
		timecode := util.TimecodeToFloat(lyricParts[0])
		if timecode == -1 {
			continue
		}
		var lyricStr string
		if len(lyricParts) != 1 {
			lyricStr = lyricParts[1]
		} else {
			lyricStr = ""
		}
		resultLyrics[i] = lyricStr
		resultTimestamps[i] = timecode
	}

	out.Lyrics = resultLyrics
	out.LyricTimestamps = resultTimestamps

	return
}
