package pkg

var CurrentConfig Config

// CONFIG
// LEVEL 0

type Config struct {
	Global    GlobalConfig    `json:"global"`
	Playerctl PlayerctlConfig `json:"playerctl"`
	Cache     CacheConfig     `json:"cache"`
	Output    OutputConfig    `json:"output"`
}

// LEVEL 1

type GlobalConfig struct {
	DisableActiveSync bool `json:"disableActiveSync"`
}

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

func (r *RomanizationConfig) IsEnabled() bool {
	return r.Japanese || r.Chinese || r.Korean
}
