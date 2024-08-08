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

	// TODO: The main functionality, the PURPOSE!
	// TODO: An actual config implementation
	// TODO: Flags implementation
	// there's definitely always more!
	// FIXME: parsing urls is done kinda wrong

	// Check if `playerctl` is installed
	err := exec.Command("playerctl", "--version").Run()

	if err != nil {
		log.Fatalln("playerctl is not found!")
	}

	fmt.Println("helo")
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	PlayLyrics()
	<-sigs
	os.Exit(0)
}
