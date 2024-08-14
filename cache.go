package main

import (
	"encoding/json"
	"fmt"
	"math"
	"os"
	"regexp"
	"strings"
	"time"
)

var badCharactersRegexp = regexp.MustCompile(`[:;|\/\\<>\.]+`)

func GetCachedLyrics(song *SongData) (LrcLibJson, bool) {
	cacheDirectory := CurrentConfig.Cache.CacheDir
	if strings.Contains(cacheDirectory, "$XDG_CACHE_DIR") && os.Getenv("$XDG_CACHE_DIR") == "" {
		cacheDirectory = strings.ReplaceAll(cacheDirectory, "$XDG_CACHE_DIR", "$HOME/.cache")
	}

	cacheDirectory = os.ExpandEnv(cacheDirectory)

	filename := getFilename(song.Song, song.Artist, song.Album, song.Duration)
	if file, err := os.ReadFile(cacheDirectory + "/" + filename + ".json"); err == nil {
		var result LrcLibJson
		err = json.Unmarshal(file, &result)
		if err != nil {
			return LrcLibJson{}, false
		}

		if CurrentConfig.Cache.DoCacheLyrics && CurrentConfig.Cache.CacheLifeSpan != 0 {
			cacheStats, _ := os.Lstat(cacheDirectory + "/" + filename + ".json")
			return result, time.Since(cacheStats.ModTime()).Hours() <= float64(CurrentConfig.Cache.CacheLifeSpan)*24
		} else {
			return result, true
		}
	} else {
		return LrcLibJson{}, false
	}
}

func StoreCachedLyrics(song *SongData, lrcData LrcLibJson) error {
	cacheDirectory := CurrentConfig.Cache.CacheDir
	if strings.Contains(cacheDirectory, "$XDG_CACHE_DIR") && os.Getenv("$XDG_CACHE_DIR") == "" {
		cacheDirectory = strings.ReplaceAll(cacheDirectory, "$XDG_CACHE_DIR", "$HOME/.cache")
	}

	cacheDirectory = os.ExpandEnv(cacheDirectory)

	if _, err := os.ReadDir(cacheDirectory); err != nil {
		os.Mkdir(cacheDirectory, 0777)
		os.Chmod(cacheDirectory, 0777)
	}

	if lrcData.PlainLyrics != "" {
		lrcData.PlainLyrics = "yes"
	}

	filename := getFilename(song.Song, song.Artist, song.Album, song.Duration)
	data, err := json.Marshal(lrcData)
	if err != nil {
		return err
	}
	if err = os.WriteFile(cacheDirectory+"/"+filename+".json", data, 0777); err != nil {
		return err
	}
	return nil
}

func getFilename(song string, artist string, album string, duration float64) string {
	return fmt.Sprintf("%v.%v.%v.%v", RemoveBadCharacters(song), RemoveBadCharacters(artist), RemoveBadCharacters(album), math.Round(duration))
}

func RemoveBadCharacters(str string) string {
	return badCharactersRegexp.ReplaceAllString(str, "_")
}
