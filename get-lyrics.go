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

// Make a URL to lrclib.net/api/get to send a GET request to
func MakeURLGet(song *SongData) url.URL {
	lrclibURL, err := url.Parse("http://lrclib.net/api/get?" + url.PathEscape(fmt.Sprintf("track_name=%v&artist_name=%v&album_name=%v&duration=%v", song.Song, song.Artist, song.Album, int(math.Ceil(song.Duration)))))
	if err != nil {
		log.Fatalln(err)
	}
	return *lrclibURL
}

// Make a URL to lrclib.net/api/search with album and duration data to send a GET request to
func MakeURLSearchWithAlbumAndDuration(song *SongData) url.URL {
	lrclibURL, err := url.Parse("http://lrclib.net/api/search?" + url.PathEscape(fmt.Sprintf("track_name=%v&artist_name=%v&album_name=%v&duration=%v", song.Song, song.Artist, song.Album, int(math.Ceil(song.Duration)))))
	if err != nil {
		log.Fatalln(err)
	}
	return *lrclibURL
}

// Make a URL to lrclib.net/api/search with album data to send a GET request to
func MakeURLSearchWithAlbum(song *SongData) url.URL {
	lrclibURL, err := url.Parse("http://lrclib.net/api/search?" + url.PathEscape(fmt.Sprintf("track_name=%v&artist_name=%v&album_name=%v", song.Song, song.Artist, song.Album)))
	if err != nil {
		log.Fatalln(err)
	}
	return *lrclibURL
}

// Make a URL to lrclib.net/api/search only with necessary data (song name and artist name) to send a GET request to
func MakeURLSearch(song *SongData) url.URL {
	lrclibURL, err := url.Parse("http://lrclib.net/api/search?" + url.PathEscape(fmt.Sprintf("track_name=%v&artist_name=%v", song.Song, song.Artist)))
	if err != nil {
		log.Fatalln(err)
	}
	return *lrclibURL
}

// Return either a slice of strings that correspond to song's lyrics and a 'false' or nil and 'true'.
// If []string is nil AND bool is false, then it's an error.
func GetSyncedLyrics(song *SongData) map[float64]string {
	lrclibURL := MakeURLGet(song)

	foundSongs, found := SendRequest(lrclibURL)

	if !found {
		lrclibURL = MakeURLSearchWithAlbumAndDuration(song)
		foundSongs, found = SendRequest(lrclibURL)
		if !found {
			lrclibURL = MakeURLSearchWithAlbum(song)
			foundSongs, found = SendRequest(lrclibURL)
			if !found {
				lrclibURL = MakeURLSearch(song)
				foundSongs, found = SendRequest(lrclibURL)
			}
		}
	}

	if !found {
		song.LyricsType = 3
		return nil
	}

	foundSong := foundSongs[0]

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

func SendRequest(link url.URL) ([]LrcLibJsonOutput, bool) {
	resp, err := http.Get(link.String())
	if err != nil || resp.StatusCode != 200 {
		return nil, false
	}

	body, err := io.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		return nil, false
	}

	var foundSongs []LrcLibJsonOutput
	json.Unmarshal(body, &foundSongs)

	return foundSongs, len(foundSongs) != 0
}
