package tui

var (
	SongInfoChanged     chan bool   = make(chan bool)
	PlayerInfoChanged   chan bool   = make(chan bool)
	CurrentLyricChanged chan int    = make(chan int)
	OverwriteReceived   chan string = make(chan string)
)
