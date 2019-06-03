package s3_test

import (
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"go.uber.org/zap"

	"os"

	p "github.com/petergtz/alexa-wikipedia/persistence"
	. "github.com/petergtz/alexa-wikipedia/s3"
)

type Credentials struct {
	AccessKeyId     string `yaml:"access_key_id"`
	SecretAccessKey string `yaml:"secret_access_key"`
}

var _ = Describe("S3", func() {
	It("can create an entry", func() {
		logger, e := zap.NewDevelopment()
		Expect(e).NotTo(HaveOccurred())

		if os.Getenv("ACCESS_KEY_ID") == "" {
			logger.Fatal("env var ACCESS_KEY_ID not provided.")
		}

		if os.Getenv("SECRET_ACCESS_KEY") == "" {
			logger.Fatal("env var SECRET_ACCESS_KEY not provided.")
		}

		persistence := NewPersistence(os.Getenv("ACCESS_KEY_ID"), os.Getenv("SECRET_ACCESS_KEY"), "alexa-wikipedia", logger.Sugar())
		defer persistence.ShutDown()
		persistence.LogDefineIntentRequest(p.LogEntry{
			UnixTimestamp: time.Now().Unix(),
			Timestamp:     time.Now(),
			SearchQuery:   "Bla",
			ActualTitle:   "blub",
			Locale:        "de-DE",
			UserID:        "userid1",
			SessionID:     "sessionid1",
		})
		persistence.LogDefineIntentRequest(p.LogEntry{
			UnixTimestamp: time.Now().Unix(),
			Timestamp:     time.Now(),
			SearchQuery:   "Bla2",
			ActualTitle:   "blub2",
			Locale:        "english",
			UserID:        "userid2",
			SessionID:     "sessionid2",
		})
		persistence.LogDefineIntentRequest(p.LogEntry{
			UnixTimestamp: time.Now().Unix(),
			Timestamp:     time.Now(),
			SearchQuery:   "Bla3",
			ActualTitle:   "blub4",
			Locale:        "de-DE",
			UserID:        "userid3",
			SessionID:     "sessionid3",
		})
	})
})
