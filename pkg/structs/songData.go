package structs

type SongData struct {
	Song       string
	Artist     string
	Album      string
	Duration   float64
	LyricsData SongLyricsData
}
