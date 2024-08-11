package main

type SongData struct {
	Song     string
	Artist   string
	Album    string
	Duration float64
	// 0 = synced, 1 = plain, 2 = instrumental, 3 = song not found, 4 = no active players, 5 = in progress, 6 = unknown
	LyricsType byte
}

type LrcLibJsonOutput struct {
	Instrumental bool   `json:"instrumental"`
	PlainLyrics  string `json:"plainLyrics"`
	SyncedLyrics string `json:"syncedLyrics"`
}
