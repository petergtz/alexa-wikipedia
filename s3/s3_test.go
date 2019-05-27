package s3_test

import (
	"io/ioutil"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"go.uber.org/zap"

	. "github.com/petergtz/alexa-wikipedia/s3"
	yaml "gopkg.in/yaml.v2"
)

type Credentials struct {
	AccessKeyId     string `yaml:"access_key_id"`
	SecretAccessKey string `yaml:"secret_access_key"`
}

var _ = Describe("S3", func() {
	It("can create an entry", func() {
		content, e := ioutil.ReadFile("../private/s3-credentials")
		Expect(e).NotTo(HaveOccurred())
		var credentials Credentials
		e = yaml.Unmarshal(content, &credentials)
		Expect(e).NotTo(HaveOccurred())
		logger, e := zap.NewDevelopment()
		Expect(e).NotTo(HaveOccurred())

		persistence := NewPersistence(credentials.AccessKeyId, credentials.SecretAccessKey, "alexa-wikipedia", logger.Sugar())
		persistence.LogDefineIntentRequest(time.Now(), "Bla", "blub", "de-DE")
	})
})
