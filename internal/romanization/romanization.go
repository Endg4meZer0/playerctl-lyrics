package romanization

import (
	"strings"
	"unicode"

	"lrcsnc/internal/pkg/global"

	jp "github.com/mochi-co/kana-tools"
	zh "github.com/mozillazg/go-pinyin"
	kr "github.com/srevinsaju/korean-romanizer-go"
)

type Language uint

var (
	Default  Language = 0
	Japanese Language = 1
	Korean   Language = 2
	Chinese  Language = 3
)

var jpUnicodeRangeTable = []*unicode.RangeTable{
	unicode.Hiragana,
	unicode.Katakana,
	unicode.Diacritic,
}

var zhUnicodeRangeTable = []*unicode.RangeTable{
	unicode.Ideographic,
	unicode.Han,
}

var krUnicodeRangeTable = []*unicode.RangeTable{
	unicode.Hangul,
}

func GetLang(lyrics []string) Language {
	if global.CurrentConfig.Lyrics.Romanization.Japanese {
		for _, l := range lyrics {
			if isChar(l, jpUnicodeRangeTable) {
				return 1
			}
		}
	}
	if global.CurrentConfig.Lyrics.Romanization.Korean {
		for _, l := range lyrics {
			if isChar(l, krUnicodeRangeTable) {
				return 2
			}
		}
	}
	if global.CurrentConfig.Lyrics.Romanization.Chinese {
		for _, l := range lyrics {
			if isChar(l, zhUnicodeRangeTable) {
				return 3
			}
		}
	}
	return 0
}

func Romanize(strs []string, lang Language) []string {
	outs := make([]string, 0, len(strs))
	for _, str := range strs {
		switch lang {
		case 1:
			out := jp.ToRomaji(str, true)
			// Kanji and zh/kr characters are coded using unicode.Ideographic/unicode.Hangul using 3 bytes.
			// So if a character did not get romanized, this should block the uppercasing (and crashing in the process)
			if !isChar(out[:3], zhUnicodeRangeTable) {
				out = strings.ToUpper(out[:1]) + out[1:]
			}
			outs = append(outs, out)
		case 2:
			r := kr.NewRomanizer(str)
			out := r.Romanize()
			if !isChar(out[:3], krUnicodeRangeTable) {
				out = strings.ToUpper(out[:1]) + out[1:]
			}
			outs = append(outs, out)
		case 3:
			out := zhCharToPinyin(str)
			if !isChar(out[:3], zhUnicodeRangeTable) {
				out = strings.ToUpper(out[:1]) + out[1:]
			}
			outs = append(outs, out)
		default:
			outs = append(outs, str)
		}

	}

	return outs
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
