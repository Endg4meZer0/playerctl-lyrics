package tui

type TUIOutputController struct{}

func (TUIOutputController) OnSongInfoChange() {
	SongInfoChanged <- true
}
func (TUIOutputController) OnPlayerInfoChange() {
	PlayerInfoChanged <- true
}
func (TUIOutputController) OnOverwriteReceived(overwrite string) {
	OverwriteReceived <- overwrite
}
func (TUIOutputController) DisplayCurrentLyric(lyricIndex int) {
	CurrentLyricChanged <- lyricIndex
}
