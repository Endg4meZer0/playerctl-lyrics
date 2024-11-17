package dto_test

import (
	"lrcsnc/internal/lyrics/dto"
	"lrcsnc/internal/lyrics/dto/lrclib"
	"lrcsnc/internal/pkg/structs"
	"slices"
	"testing"
)

func TestRemoveMismatches(t *testing.T) {
	tests := []struct {
		name       string
		song       structs.Song
		lyricsData []dto.LyricsDTO
		expected   []dto.LyricsDTO
	}{
		{
			name:       "Empty lyrics data",
			song:       structs.Song{Title: "Test Song", Duration: 300},
			lyricsData: []dto.LyricsDTO{},
			expected:   []dto.LyricsDTO{},
		},
		{
			name: "Matching title and duration",
			song: structs.Song{Title: "Test Song", Duration: 300},
			lyricsData: []dto.LyricsDTO{
				lrclib.LrcLibDTO{Title: "Test Song", Duration: 300},
				lrclib.LrcLibDTO{Title: "Test Song", Duration: 302},
			},
			expected: []dto.LyricsDTO{
				lrclib.LrcLibDTO{Title: "Test Song", Duration: 300},
				lrclib.LrcLibDTO{Title: "Test Song", Duration: 302},
			},
		},
		{
			name: "Non-matching title",
			song: structs.Song{Title: "Test Song", Duration: 300},
			lyricsData: []dto.LyricsDTO{
				lrclib.LrcLibDTO{Title: "Another Song", Duration: 300},
			},
			expected: []dto.LyricsDTO{},
		},
		{
			name: "Non-matching duration",
			song: structs.Song{Title: "Test Song", Duration: 300},
			lyricsData: []dto.LyricsDTO{
				lrclib.LrcLibDTO{Title: "Test Song", Duration: 305},
			},
			expected: []dto.LyricsDTO{},
		},
		{
			name: "Song duration is zero",
			song: structs.Song{Title: "Test Song", Duration: 0},
			lyricsData: []dto.LyricsDTO{
				lrclib.LrcLibDTO{Title: "Test Song", Duration: 300},
				lrclib.LrcLibDTO{Title: "Test Song", Duration: 305},
			},
			expected: []dto.LyricsDTO{
				lrclib.LrcLibDTO{Title: "Test Song", Duration: 300},
				lrclib.LrcLibDTO{Title: "Test Song", Duration: 305},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := dto.RemoveMismatches(tt.song, tt.lyricsData)
			if !slices.Equal(result, tt.expected) {
				t.Errorf("[tests/lyrics/dto/mismatches/%v] Returned %v, want %v", tt.name, result, tt.expected)
			}
		})
	}
}
