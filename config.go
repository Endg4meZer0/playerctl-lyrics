package main

import (
	"encoding/json"
	"os"
)

var currentConfigPath string

// LEVEL 0

type Config struct {
	Playerctl PlayerctlConfig `json:"playerctl"`
	Cache     CacheConfig     `json:"cache"`
	Output    OutputConfig    `json:"output"`
}

// LEVEL 1

type PlayerctlConfig struct {
	IncludedPlayers            []string `json:"includedPlayers"`
	ExcludedPlayers            []string `json:"excludedPlayers"`
	PlayerctlSongCheckInterval float64  `json:"playerctlMetadataCheckInterval"`
}

type CacheConfig struct {
	Enabled       bool   `json:"enabled"`
	CacheDir      string `json:"cacheDir"`
	CacheLifeSpan uint   `json:"cacheLifeSpan"`
}

type OutputConfig struct {
	TimestampOffset                         int64              `json:"timestampOffset"`
	TerminalOutputInOneLine                 bool               `json:"terminalOutputInOneLine"`
	ShowSongNotFoundWarning                 bool               `json:"showSongNotFoundWarning"`
	ShowNotSyncedLyricsWarning              bool               `json:"showNotSyncedLyricsWarning"`
	ShowGettingLyricsMessage                bool               `json:"showGettingLyricsMessage"`
	ShowRepeatedLyricsMultiplier            bool               `json:"showRepeatedLyricsMultiplier"`
	RepeatedLyricsMultiplierFormat          string             `json:"repeatedLyricsMultiplierFormat"`
	PrintRepeatedLyricsMultiplierToTheRight bool               `json:"printRepeatedLyricsMultiplierToTheRight"`
	Romanization                            RomanizationConfig `json:"romanization"`
	Instrumental                            InstrumentalConfig `json:"instrumental"`
}

// LEVEL 2

type RomanizationConfig struct {
	Japanese bool `json:"japanese"`
	Chinese  bool `json:"chinese"`
	Korean   bool `json:"korean"`
}

type InstrumentalConfig struct {
	Interval float64 `json:"interval"`
	Symbol   string  `json:"symbol"`
	MaxCount uint    `json:"maxCount"`
}

func ReadConfig(path string) error {
	configFile, err := os.ReadFile(os.ExpandEnv(path))
	if err != nil {
		return err
	}

	if err := json.Unmarshal(configFile, &CurrentConfig); err != nil {
		return err
	}

	currentConfigPath = path

	return nil
}

func ReadConfigFromDefaultPath() error {
	CurrentConfig = defaultConfig

	defaultDirectory, err := os.UserConfigDir()
	if err != nil {
		return err
	}
	defaultDirectory += "/playerctl-lyrics"

	if _, err := os.ReadDir(defaultDirectory); err != nil {
		os.Mkdir(defaultDirectory, 0777)
		os.Chmod(defaultDirectory, 0777)
	}

	if _, err := os.Lstat(defaultDirectory + "/config.json"); err != nil {
		defaultConfigJSON, err := json.MarshalIndent(defaultConfig, "", "    ")
		if err != nil {
			return err
		}
		err = os.WriteFile(defaultDirectory+"/config.json", defaultConfigJSON, 0777)
		if err != nil {
			return err
		}
	} else {
		configFile, err := os.ReadFile(defaultDirectory + "/config.json")
		if err != nil {
			return err
		}
		if err := json.Unmarshal(configFile, &CurrentConfig); err != nil {
			return err
		}
	}

	currentConfigPath = defaultDirectory + "/config.json"

	return nil
}

func UpdateConfig() {
	configFile, err := os.ReadFile(os.ExpandEnv(currentConfigPath))
	if err != nil {
		PrintOverwrite("Errors while reading config! Falling back...")
		return
	}

	if err := json.Unmarshal(configFile, &CurrentConfig); err != nil {
		PrintOverwrite("Errors while parsing config! Falling back...")
		return
	}
}

func (r *RomanizationConfig) IsEnabled() bool {
	return r.Japanese || r.Chinese || r.Korean
}

var defaultConfig = Config{
	Playerctl: PlayerctlConfig{
		IncludedPlayers:            []string{},
		ExcludedPlayers:            []string{},
		PlayerctlSongCheckInterval: 0.5,
	},
	Cache: CacheConfig{
		Enabled:       true,
		CacheDir:      "$XDG_CACHE_DIR/playerctl-lyrics",
		CacheLifeSpan: 14,
	},
	Output: OutputConfig{
		TimestampOffset:                         0,
		TerminalOutputInOneLine:                 false,
		ShowSongNotFoundWarning:                 true,
		ShowNotSyncedLyricsWarning:              true,
		ShowGettingLyricsMessage:                true,
		ShowRepeatedLyricsMultiplier:            true,
		RepeatedLyricsMultiplierFormat:          "(x%v)",
		PrintRepeatedLyricsMultiplierToTheRight: true,
		Romanization: RomanizationConfig{
			Japanese: false,
			Chinese:  false,
			Korean:   false,
		},
		Instrumental: InstrumentalConfig{
			Interval: 0.5,
			Symbol:   "♪",
			MaxCount: 3,
		},
	},
}
