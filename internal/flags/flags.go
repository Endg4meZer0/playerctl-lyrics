package flags

import (
	"flag"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"

	"lrcsnc/internal/config"
	"lrcsnc/internal/output/piped"
	"lrcsnc/internal/pkg/global"
	"lrcsnc/internal/pkg/util"
)

var helpText = []string{
	"Usage:",
	"  lrcsnc [FLAGS...]",
	"",
	"Options:",
}

var (
	configPath       = flag.String("config", "", "Sets the config file to use")
	cacheDirectory   = flag.String("cache-dir", "", "Sets the cache directory")
	isPiped          = flag.Bool("piped", false, "Pipe the output to stdout or file instead of using TUI, ignoring the set value in config.")
	outputFilePath   = flag.String("output", "", "Sets an output file to use instead of standard output when using -piped")
	clearCacheMode   = flag.Bool("clear-cache", false, "If true, searches the cache directory, removes cache files that fit the filters (-song-name, -song-artist, etc.) and exits. Only songs that contain the set patterns will be affected")
	songNameFilter   = flag.String("song-name-filter", "", "Sets the song name filter to use when -clear-cache is also set")
	artistNameFilter = flag.String("artist-name-filter", "", "Sets the artist name filter to use when -clear-cache is also set")
	albumNameFilter  = flag.String("album-name-filter", "", "Sets the album name filter to use when -clear-cache is also set")
	durationFilter   = flag.Int("duration-filter", 0, "Sets the duration filter to use when -clear-cache is also set")
	displayHelp      = flag.Bool("help", false, "Display the help message and exit")
	displayVersion   = flag.Bool("version", false, "Display the version")
)

func HandleFlags() {
	flag.Parse()

	if *displayVersion {
		fmt.Println("v0.3.0")
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
		if err = config.ReadConfig(*configPath); err != nil {
			log.Println("The set config path is not valid! Falling back to the config from default path...\nErrors: ", err.Error())
		}
	}
	if *configPath == "" || err != nil {
		if err := config.ReadConfigFromDefaultPath(); err != nil {
			log.Println("The config from default path is no valid! Falling back to the default config...\nErrors: ", err.Error())
		}
	}

	if *cacheDirectory != "" {
		if _, err := os.ReadDir(os.ExpandEnv(*cacheDirectory)); err != nil {
			os.MkdirAll(*cacheDirectory, 0777)
			os.Chmod(*cacheDirectory, 0777)
		}
		global.Config.Cache.CacheDir = *cacheDirectory
	}

	if *isPiped {
		global.Config.Global.Output = "piped"
	}

	if *outputFilePath != "" && global.Config.Global.Output == "piped" {
		t := os.ExpandEnv(*outputFilePath)
		piped.UpdateOutputDestination(t)
	}

	if *clearCacheMode {
		if *songNameFilter == "" && *artistNameFilter == "" && *albumNameFilter == "" && *durationFilter == 0 {
			log.Fatalln("The -clear-cache flag is set, but no filters are! Check -help for more information.")
		}

		currentCacheDir := global.Config.Cache.CacheDir

		if strings.Contains(currentCacheDir, "$XDG_CACHE_DIR") && os.Getenv("$XDG_CACHE_DIR") == "" {
			currentCacheDir = strings.ReplaceAll(currentCacheDir, "$XDG_CACHE_DIR", "$HOME/.cache")
		}

		cacheFiles, err := os.ReadDir(os.ExpandEnv(currentCacheDir))

		if err != nil {
			log.Println("Something is wrong with the cache directory. Not created yet? Try and launch the main process.")
		}

		deletedFiles := 0
		for _, cacheFile := range cacheFiles {
			sections := strings.Split(cacheFile.Name(), ".")
			if *songNameFilter != "" {
				if found, _ := regexp.MatchString(util.RemoveBadCharacters(*songNameFilter), sections[0]); !found {
					continue
				}
			}
			if *artistNameFilter != "" {
				if found, _ := regexp.MatchString(util.RemoveBadCharacters(*artistNameFilter), sections[1]); !found {
					continue
				}
			}
			if *albumNameFilter != "" {
				if found, _ := regexp.MatchString(util.RemoveBadCharacters(*albumNameFilter), sections[2]); !found {
					continue
				}
			}
			if *durationFilter != 0 {
				if found, _ := regexp.MatchString(fmt.Sprint(*durationFilter), sections[3]); !found {
					continue
				}
			}

			if os.Remove(os.ExpandEnv(currentCacheDir)+"/"+cacheFile.Name()) != nil {
				fmt.Printf("Couldn't delete file %v! Missing permissions?", cacheFile.Name())
				continue
			}
			deletedFiles++
		}

		fmt.Printf("Deleted %v files, exiting...", deletedFiles)

		os.Exit(0)
	}
}
