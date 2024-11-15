package lrclib

import (
	"encoding/json"
	"io"
	"net/http"
	"net/url"

	dto "lrcsnc/internal/lyrics/dto"
	lrclibdto "lrcsnc/internal/lyrics/dto/lrclib"
	"lrcsnc/internal/pkg/structs"
)

type ErrorResponse struct {
	Value string
}

func (e ErrorResponse) Error() string {
	return e.Value
}

// GetLyricsData gets the lyrics data for a song from the LrcLib API
func (l LrcLibLyricsProvider) GetLyricsDTOList(song structs.Song) ([]dto.LyricsDTO, error) {
	var getURL *url.URL
	var foundSongs []dto.LyricsDTO
	var matchedSongs []dto.LyricsDTO
	var err error

	if song.Duration != 0 {
		getURL = makeURLGet(song)
		foundSongs, err = sendRequest(getURL)
		if err == nil {
			matchedSongs = dto.RemoveMismatches(song, foundSongs)
		}
	}

	if len(matchedSongs) == 0 {
		getURL = makeURLSearchWithAlbum(song)
		foundSongs, err = sendRequest(getURL)
		if err == nil {
			matchedSongs = dto.RemoveMismatches(song, foundSongs)
		}

		if len(matchedSongs) == 0 {
			getURL = makeURLSearch(song)
			foundSongs, err = sendRequest(getURL)
			if err == nil {
				matchedSongs = dto.RemoveMismatches(song, foundSongs)
			}
		}
	}

	return matchedSongs, nil
}

func sendRequest(link *url.URL) ([]dto.LyricsDTO, error) {
	resp, err := http.Get((*link).String())
	if resp.StatusCode == 404 {
		return nil, ErrorResponse{Value: "Song is not found"}
	}
	if err != nil || resp.StatusCode != 200 {
		return nil, ErrorResponse{Value: "Unknown server error"}
	}

	body, err := io.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		return nil, ErrorResponse{Value: "Failure to read response body"}
	}

	var foundSong lrclibdto.LrcLibDTO
	if err := json.Unmarshal(body, &foundSong); err != nil {
		var foundSongs []dto.LyricsDTO = make([]dto.LyricsDTO, 0)
		err = json.Unmarshal(body, &foundSongs)
		if err != nil {
			return nil, ErrorResponse{Value: "Unmarshal error"}
		}

		if len(foundSongs) == 0 {
			return foundSongs, ErrorResponse{Value: "Song is not found"}
		}

		return foundSongs, nil
	} else {
		return []dto.LyricsDTO{foundSong}, nil
	}
}
