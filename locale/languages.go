package locale

var (
	DeDe = []byte(`
CouldNotFindExpression = "Diesen Begriff konnte ich bei Wikipedia leider nicht finden. Versuche es doch mit einem anderen Begriff."
CouldNotFindSection = "Ich konnte den angegebenen Abschnitt \"{{.SectionTitleOrNumber}}\" nicht finden."
"No?Okay" = "Nein? Okay."
Point = "punkt"
Section = "Abschnitt"
ShouldIContinue = "Soll ich noch weiterlesen?"
TableOfContents = "Inhaltsverzeichnis"
What = "Wie meinen?"
WhichSectionToJump = "Zu welchem Abschnitt möchtest Du springen?"
YouAreAtWikipediaNow = "Du befindest Dich jetzt bei Wikipedia."
HelpText = "Um einen Artikel vorgelesen zu bekommen, sage z.B. \"Suche nach Käsekuchen.\" oder \"Was ist Käsekuchen?\". Du kannst jederzeit zum Inhaltsverzeichnis springen, indem Du \"Inhaltsverzeichnis\" sagst. Oder sage \"Springe zu Abschnitt 3.2\", um direkt zu diesem Abschnitt zu springen."
`)

	EnUs = []byte(`
CouldNotFindExpression = "I couldn't find this in Wikipedia. Why don't a try a different expression?"
CouldNotFindSection = "I couldn't find the section \"{{.SectionTitleOrNumber}}\"."
"No?Okay" = "No? Okay."
Point = "point"
Section = "section"
ShouldIContinue = "Shall I continue?"
TableOfContents = "Table of contents"
What = "Excuse me?"
WhichSectionToJump = "Which section do you want to go to?"
YouAreAtWikipediaNow = "This is Wikipedia."
HelpText = "To read an article say e.g. \"What is a cheese cake?\"."
`)
)
