package providers

import (
	"lrcsnc/internal/pkg/structs"
)

type PlayerProvider interface {
	// TODO: replace with signal-based getter
	GetInfo() (structs.Player, error)
	Subscribe() <-chan structs.Player
}

var PlayerProviders map[string]PlayerProvider = map[string]PlayerProvider{
	"mpris": &MprisPlayerProvider{},
}
