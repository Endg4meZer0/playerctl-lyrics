//go:build linux

package controllers

import (
	"fmt"
	"time"

	"github.com/Pauloo27/go-mpris"
	"github.com/godbus/dbus/v5"

	"lrcsnc/internal/player/sessions"
)

type MprisPlayerController struct{}

func (m MprisPlayerController) SeekTo(pos float64) error {
	var dbusconn = sessions.MediaSessions["mpris"].GetSession()
	if dbusconn == nil {
		return fmt.Errorf("[player/providers/mpris/SeekTo] ERROR: D-Bus connection does not exist at the moment")
	}
	var player = sessions.MediaSessions["mpris"].GetPlayer()
	if player == nil {
		return fmt.Errorf("[player/providers/mpris/SeekTo] ERROR: Player handler does not exist at the moment")
	}

	// Handling MPRIS directly is not really fail-safe, so to prevent any panics from reaching the main program there is a recover defer.
	defer func() {
		_ = recover()
	}()

	metadata, err := player.(*mpris.Player).GetMetadata()
	if err != nil {
		return fmt.Errorf("[player/providers/mpris/SeekTo] ERROR: Couldn't get metadata using existing player handler")
	}
	trackId, ok := metadata["mpris:trackid"]
	if !ok {
		return fmt.Errorf("[player/providers/mpris/SeekTo] ERROR: Couldn't get track ID from metadata using existing player handler")
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
		return fmt.Errorf("[player/providers/mpris/SeekTo] ERROR: Encountered a type mismatch when trying to handle track ID from metadata received by existing player handler")
	}
}

func (m MprisPlayerController) ApplySpotifyWorkaround() error {
	var dbusconn = sessions.MediaSessions["mpris"].GetSession()
	if dbusconn == nil {
		return fmt.Errorf("[player/providers/mpris/SetRepeat] ERROR: D-Bus connection does not exist at the moment")
	}
	var player = sessions.MediaSessions["mpris"].GetPlayer().(*mpris.Player)
	if player == nil {
		return fmt.Errorf("[player/providers/mpris/SetRepeat] ERROR: Player handler does not exist at the moment")
	}

	// Handling MPRIS directly is not really fail-safe, so to prevent any panics from reaching the main program there is a recover defer.
	defer func() {
		_ = recover()
	}()

	currentLoopStatus, err := player.GetLoopStatus()
	if err != nil {
		return fmt.Errorf("[player/providers/mpris/ApplySpotifyWorkaround] ERROR: Couldn't get loop status using existing player handler")
	}
	time.Sleep(1000)

	switch currentLoopStatus {
	case mpris.LoopTrack:
		_ = player.SetLoopStatus(mpris.LoopNone)
		time.Sleep(1000)
		err = player.SetLoopStatus(currentLoopStatus)
	default:
		_ = player.SetLoopStatus(mpris.LoopTrack)
		time.Sleep(1000)
		err = player.SetLoopStatus(currentLoopStatus)
	}

	return err
}
