package player

import (
	"fmt"
	"lrcsnc/internal/pkg/global"
	"lrcsnc/internal/player/sessions"
)

func InitSession() error {
	init, ok := sessions.MediaSessions[global.CurrentConfig.Player.PlayerProvider]
	if !ok {
		return fmt.Errorf("FATAL: The specified player provider doesn't have a session handler")
	}

	return init.Init()
}

func CloseSession() error {
	provider, ok := sessions.MediaSessions[global.CurrentConfig.Player.PlayerProvider]
	if !ok {
		return fmt.Errorf("FATAL: The specified player provider doesn't have a session handler")
	}

	return provider.Close()
}
