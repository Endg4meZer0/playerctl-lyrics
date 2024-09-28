package lrclib

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math"
	"net/http"
	"net/url"
	"strings"

	"lrcsnc/internal/cache"
	"lrcsnc/internal/pkg/global"
	"lrcsnc/internal/pkg/structs"
	"lrcsnc/internal/util"
)

type LrcLibLyricsProvider struct{}

// Sets the Lyrics and LyricTimestamps properties in SongInfo object.
func (l LrcLibLyricsProvider) GetLyricsData(song structs.SongInfo) (out structs.SongLyricsData) {
	if global.CurrentConfig.Cache.Enabled {
		cachedData, isNotExpired := cache.GetCachedLyrics(song)
		if isNotExpired && !(len(cachedData.Lyrics) == 0 && cachedData.LyricsType != 3) {
			out = cachedData
			return
		}
	}

	var getURL url.URL
	var foundSongs []LrcLibJson
	var found bool = false

	if song.Duration != 0 {
		getURL := makeURLGet(song)
		foundSongs, found = sendRequest(getURL)
	}

	if !found {
		getURL = makeURLSearchWithAlbum(song)
		foundSongs, found = sendRequest(getURL)
		if !found {
			getURL = makeURLSearch(song)
			foundSongs, found = sendRequest(getURL)
		}
	}

	if !found {
		out.LyricsType = 3
		return
	}

	foundSong := foundSongs[0]

	if foundSong.Instrumental {
		out.LyricsType = 2
	} else if foundSong.PlainLyrics != "" && foundSong.SyncedLyrics == "" {
		out.Lyrics = strings.Split(foundSong.PlainLyrics, "\n")
		out.LyricsType = 1
	} else {
		out.LyricsType = 0

		syncedLyrics := strings.Split(foundSong.SyncedLyrics, "\n")

		resultLyrics := make([]string, len(syncedLyrics))
		resultTimestamps := make([]float64, len(syncedLyrics))

		for i, lyric := range syncedLyrics {
			lyricParts := strings.SplitN(lyric, " ", 2)
			timecode := util.TimecodeToFloat(lyricParts[0])
			if timecode == -1 {
				continue
			}
			var lyricStr string
			if len(lyricParts) != 1 {
				lyricStr = lyricParts[1]
			} else {
				lyricStr = ""
			}
			resultLyrics[i] = lyricStr
			resultTimestamps[i] = timecode
		}

		out.Lyrics = resultLyrics
		out.LyricTimestamps = resultTimestamps
	}

	if global.CurrentConfig.Cache.Enabled && out.LyricsType != 1 {
		if err := cache.StoreCachedLyrics(song, out); err != nil {
			log.Println("Could not save the lyrics to the cache! Are there writing perms?")
		}
	}

	return
}

// Make a URL to lrclib.net/api/get to send a GET request to
func makeURLGet(song structs.SongInfo) url.URL {
	lrclibURL, err := url.Parse("http://lrclib.net/api/get?" + url.PathEscape(fmt.Sprintf("track_name=%v&artist_name=%v&album_name=%v&duration=%v", song.Title, song.Artist, song.Album, int(math.Ceil(song.Duration)))))
	if err != nil {
		log.Fatalln(err)
	}
	return *lrclibURL
}

// Make a URL to lrclib.net/api/search with album data to send a GET request to
func makeURLSearchWithAlbum(song structs.SongInfo) url.URL {
	lrclibURL, err := url.Parse("http://lrclib.net/api/search?" + url.PathEscape(fmt.Sprintf("track_name=%v&artist_name=%v&album_name=%v", song.Title, song.Artist, song.Album)))
	if err != nil {
		log.Fatalln(err)
	}
	return *lrclibURL
}

// Make a URL to lrclib.net/api/search only with necessary data (song name and artist name) to send a GET request to
func makeURLSearch(song structs.SongInfo) url.URL {
	lrclibURL, err := url.Parse("http://lrclib.net/api/search?" + url.PathEscape(fmt.Sprintf("track_name=%v&artist_name=%v", song.Title, song.Artist)))
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
