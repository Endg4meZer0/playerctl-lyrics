package providers

import (
	"strings"

	"github.com/Pauloo27/go-mpris"
	"github.com/godbus/dbus/v5"

	"lrcsnc/internal/pkg/global"
	"lrcsnc/internal/pkg/structs"
)

type MprisPlayerProvider struct{}

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

func (m MprisPlayerProvider) GetSongInfo() (out structs.SongInfo) {
	if global.CurrentPlayer.PlayerName == "" {
		out.LyricsData.LyricsType = 4
		return
	}

	// Handling MPRIS directly is not really fail-safe, so to prevent any panics from reaching the main program there is a recover defer.
	defer func() {
		_ = recover()
	}()

	conn := InitConn()
	if conn == nil {
		out.LyricsData.LyricsType = 4
		return
	}
	defer conn.Close()

	player := mpris.New(conn, global.CurrentPlayer.PlayerName)

	metadata, err := player.GetMetadata()
	if err != nil {
		out.LyricsData.LyricsType = 4
		return
	}

	song, ok1 := metadata["xesam:title"]
	artist, ok2 := metadata["xesam:artist"]
	album, ok3 := metadata["xesam:album"]
	if !(ok1 && ok2 && ok3) {
		out.LyricsData.LyricsType = 4
		return
	}
	duration, ok := metadata["mpris:length"]

	out.Title = song.Value().(string)
	out.Artist = strings.Join(artist.Value().([]string), ", ")
	out.Album = album.Value().(string)
	out.LyricsData.LyricsType = 5

	if ok {
		out.Duration = float64(duration.Value().(uint64)) / 1000000
	}

	return
}

func (m MprisPlayerProvider) GetPlayerInfo() (out structs.PlayerInfo) {

	// Handling MPRIS directly is not really fail-safe, so to prevent any panics from reaching the main program there is a recover defer.
	defer func() {
		_ = recover()
	}()

	conn := InitConn()
	if conn == nil {
		return
	}
	defer conn.Close()

	playerNames, err := mpris.List(conn)
	if err != nil {
		return
	}
	if len(playerNames) == 0 {
		return
	}

	var player *mpris.Player
	playerName := ""

	if len(global.CurrentConfig.Player.IncludedPlayers) != 0 {
		found := false
		for _, v1 := range playerNames {
			for _, v2 := range global.CurrentConfig.Player.IncludedPlayers {
				if strings.Contains(v1, v2) {
					player = mpris.New(conn, v1)
					playerName = v1
					found = true
					break
				}
			}
			if found {
				break
			}
		}
	} else if len(global.CurrentConfig.Player.ExcludedPlayers) != 0 {
		found := false
		for _, v1 := range playerNames {
			for _, v2 := range global.CurrentConfig.Player.ExcludedPlayers {
				if !strings.Contains(v1, v2) {
					player = mpris.New(conn, v1)
					playerName = v1
					found = true
					break
				}
			}
			if found {
				break
			}
		}
	} else {
		player = mpris.New(conn, playerNames[0])
		playerName = playerNames[0]
	}

	status, err := player.GetPlaybackStatus()
	if err != nil {
		return
	}
	position, err := player.GetPosition()
	if err != nil {
		return
	}

	out.PlayerName = playerName
	out.IsPlaying = status == mpris.PlaybackPlaying
	out.Position = position

	return
}
