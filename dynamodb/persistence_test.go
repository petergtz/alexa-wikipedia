package dynamodb_test

import (
	"time"

	"os"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/petergtz/alexa-wikipedia/dynamodb"
	p "github.com/petergtz/alexa-wikipedia/persistence"
	"go.uber.org/zap"
)

var _ = Describe("Persistence", func() {
	It("works", func() {
		logger, e := zap.NewDevelopment()
		Expect(e).NotTo(HaveOccurred())
		if os.Getenv("ACCESS_KEY_ID") == "" {
			logger.Fatal("env var ACCESS_KEY_ID not provided.")
		}

		if os.Getenv("SECRET_ACCESS_KEY") == "" {
			logger.Fatal("env var SECRET_ACCESS_KEY not provided.")
		}

		persistence := dynamodb.NewPersistence(os.Getenv("ACCESS_KEY_ID"), os.Getenv("SECRET_ACCESS_KEY"), logger.Sugar())
		persistence.LogDefineIntentRequest(p.LogEntry{
			RequestID:     "req1",
			UnixTimestamp: time.Now().Unix(),
			Timestamp:     time.Now(),
			SearchQuery:   "Bla",
			ActualTitle:   "blub",
			Locale:        "de-DE",
			UserID:        "userid1",
			SessionID:     "sessionid1",
		})
		persistence.LogDefineIntentRequest(p.LogEntry{
			RequestID:     "req2",
			UnixTimestamp: time.Now().Unix(),
			Timestamp:     time.Now(),
			SearchQuery:   "Bla2",
			ActualTitle:   "blub2",
			Locale:        "english",
			UserID:        "userid2",
			SessionID:     "sessionid2",
		})
		persistence.LogDefineIntentRequest(p.LogEntry{
			RequestID:     "req3",
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
