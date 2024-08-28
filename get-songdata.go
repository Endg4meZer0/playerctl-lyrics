package main

import (
	"os/exec"
	"strconv"
	"strings"
)

func GetSongData() SongData {
	var cmd *exec.Cmd
	if len(CurrentConfig.Playerctl.IncludedPlayers) != 0 {
		includedPlayers := strings.Join(CurrentConfig.Playerctl.IncludedPlayers, ",")
		cmd = exec.Command("playerctl", "metadata", "-p", includedPlayers, "-f", "{{title}}\n{{artist}}\n{{album}}\n{{mpris:length}}")
	} else if len(CurrentConfig.Playerctl.ExcludedPlayers) != 0 {
		excludedPlayers := strings.Join(CurrentConfig.Playerctl.ExcludedPlayers, ",")
		cmd = exec.Command("playerctl", "metadata", "-i", excludedPlayers, "-f", "{{title}}\n{{artist}}\n{{album}}\n{{mpris:length}}")
	} else {
		cmd = exec.Command("playerctl", "metadata", "-f", "{{title}}\n{{artist}}\n{{album}}\n{{mpris:length}}")
	}

	output, _ := cmd.CombinedOutput()

	soutput := strings.Split(string(output), "\n")

	if len(soutput) != 4 {
		return SongData{Song: "", Artist: "", Album: "", LyricsType: 4}
	}

	song := soutput[0]
	artist := soutput[1]
	album := soutput[2]
	durationStr := soutput[3]

	if song == "" || durationStr == "" {
		return SongData{Song: "", Artist: "", Album: "", LyricsType: 4}
	}

	var duration float64
	duration2, err := strconv.ParseFloat(durationStr, 64)
	if err != nil {
		duration2 = 0.0
	} else {
		duration = duration2 / 1000000
	}

	return SongData{Song: song, Artist: artist, Album: album, Duration: duration, LyricsType: 5}
}

// Get the player's status and position
func GetPlayerData() (bool, float64) {
	var cmd *exec.Cmd
	if len(CurrentConfig.Playerctl.IncludedPlayers) != 0 {
		includedPlayers := strings.Join(CurrentConfig.Playerctl.IncludedPlayers, ",")
		cmd = exec.Command("playerctl", "metadata", "-p", includedPlayers, "-f", "{{status}}\n{{position}}")
	} else if len(CurrentConfig.Playerctl.ExcludedPlayers) != 0 {
		excludedPlayers := strings.Join(CurrentConfig.Playerctl.ExcludedPlayers, ",")
		cmd = exec.Command("playerctl", "metadata", "-i", excludedPlayers, "-f", "{{status}}\n{{position}}")
	} else {
		cmd = exec.Command("playerctl", "metadata", "-f", "{{status}}\n{{position}}")
	}
	output, _ := cmd.CombinedOutput()
	soutput := strings.Split(string(output), "\n")
	if soutput[0] != "Stopped" {
		durationInt, err := strconv.ParseInt(soutput[1], 10, 64)
		duration := float64(durationInt) / 1000000
		if err != nil {
			return soutput[0] == "Playing", 0
		} else {
			return soutput[0] == "Playing", duration
		}
	} else {
		return false, 0
	}
}
