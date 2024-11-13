package providers

import (
	dto "lrcsnc/internal/lyrics/dto"
	lrclib "lrcsnc/internal/lyrics/providers/lrclib"
	"lrcsnc/internal/pkg/structs"
)

type LyricsDataProvider interface {
	// Gets the lyrics data for a song in the form of LyricsDTO for later handling
	GetLyricsData(structs.Song) (dto.LyricsDTO, error)
	RemoveMismatches(structs.Song, []dto.LyricsDTO) []dto.LyricsDTO
}

var LyricsDataProviders map[string]LyricsDataProvider = map[string]LyricsDataProvider{
	"lrclib": lrclib.LrcLibLyricsProvider{},
}
