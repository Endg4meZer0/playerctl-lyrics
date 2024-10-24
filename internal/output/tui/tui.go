package tui

import (
	"fmt"
	"math"
	"strings"
	"time"

	"lrcsnc/internal/pkg/global"
	"lrcsnc/internal/player"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/progress"
	vp "github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	gloss "github.com/charmbracelet/lipgloss"
)

type model struct {
	ready bool

	lyricLines      [][]string
	overwrite       string
	header          string
	currentLyric    int
	cursor          int
	followSync      bool
	showTimestamps  bool
	showProgressBar bool
	// showLyricTimer  bool

	w, h int

	position, duration int

	lyricViewport vp.Model
	progressBar   progress.Model
	help          help.Model
}

func InitialModel() model {
	helpModel := help.New()
	helpModel.Styles.FullKey = gloss.NewStyle().Align(gloss.Center).Faint(false).Foreground(gloss.Color("11"))
	helpModel.Styles.FullDesc = gloss.NewStyle().Align(gloss.Center).Faint(false).Foreground(gloss.Color("15"))
	return model{
		ready: false,

		lyricLines:      [][]string{},
		overwrite:       "",
		header:          "",
		currentLyric:    -1,
		cursor:          -1,
		followSync:      true,
		showTimestamps:  false,
		showProgressBar: true,

		position: int(global.CurrentPlayer.Position),
		duration: int(global.CurrentSong.Duration),

		progressBar: progress.New(progress.WithSolidFill(global.CurrentConfig.Output.TUI.Colors.ProgressBarColor)),
		help:        helpModel,
	}
}

type animateProgressBarTick bool

func (m model) progressBarTick() tea.Cmd {
	return tea.Tick(1*time.Second, func(t time.Time) tea.Msg { return animateProgressBarTick(true) })
}

