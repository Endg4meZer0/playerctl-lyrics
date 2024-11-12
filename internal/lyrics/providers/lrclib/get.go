package lrclib

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	dto "lrcsnc/internal/lyrics/dto"
	lrclibdto "lrcsnc/internal/lyrics/dto/lrclib"
	"lrcsnc/internal/pkg/structs"
)

type LrcLibLyricsProvider struct{}

// GetLyricsData gets the lyrics data for a song from the LrcLib API
func (l LrcLibLyricsProvider) GetLyricsData(song structs.Song) (dto.LyricsDTO, error) {
	// If the song duration is 0, return an error
	if song.Duration == 0 {
		return nil, fmt.Errorf("[lyrics/providers/lrclib/get] WARNING: Song duration is 0, cannot get lyrics")
	}

	// Make a URL to the LrcLib API's `get` endpoint and send a GET request to it
	getURL := makeURLGet(song)
	foundSongs, err := sendRequest(getURL)

	if err != nil {
		getURL = makeURLSearchWithAlbum(song)
		foundSongs, err = sendRequest(getURL)
		if err != nil {
			getURL = makeURLSearch(song)
			foundSongs, err = sendRequest(getURL)
		}
	}

	return foundSongs[0], nil
}

func sendRequest(link *url.URL) ([]lrclibdto.LrcLibDTO, error) {
	resp, err := http.Get((*link).String())
	if err != nil || resp.StatusCode != 200 {
		return nil, fmt.Errorf("[lyrics/providers/lrclib/get] WARNING: Couldn't get a successful response. The track is probably missing from LrcLib's library")
	}

	body, err := io.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		return nil, fmt.Errorf("[lyrics/providers/lrclib/get] WARNING: Couldn't properly read the body of the response")
	}

	var foundSong lrclibdto.LrcLibDTO
	if json.Unmarshal(body, &foundSong) != nil {
		var foundSongs []lrclibdto.LrcLibDTO = make([]lrclibdto.LrcLibDTO, 0)
		json.Unmarshal(body, &foundSongs)

		return foundSongs, nil
	} else {
		return []lrclibdto.LrcLibDTO{foundSong}, nil
	}
}
