package wiki

import (
	"fmt"
	"strings"
)

type Page Section

type Section struct {
	Number      string
	Title       string
	Body        string
	Subsections []Section
}

type Wiki interface {
	GetPage(url string) (Page, error)
	SearchPage(url string) (Page, error)
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

var numbers = []string{"null", "eins", "zwei", "drei", "vier", "fünf", "sechs", "sieben", "acht", "neun", "zehn", "elf", "zwölf"}

func (p Page) TextAndPositionFromSectionNumber(sectionNumber string) (text string, position int) {
	return traverse2(Section(p), 0, 1, 0, sectionNumber, "")
}

func Convert(i int) string {
	return numbers[i]
}

func textFor(section Section) string {
	s := "Abschnitt " + section.Number + ". " + section.Title + ". " + section.Body
	for section.Body == "" && len(section.Subsections) > 0 {
		s += "Abschnitt " + section.Subsections[0].Number + ". " + section.Subsections[0].Title + ". " + section.Subsections[0].Body
		section = section.Subsections[0]
	}
	return s
}

func traverse2(s Section, level int, index int, cur int, sectionNumber string, prefix string) (text string, new_cur int) {
	var currentSectionNumber string
	switch level {
	case 0:
		currentSectionNumber = ""
	case 1:
		currentSectionNumber = Convert(index)
	default:
		currentSectionNumber = prefix + " punkt " + Convert(index)
	}

	if sectionNumber == currentSectionNumber {
		return textFor(s), cur
	}
	if s.Body != "" {
		cur++
	}
	for i, section := range s.Subsections {
		var ss string
		ss, cur = traverse2(section, level+1, i+1, cur, sectionNumber, currentSectionNumber)
		if ss != "" {
			return ss, cur
		}
	}
	return "", cur
}

func (p Page) TextAndPositionFromSectionName(sectionName string) (text string, position int) {
	return traverse3(Section(p), 0, 1, 0, sectionName)
}

func traverse3(s Section, level int, index int, cur int, sectionName string) (text string, new_cur int) {
	if strings.ToLower(sectionName) == strings.ToLower(s.Title) {
		return textFor(s), cur
	}
	if s.Body != "" {
		cur++
	}
	for i, section := range s.Subsections {
		var ss string
		ss, cur = traverse3(section, level+1, i+1, cur, sectionName)
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
