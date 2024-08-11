package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math"
	"os"
	"regexp"
)

var badCharactersRegexp = regexp.MustCompile(`[:;|\/\\<>\.]+`)

func GetCachedLyrics(song *SongData) LrcLibJson {
	cacheDirectory, err := os.UserCacheDir()
	if err != nil {
		log.Println("Could not get cache directory!")
		return LrcLibJson{}
	}

	filename := GetFilename(song.Song, song.Artist, song.Album, song.Duration)
	if file, err := os.ReadFile(cacheDirectory + "/playerctl-lyrics/" + filename + ".json"); err == nil {
		var result LrcLibJson
		err = json.Unmarshal(file, &result)
		if err != nil {
			return LrcLibJson{}
		}
		return result
	} else {
		return LrcLibJson{}
	}
}

func StoreCachedLyrics(song *SongData, lrcData LrcLibJson) error {
	cacheDirectory, err := os.UserCacheDir()
	if err != nil {
		log.Println("Could not get cache directory!")
		return err
	}

	os.Mkdir(cacheDirectory+"/playerctl-lyrics", 0777)

	filename := GetFilename(song.Song, song.Artist, song.Album, song.Duration)
	data, err := json.Marshal(lrcData)
	if err != nil {
		return err
	}
	if err = os.WriteFile(cacheDirectory+"/playerctl-lyrics/"+filename+".json", data, 0660); err != nil {
		return err
	}
	return nil
}

func GetFilename(song string, artist string, album string, duration float64) string {
	return fmt.Sprintf("%v.%v.%v.%v", RemoveBadCharacters(song), RemoveBadCharacters(artist), RemoveBadCharacters(album), math.Round(duration))
}

func RemoveBadCharacters(str string) string {
	return badCharactersRegexp.ReplaceAllString(str, "_")
}
