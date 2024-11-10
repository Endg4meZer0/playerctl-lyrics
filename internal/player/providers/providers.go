package providers

import (
	"lrcsnc/internal/pkg/structs"
)

type PlayerProvider interface {
	GetPlayerInfo() (structs.PlayerInfo, error)
	// TODO: replace with signal-based getter
	GetSongInfo() (structs.SongInfo, error)
}

var PlayerProviders map[string]PlayerProvider = map[string]PlayerProvider{
	"mpris":     &MprisPlayerProvider{},
	"playerctl": &PlayerctlPlayerProvider{},
}
