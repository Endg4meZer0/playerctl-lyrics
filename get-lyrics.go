package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

// Make a URL to lrclib.net to send a GET request to
func MakeURL(song *SongData) url.URL {
	lrclibURL, err := url.Parse(fmt.Sprintf("http://lrclib.net/api/search?track_name=%v&artist_name=%v", song.Song, song.Artist))
	if err != nil {
		log.Fatalln(err)
	}
	return *lrclibURL
}

// Make a URL to lrclib.net to send a GET request to
func MakeURLWithAlbum(song *SongData) url.URL {
	lrclibURL, err := url.Parse(fmt.Sprintf("http://lrclib.net/api/search?track_name=%v&artist_name=%v&album_name=%v", song.Song, song.Artist, song.Album))
	if err != nil {
		log.Fatalln(err)
	}
	return *lrclibURL
}

// Return either a slice of strings that correspond to song's lyrics and a 'false' or nil and 'true'.
// If []string is nil AND bool is false, then it's an error.
func GetSyncedLyrics(song *SongData) (map[float64]string, bool) {
	var lrclibURL url.URL
	var output []LrcLibJsonOutput
	if song.Album != "" {
		lrclibURL = MakeURLWithAlbum(song)
		output = SendLrcLibGetRequest(lrclibURL)
		if len(output) == 0 {
			lrclibURL = MakeURL(song)
			output = SendLrcLibGetRequest(lrclibURL)
		}
	} else {
		lrclibURL = MakeURL(song)
		output = SendLrcLibGetRequest(lrclibURL)
	}

	if len(output) == 0 {
		return nil, false
	}

	if output[0].Instrumental {
		return nil, true
	}

	result := map[float64]string{}

	for _, respondedSong := range output {
		if respondedSong.SyncedLyrics != "" {
			syncedLyrics := strings.Split(respondedSong.SyncedLyrics, "\n")
			for _, lyric := range syncedLyrics {
				lyricParts := strings.SplitN(lyric, " ", 2)
				timecode := TimecodeStrToFloat(lyricParts[0])
				lyricStr := lyricParts[1]
				result[timecode] = lyricStr
			}
			return result, false
		}
	}

	return nil, false
}

func SendLrcLibGetRequest(lrclibURL url.URL) []LrcLibJsonOutput {
	urls := strings.ReplaceAll(lrclibURL.String(), " ", "%20")
	resp, err := http.Get(urls)
	if err != nil {
		log.Fatalln(err)
	}
	body, err := io.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		log.Fatalln(err)
	}
	var output []LrcLibJsonOutput
	json.Unmarshal(body, &output)
	return output
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
