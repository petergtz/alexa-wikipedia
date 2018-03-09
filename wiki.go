package wiki

import (
	"fmt"
	"strings"
)

type Page Section

type Section struct {
	Title       string
	Body        string
	Subsections []Section
}

type Wiki interface {
	GetPage() Page
}

func (p Page) TextForPosition(position int) string {
	s, _ := traverse(Section(p), 0, position, "")
	return s
}

func traverse(s Section, cur int, target int, prefix string) (text string, new_cur int) {
	if s.Body != "" {
		if cur == target {
			return prefix + s.Title + ". " + s.Body, cur
		}
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

func (p Page) TextAndPositionFromSectionNumber(sectionNumber int) (text string, position int) {
	s := p.Subsections[sectionNumber].Body
	if s == "" && len(p.Subsections[sectionNumber].Subsections) > 0 {
		s = p.Subsections[sectionNumber].Subsections[0].Body
	}
	// TODO find position
	return s, 0
}

func (p Page) TextAndPositionFromSectionName(sectionName string) (text string, position int) {
	for _, section := range p.Subsections {
		if strings.ToLower(section.Title) == strings.ToLower(sectionName) {
			// TODO find position
			return section.Body, 0
		}
	}
	return "", 0
}

func (p Page) Toc() string {
	s := ""
	for i, section := range p.Subsections {
		s += fmt.Sprintf("Abschnitt %v: %v.\n", i+1, section.Title)
	}
	return s
}
