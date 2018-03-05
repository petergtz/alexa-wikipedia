package mediawiki_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestMediawiki(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Mediawiki Suite")
}
