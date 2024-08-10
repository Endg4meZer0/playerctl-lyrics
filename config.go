package main

import (
	"encoding/json"
	"log"
	"os"
)

type Config struct {
	Playerctl    Playerctl    `json:"playerctl"`
	Output       Output       `json:"output"`
	Instrumental Instrumental `json:"instrumental"`
}

type Playerctl struct {
	ExcludedPlayers            []string `json:"excludedPlayers"`
	IncludedPlayers            []string `json:"includedPlayers"`
	PlayerctlSongCheckInterval float64  `json:"playerctlMetadataCheckInterval"`
}

type Instrumental struct {
	Interval float64 `json:"interval"`
	Symbol   string  `json:"symbol"`
	MaxCount uint    `json:"maxCount"`
}

type Output struct {
	ShowSongNotFoundWarning bool `json:"showSongNotFoundWarning"`
	MaximumCharacterLength  uint `json:"maximumCharacterLength"`
	ScrollText              bool `json:"scrollText"`
}

func ReadConfig(path string) Config {
	jsonConfig, err := os.ReadFile(path)
	if err != nil {
		log.Printf("The config at %v does not exist or is not readable! Falling back to default config!", path)
	}

	var config Config
	if json.Unmarshal(jsonConfig, &config) != nil {
		log.Printf("The config at %v does not exist or is not readable! Falling back to default config!", path)
	}
	return config
}

func ValidateConfig() []error {
	panic("not implemented!!1")
}

func DefaultConfig() Config {
	panic("not implemented!!1")
}
