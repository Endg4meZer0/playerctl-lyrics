package config

import (
	"os/exec"

	"lrcsnc/pkg/structs"
)

type ConfigError string

func (ce ConfigError) Error() string {
	return string(ce)
}

func ValidateConfig(c *structs.Config) (errs []ConfigError, fatal bool) {
	// Check if playerctl is installed, rollback to
	err := exec.Command("playerctl", "--version").Run()
	if err != nil {
		c.Player.PlayerProvider = "mpris"
		errs = append(errs, `WARNING: The player data provider is set to 'playerctl', but it is not detected. Falling back to 'mpris'...`)
	}

	return
}
