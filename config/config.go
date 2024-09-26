package config

import (
	"encoding/json"
	"log"
	"os"

	"lrcsnc/internal/output"
	"lrcsnc/pkg/global"
	"lrcsnc/pkg/structs"
)

var currentConfigPath string

func ReadConfig(path string) error {
	configFile, err := os.ReadFile(os.ExpandEnv(path))
	if err != nil {
		return err
	}

	var config structs.Config

	if err := json.Unmarshal(configFile, &config); err != nil {
		return err
	}

	errs, fatal := ValidateConfig(&config)

	for _, v := range errs {
		log.Println(v)
	}

	if !fatal {
		global.CurrentConfig = config
		currentConfigPath = path
	} else {
		return ConfigError("FATAL ERRORS IN THE CONFIG WERE DETECTED! Rolling back...")
	}

	return nil
}

func ReadConfigFromDefaultPath() error {
	global.CurrentConfig = defaultConfig

	defaultDirectory, err := os.UserConfigDir()
	if err != nil {
		return err
	}
	defaultDirectory += "/lrcsnc"

	if _, err := os.ReadDir(defaultDirectory); err != nil {
		os.Mkdir(defaultDirectory, 0777)
		os.Chmod(defaultDirectory, 0777)
	}

	if _, err := os.Lstat(defaultDirectory + "/config.json"); err != nil {
		defaultConfigJSON, err := json.MarshalIndent(defaultConfig, "", "    ")
		if err != nil {
			return err
		}
		err = os.WriteFile(defaultDirectory+"/config.json", defaultConfigJSON, 0777)
		if err != nil {
			return err
		}
	} else {
		configFile, err := os.ReadFile(defaultDirectory + "/config.json")
		if err != nil {
			return err
		}

		var config structs.Config

		if err := json.Unmarshal(configFile, &config); err != nil {
			return err
		}

		errs, fatal := ValidateConfig(&config)

		for _, v := range errs {
			log.Println(v)
		}

		if !fatal {
			global.CurrentConfig = config
		} else {
			return ConfigError("FATAL ERRORS IN THE CONFIG WERE DETECTED! Rolling back...")
		}
	}

	currentConfigPath = defaultDirectory + "/config.json"

	return nil
}

func UpdateConfig() {
	configFile, err := os.ReadFile(os.ExpandEnv(currentConfigPath))
	if err != nil {
		output.PrintOverwrite("Errors while reading config! Falling back...")
		return
	}

	var config structs.Config

	if err := json.Unmarshal(configFile, &config); err != nil {
		output.PrintOverwrite("Errors while parsing config! Falling back...")
		return
	}

	errs, fatal := ValidateConfig(&config)

	for _, v := range errs {
		log.Println(v)
	}

	if !fatal {
		global.CurrentConfig = config
	} else {
		output.PrintOverwrite("Errors while parsing config! Falling back...")
		return
	}
}
