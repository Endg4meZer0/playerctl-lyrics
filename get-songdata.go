package main

import (
	"os/exec"
	"strconv"
	"strings"
)

func GetCurrentSongData() SongData {
	cmd := exec.Command("playerctl", "metadata")

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
	output, _ := exec.Command("playerctl", "status").CombinedOutput()
	return string(output) == "Playing\n"
}

func GetCurrentSongPosition() float64 {
	output, _ := exec.Command("playerctl", "position").CombinedOutput()
	soutput, _ := strings.CutSuffix(string(output), "\n")
	currentTimestamp, err := strconv.ParseFloat(soutput, 64)
	if err != nil {
		return 0.0
	}

	return currentTimestamp
}
