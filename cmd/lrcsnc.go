package cmd

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"lrcsnc/internal/config"
	"lrcsnc/internal/flags"
	"lrcsnc/internal/loop"
	"lrcsnc/internal/output/piped"
	"lrcsnc/internal/output/tui"
	"lrcsnc/internal/pkg/global"

	tea "github.com/charmbracelet/bubbletea"
)

func Start() {
	flags.HandleFlags()

	loop.SyncLoop()

	usr1Sig := make(chan os.Signal, 1)
	signal.Notify(usr1Sig, syscall.SIGUSR1)

	go func() {
		for {
			<-usr1Sig
			config.UpdateConfig()
		}
	}()

	if global.CurrentConfig.Global.Output == "piped" {
		go piped.Init()
		defer piped.CloseOutput()

		exitSigs := make(chan os.Signal, 1)
		signal.Notify(exitSigs, syscall.SIGINT, syscall.SIGTERM)

		<-exitSigs
		os.Exit(0)
	} else if global.CurrentConfig.Global.Output == "tui" {
		p := tea.NewProgram(tui.InitialModel(), tea.WithAltScreen())
		if _, err := p.Run(); err != nil {
			fmt.Printf("Error: %v", err)
			os.Exit(1)
		}
	}
}
