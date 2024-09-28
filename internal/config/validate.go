package config

import (
	"os/exec"

	"lrcsnc/internal/pkg/structs"
)

type ConfigError string

func (ce ConfigError) Error() string {
	return string(ce)
}

func ValidateConfig(c *structs.Config) (errs []ConfigError, fatal bool) {
	// Check if lrclib is set as the lyric provider
	if c.Global.LyricsProvider != "lrclib" {
		c.Global.LyricsProvider = "lrclib"
		errs = append(errs, `WARNING: For now, 'lrclib' is the only lyrics provider. So the "lyricsProvider" property will always turn to 'lrclib' until there are new lyrics providers introduced.`)
	}

	// Check if playerctl is installed, if not - rollback to direct MPRIS handler
	err := exec.Command("playerctl", "--version").Run()
	if err != nil {
		c.Player.PlayerProvider = "mpris"
		errs = append(errs, `WARNING: The player data provider is set to 'playerctl', but it is not detected. Falling back to 'mpris'...`)
	}

	return
}
