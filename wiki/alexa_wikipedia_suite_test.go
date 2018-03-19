package wiki_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestAlexaWikipedia(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "AlexaWikipedia Suite")
}
