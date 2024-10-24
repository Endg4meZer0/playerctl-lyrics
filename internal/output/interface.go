package output

import (
	"lrcsnc/internal/output/piped"
	"lrcsnc/internal/output/tui"
)

type OutputController interface {
	OnSongInfoChange()
	OnPlayerInfoChange()
	OnOverwriteReceived(overwrite string)

	DisplayCurrentLyric(lyricIndex int)
}

var OutputControllers = map[string]OutputController{
	"piped": piped.PipedOutputController{},
	"tui":   tui.TUIOutputController{},
}
