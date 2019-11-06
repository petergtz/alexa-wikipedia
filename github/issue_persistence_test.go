package github_test

import (
	"io/ioutil"
	"strings"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/petergtz/alexa-wikipedia/github"
	"github.com/petergtz/go-alexa"
)

type noOpErrorReporter struct{}

func (*noOpErrorReporter) ReportPanic(interface{}, *alexa.RequestEnvelope) {}
func (*noOpErrorReporter) ReportError(error)                               {}

var _ = Describe("GithubPersistence", func() {
	It("can use an issue comment as persistence", func() {
		token, e := ioutil.ReadFile("../private/github-access-token")
		Expect(e).NotTo(HaveOccurred())

		pp := github.NewGithubPersistence(
			"petergtz",
			"alexa-wikipedia",
			549126277,
			strings.TrimSpace(string(token)),
			&noOpErrorReporter{})
		pp.Persist([]string{"xyz"})
	})
})
