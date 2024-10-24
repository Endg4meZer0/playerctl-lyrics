package piped

import "lrcsnc/internal/pkg/global"

func lyricIndexToString(l int) string {
	if l < 0 || l >= len(global.CurrentSong.LyricsData.Lyrics) {
		return ""
	} else {
		return global.CurrentSong.LyricsData.Lyrics[l]
	}
}
