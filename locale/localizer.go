package locale

import "github.com/nicksnyder/go-i18n/v2/i18n"

type Localizer struct {
	i18n.Localizer
	lang string
}

func NewLocalizer(bundle *i18n.Bundle, lang string) *Localizer {
	return &Localizer{
		Localizer: *i18n.NewLocalizer(bundle, lang),
		lang:      lang,
	}
}

var numbers = []string{"null", "eins", "zwei", "drei", "vier", "fünf", "sechs", "sieben", "acht", "neun", "zehn", "elf", "zwölf", "dreizehn", "vierzehn", "fünfzehn", "sechzehn", "siebzehn", "achtzehn", "neunzehn", "zwanzig"}

func (l *Localizer) Spell(number int) string {
	if l.lang != "de-DE" {
		panic("Only de-DE supported for number spelling")
	}
	return numbers[number]
}

func (l *Localizer) WikiEndpoint() string {
	switch l.lang {
	case "de-DE":
		return "de.wikipedia.org"
	default:
		panic("language '" + l.lang + "' not supported")
	}
}
