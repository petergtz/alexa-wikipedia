package github_test

import (
	"io/ioutil"
	"strings"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/petergtz/alexa-wikipedia/github"
)

var _ = Describe("GithubPersistence", func() {
	It("can use an issue comment as persistence", func() {
		token, e := ioutil.ReadFile("../private/github-access-token")
		Expect(e).NotTo(HaveOccurred())

		pp := github.NewPersistence(
			"petergtz",
			"alexa-wikipedia",
			549126277,
			strings.TrimSpace(string(token)),
		)
		e = pp.Persist([]string{"xyz"})
		Expect(e).NotTo(HaveOccurred())
	})
})
