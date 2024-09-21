package cmd

import (
	"log"
	"os"
	"os/exec"
	"os/signal"
	"syscall"

	"lrcsnc/config"
	"lrcsnc/internal/flags"
	"lrcsnc/internal/loop"
	"lrcsnc/internal/output"
)

func Start() {
	// Check if playerctl is installed
	err := exec.Command("playerctl", "--version").Run()
	if err != nil {
		log.Fatalln("playerctl is not found!")
	}

	flags.HandleFlags()

	defer output.CloseOutput()

	exitSigs := make(chan os.Signal, 1)
	signal.Notify(exitSigs, syscall.SIGINT, syscall.SIGTERM)

	usr1Sig := make(chan os.Signal, 1)
	signal.Notify(usr1Sig, syscall.SIGUSR1)

	go func() {
		for {
			<-usr1Sig
			config.UpdateConfig()
		}
	}()

	loop.SyncLoop()
	<-exitSigs
	os.Exit(0)
}
