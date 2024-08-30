package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

type SpotifyAccessTokenData struct {
	Token      string        `json:"token"`
	AcquiredAt time.Time     `json:"acquiredAt"`
	ExpiredIn  time.Duration `json:"expiredAt"`
}

var CurrentSpotifyAccessTokenData SpotifyAccessTokenData = SpotifyAccessTokenData{
	Token: "testeyton",
}

var client = &http.Client{Timeout: 5 * time.Second}

func Init() {
	cacheDirectory := CurrentConfig.Cache.CacheDir
	if strings.Contains(cacheDirectory, "$XDG_CACHE_DIR") && os.Getenv("$XDG_CACHE_DIR") == "" {
		cacheDirectory = strings.ReplaceAll(cacheDirectory, "$XDG_CACHE_DIR", "$HOME/.cache")
	}

	cacheDirectory = os.ExpandEnv(cacheDirectory)

	if file, err := os.ReadFile(cacheDirectory + "/spotifyStuff.json"); err == nil {
		var cachedSpotifyStuff SpotifyAccessTokenData
		err = json.Unmarshal(file, &cachedSpotifyStuff)
		if err != nil {
			log.Println(err)
			getSpotifyAccessToken()
		} else {
			CurrentSpotifyAccessTokenData.Token = cachedSpotifyStuff.Token
			CurrentSpotifyAccessTokenData.AcquiredAt = cachedSpotifyStuff.AcquiredAt
			CurrentSpotifyAccessTokenData.ExpiredIn = cachedSpotifyStuff.ExpiredIn
		}
	} else {
		getSpotifyAccessToken()
	}
}

func GetSongBPM(song *SongData) {
	if CurrentSpotifyAccessTokenData.Token == "testeyton" {
		Init()
	} else if time.Until(CurrentSpotifyAccessTokenData.AcquiredAt.Add(CurrentSpotifyAccessTokenData.ExpiredIn)) < 0 {
		getSpotifyAccessToken()
	}

	id := searchSpotifyTracks(song)
	song.BPM = getSpotifySongBPM(id)
}

func getSpotifyAccessToken() {
	req, err := http.NewRequest("POST", "https://accounts.spotify.com/api/token", bytes.NewBuffer([]byte(`grant_type=client_credentials`)))
	if err != nil {
		return
	}

	req.Header.Set("Authorization", "Basic "+base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%v:%v", CurrentConfig.Spotify.ClientId, CurrentConfig.Spotify.ClientSecret))))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := client.Do(req)
	if err != nil || resp.StatusCode != 200 {
		return
	}

	body, err := io.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		return
	}

	var output struct {
		AccessToken string `json:"access_token"`
		ExpiresIn   int    `json:"expires_in"`
	}
	err = json.Unmarshal(body, &output)
	if err != nil {
		return
	}

	CurrentSpotifyAccessTokenData.Token = output.AccessToken
	CurrentSpotifyAccessTokenData.AcquiredAt = time.Now()
	CurrentSpotifyAccessTokenData.ExpiredIn = time.Duration(output.ExpiresIn * int(time.Second))

	cacheDirectory := CurrentConfig.Cache.CacheDir
	if strings.Contains(cacheDirectory, "$XDG_CACHE_DIR") && os.Getenv("$XDG_CACHE_DIR") == "" {
		cacheDirectory = strings.ReplaceAll(cacheDirectory, "$XDG_CACHE_DIR", "$HOME/.cache")
	}

	cacheDirectory = os.ExpandEnv(cacheDirectory)

	data, err := json.Marshal(CurrentSpotifyAccessTokenData)
	if err != nil {
		log.Println(err)
	} else {
		err = os.WriteFile(cacheDirectory+"/spotifyStuff.json", []byte(data), 0777)
		if err != nil {
			log.Println(err)
		}
	}
}

func searchSpotifyTracks(song *SongData) string {
	reqURL, err := url.Parse("https://api.spotify.com/v1/search?" + url.PathEscape(fmt.Sprintf("q=%v+%v+%v&type=track&market=US&limit=1", strings.ReplaceAll(song.Song, " ", "+"), strings.ReplaceAll(song.Artist, " ", "+"), strings.ReplaceAll(song.Album, " ", "+"))))
	if err != nil {
		return ""
	}

	req := http.Request{
		URL: reqURL,
		Header: map[string][]string{
			"Authorization": {fmt.Sprintf("Bearer %v", CurrentSpotifyAccessTokenData.Token)},
		},
		Method: "GET",
	}

	resp, err := client.Do(&req)
	if err != nil || resp.StatusCode != 200 {
		return ""
	}

	body, err := io.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		return ""
	}

	var output struct {
		Tracks struct {
			Items []struct {
				Id string `json:"id"`
			} `json:"items"`
		} `json:"tracks"`
	}
	err = json.Unmarshal(body, &output)
	if err != nil || len(output.Tracks.Items) == 0 {
		return ""
	}

	return output.Tracks.Items[0].Id
}

func getSpotifySongBPM(id string) float64 {
	reqURL, err := url.Parse("https://api.spotify.com/v1/audio-features/" + url.PathEscape(id))
	if err != nil {
		return 0
	}

	req := http.Request{
		URL: reqURL,
		Header: map[string][]string{
			"Authorization": {fmt.Sprintf("Bearer %v", CurrentSpotifyAccessTokenData.Token)},
		},
		Method: "GET",
	}

	resp, err := client.Do(&req)
	if err != nil || resp.StatusCode != 200 {
		return 0
	}

	body, err := io.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		return 0
	}

	var output struct {
		Tempo float64 `json:"tempo"`
	}
	err = json.Unmarshal(body, &output)
	if err != nil {
		return 0
	}

	return output.Tempo
}
