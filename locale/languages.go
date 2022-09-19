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
YouAreAtWikipediaNow = "Du befindest Dich jetzt bei Wikipedia. Um einen Artikel vorgelesen zu bekommen, sage z.B. \"Suche nach Käsekuchen\"."
HelpText = "Um einen Artikel vorgelesen zu bekommen, sage z.B. \"Suche nach Käsekuchen.\" oder \"Was ist Käsekuchen?\". Du kannst jederzeit zum Inhaltsverzeichnis springen, indem Du \"Inhaltsverzeichnis\" sagst. Oder sage \"Springe zu Abschnitt 3.2\", um direkt zu diesem Abschnitt zu springen."
QuickHelpText = "Suche zunächst nach einem Begriff. Sage z.B. \"Suche nach Käsekuchen.\" oder \"Was ist Käsekuchen?\"."
FurtherNavigationHints = "Zur weiteren Navigation kannst Du jederzeit zum Inhaltsverzeichnis springen indem Du \"Inhaltsverzeichnis\" oder \"nächster Abschnitt\" sagst. Soll ich zunächst einfach weiterlesen?"
FallbackText = "Meine Enzyklopädie kann hiermit nicht weiterhelfen. Aber Du kannst z.B. sagen \"Suche nach Käsekuchen\"."
CouldNotFindSpelledTerm = "Den buchstabierten Begriff {{.SpelledTerm}} konnte ich bei Wikipedia leider nicht finden. Versuche es doch mit einem anderen Begriff."
InternalError = "Es ist ein interner Fehler aufgetreten bei der Benutzung von Wikipedia."
EndOfArticle = "Oh! Wir sind bereits am Ende angelangt. Wenn Du noch einen weiteren Artikel vorgelesen kriegen möchtest, sage z.B. \"Suche nach Elefant\"."
SpellingHint = "Ich habe den Artikel, \"{{.Title}}\", gerade erst gelesen. Falls ich nicht Deinen gewünschten Artikel gefunden habe, unterbrich mich und sage: \"Alexa, Suche buchstabieren\", um Deine Suchanfrage zu buchstabieren. Hier ist der Artikel:"
`)

	EnUs = []byte(`
CouldNotFindExpression = "I couldn't find this in Wikipedia. Why don't you a try a different expression?"
CouldNotFindSection = "I couldn't find the section \"{{.SectionTitleOrNumber}}\"."
"No?Okay" = "No? Okay."
Point = "dot"
Section = "section"
ShouldIContinue = "Shall I continue?"
TableOfContents = "Table of contents"
What = "Excuse me?"
WhichSectionToJump = "Which section do you want to go to?"
YouAreAtWikipediaNow = "This is Wikipedia. To read an article say e.g. \"What is a cheese cake?\"."
HelpText = "To read an article say e.g. \"What is a cheese cake?\". You can jump to the table of contents any time by saying \"table of contents\". Or say \"Jump to section 3.2\"."
QuickHelpText = "First, search for something. E.g. say \"search for cheese cake\" or \"What is cheese cake?\"."
FurtherNavigationHints = "For further navigation you can jump to the table of contents any time or jump to a specific section. For now, shall I simply continue reading?"
FallbackText = "My encyclopedia cannot help with that. But you could e.g. search for cheese cake."
CouldNotFindSpelledTerm = "I could not find the spelled term \"{{.SpelledTerm}}\" on Wikipedia. Why don't you a try a different term?"
InternalError = "An internal error occurred while using My encyclopdia."
EndOfArticle = "Oh, we've already reached the end of the article. Why don't you try a different search. E.g. say \"What is an elephant?\""
SpellingHint = "I've just read the article \"{{.Title}}\" already. Did I not find what you were looking for? If so, interrupt me and say \"Alexa, spell search\" to spell your search query. Here's the article:"
`)
)
