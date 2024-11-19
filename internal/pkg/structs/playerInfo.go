package structs

type Player struct {
	Identity  string
	IsPlaying bool
	LoopMode  byte
	Position  float64
	Song      Song
}

type Song struct {
	Title      string
	Artist     string
	Album      string
	Duration   float64
	LyricsData LyricsData
}
