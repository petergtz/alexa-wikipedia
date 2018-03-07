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
