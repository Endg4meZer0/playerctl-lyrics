package lrclib

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
				LyricsType: 2,
			},
		},
		{
			name: "Plain Lyrics",
			dto: lrclib.LrcLibDTO{
				PlainLyrics: "Line 1\nLine 2\nLine 3",
			},
			expected: structs.LyricsData{
				LyricsType: 1,
				Lyrics:     []string{"Line 1", "Line 2", "Line 3"},
			},
		},
		{
			name: "Synced Lyrics",
			dto: lrclib.LrcLibDTO{
				SyncedLyrics: "[00:01.00] Line 1\n[00:02.00] Line 2\n[00:03.00] Line 3",
			},
			expected: structs.LyricsData{
				LyricsType:      0,
				Lyrics:          []string{"Line 1", "Line 2", "Line 3"},
				LyricTimestamps: []float64{1.0, 2.0, 3.0},
			},
		},
		{
			name: "Empty Lyrics",
			dto:  lrclib.LrcLibDTO{},
			expected: structs.LyricsData{
				LyricsType: 0,
				Lyrics:     []string{""},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.dto.ToLyricsData()
			if !reflect.DeepEqual(got, tt.expected) {
				t.Errorf("ToLyricsData() = %v, want %v", got, tt.expected)
			}
		})
	}
}
