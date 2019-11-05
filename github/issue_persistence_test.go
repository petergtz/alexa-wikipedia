package github_test

import (
	"io/ioutil"
	"strings"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/petergtz/alexa-wikipedia/github"
	"go.uber.org/zap"
)

var _ = Describe("GithubPersistence", func() {
	It("can use an issue comment as persistence", func() {
		token, e := ioutil.ReadFile("../private/github-access-token")
		Expect(e).NotTo(HaveOccurred())
		l, e := zap.NewDevelopment()
		if e != nil {
			panic(e)
		}
		defer l.Sync()
		log := l.Sugar()

		pp := github.NewGithubPersistence("petergtz", "alexa-wikipedia", 549126277, strings.TrimSpace(string(token)), log)
		pp.Persist([]string{"xyz"})
	})
})
