package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"syscall"
)

func main() {

	// Check if `playerctl` is installed
	err := exec.Command("playerctl", "--version").Run()

	if err != nil {
		log.Fatalln("playerctl is not found! Please, install playerctl.")
	}

	fmt.Println("helo")
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	var currentSong SongData
	//var currentLyrics map[float64]string
	//var currentlyInstrumental bool

	go func() {
		song := GetCurrentSongData()
		if song != currentSong {
			//currentLyrics, currentlyInstrumental := GetSyncedLyrics(song)
		}
	}()

}
