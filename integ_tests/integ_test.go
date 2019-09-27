package main_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"math/rand"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"time"

	"go.uber.org/zap"

	"github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo"
	"github.com/onsi/gomega"
	. "github.com/onsi/gomega"
	"github.com/petergtz/alexa-wikipedia/cmd/skill/factory"
	"github.com/petergtz/go-alexa"
)

func TestEndToEnd(t *testing.T) {
	rand.Seed(time.Now().Unix())
	gomega.RegisterFailHandler(ginkgo.Fail)
	ginkgo.RunSpecs(t, "EndToEnd")
}

var skill alexa.Skill

var _ = Describe("Skill", func() {
	BeforeSuite(func() {
		logger, e := zap.NewDevelopment()
		Expect(e).NotTo(HaveOccurred())
		skill = factory.CreateSkill(logger.Sugar())
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
					c, e := ioutil.ReadAll(readerFrom(fileInfo.Name()))
					Expect(e).NotTo(HaveOccurred())

					var requestEnv alexa.RequestEnvelope
					e = json.Unmarshal(c, &requestEnv)
					Expect(e).NotTo(HaveOccurred())

					Expect(json.MarshalIndent(skill.ProcessRequest(&requestEnv), "", "  ")).To(MatchJSON(stringFrom(strings.Replace(fileInfo.Name(), "json", "response.json", -1))))
				})
			})
		}

		for loopI := len(parts) - 2; loopI > 0; loopI-- {
			i := loopI
			d[i] = func() { Describe(parts[i], d[i+1]) }
		}
		Describe(parts[0], d[1])
	}
})

func readerFrom(fixturename string) io.Reader {
	_, filename, _, _ := runtime.Caller(0)
	buf, e := ioutil.ReadFile(filepath.Join(filepath.Dir(filename), "fixtures", fixturename))
	Expect(e).NotTo(HaveOccurred())
	buf = bytes.Replace(buf, []byte("TIMESTAMP"), []byte(time.Now().UTC().Format("2006-01-02T15:04:05Z")), -1)
	// This is necessary to avoid the skill to ask the user if Alexa understood her right, due to repeated same queries:
	buf = bytes.Replace(buf, []byte(`"userId": "xxx"`), []byte(fmt.Sprintf(`"userId": "%v"`, rand.Int())), -1)
	return bytes.NewReader(buf)
}

func stringFrom(fixturename string) string {
	_, filename, _, _ := runtime.Caller(0)
	buf, e := ioutil.ReadFile(filepath.Join(filepath.Dir(filename), "fixtures", fixturename))
	Expect(e).NotTo(HaveOccurred())
	return string(bytes.Replace(buf, []byte("TIMESTAMP"), []byte(time.Now().UTC().Format("2006-01-02T15:04:05Z")), -1))
}
