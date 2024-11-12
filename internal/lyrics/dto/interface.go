package dto

import (
	lrclib "lrcsnc/internal/lyrics/dto/lrclib"
	"lrcsnc/internal/pkg/structs"
)

type LyricsDTO interface {
	ToLyricsData() structs.LyricsData
}

var LyricsDTOList map[string]LyricsDTO = map[string]LyricsDTO{
	"lrclib": lrclib.LrcLibDTO{},
}
