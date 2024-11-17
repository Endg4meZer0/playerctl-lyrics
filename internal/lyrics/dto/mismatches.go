package dto

import (
	lrclibdto "lrcsnc/internal/lyrics/dto/lrclib"
	"lrcsnc/internal/pkg/structs"
	"math"
	"strings"
)

func RemoveMismatches(song structs.Song, lyricsData []LyricsDTO) []LyricsDTO {
	if len(lyricsData) == 0 {
		return lyricsData
	}

	var matchingLyricsData []LyricsDTO = make([]LyricsDTO, 0, len(lyricsData))

	switch lyricsData[0].(type) {
	case lrclibdto.LrcLibDTO:
		for _, lyrics := range lyricsData {
			if strings.EqualFold(lyrics.(lrclibdto.LrcLibDTO).Title, song.Title) &&
				// If player doesn't provide the song's duration, ignore the duration check
				// Otherwise, do a check that prevents different versions of a song of messing up the response
				((song.Duration != 0) == (math.Abs(float64(lyrics.(lrclibdto.LrcLibDTO).Duration)-song.Duration) <= 2)) {
				matchingLyricsData = append(matchingLyricsData, lyrics)
			}
		}
	default:
		return lyricsData
	}

	return matchingLyricsData
}
