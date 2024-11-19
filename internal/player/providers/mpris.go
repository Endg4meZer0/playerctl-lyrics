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

func (m MprisPlayerProvider) GetInfo() (structs.Player, error) {
	var out structs.Player
	var err error
	var dbusconn = sessions.MediaSessions["mpris"].GetSession()
	if dbusconn == nil {
		out.Song.LyricsData.LyricsType = 4
		return out, fmt.Errorf("[player/providers/mpris] ERROR: D-Bus connection does not exist at the moment")
	}

	var player *mpris.Player = sessions.MediaSessions["mpris"].GetPlayer().(*mpris.Player)
	if player != nil {
		out.Identity, err = player.GetIdentity()
		if err == nil {
			playing, err := player.GetPlaybackStatus()
			if err != nil {
				return structs.Player{}, fmt.Errorf("[player/providers/mpris] ERROR: Failed to get playback status using existing player handler")
			}
			out.IsPlaying = playing == mpris.PlaybackPlaying
			pos, err := player.GetPosition()
			if err != nil {
				return structs.Player{}, fmt.Errorf("[player/providers/mpris] ERROR: Failed to get player position using existing player handler")
			}
			out.Position = pos
		}
	} else {
		playerNames, err := mpris.List(dbusconn.(*dbus.Conn))
		if err != nil {
			return structs.Player{}, fmt.Errorf("[player/providers/mpris/GetPlayerInfo] ERROR: Failed to get list of players using existing D-Bus connection")
		}
		if len(playerNames) == 0 {
			return structs.Player{}, nil
		}

		if len(global.Config.Player.IncludedPlayers) != 0 {
			found := false
			for _, v1 := range playerNames {
				for _, v2 := range global.Config.Player.IncludedPlayers {
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
		} else if len(global.Config.Player.ExcludedPlayers) != 0 {
			found := false
			for _, v1 := range playerNames {
				for _, v2 := range global.Config.Player.ExcludedPlayers {
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
	}

	status, err := player.GetPlaybackStatus()
	if err != nil {
		return structs.Player{}, fmt.Errorf("[player/providers/mpris/GetInfo] ERROR: Failed to get playback status using newly found player handler")
	}
	position, err := player.GetPosition()
	if err != nil {
		return structs.Player{}, fmt.Errorf("[player/providers/mpris/GetInfo] ERROR: Failed to get player position using newly found player handler")
	}
	loopMode, err := player.GetLoopStatus()
	if err != nil {
		return structs.Player{}, fmt.Errorf("[player/providers/mpris/GetInfo] ERROR: Failed to get loop mode using newly found player handler")
	}

	out.IsPlaying = status == mpris.PlaybackPlaying
	out.Position = position
	switch loopMode {
	case mpris.LoopTrack:
		out.LoopMode = 1
	case mpris.LoopPlaylist:
		out.LoopMode = 2
	default:
		out.LoopMode = 0
	}

	err = sessions.MediaSessions["mpris"].SetPlayer(player)
	if err != nil {
		return out, fmt.Errorf("[player/providers/mpris/GetInfo] ERROR: Couldn't save the player handler in session struct")
	}

	metadata, err := player.GetMetadata()
	if err != nil {
		out.Song.LyricsData.LyricsType = 4
		return out, fmt.Errorf("[player/providers/mpris/GetInfo] ERROR: Failed to get metadata from existing player handler")
	}

	song, ok1 := metadata["xesam:title"]
	artist, ok2 := metadata["xesam:artist"]
	album, ok3 := metadata["xesam:album"]
	if !(ok1 && ok2 && ok3) {
		out.Song.LyricsData.LyricsType = 4
		return out, nil
	}
	duration, ok := metadata["mpris:length"]

	out.Song.Title = song.Value().(string)
	out.Song.Artist = strings.Join(artist.Value().([]string), ", ")
	out.Song.Album = album.Value().(string)
	out.Song.LyricsData.LyricsType = 5

	if ok {
		out.Song.Duration = float64(getU64(duration.Value())) / 1000000
	}

	return out, nil
}

func (m MprisPlayerProvider) Subscribe() <-chan structs.Player {
	out := make(chan structs.Player)
	changesCh := make(chan *dbus.Signal)
	player := sessions.MediaSessions["mpris"].GetPlayer().(*mpris.Player)
	if player == nil {
		return nil
	}
	err := player.OnSignal(changesCh)
	if err != nil {
		// TODO: logger :)
		return nil
	}
	go func() {
		for {
			<-changesCh
			player, err := m.GetInfo()
			if err != nil {
				// TODO: logger :)
				continue
			}
			out <- player
		}
	}()
	return out
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
