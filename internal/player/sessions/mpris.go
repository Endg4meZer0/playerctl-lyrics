package sessions

import (
	"fmt"

	"github.com/Pauloo27/go-mpris"
	"github.com/godbus/dbus/v5"
)

type DBusSession struct {
	conn   *dbus.Conn
	player *mpris.Player
}

func (s *DBusSession) Init() error {
	var err error
	s.conn, err = dbus.SessionBusPrivate()
	if err != nil {
		s.conn.Close()
		return err
	}
	err = s.conn.Auth(nil)
	if err != nil {
		s.conn.Close()
		return err
	}
	err = s.conn.Hello()
	if err != nil {
		s.conn.Close()
		return err
	}

	return nil
}

func (s *DBusSession) GetSession() interface{} {
	return s.conn
}

func (s *DBusSession) GetPlayer() interface{} {
	return s.player
}

func (s *DBusSession) SetPlayer(player interface{}) error {
	switch player := player.(type) {
	case *mpris.Player:
		s.player = player
		return nil
	default:
		return fmt.Errorf("[player/sessions/mpris] ERROR: Failed to set the new player in session struct because of type mismatch")
	}
}

func (s *DBusSession) Close() error {
	return s.conn.Close()
}
