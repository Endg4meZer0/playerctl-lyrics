package cache_test

import (
	"lrcsnc/internal/cache"
	"lrcsnc/internal/pkg/global"
	"lrcsnc/internal/pkg/structs"
	"testing"
)

func TestStoreGetCycle(t *testing.T) {
	global.CurrentConfig.Cache.CacheDir = "$XDG_CACHE_DIR/lrcsnc"
	testSong := structs.Song{
		Title:    "Is This A Test?",
		Artist:   "Endg4me_",
		Album:    "lrcsnc",
		Duration: 12.12,
		LyricsData: structs.LyricsData{
			Lyrics: []string{
				"Pam-pam-pampararam",
				"Pam-pam-pam-param-pamparam",
			},
			LyricTimestamps: []float64{
				4.12,
				7.54,
			},
			LyricsType: 0,
		},
	}
	err := cache.StoreCachedLyrics(testSong)
	if err != nil {
		t.Errorf("[tests/cache/TestStoreGetCycle] %v", err)
	}

	global.CurrentConfig.Cache.Enabled = false
	answerDisabled, cacheStateDisabled := cache.GetCachedLyrics(testSong)

	global.CurrentConfig.Cache.Enabled = true
	answerInfLifeSpan, cacheStateInfLifeSpan := cache.GetCachedLyrics(testSong)

	if len(answerDisabled.Lyrics) != 0 || answerDisabled.LyricsType != 0 || len(answerDisabled.LyricTimestamps) != 0 || cacheStateDisabled != cache.CacheStateDisabled {
		t.Errorf("[tests/cache/TestStoreGetCycle] ERROR: Disabled caching in config doesn't stop getter from getting cached data")
	}

	if len(answerInfLifeSpan.Lyrics) != 2 || answerInfLifeSpan.LyricsType != 0 || len(answerInfLifeSpan.LyricTimestamps) != 2 || cacheStateInfLifeSpan != cache.CacheStateActive {
		t.Errorf("[tests/cache/TestStoreGetCycle] ERROR: Received wrong cached data: expected [%v]string, [%v]float64, %v and %v, received [%v]string, [%v]float64, %v and %v",
			len(testSong.LyricsData.Lyrics), len(testSong.LyricsData.LyricTimestamps), testSong.LyricsData.LyricsType, cache.CacheStateActive,
			len(answerInfLifeSpan.Lyrics), len(answerInfLifeSpan.LyricTimestamps), answerInfLifeSpan.LyricsType, cacheStateInfLifeSpan,
		)
	}

	cache.RemoveCachedLyrics(testSong)
}
