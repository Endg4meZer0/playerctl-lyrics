package tui

import "github.com/charmbracelet/bubbles/key"

type keyMap struct {
	Up              key.Binding
	Down            key.Binding
	Follow          key.Binding
	ShowTimestamps  key.Binding
	ShowProgressBar key.Binding
	SeekToLyric     key.Binding
	Help            key.Binding
	Quit            key.Binding
}

func (k keyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Help, k.Quit}
}

func (k keyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Up, k.Down, k.Follow, k.ShowTimestamps},
		{k.ShowProgressBar, k.SeekToLyric, k.Help, k.Quit},
	}
}

var keys = keyMap{
	Up: key.NewBinding(
		key.WithKeys("up", "k"),
		key.WithHelp("↑/k", "move cursor up"),
	),
	Down: key.NewBinding(
		key.WithKeys("down", "j"),
		key.WithHelp("↓/j", "move cursor down"),
	),
	Follow: key.NewBinding(
		key.WithKeys("f"),
		key.WithHelp("f", "toggle following the active lyric"),
	),
	ShowProgressBar: key.NewBinding(
		key.WithKeys("p"),
		key.WithHelp("p", "toggle progress bar"),
	),
	ShowTimestamps: key.NewBinding(
		key.WithKeys("t"),
		key.WithHelp("t", "toggle timestamps"),
	),
	SeekToLyric: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "seek to selected lyric"),
	),
	Help: key.NewBinding(
		key.WithKeys("?", "h"),
		key.WithHelp("?/h", "toggle help"),
	),
	Quit: key.NewBinding(
		key.WithKeys("q", "ctrl+c", "esc"),
		key.WithHelp("q", "quit"),
	),
}
