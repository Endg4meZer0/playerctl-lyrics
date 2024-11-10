package providers

import (
	"fmt"
	"strings"

	"github.com/Pauloo27/go-mpris"
	"github.com/godbus/dbus/v5"

	"lrcsnc/internal/pkg/global"
	"lrcsnc/internal/pkg/structs"
	"lrcsnc/internal/player/sessions"
)

type MprisPlayerProvider struct{}

func (m *MprisPlayerProvider) GetSongInfo() (structs.SongInfo, error) {
	var out structs.SongInfo
	var dbusconn = sessions.MediaSessions["mpris"].GetSession()
	if dbusconn == nil {
		out.LyricsData.LyricsType = 4
		return out, fmt.Errorf("[player/providers/mpris/GetSongInfo] ERROR: D-Bus connection does not exist at the moment")
	}
	var player = sessions.MediaSessions["mpris"].GetPlayer()
	if player == nil {
		out.LyricsData.LyricsType = 4
		return out, fmt.Errorf("[player/providers/mpris/GetSongInfo] ERROR: Player handler does not exist at the moment")
	}

	// FIXME: Handling MPRIS directly is not really fail-safe yet,
	// so to prevent any panics from crashing the whole app there is a recover defer.
	// That should be changed ASAP, but it requires a significant rewrite of player listener's logic.
	defer func() {
		_ = recover()
	}()

	metadata, err := player.(*mpris.Player).GetMetadata()
	if err != nil {
		out.LyricsData.LyricsType = 4
		return out, fmt.Errorf("[player/providers/mpris/GetSongInfo] ERROR: Failed to get metadata from existing player handler")
	}

	song, ok1 := metadata["xesam:title"]
	artist, ok2 := metadata["xesam:artist"]
	album, ok3 := metadata["xesam:album"]
	if !(ok1 && ok2 && ok3) {
		out.LyricsData.LyricsType = 4
		return out, nil
	}
	duration, ok := metadata["mpris:length"]

	out.Title = song.Value().(string)
	out.Artist = strings.Join(artist.Value().([]string), ", ")
	out.Album = album.Value().(string)
	out.LyricsData.LyricsType = 5

	if ok {
		out.Duration = float64(getU64(duration.Value())) / 1000000
	}

	return out, nil
}

func (m *MprisPlayerProvider) GetPlayerInfo() (structs.PlayerInfo, error) {
	var out structs.PlayerInfo
	var dbusconn = sessions.MediaSessions["mpris"].GetSession()
	if dbusconn.(*dbus.Conn) == nil {
		return out, fmt.Errorf("[player/providers/mpris/GetPlayerInfo] ERROR: D-Bus connection does not exist at the moment")
	}

	var connPlayer = sessions.MediaSessions["mpris"].GetPlayer()
	if connPlayer.(*mpris.Player) != nil {
		identity, err := connPlayer.(*mpris.Player).GetIdentity()
		if err == nil {
			out.PlayerName = identity
			playing, err := connPlayer.(*mpris.Player).GetPlaybackStatus()
			if err != nil {
				return structs.PlayerInfo{}, fmt.Errorf("[player/providers/mpris/GetPlayerInfo] ERROR: Failed to get playback status using existing player handler")
			}
			out.IsPlaying = playing == mpris.PlaybackPlaying
			pos, err := connPlayer.(*mpris.Player).GetPosition()
			if err != nil {
				return structs.PlayerInfo{}, fmt.Errorf("[player/providers/mpris/GetPlayerInfo] ERROR: Failed to get player position using existing player handler")
			}
			out.Position = pos
			return out, nil
		}
	}

	// FIXME: Handling MPRIS directly is not really fail-safe yet,
	// so to prevent any panics from crashing the whole app there is a recover defer.
	// That should be changed ASAP, but it requires a significant rewrite of player listener's logic.
	defer func() {
		_ = recover()
	}()

	playerNames, err := mpris.List(dbusconn.(*dbus.Conn))
	if err != nil {
		return structs.PlayerInfo{}, fmt.Errorf("[player/providers/mpris/GetPlayerInfo] ERROR: Failed to get list of players using existing D-Bus connection")
	}
	if len(playerNames) == 0 {
		return structs.PlayerInfo{}, nil
	}

	var player *mpris.Player = nil

	if len(global.CurrentConfig.Player.IncludedPlayers) != 0 {
		found := false
		for _, v1 := range playerNames {
			for _, v2 := range global.CurrentConfig.Player.IncludedPlayers {
				if strings.Contains(v1, v2) {
					player = mpris.New(dbusconn.(*dbus.Conn), v1)
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
					player = mpris.New(dbusconn.(*dbus.Conn), v1)
					found = true
					break
				}
			}
			if found {
				break
			}
		}
	} else {
		player = mpris.New(dbusconn.(*dbus.Conn), playerNames[0])
	}

	identity, err := player.GetIdentity()
	if err != nil {
		return structs.PlayerInfo{}, fmt.Errorf("[player/providers/mpris/GetPlayerInfo] ERROR: Failed to get player identity using newly found player handler")
	}
	status, err := player.GetPlaybackStatus()
	if err != nil {
		return structs.PlayerInfo{}, fmt.Errorf("[player/providers/mpris/GetPlayerInfo] ERROR: Failed to get playback status using newly found player handler")
	}
	position, err := player.GetPosition()
	if err != nil {
		return structs.PlayerInfo{}, fmt.Errorf("[player/providers/mpris/GetPlayerInfo] ERROR: Failed to get player position using newly found player handler")
	}

	out.PlayerName = identity
	out.IsPlaying = status == mpris.PlaybackPlaying
	out.Position = position

	err = sessions.MediaSessions["mpris"].SetPlayer(player)
	if err != nil {
		return out, fmt.Errorf("[player/providers/mpris/GetPlayerInfo] ERROR: Couldn't save the player handler in session struct")
	}

	return out, nil
}

func getU64(duration interface{}) uint64 {
	switch x := duration.(type) {
	case int:
		return max(0, uint64(x))
	case int8:
		return max(0, uint64(x))
	case int16:
		return max(0, uint64(x))
	case int32:
		return max(0, uint64(x))
	case int64:
		return max(0, uint64(x))
	case uint:
		return uint64(x)
	case uint8:
		return uint64(x)
	case uint16:
		return uint64(x)
	case uint32:
		return uint64(x)
	case uint64:
		return x
	default:
		return 0
	}

}
