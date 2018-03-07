package wiki

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
	return traverse(Section(p), 0, position)
}

func traverse(s Section, cur int, target int) string {
	if cur == target && s.Body != "" {
		return s.Title + ". " + s.Body
	}
	for _, section := range s.Subsections {
		s := traverse(section, cur+1, target)
		if s != "" {
			return s
		}
	}
	return ""
}
