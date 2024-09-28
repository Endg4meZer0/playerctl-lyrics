package controllers

import (
	"github.com/Pauloo27/go-mpris"
	"github.com/godbus/dbus/v5"

	"lrcsnc/internal/pkg/global"
)

type MprisPlayerController struct{}

func InitConn() (conn *dbus.Conn) {
	var err error
	conn, err = dbus.SessionBusPrivate()
	if err != nil {
		conn.Close()
		conn = nil
		return
	}
	err = conn.Auth(nil)
	if err != nil {
		conn.Close()
		conn = nil
		return
	}
	err = conn.Hello()
	if err != nil {
		conn.Close()
		conn = nil
		return
	}
	return
}

func (m MprisPlayerController) SeekTo(pos float64) bool {
	if global.CurrentPlayer.PlayerName == "" {
		return false
	}

	// Handling MPRIS directly is not fail-safe, so to prevent various crashes there is a recover defer.
	defer func() {
		if r := recover(); r != nil {
		}
	}()

	conn := InitConn()
	if conn == nil {
		return false
	}
	defer conn.Close()

	player := mpris.New(conn, global.CurrentPlayer.PlayerName)

	err := player.SetPosition(pos)

	return err == nil
}
