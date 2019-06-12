package locale

import (
	"strings"

	"github.com/nicksnyder/go-i18n/v2/i18n"
	"go.uber.org/zap"
)

type Localizer struct {
	i18n.Localizer
	lang   string
	logger *zap.SugaredLogger
}

var (
	numbers = map[string][]string{
		"de-DE": []string{"null", "eins", "zwei", "drei", "vier", "fünf", "sechs", "sieben", "acht", "neun", "zehn", "elf", "zwölf", "dreizehn", "vierzehn", "fünfzehn", "sechzehn", "siebzehn", "achtzehn", "neunzehn", "zwanzig"},
		"en-US": []string{"zero", "one", "two", "three", "four", "five", "six", "seven", "eight", "nine", "ten", "eleven", "twelf", "thirteen", "fourteen", "fifteen", "sixteen", "seventeen", "eighteen", "nineteen", "twenty"},
		"en-GB": []string{"zero", "one", "two", "three", "four", "five", "six", "seven", "eight", "nine", "ten", "eleven", "twelf", "thirteen", "fourteen", "fifteen", "sixteen", "seventeen", "eighteen", "nineteen", "twenty"},
		"en-IN": []string{"zero", "one", "two", "three", "four", "five", "six", "seven", "eight", "nine", "ten", "eleven", "twelf", "thirteen", "fourteen", "fifteen", "sixteen", "seventeen", "eighteen", "nineteen", "twenty"},
		"en-AU": []string{"zero", "one", "two", "three", "four", "five", "six", "seven", "eight", "nine", "ten", "eleven", "twelf", "thirteen", "fourteen", "fifteen", "sixteen", "seventeen", "eighteen", "nineteen", "twenty"},
		"en-CA": []string{"zero", "one", "two", "three", "four", "five", "six", "seven", "eight", "nine", "ten", "eleven", "twelf", "thirteen", "fourteen", "fifteen", "sixteen", "seventeen", "eighteen", "nineteen", "twenty"},
	}
	endpoints = map[string]string{
		"de-DE": "de.wikipedia.org",
		"en-US": "en.wikipedia.org",
		"en-GB": "en.wikipedia.org",
		"en-IN": "en.wikipedia.org",
		"en-AU": "en.wikipedia.org",
		"en-CA": "en.wikipedia.org",
	}
)

func NewLocalizer(bundle *i18n.Bundle, lang string, logger *zap.SugaredLogger) *Localizer {
	if _, langExist := numbers[lang]; !langExist {
		panic("language '" + lang + "' not supported for number spelling")
	}
	if _, langExist := endpoints[lang]; !langExist {
		panic("language '" + lang + "' not supported for wikipedia endpoint")
	}
	return &Localizer{
		Localizer: *i18n.NewLocalizer(bundle, lang),
		lang:      lang,
		logger:    logger,
	}
}

func (l *Localizer) Spell(number int) string {
	if number >= len(numbers[l.lang]) {
		l.logger.Errorw("Tried to spell a too big number", "number", number, "language", l.lang)
		return "whatever"
	}
	return numbers[l.lang][number]
}

func (l *Localizer) WikiEndpoint() string {
	return endpoints[l.lang]
}

func (l *Localizer) AssembleTermFromSpelling(spelledTerm string) string {
	if l.lang == "de-DE" {
		return strings.Join(strings.Split(strings.ReplaceAll(spelledTerm, "leerzeichen", " "), ". "), "")
	}
	return strings.ToLower(strings.ReplaceAll(spelledTerm, "space", " "))
}
