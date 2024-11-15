package lyrics

import (
	"fmt"
	"lrcsnc/internal/cache"
	"lrcsnc/internal/lyrics/providers"
	"lrcsnc/internal/pkg/global"
	"lrcsnc/internal/pkg/structs"
)

func GetLyricsData(song structs.Song) (structs.LyricsData, error) {
	if song.Duration == 0 {
		return structs.LyricsData{LyricsType: structs.LyricsStateNotFound}, fmt.Errorf("[lyrics/providers/get] WARNING: Song duration is 0, cannot get lyrics")
	}

	if global.CurrentConfig.Cache.Enabled {
		cachedData, cacheState := cache.GetCachedLyrics(song)
		if cacheState == cache.CacheStateActive {
			return cachedData, nil
		}
	}

	dtoList, err := providers.LyricsDataProviders[global.CurrentConfig.Global.LyricsProvider].GetLyricsDTOList(song)
	if err != nil && err.Error() != "Song is not found" {
		// TODO: logger :)
		return structs.LyricsData{LyricsType: structs.LyricsStateUnknown}, fmt.Errorf("[lyrics/providers/get] WARNING: Error while getting lyrics data: %v", err)
	}

	if len(dtoList) == 0 {
		// TODO: logger :)
		return structs.LyricsData{LyricsType: structs.LyricsStateNotFound}, fmt.Errorf("[lyrics/providers/get] WARNING: Couldn't find any matching songs")
	}

	res := dtoList[0].ToLyricsData()

	if global.CurrentConfig.Cache.Enabled && res.LyricsType != 1 {
		defer func() {
			song.LyricsData = res
			if cache.StoreCachedLyrics(song) != nil {
				// TODO: logger :)
				// log.Println("Could not save the lyrics to the cache! Is there an issue with perms?")
			}
		}()
	}

	return res, nil
}
