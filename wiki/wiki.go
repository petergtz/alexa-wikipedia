package wiki

import (
	"fmt"
	"strings"

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
	return traverse2(Section(p), 0, 1, 0, sectionNumber, "", localizer)
}

func textFor(section Section, localizer *locale.Localizer) string {
	s := "Abschnitt " + section.Number + ". " + section.Title + ". " + section.Body
	for section.Body == "" && len(section.Subsections) > 0 {
		s += "Abschnitt " + section.Subsections[0].Number + ". " + section.Subsections[0].Title + ". " + section.Subsections[0].Body
		section = section.Subsections[0]
	}
	return s
}

func traverse2(s Section, level int, index int, cur int, sectionNumber string, prefix string, localizer *locale.Localizer) (text string, new_cur int) {
	var currentSectionNumber string
	switch level {
	case 0:
		currentSectionNumber = ""
	case 1:
		currentSectionNumber = localizer.Spell(index)
	default:
		currentSectionNumber = prefix + " punkt " + localizer.Spell(index)
	}

	if sectionNumber == currentSectionNumber {
		return textFor(s, localizer), cur
	}
	if s.Body != "" {
		cur++
	}
	for i, section := range s.Subsections {
		var ss string
		ss, cur = traverse2(section, level+1, i+1, cur, sectionNumber, currentSectionNumber, localizer)
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

func (p Page) Toc() string {
	s := "Inhaltsverzeichnis. "
	for i, section := range p.Subsections {
		s += fmt.Sprintf("Abschnitt %v: %v.\n", i+1, section.Title)
	}
	return s
}
