package output

import (
	"lrcsnc/internal/output/piped"
	"lrcsnc/internal/output/tui"
	"lrcsnc/internal/pkg/global"
)

func SongInfoChangedToOutput() {
	switch global.CurrentConfig.Global.Output {
	case "tui":
		tui.SongInfoChanged <- true
	}
}

func PlayerInfoChangedToOutput() {
	switch global.CurrentConfig.Global.Output {
	case "tui":
		tui.PlayerInfoChanged <- true
	}
}

func CurrentLyricToOutput(l int) {
	switch global.CurrentConfig.Global.Output {
	case "piped":
		piped.PrintLyric(lyricIndexToString(l))
	case "tui":
		tui.CurrentLyricChanged <- l
	}
}

func OverwriteToOutput(s string) {
	switch global.CurrentConfig.Global.Output {
	case "piped":
		piped.PrintOverwrite(s)
	case "tui":
		tui.OverwriteReceived <- s
	}
}

// util

func lyricIndexToString(l int) string {
	if l < 0 || l >= len(global.CurrentSong.LyricsData.Lyrics) {
		return ""
	} else {
		return global.CurrentSong.LyricsData.Lyrics[l]
	}
}
