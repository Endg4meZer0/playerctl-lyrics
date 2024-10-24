package piped

type PipedOutputController struct{}

func (PipedOutputController) OnSongInfoChange()   {}
func (PipedOutputController) OnPlayerInfoChange() {}
func (PipedOutputController) OnOverwriteReceived(overwrite string) {
	PrintOverwrite(overwrite)
}
func (PipedOutputController) DisplayCurrentLyric(lyricIndex int) {
	PrintLyric(lyricIndexToString(lyricIndex))
}
