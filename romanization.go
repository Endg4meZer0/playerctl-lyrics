package main

import (
	"strings"
	"unicode"

	jp "github.com/mochi-co/kana-tools"
	zh "github.com/mozillazg/go-pinyin"
	kr "github.com/srevinsaju/korean-romanizer-go"
)

var supportedAsianLangsUnicodeRangeTable = []*unicode.RangeTable{
	unicode.Ideographic, // jp kanji and some zh characters
	unicode.Hiragana,    // jp
	unicode.Katakana,    // jp
	unicode.Diacritic,   // jp (?)
	unicode.Han,         // zh
	unicode.Hangul,      // kr
}

func IsSupportedAsianLang(str string) bool {
	return isChar(str, supportedAsianLangsUnicodeRangeTable)
}

func Romanize(str string) (out string) {
	out = str

	// Sometimes there are japanese lyrics that consist only of kanji (unromanazible without a dictionary) characters
	// They fall down to chinese romanization and that sometimes causes trouble.
	// Until I know how to fix it properly, there will be a recover function
	defer func() {
		if r := recover(); r != nil {
			out = str
		}
	}()

	if CurrentConfig.Output.Romanization.Japanese {
		out = jp.ToRomaji(str, true)
		if out != strings.ToLower(str) {
			out = strings.ToUpper(out[:1]) + out[1:]
			return
		}
	}

	if CurrentConfig.Output.Romanization.Chinese {
		out = zhCharToPinyin(str)
		if out != str {
			out = strings.ToUpper(out[:1]) + out[1:]
			return
		}
	}

	if CurrentConfig.Output.Romanization.Korean {
		r := kr.NewRomanizer(str)
		out = r.Romanize()
		if out != str {
			out = strings.ToUpper(out[:1]) + out[1:]
			return
		}
	}

	return
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
