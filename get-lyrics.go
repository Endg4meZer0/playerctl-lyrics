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

// Sets the Lyrics and LyricTimestamps properties in SongData object.
func GetSyncedLyrics(song *SongData) {
	var foundSong LrcLibJson
	cachedLyrics, isNotExpired := GetCachedLyrics(song)
	if cachedLyrics == "" || !isNotExpired {
		lrclibURL := makeURLGet(song)

		foundSongs, found := sendRequest(lrclibURL)

		if !found {
			lrclibURL = makeURLSearchWithAlbum(song)
			foundSongs, found = sendRequest(lrclibURL)
			if !found {
				lrclibURL = makeURLSearch(song)
				foundSongs, found = sendRequest(lrclibURL)
			}
		}

		if !found {
			song.LyricsType = 3
			return
		}

		foundSong = foundSongs[0]

		if foundSong.Instrumental {
			song.LyricsType = 2
			return
		}

		if foundSong.PlainLyrics != "" && foundSong.SyncedLyrics == "" {
			song.LyricsType = 1
			return
		}

		if CurrentConfig.Cache.Enabled {
			if err := StoreCachedLyrics(*song, foundSong.SyncedLyrics); err != nil {
				log.Println("Could not save the lyrics to the cache! Are there writing perms?")
			}
		}
	} else {
		foundSong.SyncedLyrics = cachedLyrics
	}

	song.LyricsType = 0

	resultLyrics := []string{}
	resultTimestamps := []float64{}

	syncedLyrics := strings.Split(foundSong.SyncedLyrics, "\n")
	for _, lyric := range syncedLyrics {
		lyricParts := strings.SplitN(lyric, " ", 2)
		timecode := timecodeStrToFloat(lyricParts[0])
		if timecode == -1 {
			continue
		}
		var lyricStr string
		if len(lyricParts) != 1 {
			lyricStr = lyricParts[1]
		} else {
			lyricStr = ""
		}
		resultLyrics = append(resultLyrics, lyricStr)
		resultTimestamps = append(resultTimestamps, timecode)
	}

	song.Lyrics = resultLyrics
	song.LyricTimestamps = resultTimestamps
}

func timecodeStrToFloat(timecode string) float64 {
	// [00:00.00]
	if len(timecode) != 10 {
		return -1
	}
	minutes, err := strconv.ParseFloat(timecode[1:3], 64)
	if err != nil {
		return -1
	}
	seconds, err := strconv.ParseFloat(timecode[4:9], 64)
	if err != nil {
		return -1
	}
	return minutes*60.0 + seconds
}

// Make a URL to lrclib.net/api/get to send a GET request to
func makeURLGet(song *SongData) url.URL {
	lrclibURL, err := url.Parse("http://lrclib.net/api/get?" + url.PathEscape(fmt.Sprintf("track_name=%v&artist_name=%v&album_name=%v&duration=%v", song.Song, song.Artist, song.Album, int(math.Ceil(song.Duration)))))
	if err != nil {
		log.Fatalln(err)
	}
	return *lrclibURL
}

// Make a URL to lrclib.net/api/search with album data to send a GET request to
func makeURLSearchWithAlbum(song *SongData) url.URL {
	lrclibURL, err := url.Parse("http://lrclib.net/api/search?" + url.PathEscape(fmt.Sprintf("track_name=%v&artist_name=%v&album_name=%v", song.Song, song.Artist, song.Album)))
	if err != nil {
		log.Fatalln(err)
	}
	return *lrclibURL
}

// Make a URL to lrclib.net/api/search only with necessary data (song name and artist name) to send a GET request to
func makeURLSearch(song *SongData) url.URL {
	lrclibURL, err := url.Parse("http://lrclib.net/api/search?" + url.PathEscape(fmt.Sprintf("track_name=%v&artist_name=%v", song.Song, song.Artist)))
	if err != nil {
		log.Fatalln(err)
	}
	return *lrclibURL
}

func sendRequest(link url.URL) ([]LrcLibJson, bool) {
	resp, err := http.Get(link.String())
	if err != nil || resp.StatusCode != 200 {
		return nil, false
	}

	body, err := io.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		return nil, false
	}

	var foundSong LrcLibJson
	if json.Unmarshal(body, &foundSong) != nil {
		var foundSongs []LrcLibJson
		json.Unmarshal(body, &foundSongs)

		return foundSongs, len(foundSongs) != 0
	} else {
		return []LrcLibJson{foundSong}, true
	}
}
