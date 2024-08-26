package main

import (
	"fmt"
	"math"
	"os"
	"regexp"
	"strings"
	"time"
)

var badCharactersRegexp = regexp.MustCompile(`[:;|\/\\<>\.]+`)

func GetCachedLyrics(song *SongData) (string, bool) {
	cacheDirectory := CurrentConfig.Cache.CacheDir
	if strings.Contains(cacheDirectory, "$XDG_CACHE_DIR") && os.Getenv("$XDG_CACHE_DIR") == "" {
		cacheDirectory = strings.ReplaceAll(cacheDirectory, "$XDG_CACHE_DIR", "$HOME/.cache")
	}

	cacheDirectory = os.ExpandEnv(cacheDirectory)

	filename := getFilename(song.Song, song.Artist, song.Album, song.Duration)
	fullPath := cacheDirectory + "/" + filename + ".lrc"

	if file, err := os.ReadFile(fullPath); err == nil {
		if CurrentConfig.Cache.Enabled && CurrentConfig.Cache.CacheLifeSpan != 0 {
			cacheStats, _ := os.Lstat(fullPath)
			return string(file), time.Since(cacheStats.ModTime()).Hours() <= float64(CurrentConfig.Cache.CacheLifeSpan)*24
		} else {
			return string(file), true
		}
	} else {
		return "", false
	}
}

func StoreCachedLyrics(song *SongData, lrcData string) error {
	cacheDirectory := CurrentConfig.Cache.CacheDir
	if strings.Contains(cacheDirectory, "$XDG_CACHE_DIR") && os.Getenv("$XDG_CACHE_DIR") == "" {
		cacheDirectory = strings.ReplaceAll(cacheDirectory, "$XDG_CACHE_DIR", "$HOME/.cache")
	}

	cacheDirectory = os.ExpandEnv(cacheDirectory)

	if _, err := os.ReadDir(cacheDirectory); err != nil {
		os.Mkdir(cacheDirectory, 0777)
		os.Chmod(cacheDirectory, 0777)
	}

	filename := getFilename(song.Song, song.Artist, song.Album, song.Duration)
	fullPath := cacheDirectory + "/" + filename + ".lrc"

	if err := os.WriteFile(fullPath, []byte(lrcData), 0777); err != nil {
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
