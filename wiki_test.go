package wiki_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	wiki "github.com/petergtz/alexa-wikipedia"
)

var page = wiki.Page{
	Title: "Main Title",
	Body:  "Intro",
	Subsections: []wiki.Section{
		wiki.Section{
			Title: "A",
			Body:  "Body A",
		},
		wiki.Section{
			Title: "B",
			Body:  "Body B",
		},
		wiki.Section{
			Title: "C",
			Body:  "",
			Subsections: []wiki.Section{
				wiki.Section{
					Title: "C.A",
					Body:  "Body C.A",
					Subsections: []wiki.Section{
						wiki.Section{
							Title: "C.A.A",
							Body:  "Body C.A.A",
						},
						wiki.Section{
							Title: "C.A.B",
							Body:  "Body C.A.B",
						},
						wiki.Section{
							Title: "C.A.C",
							Body:  "Body C.A.C",
						},
					},
				},
				wiki.Section{
					Title: "C.B",
					Body:  "Body C.B",
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

})