func (m model) Init() tea.Cmd {
	return tea.Sequence(
		tea.SetWindowTitle("lrcsnc"),
		tea.Batch(watchSongInfoChanges(), watchPlayerInfoChanges(), watchCurrentLyricChanges(), watchReceivedOverwrites(), m.progressBarTick()),
	)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case songInfoChanged:
		if global.CurrentSong.LyricsData.LyricsType >= 2 {
			m.lyricLines = [][]string{{"───"}}
		} else {
			m.lyricLines = make([][]string, 0, len(global.CurrentSong.LyricsData.Lyrics))
			for _, s := range global.CurrentSong.LyricsData.Lyrics {
				m.lyricLines = append(m.lyricLines, m.lyricWrap(s))
			}
		}
		switch global.CurrentSong.LyricsData.LyricsType {
		case 1:
			m.header = " (unsynced)"
		case 6:
			m.header = " (unknown error!!!)"
		default:
			m.header = ""
		}
		m.lyricViewport.SetContent(gloss.PlaceHorizontal(m.w, gloss.Center, m.lyricsView()))
		m.lyricViewport.SetYOffset(m.lyricViewport.Height)
		if global.CurrentSong.LyricsData.LyricsType != 0 {
			m.followSync = false
		} else {
			m.followSync = true
		}
		m.currentLyric = -1
		return m, watchSongInfoChanges()

	case playerInfoChanged:
		m.position = int(global.CurrentPlayer.Position)
		m.duration = int(global.CurrentSong.Duration)
		return m, tea.Batch(watchPlayerInfoChanges(), m.progressBar.SetPercent(float64(m.position)/float64(m.duration)))

	case animateProgressBarTick:
		if global.CurrentPlayer.IsPlaying {
			m.position++
		}
		return m, m.progressBarTick()

	case currentLyricChanged:
		m.currentLyric = int(msg)
		if m.followSync {
			m.cursor = m.currentLyric
		}
		m.lyricViewport.SetContent(gloss.PlaceHorizontal(m.w, gloss.Center, m.lyricsView()))
		m.lyricViewport.YOffset = min(m.lyricViewport.TotalLineCount()-m.lyricViewport.Height, max(0, m.calcYOffset(m.cursor)))
		return m, watchCurrentLyricChanges()

	case overwriteReceived:
		m.overwrite = string(msg)
		return m, tea.Batch(watchReceivedOverwrites(), tea.Tick(5*time.Second, func(t time.Time) tea.Msg { return overwriteReceived("") }))

	case tea.WindowSizeMsg:
		m.w = msg.Width
		m.h = msg.Height
		m.lyricLines = make([][]string, 0, len(global.CurrentSong.LyricsData.Lyrics))
		for _, s := range global.CurrentSong.LyricsData.Lyrics {
			m.lyricLines = append(m.lyricLines, m.lyricWrap(s))
		}
		headerHeight := gloss.Height(m.headerView())
		footerHeight := gloss.Height(m.footerView())
		verticalMarginHeight := headerHeight + footerHeight
		if !m.ready {
			m.ready = true
			m.lyricViewport = vp.New(m.w, m.h-verticalMarginHeight)
			m.lyricViewport.Style = m.lyricViewport.Style.AlignHorizontal(gloss.Center).AlignVertical(gloss.Center)
		} else {
			m.lyricViewport.Width = m.w
			m.lyricViewport.Height = m.h - verticalMarginHeight
		}
		m.lyricViewport.SetContent(gloss.PlaceHorizontal(m.w, gloss.Center, m.lyricsView()))
		m.lyricViewport.SetYOffset(min(m.lyricViewport.TotalLineCount()-m.lyricViewport.Height, max(0, m.calcYOffset(m.cursor))))

	case tea.KeyMsg:
		if m.help.ShowAll {
			m.help.ShowAll = false
			return m, nil
		}

		switch msg.String() {

		case "ctrl+c", "q":
			return m, tea.Quit

		case "?", "h":
			m.help.ShowAll = true
			return m, nil

		case "enter":
			if global.CurrentSong.LyricsData.LyricsType != 0 || m.currentLyric == -1 {
				return m, nil
			}

			if !player.PlayerInfoControllers[global.CurrentConfig.Player.PlayerProvider].SeekTo(global.CurrentSong.LyricsData.LyricTimestamps[m.cursor]) {
				OverwriteReceived <- "Couldn't seek to the lyric!"
			}

			return m, nil

		case "up", "k":
			if m.followSync {
				m.followSync = false
			}
			m.cursor--
			if m.cursor < 0 {
				m.cursor = 0
			}
			m.lyricViewport.SetContent(gloss.PlaceHorizontal(m.w, gloss.Center, m.lyricsView()))
			m.lyricViewport.SetYOffset(min(m.lyricViewport.TotalLineCount()-m.lyricViewport.Height, max(0, m.calcYOffset(m.cursor))))

		case "down", "j":
			if m.followSync {
				m.followSync = false
			}
			m.cursor++
			if m.cursor > len(m.lyricLines)-1 {
				m.cursor = len(m.lyricLines) - 1
			}
			m.lyricViewport.SetContent(gloss.PlaceHorizontal(m.w, gloss.Center, m.lyricsView()))
			m.lyricViewport.SetYOffset(min(m.lyricViewport.TotalLineCount()-m.lyricViewport.Height, max(0, m.calcYOffset(m.cursor))))

		case "f":
			if global.CurrentSong.LyricsData.LyricsType == 0 {
				m.followSync = !m.followSync
				if !m.followSync {
					if m.cursor == -1 {
						m.cursor = 0
					}
				} else {
					m.cursor = m.currentLyric
				}
			}
			m.lyricViewport.SetContent(gloss.PlaceHorizontal(m.w, gloss.Center, m.lyricsView()))
			m.lyricViewport.SetYOffset(min(m.lyricViewport.TotalLineCount()-m.lyricViewport.Height, max(0, m.calcYOffset(m.cursor))))

		case "t":
			m.showTimestamps = !m.showTimestamps
			m.lyricLines = make([][]string, 0, len(global.CurrentSong.LyricsData.Lyrics))
			for _, s := range global.CurrentSong.LyricsData.Lyrics {
				m.lyricLines = append(m.lyricLines, m.lyricWrap(s))
			}
			m.lyricViewport.SetContent(gloss.PlaceHorizontal(m.w, gloss.Center, m.lyricsView()))

		case "p":
			m.showProgressBar = !m.showProgressBar
			if m.showProgressBar {
				m.lyricViewport.Height -= 1
			} else {
				m.lyricViewport.Height += 1
			}
		}
	}

	return m, nil
}

