package romanization_test

import (
	"lrcsnc/internal/pkg/global"
	"lrcsnc/internal/romanization"
	"slices"
	"testing"
)

func TestGetLang(t *testing.T) {
	global.Config.Lyrics.Romanization.Japanese = true
	global.Config.Lyrics.Romanization.Chinese = true
	global.Config.Lyrics.Romanization.Korean = true
	answerJapanese := romanization.GetLang(
		[]string{
			"something something Lorem Ipsum",
			"ああ？私に近づいてるの？",
			"кириллица",
		},
	)
	answerKorean := romanization.GetLang(
		[]string{
			"something something Lorem Ipsum",
			"어? 나한테 다가오니?",
			"кириллица",
		},
	)
	answerChinese := romanization.GetLang(
		[]string{
			"something something Lorem Ipsum",
			"哦？你在接近我吗？",
			"кириллица",
		},
	)
	answerDefault := romanization.GetLang(
		[]string{
			"something something Lorem Ipsum",
			"france?!?",
			"кириллица",
		},
	)

	if answerJapanese != romanization.Japanese ||
		answerKorean != romanization.Korean ||
		answerChinese != romanization.Chinese ||
		answerDefault != romanization.Default {
		t.Errorf(
			"[tests/romanization/TestGetLang] ERROR: Expected %v, %v, %v and %v; received %v, %v, %v and %v",
			romanization.Japanese, romanization.Korean, romanization.Chinese, romanization.Default,
			answerJapanese, answerKorean, answerChinese, answerDefault)
	}
}

func TestRomanize(t *testing.T) {
	global.Config.Lyrics.Romanization.Japanese = true
	global.Config.Lyrics.Romanization.Chinese = true
	global.Config.Lyrics.Romanization.Korean = true
	answerJapanese := romanization.Romanize(
		[]string{
			"ああ？私に近づいてるの？",
		},
		romanization.Japanese,
	)
	answerKorean := romanization.Romanize(
		[]string{
			"어? 나한테 다가오니?",
		},
		romanization.Korean,
	)
	answerChinese := romanization.Romanize(
		[]string{
			"哦？你在接近我吗？",
		},
		romanization.Chinese,
	)
	answerDefault := romanization.Romanize(
		[]string{
			"france?!?",
		},
		romanization.Default,
	)

	rightAnswerJapanese := []string{"Aa？私ni近zuiteruno？"}
	rightAnswerKorean := []string{"Eo? nahante dagaoni?"}
	rightAnswerChinese := []string{"O？nizaijiejinwoma？"}
	rightAnswerDefault := []string{"france?!?"}

	if slices.Compare(answerJapanese, rightAnswerJapanese) != 0 ||
		slices.Compare(answerKorean, rightAnswerKorean) != 0 ||
		slices.Compare(answerChinese, rightAnswerChinese) != 0 ||
		slices.Compare(answerDefault, rightAnswerDefault) != 0 {
		t.Errorf(
			"[tests/romanization/TestRomanize] ERROR: Expected \"%v\", \"%v\", \"%v\" and \"%v\"; received \"%v\", \"%v\", \"%v\" and \"%v\"",
			rightAnswerJapanese[0], rightAnswerKorean[0], rightAnswerChinese[0], rightAnswerDefault[0],
			answerJapanese[0], answerKorean[0], answerChinese[0], answerDefault[0])
	}
}
