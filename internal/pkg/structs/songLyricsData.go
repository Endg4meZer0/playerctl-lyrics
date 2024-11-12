package structs

type LyricsData struct {
	Lyrics          []string
	LyricTimestamps []float64
	LyricsType      LyricsState
}

/*
LyricsState represents the state of the lyrics

Possible states:

- LyricsStateSynced = 0: Lyrics are synced with the song

- LyricsStatePlain = 1: Lyrics are plain text

- LyricsStateInstrumental = 2: There are no lyrics since the track is instrumental

- LyricsStateNotFound = 3: Lyrics are not found

- LyricsStateInProgress = 4: Lyrics are being fetched

- LyricsStateUnknown = 5: Lyrics state is unknown
*/
type LyricsState byte

const (
	LyricsStateSynced       LyricsState = 0
	LyricsStatePlain        LyricsState = 1
	LyricsStateInstrumental LyricsState = 2
	LyricsStateNotFound     LyricsState = 3
	LyricsStateInProgress   LyricsState = 4
	LyricsStateUnknown      LyricsState = 5
)