func (m model) View() (s string) {
	switch {
	case m.help.ShowAll:
		return gloss.PlaceHorizontal(m.w, gloss.Center, gloss.PlaceVertical(m.h, gloss.Center, m.help.View(keys)))
	case m.ready:
		return gloss.JoinVertical(gloss.Center, m.headerView(), m.lyricViewport.View(), m.footerView())
	default:
		return gloss.NewStyle().AlignVertical(gloss.Center).AlignHorizontal(gloss.Center).Render("Loading...")
	}
}

func (m model) headerView() string {
	var title string
	if m.overwrite != "" {
		title = gloss.NewStyle().Foreground(gloss.Color("11")).AlignHorizontal(gloss.Center).Render(m.overwrite)
	} else {
		title = gloss.NewStyle().Foreground(gloss.Color("15")).AlignHorizontal(gloss.Center).Render(fmt.Sprintf("%s - %s%s", global.CurrentSong.Artist, global.CurrentSong.Title, m.header))
	}
	return gloss.NewStyle().BorderBottom(true).BorderStyle(gloss.NormalBorder()).Render(gloss.PlaceHorizontal(m.w-4, gloss.Center, title))
}

func (m model) footerView() string {
	if !m.showProgressBar {
		return ""
	}
	position := gloss.NewStyle().AlignHorizontal(gloss.Left).Render(positionIntoString(m.position) + " ")
	duration := gloss.NewStyle().AlignHorizontal(gloss.Right).Render(" " + positionIntoString(m.duration))
	m.progressBar.Empty = ' '
	m.progressBar.Full = '█'
	m.progressBar.Width = max(0, m.w-gloss.Width(position)-gloss.Width(duration)-4)
	m.progressBar.ShowPercentage = false
	return gloss.NewStyle().BorderTop(true).BorderStyle(gloss.NormalBorder()).Margin(0, 2).Render(gloss.JoinHorizontal(gloss.Center, position, m.progressBar.ViewAs(float64(m.position)/float64(m.duration)), duration))
}

