package main

import (
	"log"
	"os"
	"os/exec"
	"os/signal"
	"syscall"
)

func main() {

	// TODO: The main functionality, the PURPOSE!
	// TODO: An actual config implementation
	// TODO: Flags implementation
	// TODO: At least some sort of caching system (using .lrc files from LrcLib should suffice)
	// there's definitely always more!
	// FIXME: parsing urls is done kinda wrong, so e.g. Panchiko - D>E>A>T>H>M>E>T>A>L fails horribly

	// Check if `playerctl` is installed
	err := exec.Command("playerctl", "--version").Run()

	if err != nil {
		log.Fatalln("playerctl is not found!")
	}

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	SyncLoop()
	<-sigs
	os.Exit(0)
}
