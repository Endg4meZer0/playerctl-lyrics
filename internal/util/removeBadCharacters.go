package util

import "regexp"

func RemoveBadCharacters(str string) string {
	return regexp.MustCompile(`[:;!?,\.\[\]<>\/\\*|]+`).ReplaceAllString(str, "_")
}
