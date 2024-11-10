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
	"lrcsnc/internal/player"

	tea "github.com/charmbracelet/bubbletea"
)

func Start() {
	// Handle the -(-) flags
	flags.HandleFlags()

	// Start the USR1 signal listener for config updates
	// TODO: replace with live file watcher
	usr1Sig := make(chan os.Signal, 1)
	signal.Notify(usr1Sig, syscall.SIGUSR1)

	go func() {
		for {
			<-usr1Sig
			config.UpdateConfig()
		}
	}()

	// Initialize the player listener session
	player.InitSession()
	defer player.CloseSession()

	// Start the main loop
	loop.SyncLoop()

	// Initialize the output
	switch global.CurrentConfig.Global.Output {
	case "piped":
		go piped.Init()
		defer piped.CloseOutput()

		exitSigs := make(chan os.Signal, 1)
		signal.Notify(exitSigs, syscall.SIGINT, syscall.SIGTERM)

		<-exitSigs
		os.Exit(0)
	// TUI also falls under the default handler
	default:
		p := tea.NewProgram(tui.InitialModel(), tea.WithAltScreen())
		if _, err := p.Run(); err != nil {
			fmt.Printf("Error: %v", err)
			os.Exit(1)
		}
	}
}
