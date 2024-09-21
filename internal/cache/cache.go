package cache

import (
	"encoding/json"
	"fmt"
	"log"
	"math"
	"os"
	"strings"
	"time"

	"lrcsnc/pkg"
	"lrcsnc/util"
)

type Cache struct {
	LyricTimestamps []float64 `json:"lyricTimestamps"`
	Lyrics          []string  `json:"lyrics"`
	Instrumental    bool      `json:"instrumental"`
}

func GetCachedLyrics(song *pkg.SongData) (Cache, bool) {
	if !pkg.CurrentConfig.Cache.Enabled {
		return Cache{}, true
	}
	cacheDirectory := pkg.CurrentConfig.Cache.CacheDir
	if strings.Contains(cacheDirectory, "$XDG_CACHE_DIR") && os.Getenv("$XDG_CACHE_DIR") == "" {
		cacheDirectory = strings.ReplaceAll(cacheDirectory, "$XDG_CACHE_DIR", "$HOME/.cache")
	}

	cacheDirectory = os.ExpandEnv(cacheDirectory)

	filename := getFilename(song.Song, song.Artist, song.Album, song.Duration)
	fullPath := cacheDirectory + "/" + filename + ".json"

	if file, err := os.ReadFile(fullPath); err == nil {
		var cachedData Cache
		err = json.Unmarshal(file, &cachedData)
		if err != nil {
			log.Println(err)
			return Cache{}, false
		}

		if pkg.CurrentConfig.Cache.CacheLifeSpan != 0 {
			cacheStats, _ := os.Lstat(fullPath)
			return cachedData, time.Since(cacheStats.ModTime()).Hours() <= float64(pkg.CurrentConfig.Cache.CacheLifeSpan)*24
		} else {
			return cachedData, true
		}
	} else {
		return Cache{}, false
	}
}

func StoreCachedLyrics(song pkg.SongData, data Cache) error {
	cacheDirectory := pkg.CurrentConfig.Cache.CacheDir
	if strings.Contains(cacheDirectory, "$XDG_CACHE_DIR") && os.Getenv("$XDG_CACHE_DIR") == "" {
		cacheDirectory = strings.ReplaceAll(cacheDirectory, "$XDG_CACHE_DIR", "$HOME/.cache")
	}

	cacheDirectory = os.ExpandEnv(cacheDirectory)

	if _, err := os.ReadDir(cacheDirectory); err != nil {
		os.Mkdir(cacheDirectory, 0777)
		os.Chmod(cacheDirectory, 0777)
	}

	filename := getFilename(song.Song, song.Artist, song.Album, song.Duration)
	fullPath := cacheDirectory + "/" + filename + ".json"

	encodedData, err := json.Marshal(data)
	if err != nil {
		return err
	}

	if err := os.WriteFile(fullPath, []byte(encodedData), 0777); err != nil {
		return err
	}
	return nil
}

func getFilename(song string, artist string, album string, duration float64) string {
	return fmt.Sprintf("%v.%v.%v.%v", util.RemoveBadCharacters(song), util.RemoveBadCharacters(artist), util.RemoveBadCharacters(album), math.Round(duration))
}
