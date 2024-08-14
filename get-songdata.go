package main

import (
	"os/exec"
	"strconv"
	"strings"
)

func GetCurrentSongData() SongData {
	var cmd *exec.Cmd
	if len(CurrentConfig.Playerctl.IncludedPlayers) != 0 {
		includedPlayers := strings.Join(CurrentConfig.Playerctl.IncludedPlayers, ",")
		cmd = exec.Command("playerctl", "metadata", "-p", includedPlayers)
	} else if len(CurrentConfig.Playerctl.ExcludedPlayers) != 0 {
		excludedPlayers := strings.Join(CurrentConfig.Playerctl.ExcludedPlayers, ",")
		cmd = exec.Command("playerctl", "metadata", "-i", excludedPlayers)
	} else {
		cmd = exec.Command("playerctl", "metadata")
	}

	output, _ := cmd.CombinedOutput()

	soutput := string(output)
	if !strings.Contains(soutput, "xesam:title") || !strings.Contains(soutput, "mpris:length") {
		return SongData{Song: "", Artist: "", Album: ""}
	}
	_, songTitleExtra, _ := strings.Cut(soutput, "xesam:title")
	song := strings.TrimLeft(strings.Split(songTitleExtra, "\n")[0], " ")
	_, artistNameExtra, _ := strings.Cut(soutput, "xesam:artist")
	artist := strings.TrimLeft(strings.Split(artistNameExtra, "\n")[0], " ")
	_, albumNameExtra, albumFound := strings.Cut(soutput, "xesam:album")
	album := ""
	if albumFound {
		album = strings.TrimLeft(strings.Split(albumNameExtra, "\n")[0], " ")
	}
	_, durationExtra, durationFound := strings.Cut(soutput, "mpris:length")
	var duration float64
	if durationFound {
		duration2, err := strconv.ParseFloat(strings.TrimLeft(strings.Split(durationExtra, "\n")[0], " "), 64)
		if err != nil {
			duration2 = 0.0
		} else {
			duration = duration2 / 1000000
		}
	}

	return SongData{Song: song, Artist: artist, Album: album, Duration: duration, LyricsType: 5}
}

func GetCurrentSongStatus() bool {
	var cmd *exec.Cmd
	if len(CurrentConfig.Playerctl.IncludedPlayers) != 0 {
		includedPlayers := strings.Join(CurrentConfig.Playerctl.IncludedPlayers, ",")
		cmd = exec.Command("playerctl", "status", "-p", includedPlayers)
	} else if len(CurrentConfig.Playerctl.ExcludedPlayers) != 0 {
		excludedPlayers := strings.Join(CurrentConfig.Playerctl.ExcludedPlayers, ",")
		cmd = exec.Command("playerctl", "status", "-i", excludedPlayers)
	} else {
		cmd = exec.Command("playerctl", "status")
	}
	output, _ := cmd.CombinedOutput()
	return string(output) == "Playing\n"
}

func GetCurrentSongPosition() float64 {
	var cmd *exec.Cmd
	if len(CurrentConfig.Playerctl.IncludedPlayers) != 0 {
		includedPlayers := strings.Join(CurrentConfig.Playerctl.IncludedPlayers, ",")
		cmd = exec.Command("playerctl", "position", "-p", includedPlayers)
	} else if len(CurrentConfig.Playerctl.ExcludedPlayers) != 0 {
		excludedPlayers := strings.Join(CurrentConfig.Playerctl.ExcludedPlayers, ",")
		cmd = exec.Command("playerctl", "position", "-i", excludedPlayers)
	} else {
		cmd = exec.Command("playerctl", "position")
	}
	output, _ := cmd.CombinedOutput()
	soutput, _ := strings.CutSuffix(string(output), "\n")
	currentTimestamp, err := strconv.ParseFloat(soutput, 64)
	if err != nil {
		return 0.0
	}

	return currentTimestamp
}
