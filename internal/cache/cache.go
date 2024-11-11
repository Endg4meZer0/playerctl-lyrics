package cache

import (
	"encoding/json"
	"fmt"
	"math"
	"os"
	"strings"
	"time"

	"lrcsnc/internal/pkg/global"
	"lrcsnc/internal/pkg/structs"
	"lrcsnc/internal/pkg/util"
)

type CacheState byte

var (
	CacheStateActive      CacheState = 0
	CacheStateExpired     CacheState = 1
	CacheStateNonExistant CacheState = 1
	CacheStateDisabled    CacheState = 2
)

// Returns the specified song's lyrics data from cache. The returned boolean is true if the cached data exists, is not expired and the function didn't end up with error.
func GetCachedLyrics(song structs.SongInfo) (structs.SongLyricsData, CacheState) {
	if !global.CurrentConfig.Cache.Enabled {
		return structs.SongLyricsData{}, CacheStateDisabled
	}
	cacheDirectory := GetCacheDir()

	filename := getFilename(song.Title, song.Artist, song.Album, song.Duration)
	fullPath := cacheDirectory + "/" + filename + ".json"

	if file, err := os.ReadFile(fullPath); err == nil {
		var cachedData structs.SongLyricsData
		err = json.Unmarshal(file, &cachedData)
		if err != nil {
			return structs.SongLyricsData{}, CacheStateNonExistant
		}

		if global.CurrentConfig.Cache.CacheLifeSpan != 0 {
			cacheStats, _ := os.Lstat(fullPath)
			isExpired := time.Since(cacheStats.ModTime()).Hours() <= float64(global.CurrentConfig.Cache.CacheLifeSpan)*24
			if isExpired {
				return cachedData, CacheStateExpired
			} else {
				return cachedData, CacheStateActive
			}
		} else {
			return cachedData, CacheStateActive
		}
	} else {
		return structs.SongLyricsData{}, CacheStateNonExistant
	}
}

// Stores the specified song's data to cache
func StoreCachedLyrics(song structs.SongInfo) error {
	cacheDirectory := GetCacheDir()
	if _, err := os.ReadDir(cacheDirectory); err != nil {
		os.Mkdir(cacheDirectory, 0777)
		os.Chmod(cacheDirectory, 0777)
	}

	filename := getFilename(song.Title, song.Artist, song.Album, song.Duration)
	fullPath := cacheDirectory + "/" + filename + ".json"

	encodedData, err := json.Marshal(song.LyricsData)
	if err != nil {
		return err
	}

	if err := os.WriteFile(fullPath, []byte(encodedData), 0777); err != nil {
		return err
	}
	return nil
}

// Delete the specified song's cached data
func RemoveCachedLyrics(song structs.SongInfo) error {
	cacheDirectory := GetCacheDir()
	if _, err := os.ReadDir(cacheDirectory); err != nil {
		os.Mkdir(cacheDirectory, 0777)
		os.Chmod(cacheDirectory, 0777)
	}

	filename := getFilename(song.Title, song.Artist, song.Album, song.Duration)
	fullPath := cacheDirectory + "/" + filename + ".json"

	if err := os.Remove(fullPath); err != nil {
		return fmt.Errorf("[internal/cache/RemoveCachedLyrics] ERROR: Couldn't delete the specified song's (%v) cached data. Maybe the data didn't exist in the first place?", filename)
	}
	return nil
}

func GetCacheDir() string {
	cacheDirectory := global.CurrentConfig.Cache.CacheDir
	if strings.Contains(cacheDirectory, "$XDG_CACHE_DIR") && os.Getenv("$XDG_CACHE_DIR") == "" {
		cacheDirectory = strings.ReplaceAll(cacheDirectory, "$XDG_CACHE_DIR", "$HOME/.cache")
	}

	cacheDirectory = os.ExpandEnv(cacheDirectory)
	return cacheDirectory
}

func getFilename(song string, artist string, album string, duration float64) string {
	return fmt.Sprintf("%v.%v.%v.%v",
		util.RemoveBadCharacters(song),
		util.RemoveBadCharacters(artist),
		util.RemoveBadCharacters(album),
		math.Round(duration),
	)
}
