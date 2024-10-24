package romanization

import (
	"strings"
	"unicode"

	"lrcsnc/internal/pkg/global"

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

	if global.CurrentConfig.Lyrics.Romanization.Japanese {
		out = jp.ToRomaji(str, true)
		if out != strings.ToLower(str) {
			// Kanji and zh/kr characters are coded using 3 bytes.
			// So if a character did not get romanized, this should block the uppercasing (and failing in the process)
			if !isChar(out[:3], supportedAsianLangsUnicodeRangeTable) {
				out = strings.ToUpper(out[:1]) + out[1:]
			}
			return
		} else {
			out = str
		}
	}

	if global.CurrentConfig.Lyrics.Romanization.Chinese {
		out = zhCharToPinyin(str)
		if out != str {
			if !isChar(out[:3], supportedAsianLangsUnicodeRangeTable) {
				out = strings.ToUpper(out[:1]) + out[1:]
			}
			return
		} else {
			out = str
		}
	}

	if global.CurrentConfig.Lyrics.Romanization.Korean {
		r := kr.NewRomanizer(str)
		out = r.Romanize()
		if out != str {
			if !isChar(out[:3], supportedAsianLangsUnicodeRangeTable) {
				out = strings.ToUpper(out[:1]) + out[1:]
			}
			return
		} else {
			out = str
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
