package player

import (
	"lrcsnc/internal/pkg/structs"
	p "lrcsnc/internal/player/providers"
)

type PlayerInfoProvider interface {
	GetPlayerInfo() structs.PlayerInfo
	GetSongInfo() structs.SongInfo
}

var PlayerInfoProviders map[string]PlayerInfoProvider = map[string]PlayerInfoProvider{
	"mpris":     p.MprisPlayerProvider{},
	"playerctl": p.PlayerctlPlayerProvider{},
}
