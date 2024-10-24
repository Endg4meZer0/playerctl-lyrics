package cache

import (
	"encoding/json"
	"fmt"
	"log"
	"math"
	"os"
	"strings"
	"time"

	"lrcsnc/internal/pkg/global"
	"lrcsnc/internal/pkg/structs"
	"lrcsnc/internal/util"
)

func GetCachedLyrics(song structs.SongInfo) (structs.SongLyricsData, bool) {
	if !global.CurrentConfig.Cache.Enabled {
		return structs.SongLyricsData{}, true
	}
	cacheDirectory := global.CurrentConfig.Cache.CacheDir
	if strings.Contains(cacheDirectory, "$XDG_CACHE_DIR") && os.Getenv("$XDG_CACHE_DIR") == "" {
		cacheDirectory = strings.ReplaceAll(cacheDirectory, "$XDG_CACHE_DIR", "$HOME/.cache")
	}

	cacheDirectory = os.ExpandEnv(cacheDirectory)

	filename := getFilename(song.Title, song.Artist, song.Album, song.Duration)
	fullPath := cacheDirectory + "/" + filename + ".json"

	if file, err := os.ReadFile(fullPath); err == nil {
		var cachedData structs.SongLyricsData
		err = json.Unmarshal(file, &cachedData)
		if err != nil {
			log.Println(err)
			return structs.SongLyricsData{}, false
		}

		if global.CurrentConfig.Cache.CacheLifeSpan != 0 {
			cacheStats, _ := os.Lstat(fullPath)
			return cachedData, time.Since(cacheStats.ModTime()).Hours() <= float64(global.CurrentConfig.Cache.CacheLifeSpan)*24
		} else {
			return cachedData, true
		}
	} else {
		return structs.SongLyricsData{}, false
	}
}

func StoreCachedLyrics(song structs.SongInfo, data structs.SongLyricsData) error {
	cacheDirectory := global.CurrentConfig.Cache.CacheDir
	if strings.Contains(cacheDirectory, "$XDG_CACHE_DIR") && os.Getenv("$XDG_CACHE_DIR") == "" {
		cacheDirectory = strings.ReplaceAll(cacheDirectory, "$XDG_CACHE_DIR", "$HOME/.cache")
	}

	cacheDirectory = os.ExpandEnv(cacheDirectory)

	if _, err := os.ReadDir(cacheDirectory); err != nil {
		os.Mkdir(cacheDirectory, 0777)
		os.Chmod(cacheDirectory, 0777)
	}

	filename := getFilename(song.Title, song.Artist, song.Album, song.Duration)
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
