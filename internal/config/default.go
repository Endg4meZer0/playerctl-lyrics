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
	Cache: structs.CacheConfig{
		Enabled:       true,
		CacheDir:      "$XDG_CACHE_DIR/lrcsnc",
		CacheLifeSpan: 14,
	},
	Output: structs.OutputConfig{
		TimestampOffset:                         0,
		TerminalOutputInOneLine:                 false,
		ShowSongNotFoundWarning:                 true,
		ShowNotSyncedLyricsWarning:              true,
		ShowGettingLyricsMessage:                true,
		ShowRepeatedLyricsMultiplier:            true,
		RepeatedLyricsMultiplierFormat:          "(x%v)",
		PrintRepeatedLyricsMultiplierToTheRight: true,
		Romanization: structs.RomanizationConfig{
			Japanese: false,
			Chinese:  false,
			Korean:   false,
		},
		Instrumental: structs.InstrumentalConfig{
			Interval: 0.5,
			Symbol:   "♪",
			MaxCount: 3,
		},
	},
}
