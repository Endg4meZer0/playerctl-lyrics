package config

import "lrcsnc/internal/pkg/structs"

var defaultConfig = structs.Config{
	Global: structs.GlobalConfig{
		Output:           "tui",
		LyricsProvider:   "lrclib",
		EnableActiveSync: false,
	},
	Player: structs.PlayerConfig{
		PlayerProvider:    "playerctl",
		IncludedPlayers:   []string{},
		ExcludedPlayers:   []string{},
		SongCheckInterval: 0.5,
	},
	Lyrics: structs.LyricsConfig{
		TimestampOffset: 0.0,
		Romanization: structs.RomanizationConfig{
			Japanese: false,
			Chinese:  false,
			Korean:   false,
		},
	},
	Cache: structs.CacheConfig{
		Enabled:       true,
		CacheDir:      "$XDG_CACHE_DIR/lrcsnc",
		CacheLifeSpan: 14,
	},
	Output: structs.OutputConfig{
		Piped: structs.PipedOutputConfig{
			ShowSongNotFoundWarning:                 true,
			ShowNotSyncedLyricsWarning:              true,
			ShowGettingLyricsMessage:                true,
			ShowRepeatedLyricsMultiplier:            true,
			RepeatedLyricsMultiplierFormat:          "(x%v)",
			PrintRepeatedLyricsMultiplierToTheRight: true,
			Instrumental: structs.InstrumentalConfig{
				Interval: 0.5,
				Symbol:   "♪",
				MaxCount: 3,
			},
		},
		TUI: structs.TUIOutputConfig{
			Colors: structs.TUIColorsConfig{
				LyricBefore:      "15",
				LyricCurrent:     "11",
				LyricAfter:       "15", // faint
				LyricCursor:      "3",
				BorderCursor:     "3",
				Timestamp:        "8",
				TimestampCurrent: "3",
				TimestampCursor:  "3", // faint
				ProgressBarColor: "10",
			},
			ShowTimestamps:  false,
			ShowProgressBar: true,
		},
	},
}
