package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"regexp"
	"strings"
	"syscall"
)

func main() {

	// TODO: An actual config implementation
	// TODO: A better options implementation
	// there's definitely always more!

	// Check if playerctl is installed
	err := exec.Command("playerctl", "--version").Run()
	if err != nil {
		log.Fatalln("playerctl is not found!")
	}

	HandleFlags()

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	SyncLoop()
	<-sigs
	os.Exit(0)
}

// Handling flags

var helpText = []string{
	"Usage:",
	"  playerctl-lyrics [FLAGS...]",
	"  If certain flags are not set, the command will start the main process that gets lyrics, syncs them with playerctl and prints them to stdout.",
	"",
	"Flags:",
}

func HandleFlags() {
	configPath := flag.String("config", "", "Sets the config file to use")
	clearCacheMode := flag.Bool("clear-cache", false, "If true, searches the cache directory, removes cache files that fit the filters (-song-name, -song-artist, etc.) and exits. Only songs that contain the set patterns will be affected")
	cacheDirectory := flag.String("cache-dir", DefaultConfig().Cache.CacheDir, "Sets the cache directory")
	songNameFilter := flag.String("song-name-filter", "", "Sets the song name filter to use when -clear-cache is also set")
	artistNameFilter := flag.String("artist-name-filter", "", "Sets the artist name filter to use when -clear-cache is also set")
	albumNameFilter := flag.String("album-name-filter", "", "Sets the album name filter to use when -clear-cache is also set")
	durationFilter := flag.Int("duration-filter", 0, "Sets the duration filter to use when -clear-cache is also set")
	displayHelp := flag.Bool("help", false, "Display the help message and exit")
	displayVersion := flag.Bool("version", false, "Display the version")
	flag.Parse()

	if *displayVersion {
		fmt.Println("v0.0.1-beta")
		os.Exit(0)
	}

	if *displayHelp {
		for _, s := range helpText {
			fmt.Println(s)
		}
		flag.PrintDefaults()
		os.Exit(0)
	}

	var err error

	if *configPath != "" {
		if err = ReadConfig(*configPath); err != nil {
			fmt.Println("The set config path is not valid! Falling back to the config from default path...\nErrors: ", err.Error())
		}
	}
	if *configPath == "" || err != nil {
		if err := ReadConfigFromDefaultPath(); err != nil {
			fmt.Println("The config from default path is no valid! Falling back to the default config...\nErrors: ", err.Error())
			CurrentConfig = DefaultConfig()
		}
	}

	if *cacheDirectory != DefaultConfig().Cache.CacheDir {
		if _, err := os.ReadDir(*cacheDirectory); err != nil {
			os.MkdirAll(*cacheDirectory, 0664)
		}
		CurrentConfig.Cache.CacheDir = *cacheDirectory
	}

	if *clearCacheMode {
		if *songNameFilter == "" && *artistNameFilter == "" && *albumNameFilter == "" && *durationFilter == 0 {
			fmt.Println("The -clear-cache flag is set, but no filters are! Check -help for more information.")
			os.Exit(1)
		}

		cacheFiles, err := os.ReadDir(os.ExpandEnv(CurrentConfig.Cache.CacheDir))
		deletedFiles := 0

		if err != nil {
			fmt.Println("Something is wrong with the cache directory. Not created yet? Try and launch the main process.")
		}

		for _, cacheFile := range cacheFiles {
			sections := strings.Split(cacheFile.Name(), ".")
			if *songNameFilter != "" {
				if found, _ := regexp.MatchString(RemoveBadCharacters(*songNameFilter), sections[0]); !found {
					continue
				}
			}
			if *artistNameFilter != "" {
				if found, _ := regexp.MatchString(RemoveBadCharacters(*artistNameFilter), sections[1]); !found {
					continue
				}
			}
			if *albumNameFilter != "" {
				if found, _ := regexp.MatchString(RemoveBadCharacters(*albumNameFilter), sections[2]); !found {
					continue
				}
			}
			if *durationFilter != 0 {
				if found, _ := regexp.MatchString(string(*durationFilter), sections[3]); !found {
					continue
				}
			}

			if os.Remove(os.ExpandEnv(CurrentConfig.Cache.CacheDir)+"/"+cacheFile.Name()) != nil {
				fmt.Printf("Couldn't delete file %v! Missing permissions?", cacheFile.Name())
				continue
			}
			deletedFiles++
		}

		os.Exit(0)
	}
}
