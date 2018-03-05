package mediawiki_test

import (
	. "github.com/onsi/ginkgo"
)

var _ = Describe("Mediawiki", func() {
	Describe("action=query prop=extract", func() {
		It("works", func() {
			// r, e := http.Get("https://de.wikipedia.org/w/api.php?format=json&action=query&prop=extracts&titles=" + "käsekuchen" + "&redirects=true&formatversion=2")

		})
	})
	Describe("action=query prop=revision", func() {

	})
	Describe("action=parse", func() {

	})
})
