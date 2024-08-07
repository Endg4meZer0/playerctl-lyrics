package main

type Config struct {
	ExcludedPlayers                []string `json:"excludedPlayers"`
	IncludedPlayers                []string `json:"includedPlayers"`
	PlayerctlMetadataCheckInterval float64  `json:"playerctlMetadataCheckInterval"`
	MaximumCharacterLength         uint     `json:"maximumCharacterLength"`
}

func ReadConfig() Config {
	panic("not implemented!!1")
}
