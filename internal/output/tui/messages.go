package tui

import (
	tea "github.com/charmbracelet/bubbletea"
)

type (
	songInfoChanged     struct{}
	playerInfoChanged   struct{}
	currentLyricChanged int
	overwriteReceived   string
)

func watchSongInfoChanges() tea.Cmd {
	return func() tea.Msg {
		<-SongInfoChanged
		return songInfoChanged{}
	}
}

func watchPlayerInfoChanges() tea.Cmd {
	return func() tea.Msg {
		<-PlayerInfoChanged
		return playerInfoChanged{}
	}
}

func watchCurrentLyricChanges() tea.Cmd {
	return func() tea.Msg {
		return currentLyricChanged(<-CurrentLyricChanged)
	}
}

func watchReceivedOverwrites() tea.Cmd {
	return func() tea.Msg {
		return overwriteReceived(<-OverwriteReceived)
	}
}
