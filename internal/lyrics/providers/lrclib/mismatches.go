package lrclib

import (
	"lrcsnc/internal/lyrics/dto"
	lrclibdto "lrcsnc/internal/lyrics/dto/lrclib"
	"lrcsnc/internal/pkg/structs"
	"math"
)

func (l LrcLibLyricsProvider) RemoveMismatches(song structs.Song, lyricsData []dto.LyricsDTO) []dto.LyricsDTO {
	if len(lyricsData) == 0 {
		return lyricsData
	}

	var matchingLyricsData []dto.LyricsDTO = make([]dto.LyricsDTO, 0, len(lyricsData))

	for _, lyrics := range lyricsData {
		if lyrics.(lrclibdto.LrcLibDTO).Title == song.Title &&
			// If player doesn't provide the song's duration, ignore the duration check
			((song.Duration != 0) == (math.Abs(float64(lyrics.(lrclibdto.LrcLibDTO).Duration)-song.Duration) <= 2)) {
			matchingLyricsData = append(matchingLyricsData, lyrics)
		}
	}

	return matchingLyricsData
}
