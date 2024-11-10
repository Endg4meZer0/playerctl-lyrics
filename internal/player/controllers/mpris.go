//go:build linux

package controllers

import (
	"fmt"

	"github.com/Pauloo27/go-mpris"
	"github.com/godbus/dbus/v5"

	"lrcsnc/internal/player/sessions"
)

type MprisPlayerController struct{}

func (m *MprisPlayerController) SeekTo(pos float64) error {
	var dbusconn = sessions.MediaSessions["mpris"].GetSession()
	if dbusconn == nil {
		return fmt.Errorf("[player/providers/mpris/GetSongInfo] ERROR: D-Bus connection does not exist at the moment")
	}
	var player = sessions.MediaSessions["mpris"].GetPlayer()
	if player == nil {
		return fmt.Errorf("[player/providers/mpris/GetSongInfo] ERROR: Player handler does not exist at the moment")
	}

	// Handling MPRIS directly is not really fail-safe, so to prevent any panics from reaching the main program there is a recover defer.
	defer func() {
		_ = recover()
	}()

	metadata, err := player.(*mpris.Player).GetMetadata()
	if err != nil {
		return fmt.Errorf("[player/providers/mpris/GetSongInfo] ERROR: Couldn't get metadata using existing player handler")
	}
	trackId, ok := metadata["mpris:trackid"]
	if !ok {
		return fmt.Errorf("[player/providers/mpris/GetSongInfo] ERROR: Couldn't get track ID from metadata using existing player handler")
	}

	switch tId := trackId.Value().(type) {
	case *dbus.ObjectPath:
		return player.(*mpris.Player).SetTrackPosition(tId, pos)
	case dbus.ObjectPath:
		return player.(*mpris.Player).SetTrackPosition(&tId, pos)
	case string:
		objPath := dbus.ObjectPath(tId)
		return player.(*mpris.Player).SetTrackPosition(&objPath, pos)
	default:
		return fmt.Errorf("[player/providers/mpris/GetSongInfo] ERROR: Encountered a type mismatch when trying to handle track ID from metadata received by existing player handler")
	}
}