func (m model) lyricsView() string {
	switch global.CurrentSong.LyricsData.LyricsType {
	case 2:
		return "Listen to the instrumental... ♪♪♪"
	case 3:
		return "This song was not found on LrcLib!"
	case 4:
		return ""
	case 5:
		return "Loading lyrics..."
	case 6:
		return ":c"
	case 0:
		lines := make([]string, 0, len(m.lyricLines))

		// Iterate over our lyrics
		for i, lyricLine := range m.lyricLines {
			line := ""
			stylizedLyrics := make([]string, 0, len(lyricLine))

			if len(lyricLine) == 1 && lyricLine[0] == "" {
				lyricLine[0] = "───"
			}

			if i < m.currentLyric {
				for _, l := range lyricLine {
					stylizedLyrics = append(stylizedLyrics, styleBefore().Render(l))
				}
				line = gloss.JoinVertical(gloss.Center, stylizedLyrics...)
			} else if i == m.currentLyric {
				for _, l := range lyricLine {
					stylizedLyrics = append(stylizedLyrics, styleCurrent().Render(l))
				}
				line = gloss.NewStyle().Margin(1, 0).Render(gloss.JoinVertical(gloss.Center, stylizedLyrics...))
			} else if i > m.currentLyric {
				for _, l := range lyricLine {
					stylizedLyrics = append(stylizedLyrics, styleAfter().Render(l))
				}
				line = gloss.JoinVertical(gloss.Center, stylizedLyrics...)
			}

			if !m.followSync {
				if i == m.cursor {
					stylizedLyrics = make([]string, 0, len(lyricLine))
					if i == m.currentLyric && global.CurrentSong.LyricsData.LyricsType == 0 {
						for _, l := range lyricLine {
							stylizedLyrics = append(stylizedLyrics, styleCurrent().Render(l))
						}
						line = gloss.JoinVertical(gloss.Center, stylizedLyrics...)
					} else {
						for _, l := range lyricLine {
							stylizedLyrics = append(stylizedLyrics, styleCursor().Render(l))
						}
						line = gloss.JoinVertical(gloss.Center, stylizedLyrics...)
					}
				}
			}

			var timestampView string = ""
			if m.showTimestamps {
				style := styleTimestamp()
				if i == m.cursor && i == m.currentLyric {
					style = styleTimestampCurrent()
				} else if i == m.cursor {
					style = styleTimestampCursor()
				}
				timestampView = style.Render(timestampIntoString(global.CurrentSong.LyricsData.LyricTimestamps[i])) + " "
			}

			line = gloss.JoinHorizontal(gloss.Center, timestampView, line)

			if i == m.currentLyric && i == m.cursor && !m.followSync && global.CurrentSong.LyricsData.LyricsType == 0 {
				line = styleBorderCursor().Render(line)
			}

			lines = append(lines, line)
		}

		return gloss.JoinVertical(gloss.Center, lines...)
	case 1:
		lines := make([]string, 0, len(m.lyricLines))

		// Iterate over our lyrics
		for i, lyricLine := range m.lyricLines {
			line := ""
			stylizedLyrics := make([]string, 0, len(lyricLine))

			if len(lyricLine) == 1 && lyricLine[0] == "" {
				lyricLine[0] = "───"
			}

			for _, l := range lyricLine {
				stylizedLyrics = append(stylizedLyrics, styleBefore().Render(l))
			}

			if i == m.cursor {
				stylizedLyrics = make([]string, 0, len(lyricLine))
				for _, l := range lyricLine {
					stylizedLyrics = append(stylizedLyrics, styleCursor().Render(l))
				}
				line = gloss.JoinVertical(gloss.Center, stylizedLyrics...)
			} else {
				line = gloss.JoinVertical(gloss.Center, stylizedLyrics...)
			}

			lines = append(lines, line)
		}

		return gloss.JoinVertical(gloss.Center, lines...)
	default:
		return ""
	}
}

// func main() {
// 	t := time.NewTimer(250 * time.Millisecond)
// 	t1 := time.NewTimer(4 * time.Second)
// 	t2 := time.NewTimer(7500 * time.Millisecond)
// 	go func() {
// 		<-t.C
// 		LyricsDataChan <- true
// 		CurrentLyricDataChan <- 0
// 		<-t1.C
// 		CurrentLyricDataChan <- 6
// 		<-t2.C
// 		CurrentLyricDataChan <- len(Lyrics) - 2
// 	}()
// }

func timestampIntoString(t float64) string {
	i, f := math.Modf(t)
	return fmt.Sprintf("%02d:%02d.%02d", int(i)/60, int(i)%60, int(math.Round(f*100)))
}

func positionIntoString(t int) string {
	return fmt.Sprintf("%02d:%02d", t/60, t%60)
}

func (m model) calcYOffset(l int) (res int) {
	if l < 0 || l >= len(m.lyricLines) {
		return 0
	}

	for i := 0; i < l; i++ {
		res += len(m.lyricLines[i])
		if m.currentLyric == i {
			res += 2
		}
	}
	if m.cursor == m.currentLyric {
		res += 1
	}
	res = max(0, res-int(math.Ceil(float64(m.lyricViewport.Height)/2)-1))
	return
}

func (m model) lyricWrap(l string) (res []string) {
	if gloss.Width(styleCurrent().Render(l)) < m.lyricViewport.Width*3/4 {
		return []string{l}
	} else {
		res = make([]string, 0)
		words := strings.Split(l, " ")
		s := strings.Builder{}
		for i := 0; i < len(words); i++ {
			if gloss.Width(styleCurrent().Render(s.String()+" "+words[i])) >= m.lyricViewport.Width*2/3 {
				res = append(res, s.String())
				s.Reset()
			}

			s.WriteString(words[i])
			if i != len(words)-1 {
				s.WriteString(" ")
			}
		}
		res = append(res, s.String())
		return res
	}
}
