package lrclib_test

import (
	"reflect"
	"testing"

	"lrcsnc/internal/lyrics/dto/lrclib"
	"lrcsnc/internal/pkg/structs"
)

func TestToLyricsData(t *testing.T) {
	tests := []struct {
		name     string
		dto      lrclib.LrcLibDTO
		expected structs.LyricsData
	}{
		{
			name: "Instrumental",
			dto: lrclib.LrcLibDTO{
				Instrumental: true,
			},
			expected: structs.LyricsData{
				LyricsType: structs.LyricsStateInstrumental,
			},
		},
		{
			name: "Plain Lyrics",
			dto: lrclib.LrcLibDTO{
				PlainLyrics: "Line 1\nLine 2\nLine 3",
			},
			expected: structs.LyricsData{
				LyricsType: structs.LyricsStatePlain,
				Lyrics:     []string{"Line 1", "Line 2", "Line 3"},
			},
		},
		{
			name: "Synced Lyrics",
			dto: lrclib.LrcLibDTO{
				PlainLyrics:  "Line 1\nLine 2\nLine 3",
				SyncedLyrics: "[00:01.00] Line 1\n[00:02.00] Line 2\n[00:03.00] Line 3",
			},
			expected: structs.LyricsData{
				LyricsType:      structs.LyricsStateSynced,
				Lyrics:          []string{"Line 1", "Line 2", "Line 3"},
				LyricTimestamps: []float64{1.0, 2.0, 3.0},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.dto.ToLyricsData()
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("[tests/lyrics/dto/lrclib/dto/%v] Returned %v, want %v", tt.name, result, tt.expected)
			}
		})
	}
}
