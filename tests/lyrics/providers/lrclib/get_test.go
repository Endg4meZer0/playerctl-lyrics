package lrclib

import (
	"slices"
	"testing"

	dto "lrcsnc/internal/lyrics/dto"
	lrclibdto "lrcsnc/internal/lyrics/dto/lrclib"
	provider "lrcsnc/internal/lyrics/providers/lrclib"
	"lrcsnc/internal/pkg/structs"
)

type Response struct {
	StatusCode int
	Body       string
}

// The most simple tests for LrcLib provider
func TestGetLyricsDTOList(t *testing.T) {
	tests := []struct {
		name    string
		song    structs.Song
		dtoList []dto.LyricsDTO
	}{
		{
			name: "Existing song",
			song: structs.Song{Title: "Earthless", Artist: "Night Verses", Album: "From the Gallery of Sleep", Duration: 383},
			dtoList: []dto.LyricsDTO{
				lrclibdto.LrcLibDTO{Title: "Earthless", Artist: "Night Verses", Album: "From the Gallery of Sleep", Duration: 382, Instrumental: false, PlainLyrics: "\"He is the one who gave me the horse\nSo I could ride into the desert and see\nThe future.\"\n\n\"He is the one who gave me the horse\nSo I could ride into the desert and see\nThe future.\"\n", SyncedLyrics: "[05:44.18] \"He is the one who gave me the horse\n[05:46.74] So I could ride into the desert and see\n[05:50.77] The future.\"\n[05:51.41] \n[05:58.75] \"He is the one who gave me the horse\n[06:01.39] So I could ride into the desert and see\n[06:05.29] The future.\"\n[06:06.07] "},
			},
		},
		{
			name:    "Not-existing song",
			song:    structs.Song{Title: "Moonmore", Artist: "Day Choruses", Album: "From the Gallery of Minecraft Pictures idk", Duration: 283},
			dtoList: []dto.LyricsDTO{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := provider.LrcLibLyricsProvider{}.GetLyricsDTOList(tt.song)
			if err != nil {
				t.Errorf("[tests/lyrics/providers/lrclib/get/%v] Error: %v", tt.name, err)
				return
			}
			if !slices.Equal(got, tt.dtoList) {
				t.Errorf("[tests/lyrics/providers/lrclib/get/%v] Received %v, want %v", tt.name, got, tt.dtoList)
			}
		})
	}
}
