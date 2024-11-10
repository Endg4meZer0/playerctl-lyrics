package lyrics

import (
	lrclib "lrcsnc/internal/lyrics/providers/lrclib"
	"lrcsnc/internal/pkg/structs"
)

type LyricsDataProvider interface {
	GetLyricsData(structs.SongInfo) structs.SongLyricsData
}

var LyricsDataProviders map[string]LyricsDataProvider = map[string]LyricsDataProvider{
	"lrclib": lrclib.LrcLibLyricsProvider{},
}

// TODO: separate getter from DTO-controller
