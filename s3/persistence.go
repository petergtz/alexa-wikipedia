package s3

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math/rand"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"go.uber.org/zap"
)

type Persistence struct {
	s3Client *s3.S3
	bucket   string
	logger   *zap.SugaredLogger
}

func NewPersistence(accessKeyID, secretAccessKey, bucket string, logger *zap.SugaredLogger) *Persistence {
	rand.Seed(time.Now().Unix())
	return &Persistence{
		s3Client: s3.New(session.Must(session.NewSession(&aws.Config{
			Region:      aws.String("eu-central-1"),
			Credentials: credentials.NewStaticCredentials(accessKeyID, secretAccessKey, ""),
		}))),
		bucket: bucket,
		logger: logger,
	}
}

func (p *Persistence) LogDefineIntentRequest(timestamp time.Time, searchQuery string, actualTitle string, locale string) {
	b, e := json.Marshal(map[string]interface{}{
		"timestamp":    timestamp.Unix(),
		"search_query": searchQuery,
		"actual_title": actualTitle,
		"locale":       locale,
	})
	if e != nil {
		p.logger.Errorw("Error while trying to marshal DefineIntent request data",
			"bucket", p.bucket,
			"timestamp", timestamp,
			"search-query", searchQuery,
			"actual-title", actualTitle,
			"locale", locale,
			"error", e)
	}
	_, e = p.s3Client.PutObject(&s3.PutObjectInput{
		Bucket: aws.String(p.bucket),
		Key:    aws.String(locale + "_" + timestamp.Format(time.RFC3339Nano) + "_" + fmt.Sprintf("%v", rand.Intn(10000)) + ".json"),
		Body:   bytes.NewReader(b),
	})
	if e != nil {
		p.logger.Errorw("Error while trying to upload DefineIntent request data",
			"bucket", p.bucket,
			"timestamp", timestamp,
			"search-query", searchQuery,
			"actual-title", actualTitle,
			"locale", locale,
			"error", e)
	}
}
