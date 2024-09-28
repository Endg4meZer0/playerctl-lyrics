package structs

type SongLyricsData struct {
	Lyrics          []string
	LyricTimestamps []float64
	// 0 = synced, 1 = plain, 2 = instrumental, 3 = song not found, 4 = no active players, 5 = retrieval in progress, 6 = unknown
	LyricsType byte
}
