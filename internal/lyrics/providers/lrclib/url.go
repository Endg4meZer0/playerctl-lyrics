package lrclib

import (
	"fmt"
	"lrcsnc/internal/pkg/structs"
	"math"
	"net/url"
)

// Make a URL to lrclib.net/api/get to send a GET request to
func makeURLGet(song structs.Song) *url.URL {
	lrclibURL, err := url.Parse("http://lrclib.net/api/get?" + url.PathEscape(fmt.Sprintf("track_name=%v&artist_name=%v&album_name=%v&duration=%v", song.Title, song.Artist, song.Album, int(math.Ceil(song.Duration)))))
	if err != nil {
		// TODO: logger :)
	}
	return lrclibURL
}

// Make a URL to lrclib.net/api/search with album data to send a GET request to
func makeURLSearchWithAlbum(song structs.Song) *url.URL {
	lrclibURL, err := url.Parse("http://lrclib.net/api/search?" + url.PathEscape(fmt.Sprintf("track_name=%v&artist_name=%v&album_name=%v", song.Title, song.Artist, song.Album)))
	if err != nil {
		// TODO: logger :)
	}
	return lrclibURL
}

// Make a URL to lrclib.net/api/search only with necessary data (song name and artist name) to send a GET request to
func makeURLSearch(song structs.Song) *url.URL {
	lrclibURL, err := url.Parse("http://lrclib.net/api/search?" + url.PathEscape(fmt.Sprintf("track_name=%v&artist_name=%v", song.Title, song.Artist)))
	if err != nil {
		// TODO: logger :)
	}
	return lrclibURL
}
