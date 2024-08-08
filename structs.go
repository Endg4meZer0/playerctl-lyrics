package main

type SongData struct {
	Song     string
	Artist   string
	Album    string
	Duration float64
}

type LrcLibJsonOutput struct {
	Duration     float64 `json:"duration"`
	Instrumental bool    `json:"instrumental"`
	PlainLyrics  string  `json:"plainLyrics"`
	SyncedLyrics string  `json:"syncedLyrics"`
}
