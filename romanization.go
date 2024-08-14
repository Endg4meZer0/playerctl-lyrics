package main

import (
	"strings"
	"unicode"

	jp "github.com/mochi-co/kana-tools"
	zh "github.com/mozillazg/go-pinyin"
	kr "github.com/srevinsaju/korean-romanizer-go"
)

var supportedAsianLangsUnicodeRangeTable = []*unicode.RangeTable{
	unicode.Ideographic,
	unicode.Hiragana,
	unicode.Katakana,
	unicode.Diacritic,
	unicode.Han,
	unicode.Hangul,
}

func IsSupportedAsianLang(str string) bool {
	return isChar(str, supportedAsianLangsUnicodeRangeTable)
}

func Romanize(str string) string {
	jpRomanized := jp.ToRomaji(str, true)
	if CurrentConfig.Output.Romanization.Japanese && jpRomanized != strings.ToLower(str) {
		return jpRomanized
	}
	zhRomanized := zhCharToPinyin(str)
	if CurrentConfig.Output.Romanization.Chinese && zhRomanized != str {
		return zhRomanized
	}
	r := kr.NewRomanizer(str)
	krRomanized := r.Romanize()
	if CurrentConfig.Output.Romanization.Korean && krRomanized != str {
		return krRomanized
	}

	return str
}

func isChar(s string, rangeTable []*unicode.RangeTable) bool {
	for _, r := range s {
		if unicode.IsOneOf(rangeTable, r) {
			return true
		}
	}
	return false
}

func zhCharToPinyin(p string) (s string) {
	for _, r := range p {
		if unicode.Is(unicode.Han, r) {
			s += string(zh.Pinyin(string(r), zh.NewArgs())[0][0])
		} else {
			s += string(r)
		}
	}
	return
}
