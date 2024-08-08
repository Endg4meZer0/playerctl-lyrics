package main

import (
	"log"
	"os/exec"
	"strconv"
	"strings"
)

func GetCurrentSongData() SongData {
	cmd := exec.Command("playerctl", "metadata")

	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatalln(err)
	}

	soutput := string(output)
	if !strings.Contains(soutput, "xesam:title               ") {
		return SongData{Duration: -1}
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
	duration := 0.0
	if durationFound {
		duration, err = strconv.ParseFloat(strings.TrimLeft(strings.Split(durationExtra, "\n")[0], " "), 64)
		if err != nil {
			duration = 0.0
		}
		duration = duration / 1000000
	}

	return SongData{Song: song, Artist: artist, Album: album, Duration: duration}
}

func GetCurrentSongStatus() bool {
	output, err := exec.Command("playerctl", "status").Output()
	if err != nil {
		log.Fatalln(err)
	}
	return string(output) == "Playing\n"
}

func GetCurrentSongPosition() float64 {
	output, err := exec.Command("playerctl", "position").Output()
	if err != nil {
		log.Fatalln(err)
	}
	soutput, _ := strings.CutSuffix(string(output), "\n")
	currentTimestamp, err := strconv.ParseFloat(soutput, 64)
	if err != nil {
		log.Fatalln(err)
	}

	return currentTimestamp
}
