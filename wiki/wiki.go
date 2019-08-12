package wiki

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/petergtz/alexa-wikipedia/locale"
)

type Page Section

type Section struct {
	Number      string
	Title       string
	Body        string
	Subsections []Section
	// Locale      string
}

type Wiki interface {
	GetPage(url string, localizer *locale.Localizer) (Page, error)
	SearchPage(url string, localizer *locale.Localizer) (Page, error)
}

func (p Page) TextForPosition(position int) string {
	s, _ := traverse(Section(p), 0, position, "")
	return s
}

func traverse(s Section, cur int, target int, prefix string) (text string, new_cur int) {
	if s.Body != "" && cur == target {
		return prefix + s.Title + ". " + s.Body, cur
	}
	if s.Body != "" {
		cur++
	}
	for i, section := range s.Subsections {
		var ss string
		if i == 0 && s.Body == "" {
			ss, cur = traverse(section, cur, target, prefix+s.Title+". ")
		} else {
			ss, cur = traverse(section, cur, target, "")
		}
		if ss != "" {
			return ss, cur
		}
	}
	return "", cur
}

func (p Page) TextAndPositionFromSectionNumber(sectionNumber string, localizer *locale.Localizer) (text string, position int) {
	return traverse2(Section(p), 0, 1, 0, sectionNumber, "", "", localizer)
}

func textFor(section Section, localizer *locale.Localizer) string {
	s := localizer.MustLocalize(&i18n.LocalizeConfig{DefaultMessage: &i18n.Message{
		ID: "Section", Other: "Abschnitt",
	}}) + " " + section.Number + ". " + section.Title + ". " + section.Body
	for section.Body == "" && len(section.Subsections) > 0 {
		s += localizer.MustLocalize(&i18n.LocalizeConfig{DefaultMessage: &i18n.Message{
			ID: "Section", Other: "Abschnitt",
		}}) + " " + section.Subsections[0].Number + ". " + section.Subsections[0].Title + ". " + section.Subsections[0].Body
		section = section.Subsections[0]
	}
	return s
}

func traverse2(s Section, level int, index int, cur int, sectionNumber string, prefix string, prefixDigits string, localizer *locale.Localizer) (text string, new_cur int) {
	// Behavior of Alexa ASK is different for different locales.
	// In DE "zwei" is kept as "zwei" in slot value.
	// In EN "two" is converted to "2" in slot value.
	// To account for this, we maintain 2 versions of the section number string
	var (
		currentSectionNumber       string
		currentSectionNumberDigits string
	)
	switch level {
	case 0:
		currentSectionNumber = ""
		currentSectionNumberDigits = ""
	case 1:
		currentSectionNumber = localizer.Spell(index)
		currentSectionNumberDigits = strconv.Itoa(index)
	default:
		currentSectionNumber = prefix + " " + localizer.MustLocalize(&i18n.LocalizeConfig{DefaultMessage: &i18n.Message{
			ID: "Point", Other: "punkt",
		}}) + " " + localizer.Spell(index)
		currentSectionNumberDigits = prefixDigits + " " + localizer.MustLocalize(&i18n.LocalizeConfig{DefaultMessage: &i18n.Message{
			ID: "Point", Other: "punkt",
		}}) + " " + strconv.Itoa(index)
	}

	if sectionNumber == currentSectionNumber || sectionNumber == currentSectionNumberDigits {
		return textFor(s, localizer), cur
	}
	if s.Body != "" {
		cur++
	}
	for i, section := range s.Subsections {
		var ss string
		ss, cur = traverse2(section, level+1, i+1, cur, sectionNumber, currentSectionNumber, currentSectionNumberDigits, localizer)
		if ss != "" {
			return ss, cur
		}
	}
	return "", cur
}

func (p Page) TextAndPositionFromSectionName(sectionName string, localizer *locale.Localizer) (text string, position int) {
	return traverse3(Section(p), 0, 1, 0, sectionName, localizer)
}

func traverse3(s Section, level int, index int, cur int, sectionName string, localizer *locale.Localizer) (text string, new_cur int) {
	if strings.ToLower(sectionName) == strings.ToLower(s.Title) {
		return textFor(s, localizer), cur
	}
	if s.Body != "" {
		cur++
	}
	for i, section := range s.Subsections {
		var ss string
		ss, cur = traverse3(section, level+1, i+1, cur, sectionName, localizer)
		if ss != "" {
			return ss, cur
		}
	}
	return "", cur
}

func (p Page) Toc(localizer *locale.Localizer) string {
	s := localizer.MustLocalize(&i18n.LocalizeConfig{DefaultMessage: &i18n.Message{
		ID: "TableOfContents", Other: "Inhaltsverzeichnis",
	}}) + ". "
	for i, section := range p.Subsections {
		s += fmt.Sprintf(localizer.MustLocalize(&i18n.LocalizeConfig{DefaultMessage: &i18n.Message{
			ID: "Section", Other: "Abschnitt",
		}})+" %v: %v.\n", i+1, section.Title)
	}
	return s
}
