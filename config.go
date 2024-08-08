package main

type Config struct {
	ExcludedPlayers                []string `json:"excludedPlayers"`
	IncludedPlayers                []string `json:"includedPlayers"`
	PlayerctlMetadataCheckInterval float64  `json:"playerctlMetadataCheckInterval"`
	PlayerctlPositionCheckInterval float64  `json:"playerctlPositionCheckInterval"`
	MaximumCharacterLength         uint     `json:"maximumCharacterLength"`
	ScrollText                     bool     `json:"scrollText"`
}

func ReadConfig() Config {
	panic("not implemented!!1")
}
