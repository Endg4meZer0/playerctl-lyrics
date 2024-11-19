package piped

import "lrcsnc/internal/pkg/global"

func lyricIndexToString(l int) string {
	if l < 0 || l >= len(global.Player.Song.LyricsData.Lyrics) {
		return ""
	} else {
		return global.Player.Song.LyricsData.Lyrics[l]
	}
}
