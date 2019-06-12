package main_test

import (
	"bytes"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"time"

	"github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo"
	"github.com/onsi/gomega"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"
)

func TestEndToEnd(t *testing.T) {
	gomega.RegisterFailHandler(ginkgo.Fail)
	ginkgo.RunSpecs(t, "EndToEnd")
}

var _ = Describe("Skill", func() {

	var (
		session *gexec.Session
		client  *http.Client
	)

	BeforeSuite(func() {
		pathToWebserver, err := gexec.Build("github.com/petergtz/alexa-wikipedia")
		Ω(err).ShouldNot(HaveOccurred())

		os.Setenv("PORT", "4443")
		os.Setenv("SKILL_ADDR", "127.0.0.1")
		os.Setenv("SKILL_SKIP_REQUEST_VALIDATION", "true")
		os.Setenv("APPLICATION_ID", "xxx")

		session, err = gexec.Start(exec.Command(pathToWebserver), GinkgoWriter, GinkgoWriter)
		Ω(err).ShouldNot(HaveOccurred())
		time.Sleep(200 * time.Millisecond)
		Expect(session.ExitCode()).To(Equal(-1), "Webserver error message: %s", string(session.Err.Contents()))

		client = &http.Client{}
	})

	AfterSuite(func() {
		if session != nil {
			session.Kill()
		}
		gexec.CleanupBuildArtifacts()
	})

	_, filename, _, _ := runtime.Caller(0)
	fileInfos, e := ioutil.ReadDir(filepath.Join(filepath.Dir(filename), "fixtures"))
	if e != nil {
		panic(e)
	}

	for _, loopFileInfo := range fileInfos {
		fileInfo := loopFileInfo
		if strings.Contains(fileInfo.Name(), "response") {
			continue
		}

		parts := strings.Split(fileInfo.Name(), ".")
		parts = parts[:len(parts)-1]
		d := make([]func(), len(parts))
		d[len(parts)-1] = func() {
			Describe(parts[len(parts)-1], func() {
				It("works", func() {
					response, e := client.Post("http://127.0.0.1:4443/", "", readerFrom(fileInfo.Name()))
					Expect(e).NotTo(HaveOccurred())
					Expect(response.StatusCode).To(Equal(http.StatusOK))
					Expect(ioutil.ReadAll(response.Body)).To(MatchJSON(stringFrom(strings.Replace(fileInfo.Name(), "json", "response.json", -1))))
				})
			})
		}

		for loopI := len(parts) - 2; loopI > 0; loopI-- {
			i := loopI
			d[i] = func() { Describe(parts[i], d[i+1]) }
		}
		Describe(parts[0], d[1])
	}

	Context("Invalid body", func() {
		It("returns a StatusBadRequest", func() {
			response, e := client.Post("http://127.0.0.1:4443/", "", strings.NewReader("Hello"))
			Expect(e).NotTo(HaveOccurred())

			Expect(response.StatusCode).To(Equal(http.StatusBadRequest))
		})
	})

})

func readerFrom(fixturename string) io.Reader {
	_, filename, _, _ := runtime.Caller(0)
	buf, e := ioutil.ReadFile(filepath.Join(filepath.Dir(filename), "fixtures", fixturename))
	Expect(e).NotTo(HaveOccurred())
	return bytes.NewReader(bytes.Replace(buf, []byte("TIMESTAMP"), []byte(time.Now().UTC().Format("2006-01-02T15:04:05Z")), -1))
}

func stringFrom(fixturename string) string {
	_, filename, _, _ := runtime.Caller(0)
	buf, e := ioutil.ReadFile(filepath.Join(filepath.Dir(filename), "fixtures", fixturename))
	Expect(e).NotTo(HaveOccurred())
	return string(bytes.Replace(buf, []byte("TIMESTAMP"), []byte(time.Now().UTC().Format("2006-01-02T15:04:05Z")), -1))
}
