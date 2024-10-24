package structs

// LEVEL 0

type Config struct {
	Global GlobalConfig `json:"global"`
	Player PlayerConfig `json:"player"`
	Lyrics LyricsConfig `json:"lyrics"`
	Cache  CacheConfig  `json:"cache"`
	Output OutputConfig `json:"output"`
}

// LEVEL 1

type GlobalConfig struct {
	Output           string `json:"output"`
	LyricsProvider   string `json:"lyricsProvider"`
	EnableActiveSync bool   `json:"enableActiveSync"`
}

type PlayerConfig struct {
	PlayerProvider    string   `json:"playerProvider"`
	IncludedPlayers   []string `json:"includedPlayers"`
	ExcludedPlayers   []string `json:"excludedPlayers"`
	SongCheckInterval float64  `json:"songCheckInterval"`
}

type LyricsConfig struct {
	TimestampOffset float64            `json:"timestampOffset"`
	Romanization    RomanizationConfig `json:"romanization"`
}

type CacheConfig struct {
	Enabled       bool   `json:"enabled"`
	CacheDir      string `json:"cacheDir"`
	CacheLifeSpan uint   `json:"cacheLifeSpan"`
}

type OutputConfig struct {
	Piped PipedOutputConfig `json:"piped"`
	TUI   TUIOutputConfig   `json:"tui"`
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

type PipedOutputConfig struct {
	ShowSongNotFoundWarning                 bool               `json:"showSongNotFoundWarning"`
	ShowNotSyncedLyricsWarning              bool               `json:"showNotSyncedLyricsWarning"`
	ShowGettingLyricsMessage                bool               `json:"showGettingLyricsMessage"`
	ShowRepeatedLyricsMultiplier            bool               `json:"showRepeatedLyricsMultiplier"`
	RepeatedLyricsMultiplierFormat          string             `json:"repeatedLyricsMultiplierFormat"`
	PrintRepeatedLyricsMultiplierToTheRight bool               `json:"printRepeatedLyricsMultiplierToTheRight"`
	Instrumental                            InstrumentalConfig `json:"instrumental"`
}

type TUIOutputConfig struct {
	Colors          TUIColorsConfig `json:"colors"`
	ShowTimestamps  bool            `json:"showTimestamps"`
	ShowProgressBar bool            `json:"showProgressBar"`
}

// LEVEL 3

type TUIColorsConfig struct {
	LyricBefore      string `json:"lyricBefore"`
	LyricCurrent     string `json:"lyricCurrent"`
	LyricAfter       string `json:"lyricAfter"`
	LyricCursor      string `json:"lyricCursor"`
	BorderCursor     string `json:"borderCursor"`
	Timestamp        string `json:"timestamp"`
	TimestampCurrent string `json:"timestampCurrent"`
	TimestampCursor  string `json:"timestampCursor"`
	ProgressBarColor string `json:"progressBarColor"`
}

// METHODS

func (r *RomanizationConfig) IsEnabled() bool {
	return r.Japanese || r.Chinese || r.Korean
}
