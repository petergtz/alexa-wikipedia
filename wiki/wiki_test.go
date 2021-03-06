package wiki_test

import (
	"github.com/BurntSushi/toml"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"go.uber.org/zap"
	"golang.org/x/text/language"

	"github.com/petergtz/alexa-wikipedia/locale"
	"github.com/petergtz/alexa-wikipedia/wiki"
)

var logger, _ = zap.NewDevelopment()

var _ = Describe("Wiki", func() {
	var (
		localizer  *locale.Localizer
		page       wiki.Page
		i18nBundle *i18n.Bundle
	)

	BeforeEach(func() {
		i18nBundle = i18n.NewBundle(language.English)
		i18nBundle.RegisterUnmarshalFunc("toml", toml.Unmarshal)
		i18nBundle.MustParseMessageFileBytes(locale.DeDe, "active.de.toml")
		i18nBundle.MustParseMessageFileBytes(locale.EnUs, "active.en.toml")
	})

	JustBeforeEach(func() {
		page = wiki.Page{
			Title: "Main Title",
			Body:  "Intro",
			Subsections: []wiki.Section{
				wiki.Section{
					Number: localizer.Spell(1),
					Title:  "A",
					Body:   "Body A",
				},
				wiki.Section{
					Number: localizer.Spell(2),
					Title:  "B",
					Body:   "Body B",
				},
				wiki.Section{
					Number: localizer.Spell(3),
					Title:  "C",
					Body:   "",
					Subsections: []wiki.Section{
						wiki.Section{
							Number: localizer.Spell(3) + "." + localizer.Spell(1),
							Title:  "C.A",
							Body:   "Body C.A",
							Subsections: []wiki.Section{
								wiki.Section{
									Number: localizer.Spell(3) + "." + localizer.Spell(1) + "." + localizer.Spell(1),
									Title:  "C.A.A",
									Body:   "Body C.A.A",
								},
								wiki.Section{
									Number: localizer.Spell(3) + "." + localizer.Spell(1) + "." + localizer.Spell(2),
									Title:  "C.A.B",
									Body:   "Body C.A.B",
								},
								wiki.Section{
									Number: localizer.Spell(3) + "." + localizer.Spell(1) + "." + localizer.Spell(3),
									Title:  "C.A.C",
									Body:   "Body C.A.C",
								},
							},
						},
						wiki.Section{
							Number: localizer.Spell(3) + "." + localizer.Spell(2),
							Title:  "C.B",
							Body:   "Body C.B",
						},
					},
				},
			},
		}

	})

	Context("German", func() {
		BeforeEach(func() { localizer = locale.NewLocalizer(i18nBundle, "de-DE", logger.Sugar()) })

		It("Can get position 0", func() {
			Expect(page.TextForPosition(0)).To(Equal("Main Title. Intro"))
		})
		It("Can get position 1", func() {
			Expect(page.TextForPosition(1)).To(Equal("A. Body A"))
		})

		It("Can get position 2", func() {
			Expect(page.TextForPosition(2)).To(Equal("B. Body B"))
		})

		It("Can get position 3", func() {
			Expect(page.TextForPosition(3)).To(Equal("C. C.A. Body C.A"))
		})
		It("Can get position 4", func() {
			Expect(page.TextForPosition(4)).To(Equal("C.A.A. Body C.A.A"))
		})
		It("Can get position 5", func() {
			Expect(page.TextForPosition(5)).To(Equal("C.A.B. Body C.A.B"))
		})
		It("Can get position 6", func() {
			Expect(page.TextForPosition(6)).To(Equal("C.A.C. Body C.A.C"))
		})
		It("Can get position 7", func() {
			Expect(page.TextForPosition(7)).To(Equal("C.B. Body C.B"))
		})

		It("Can get section 3", func() {
			text, _ := page.TextAndPositionFromSectionNumber("drei", localizer)
			Expect(text).To(Equal("Abschnitt drei. C. Abschnitt drei.eins. C.A. Body C.A"))
		})

		It("Can get section 3.2", func() {
			text, _ := page.TextAndPositionFromSectionNumber("drei punkt zwei", localizer)
			Expect(text).To(Equal("Abschnitt drei.zwei. C.B. Body C.B"))
		})
	})

	Context("English", func() {
		BeforeEach(func() { localizer = locale.NewLocalizer(i18nBundle, "en-US", logger.Sugar()) })

		It("Can get section two", func() {
			text, position := page.TextAndPositionFromSectionNumber("two", localizer)
			Expect(position).To(Equal(2))
			Expect(text).To(Equal("section two. B. Body B"))
		})

		It("Can get section 2", func() {
			text, position := page.TextAndPositionFromSectionNumber("2", localizer)
			Expect(position).To(Equal(2))
			Expect(text).To(Equal("section two. B. Body B"))
		})

		It("Can get section 3 dot 2", func() {
			text, position := page.TextAndPositionFromSectionNumber("3 dot 2", localizer)
			Expect(position).To(Equal(7))
			Expect(text).To(Equal("section three.two. C.B. Body C.B"))
		})
	})
})
