package main

import (
	"fmt"
	"log"
	"os"
	"strings"
)

var helpText = []string{
	"Usage:",
	"  playerctl-lyrics [OPTION]",
	"",
	"Options:",
	"  -h, --help:													print this message and exit",
	"  -v, --version: 												print version and exit",
	"  --clear-cache SONGNAME ARTISTNAME [ALBUMNAME] [DURATION]:	clear cache files of the matching songs and exit",
	"  --clear-cache-dir:											clear cache directory and exit",
}

// The option handling will 100% need to be rewritten later on. Sure of it.
func HandleOptions(args []string) {
	for _, arg := range args {
		switch arg {
		case "--clear-cache":
			cacheDirectory, err := os.UserCacheDir()
			if err != nil {
				log.Fatalln("Couldn't find the user cache directory!")
			}
			if len(args) < 3 {
				log.Fatalln("Not enough arguments! Check --help for more information.")
			}

			entries, err := os.ReadDir(cacheDirectory + "/playerctl-lyrics")
			if err != nil {
				log.Fatalln("Couldn't read the user cache directory! Are there reading permissions?")
			}

			doAlbumCheck := len(args) >= 4
			doDurationCheck := len(args) == 5

			deletedCount := 0

			for _, entry := range entries {
				if entry.IsDir() {
					continue
				}

				nameParts := strings.Split(entry.Name(), ".")

				songMatched := RemoveBadCharacters(args[1]) == nameParts[0]
				if !songMatched {
					continue
				}

				artistMatched := RemoveBadCharacters(args[2]) == nameParts[1]
				if !artistMatched {
					continue
				}

				if doAlbumCheck {
					albumMatched := RemoveBadCharacters(args[3]) == nameParts[2]
					if !albumMatched {
						continue
					}
				}

				if doDurationCheck {
					durationMatched := RemoveBadCharacters(args[4]) == nameParts[3]
					if !durationMatched {
						continue
					}
				}

				os.Remove(cacheDirectory + "/playerctl-lyrics/" + entry.Name())
				deletedCount++
			}
			log.Printf("Deleted %v files. Exiting...", deletedCount)
			os.Exit(0)
		case "--clear-cache-dir":
			cacheDirectory, err := os.UserCacheDir()
			if err != nil {
				log.Fatalln("Couldn't find the user cache directory!")
			}
			os.RemoveAll(cacheDirectory + "/playerctl-lyrics")
			os.Mkdir(cacheDirectory+"/playerctl-lyrics", 0777)
			os.Exit(0)
		case "-h", "--help":
			for _, s := range helpText {
				fmt.Println(s)
			}
			os.Exit(0)
		case "-v", "--version":
			fmt.Println("v0.0.1-alpha")
			os.Exit(0)
		default:
			log.Fatalln("Unknown option: ", arg)
			os.Exit(0)
		}
	}
}
