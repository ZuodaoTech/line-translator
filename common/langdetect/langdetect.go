package langdetect

import (
	"math"
	"strings"

	"github.com/pemistahl/lingua-go"
)

type (
	Detector struct {
		detector lingua.LanguageDetector
	}
	Config struct {
	}
)

func New() *Detector {
	languages := []lingua.Language{
		lingua.English,
		lingua.Chinese,
		lingua.Japanese,
	}

	detector := lingua.NewLanguageDetectorBuilder().
		FromLanguages(languages...).
		WithMinimumRelativeDistance(0.9).
		Build()

	return &Detector{
		detector: detector,
	}
}

func (d *Detector) Detect(text string) (string, bool) {
	language, exists := d.detector.DetectLanguageOf(text)
	if exists {
		lang := d.FormalizeName(language.String(), text)
		return lang, exists
	}
	return "", exists
}

func (d *Detector) FormalizeName(lang, text string) string {
	lang = strings.ToLower(lang)
	switch lang {
	case "chinese":
		return detectZhOrJa(text)
	case "japanese":
		return detectZhOrJa(text)
	case "english":
		return "en"
	default:
		return ""
	}
}

func detectZhOrJa(text string) string {
	// get top 100 runes from string
	allRunes := []rune(text)
	top100Runes := allRunes[:int(math.Min(100, float64(len(allRunes))))]

	// count chinese and japanese runes
	// return the language with more runes
	zhCount := 0
	jaCount := 0
	for _, r := range top100Runes {
		if isJapanese(r) {
			jaCount++
		} else if isChinese(r) {
			zhCount++
		}
	}
	// fmt.Printf("zhCount: %v\n", zhCount)
	// fmt.Printf("jaCount: %v\n", jaCount)
	// if there are more than twice as many Chinese characters as Japanese characters, it is Chinese
	if zhCount > jaCount*5 {
		return "zh"
	}
	return "ja"
}

func isChinese(c rune) bool {
	if (c >= '\u3400' && c <= '\u4db5') || // CJK Unified Ideographs Extension A
		(c >= '\u4e00' && c <= '\u9fed') || // CJK Unified Ideographs
		(c >= '\uf900' && c <= '\ufaff') { // CJK Compatibility Ideographs
		return true
	}

	return false
}

func isJapanese(c rune) bool {
	if (c >= '\u3021' && c <= '\u3029') || // Japanese Hanzi
		(c >= '\u3040' && c <= '\u309f') || // Hiragana
		(c >= '\u30a0' && c <= '\u30ff') || // Katakana
		(c >= '\u31f0' && c <= '\u31ff') || // Katakana Phonetic Extension
		(c >= '\uf900' && c <= '\ufaff') { // CJK Compatibility Ideographs
		return true
	}

	return false
}
