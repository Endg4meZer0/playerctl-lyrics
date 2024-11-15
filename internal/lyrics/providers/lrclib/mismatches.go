package lrclib

import (
	"lrcsnc/internal/lyrics/dto"
	lrclibdto "lrcsnc/internal/lyrics/dto/lrclib"
	"lrcsnc/internal/pkg/structs"
	"math"
	"strings"
)

func (l LrcLibLyricsProvider) RemoveMismatches(song structs.Song, lyricsData []dto.LyricsDTO) []dto.LyricsDTO {
	if len(lyricsData) == 0 {
		return lyricsData
	}

	var matchingLyricsData []dto.LyricsDTO = make([]dto.LyricsDTO, 0, len(lyricsData))

	for _, lyrics := range lyricsData {
		if strings.EqualFold(lyrics.(lrclibdto.LrcLibDTO).Title, song.Title) &&
			// If player doesn't provide the song's duration, ignore the duration check
			// Otherwise, do a check that prevents different versions of a song of messing up the response
			((song.Duration != 0) == (math.Abs(float64(lyrics.(lrclibdto.LrcLibDTO).Duration)-song.Duration) <= 2)) {
			matchingLyricsData = append(matchingLyricsData, lyrics)
		}
	}

	return matchingLyricsData
}
