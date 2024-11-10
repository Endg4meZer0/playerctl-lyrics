// Playerctl support is deprecated as of 0.3.0 and right now exists purely for compatibility purposes

package controllers

import (
	"fmt"
	"os/exec"

	"lrcsnc/internal/pkg/global"
)

type PlayerctlPlayerController struct{}

func (p PlayerctlPlayerController) SeekTo(pos float64) error {
	if global.CurrentPlayer.PlayerName == "" {
		return fmt.Errorf("[player/controllers/playerctl/SeekTo] ERROR: Tried to seek to lyric, but no player name is present in global config")
	}

	cmd := exec.Command("playerctl", "position", "-p", global.CurrentPlayer.PlayerName, fmt.Sprintf("%.2f", pos))
	_, err := cmd.CombinedOutput()

	return err
}
