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

type ResponseStatus struct {
	Status byte
	Error  error
}

const (
	Success byte = iota
	NotFound
	ClientError
	ServerError
)

// GetLyricsData gets the lyrics data for a song from the LrcLib API
func (l LrcLibLyricsProvider) GetLyricsData(song structs.Song) (dto.LyricsDTO, error) {
	if song.Duration == 0 {
		// TODO: logger :)
		return nil, fmt.Errorf("[lyrics/providers/lrclib/get] WARNING: Song duration is 0, cannot get lyrics")
	}

	getURL := makeURLGet(song)
	foundSongs, status := sendRequest(getURL)
	if status.Status == ClientError || status.Status == ServerError {
		return nil, status.Error
	}
	matchedSongs := l.RemoveMismatches(song, foundSongs)

	if status.Status == NotFound || len(matchedSongs) == 0 {
		getURL = makeURLSearchWithAlbum(song)
		foundSongs, status = sendRequest(getURL)

		if status.Status == ClientError || status.Status == ServerError {
			return nil, status.Error
		}
		matchedSongs = l.RemoveMismatches(song, foundSongs)

		if status.Status == NotFound || len(matchedSongs) == 0 {
			getURL = makeURLSearch(song)
			foundSongs, status = sendRequest(getURL)

			if status.Status == ClientError || status.Status == ServerError {
				return nil, status.Error
			}
			matchedSongs = l.RemoveMismatches(song, foundSongs)
		}
	}

	if len(matchedSongs) == 0 {
		// TODO: logger :)
		return nil, fmt.Errorf("[lyrics/providers/lrclib/get] WARNING: Couldn't find any matching songs")
	}

	return matchedSongs[0], nil
}

func sendRequest(link *url.URL) ([]dto.LyricsDTO, ResponseStatus) {
	resp, err := http.Get((*link).String())
	if resp.StatusCode == 404 {
		return nil, ResponseStatus{Status: NotFound}
	}
	if err != nil || resp.StatusCode != 200 {
		return nil, ResponseStatus{Status: ServerError, Error: fmt.Errorf("[lyrics/providers/lrclib/get] WARNING: Couldn't get a successful response: %v", err)}
	}

	body, err := io.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		return nil, ResponseStatus{Status: ClientError, Error: fmt.Errorf("[lyrics/providers/lrclib/get] WARNING: Couldn't read response body: %v", err)}
	}

	var foundSong lrclibdto.LrcLibDTO
	if err := json.Unmarshal(body, &foundSong); err != nil {
		var foundSongs []dto.LyricsDTO = make([]dto.LyricsDTO, 0)
		err = json.Unmarshal(body, &foundSongs)
		if err != nil {
			return nil, ResponseStatus{Status: ClientError, Error: fmt.Errorf("[lyrics/providers/lrclib/get] WARNING: An error occured while unmarshalling the found songs")}
		}

		return foundSongs, ResponseStatus{Status: Success}
	} else {
		return []dto.LyricsDTO{foundSong}, ResponseStatus{Status: Success}
	}
}
