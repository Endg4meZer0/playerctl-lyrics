package lyrics

import (
	lrclib "lrcsnc/internal/lyrics/providers/lrclib"
	"lrcsnc/pkg/structs"
)

type LyricsDataProvider interface {
	GetLyricsData(structs.SongData) structs.SongLyricsData
}

var LyricsDataProviders map[string]LyricsDataProvider = map[string]LyricsDataProvider{
	"lrclib": lrclib.LrcLibLyricsProvider{},
}
