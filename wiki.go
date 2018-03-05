package wiki

type Page struct {
	Sections []Section
}

type Section struct {
	Title string
	Body  string
}

type Wiki interface {
	GetPage() Page
}
