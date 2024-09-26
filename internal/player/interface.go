package player

import (
	p "lrcsnc/internal/player/providers"
	"lrcsnc/pkg/structs"
)

type PlayerDataProvider interface {
	GetPlayerData() structs.PlayerData
	GetSongData() structs.SongData
}

var PlayerDataProviders map[string]PlayerDataProvider = map[string]PlayerDataProvider{
	"mpris":     p.MprisPlayerProvider{},
	"playerctl": p.PlayerctlPlayerProvider{},
}
