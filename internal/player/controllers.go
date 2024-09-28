package player

import (
	c "lrcsnc/internal/player/controllers"
)

type PlayerInfoController interface {
	SeekTo(float64) bool
}

var PlayerInfoControllers map[string]PlayerInfoController = map[string]PlayerInfoController{
	"mpris":     c.MprisPlayerController{},
	"playerctl": c.PlayerctlPlayerController{},
}
