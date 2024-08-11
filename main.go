package main

import (
	"log"
	"os"
	"os/exec"
	"os/signal"
	"syscall"
)

func main() {

	// TODO: An actual config implementation
	// TODO: Flags implementation
	// there's definitely always more!

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
