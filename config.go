package main

import (
	"encoding/json"
	"os"
)

// LEVEL 0

type Config struct {
	Playerctl    Playerctl    `json:"playerctl"`
	Cache        Cache        `json:"cache"`
	Output       Output       `json:"output"`
	Instrumental Instrumental `json:"instrumental"`
}

// LEVEL 1

type Playerctl struct {
	IncludedPlayers            []string `json:"includedPlayers"`
	ExcludedPlayers            []string `json:"excludedPlayers"`
	PlayerctlSongCheckInterval float64  `json:"playerctlMetadataCheckInterval"`
}

type Cache struct {
	DoCacheLyrics bool   `json:"doCacheLyrics"`
	CacheDir      string `json:"cacheDir"`
	CacheLifeSpan uint   `json:"cacheLifeSpan"`
}

type Output struct {
	ShowSongNotFoundWarning         bool         `json:"showSongNotFoundWarning"`
	ShowNotSyncedLyricsWarning      bool         `json:"showNotSyncedLyricsWarning"`
	ShowGettingLyricsMessage        bool         `json:"showGettingLyricsMessage"`
	ShowRepeatedLyricsMultiplicator bool         `json:"showRepeatedLyricsMultiplicator"`
	Romanization                    Romanization `json:"romanization"`
}

type Instrumental struct {
	Interval float64 `json:"interval"`
	Symbol   string  `json:"symbol"`
	MaxCount uint    `json:"maxCount"`
}

// LEVEL 2

type Romanization struct {
	Japanese bool `json:"japanese"`
	Chinese  bool `json:"chinese"`
	Korean   bool `json:"korean"`
}

var CurrentConfig Config

func ReadConfig(path string) error {
	jsonConfig, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(jsonConfig, &CurrentConfig); err != nil {
		return err
		// "The config at %v does not exist or is not readable! Falling back to the default config!"
		// "The config at %v is not formatted as JSON! Falling back to the default config!"
		// "The config at %v is not valid! Error at %v. Falling back to the default config!"
	}

	return nil
}

func DefaultConfig() Config {
	return Config{
		Playerctl: Playerctl{
			IncludedPlayers:            []string{},
			ExcludedPlayers:            []string{},
			PlayerctlSongCheckInterval: 1.0,
		},
		Cache: Cache{
			DoCacheLyrics: true,
			CacheDir:      "$XDG_CACHE_DIR/playerctl-lyrics",
			CacheLifeSpan: 0,
		},
		Output: Output{
			ShowSongNotFoundWarning:         true,
			ShowNotSyncedLyricsWarning:      true,
			ShowGettingLyricsMessage:        true,
			ShowRepeatedLyricsMultiplicator: true,
			Romanization: Romanization{
				Japanese: false,
				Chinese:  false,
				Korean:   false,
			},
		},
		Instrumental: Instrumental{
			Interval: 0.5,
			Symbol:   "♪",
			MaxCount: 3,
		},
	}
}
