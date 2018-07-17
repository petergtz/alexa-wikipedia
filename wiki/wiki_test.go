package wiki_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/petergtz/alexa-wikipedia/wiki"
)

var page = wiki.Page{
	Title: "Main Title",
	Body:  "Intro",
	Subsections: []wiki.Section{
		wiki.Section{
			Number: wiki.Convert(1),
			Title:  "A",
			Body:   "Body A",
		},
		wiki.Section{
			Number: wiki.Convert(2),
			Title:  "B",
			Body:   "Body B",
		},
		wiki.Section{
			Number: wiki.Convert(3),
			Title:  "C",
			Body:   "",
			Subsections: []wiki.Section{
				wiki.Section{
					Number: wiki.Convert(3) + "." + wiki.Convert(1),
					Title:  "C.A",
					Body:   "Body C.A",
					Subsections: []wiki.Section{
						wiki.Section{
							Number: wiki.Convert(3) + "." + wiki.Convert(1) + "." + wiki.Convert(1),
							Title:  "C.A.A",
							Body:   "Body C.A.A",
						},
						wiki.Section{
							Number: wiki.Convert(3) + "." + wiki.Convert(1) + "." + wiki.Convert(2),
							Title:  "C.A.B",
							Body:   "Body C.A.B",
						},
						wiki.Section{
							Number: wiki.Convert(3) + "." + wiki.Convert(1) + "." + wiki.Convert(3),
							Title:  "C.A.C",
							Body:   "Body C.A.C",
						},
					},
				},
				wiki.Section{
					Number: wiki.Convert(3) + "." + wiki.Convert(2),
					Title:  "C.B",
					Body:   "Body C.B",
				},
			},
		},
	},
}

var _ = Describe("Wiki", func() {
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
		text, _ := page.TextAndPositionFromSectionNumber("drei")
		Expect(text).To(Equal("Abschnitt drei. C. Abschnitt drei.eins. C.A. Body C.A"))
	})

	It("Can get section 3.2", func() {
		text, _ := page.TextAndPositionFromSectionNumber("drei punkt zwei")
		Expect(text).To(Equal("Abschnitt drei.zwei. C.B. Body C.B"))
	})

})
