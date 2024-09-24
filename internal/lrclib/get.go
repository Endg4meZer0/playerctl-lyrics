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
	"lrcsnc/internal/util"
	"lrcsnc/pkg"
)

// Sets the Lyrics and LyricTimestamps properties in SongData object.
func GetSyncedLyrics(song *pkg.SongData) {
	if pkg.CurrentConfig.Cache.Enabled {
		cachedData, isNotExpired := cache.GetCachedLyrics(song)
		if isNotExpired && !(len(cachedData.Lyrics) == 0 && !cachedData.Instrumental) {
			song.LyricTimestamps = cachedData.LyricTimestamps
			song.Lyrics = cachedData.Lyrics
			if cachedData.Instrumental {
				song.LyricsType = 2
			} else {
				song.LyricsType = 0
			}
			return
		}
	}

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

	foundSong := foundSongs[0]

	resultLyrics := []string{}
	resultTimestamps := []float64{}

	if foundSong.Instrumental {
		song.LyricsType = 2
	} else if foundSong.PlainLyrics != "" && foundSong.SyncedLyrics == "" {
		song.LyricsType = 1
	} else {
		song.LyricsType = 0

		syncedLyrics := strings.Split(foundSong.SyncedLyrics, "\n")

		resultLyrics = make([]string, len(syncedLyrics))
		resultTimestamps = make([]float64, len(syncedLyrics))

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

		song.Lyrics = resultLyrics
		song.LyricTimestamps = resultTimestamps
	}

	if pkg.CurrentConfig.Cache.Enabled && song.LyricsType != 1 {
		var dataToBeCached cache.Cache = cache.Cache{
			LyricTimestamps: resultTimestamps,
			Lyrics:          resultLyrics,
			Instrumental:    false,
		}
		if song.LyricsType == 2 {
			dataToBeCached.Instrumental = true
		}
		if err := cache.StoreCachedLyrics(*song, dataToBeCached); err != nil {
			log.Println("Could not save the lyrics to the cache! Are there writing perms?")
		}
	}
}

// Make a URL to lrclib.net/api/get to send a GET request to
func makeURLGet(song *pkg.SongData) url.URL {
	lrclibURL, err := url.Parse("http://lrclib.net/api/get?" + url.PathEscape(fmt.Sprintf("track_name=%v&artist_name=%v&album_name=%v&duration=%v", song.Song, song.Artist, song.Album, int(math.Ceil(song.Duration)))))
	if err != nil {
		log.Fatalln(err)
	}
	return *lrclibURL
}

// Make a URL to lrclib.net/api/search with album data to send a GET request to
func makeURLSearchWithAlbum(song *pkg.SongData) url.URL {
	lrclibURL, err := url.Parse("http://lrclib.net/api/search?" + url.PathEscape(fmt.Sprintf("track_name=%v&artist_name=%v&album_name=%v", song.Song, song.Artist, song.Album)))
	if err != nil {
		log.Fatalln(err)
	}
	return *lrclibURL
}

// Make a URL to lrclib.net/api/search only with necessary data (song name and artist name) to send a GET request to
func makeURLSearch(song *pkg.SongData) url.URL {
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
