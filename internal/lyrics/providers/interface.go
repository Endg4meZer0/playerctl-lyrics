package providers

import (
	dto "lrcsnc/internal/lyrics/dto"
	lrclib "lrcsnc/internal/lyrics/providers/lrclib"
	"lrcsnc/internal/pkg/structs"
)

type LyricsDataProvider interface {
	GetLyricsData(structs.Song) (dto.LyricsDTO, error)
}

var LyricsDataProviders map[string]LyricsDataProvider = map[string]LyricsDataProvider{
	"lrclib": lrclib.LrcLibLyricsProvider{},
}
