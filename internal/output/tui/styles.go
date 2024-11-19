package tui

import (
	"lrcsnc/internal/pkg/global"

	gloss "github.com/charmbracelet/lipgloss"
)

var (
	styleLyric = func() gloss.Style {
		return gloss.NewStyle().AlignHorizontal(gloss.Center).Bold(true)
	}
	styleBefore = func() gloss.Style {
		return styleLyric().Foreground(gloss.Color(global.Config.Output.TUI.Colors.LyricBefore))
	}
	styleCurrent = func() gloss.Style {
		return styleLyric().Foreground(gloss.Color(global.Config.Output.TUI.Colors.LyricCurrent))
	}
	styleAfter = func() gloss.Style {
		return styleLyric().Foreground(gloss.Color(global.Config.Output.TUI.Colors.LyricAfter)).Faint(true)
	}
	styleCursor = func() gloss.Style {
		return styleLyric().Foreground(gloss.Color(global.Config.Output.TUI.Colors.LyricCursor))
	}
	styleBorderCursor = func() gloss.Style {
		return gloss.NewStyle().Border(gloss.ThickBorder(), true, false).BorderForeground(gloss.Color(global.Config.Output.TUI.Colors.BorderCursor))
	}
	styleTimestamp = func() gloss.Style {
		return gloss.NewStyle().Foreground(gloss.Color(global.Config.Output.TUI.Colors.Timestamp))
	}
	styleTimestampCurrent = func() gloss.Style {
		return gloss.NewStyle().Foreground(gloss.Color(global.Config.Output.TUI.Colors.TimestampCurrent))
	}
	styleTimestampCursor = func() gloss.Style {
		return gloss.NewStyle().Foreground(gloss.Color(global.Config.Output.TUI.Colors.TimestampCursor)).Faint(true)
	}
)
