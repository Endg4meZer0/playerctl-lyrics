package providers

import (
	"os/exec"
	"strconv"
	"strings"

	"lrcsnc/internal/pkg/global"
	"lrcsnc/internal/pkg/structs"
)

type PlayerctlPlayerProvider struct{}

func (p PlayerctlPlayerProvider) GetSongInfo() (out structs.SongInfo) {
	cmd := exec.Command("playerctl", "metadata", "-p", global.CurrentPlayer.PlayerName, "-f", "{{title}}\n{{artist}}\n{{album}}\n{{mpris:length}}")
	output, _ := cmd.CombinedOutput()
	soutput := strings.Split(string(output), "\n")
	if len(soutput) != 5 {
		out.LyricsData.LyricsType = 4
		return
	}

	if soutput[0] == "" {
		out.LyricsData.LyricsType = 4
		return
	}

	out.Title = soutput[0]
	out.Artist = soutput[1]
	out.Album = soutput[2]
	out.LyricsData.LyricsType = 5
	out.Duration = 0

	if soutput[3] != "" {
		var duration float64 = 0
		duration2, err := strconv.ParseFloat(soutput[3], 64)
		if err == nil {
			duration = duration2 / 1000000
		}
		out.Duration = duration
	}

	return
}

func (p PlayerctlPlayerProvider) GetPlayerInfo() (out structs.PlayerInfo) {
	var cmd *exec.Cmd
	if len(global.CurrentConfig.Player.IncludedPlayers) != 0 {
		includedPlayers := strings.Join(global.CurrentConfig.Player.IncludedPlayers, ",")
		cmd = exec.Command("playerctl", "metadata", "-p", includedPlayers, "-f", "{{playerName}}\n{{status}}\n{{position}}")
	} else if len(global.CurrentConfig.Player.ExcludedPlayers) != 0 {
		excludedPlayers := strings.Join(global.CurrentConfig.Player.ExcludedPlayers, ",")
		cmd = exec.Command("playerctl", "metadata", "-i", excludedPlayers, "-f", "{{playerName}}\n{{status}}\n{{position}}")
	} else {
		cmd = exec.Command("playerctl", "metadata", "-f", "{{playerName}}\n{{status}}\n{{position}}")
	}
	output, _ := cmd.CombinedOutput()
	soutput := strings.Split(string(output), "\n")
	if soutput[1] != "Stopped" {
		out.PlayerName = soutput[0]
		durationInt, err := strconv.ParseInt(soutput[2], 10, 64)
		if err == nil {
			duration := float64(durationInt) / 1000000
			out.IsPlaying = soutput[1] == "Playing"
			out.Position = duration
		} else {
			out.IsPlaying = soutput[1] == "Playing"
			out.Position = 0
		}
	} else {
		out.IsPlaying = false
		out.PlayerName = ""
		out.Position = 0
	}

	return
}
