package controllers

import (
	"fmt"
	"os/exec"

	"lrcsnc/internal/pkg/global"
)

type PlayerctlPlayerController struct{}

func (p PlayerctlPlayerController) SeekTo(pos float64) bool {
	if global.CurrentPlayer.PlayerName == "" {
		return false
	}

	cmd := exec.Command("playerctl", "position", "-p", global.CurrentPlayer.PlayerName, fmt.Sprintf("%.2f", pos))
	_, err := cmd.CombinedOutput()

	return err == nil
}
