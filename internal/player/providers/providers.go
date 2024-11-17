package providers

import (
	"lrcsnc/internal/pkg/structs"
)

type PlayerProvider interface {
	GetPlayerInfo() (structs.Player, error)
	// TODO: replace with signal-based getter
	GetSongInfo() (structs.Song, error)
}

var PlayerProviders map[string]PlayerProvider = map[string]PlayerProvider{
	"mpris": &MprisPlayerProvider{},
}
