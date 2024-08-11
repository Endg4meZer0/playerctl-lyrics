package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

// Make a URL to lrclib.net to send a GET request to
func MakeURL(song *SongData) url.URL {
	lrclibURL, err := url.Parse("http://lrclib.net/api/get?" + url.PathEscape(fmt.Sprintf("track_name=%v&artist_name=%v&album_name=%v&duration=%v", song.Song, song.Artist, song.Album, int(math.Ceil(song.Duration)))))
	if err != nil {
		log.Fatalln(err)
	}
	return *lrclibURL
}

// Return either a slice of strings that correspond to song's lyrics and a 'false' or nil and 'true'.
// If []string is nil AND bool is false, then it's an error.
func GetSyncedLyrics(song *SongData) map[float64]string {
	lrclibURL := MakeURL(song)

	resp, err := http.Get(lrclibURL.String())
	if err != nil {
		fmt.Println("Could not make a GET request to LrcLib!")
		return nil
	}

	if resp.StatusCode != 200 {
		song.LyricsType = 3
		return nil
	}

	body, err := io.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		log.Fatalln(err)
	}
	var foundSong LrcLibJsonOutput
	json.Unmarshal(body, &foundSong)

	if foundSong.Instrumental {
		song.LyricsType = 2
		return nil
	}

	if foundSong.PlainLyrics != "" && foundSong.SyncedLyrics == "" {
		song.LyricsType = 1
		return nil
	}

	song.LyricsType = 0

	result := map[float64]string{}

	syncedLyrics := strings.Split(foundSong.SyncedLyrics, "\n")
	for _, lyric := range syncedLyrics {
		lyricParts := strings.SplitN(lyric, " ", 2)
		timecode := TimecodeStrToFloat(lyricParts[0])
		lyricStr := lyricParts[1]
		result[timecode] = lyricStr
	}
	return result
}

func TimecodeStrToFloat(timecode string) float64 {
	// [00:00.00]
	minutes, err := strconv.ParseFloat(timecode[1:3], 64)
	if err != nil {
		panic(err)
	}
	seconds, err := strconv.ParseFloat(timecode[4:9], 64)
	if err != nil {
		panic(err)
	}
	return minutes*60.0 + seconds
}
