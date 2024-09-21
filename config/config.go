package config

import (
	"encoding/json"
	"os"

	"lrcsnc/internal/output"
	"lrcsnc/pkg"
)

var currentConfigPath string

func ReadConfig(path string) error {
	configFile, err := os.ReadFile(os.ExpandEnv(path))
	if err != nil {
		return err
	}

	if err := json.Unmarshal(configFile, &pkg.CurrentConfig); err != nil {
		return err
	}

	currentConfigPath = path

	return nil
}

func ReadConfigFromDefaultPath() error {
	pkg.CurrentConfig = defaultConfig

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
		if err := json.Unmarshal(configFile, &pkg.CurrentConfig); err != nil {
			return err
		}
	}

	currentConfigPath = defaultDirectory + "/config.json"

	return nil
}

func UpdateConfig() {
	configFile, err := os.ReadFile(os.ExpandEnv(currentConfigPath))
	if err != nil {
		output.PrintOverwrite("Errors while reading config! Falling back...")
		return
	}

	if err := json.Unmarshal(configFile, &pkg.CurrentConfig); err != nil {
		output.PrintOverwrite("Errors while parsing config! Falling back...")
		return
	}
}

var defaultConfig = pkg.Config{
	Global: pkg.GlobalConfig{
		DisableActiveSync: false,
	},
	Playerctl: pkg.PlayerctlConfig{
		IncludedPlayers:            []string{},
		ExcludedPlayers:            []string{},
		PlayerctlSongCheckInterval: 0.5,
	},
	Cache: pkg.CacheConfig{
		Enabled:       true,
		CacheDir:      "$XDG_CACHE_DIR/playerctl-lyrics",
		CacheLifeSpan: 14,
	},
	Output: pkg.OutputConfig{
		TimestampOffset:                         0,
		TerminalOutputInOneLine:                 false,
		ShowSongNotFoundWarning:                 true,
		ShowNotSyncedLyricsWarning:              true,
		ShowGettingLyricsMessage:                true,
		ShowRepeatedLyricsMultiplier:            true,
		RepeatedLyricsMultiplierFormat:          "(x%v)",
		PrintRepeatedLyricsMultiplierToTheRight: true,
		Romanization: pkg.RomanizationConfig{
			Japanese: false,
			Chinese:  false,
			Korean:   false,
		},
		Instrumental: pkg.InstrumentalConfig{
			Interval: 0.5,
			Symbol:   "♪",
			MaxCount: 3,
		},
	},
}
