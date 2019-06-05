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
	"github.com/facebookarchive/muster"
	"github.com/petergtz/go-alexa"
	"go.uber.org/zap"
)

type Persistence struct {
	musterClient *muster.Client
}

func NewPersistence(accessKeyID, secretAccessKey, bucket string, logger *zap.SugaredLogger) *Persistence {
	rand.Seed(time.Now().Unix())
	s3Client := s3.New(session.Must(session.NewSession(&aws.Config{
		Region:      aws.String("eu-central-1"),
		Credentials: credentials.NewStaticCredentials(accessKeyID, secretAccessKey, ""),
	})))
	result := &Persistence{
		musterClient: &muster.Client{
			MaxBatchSize:        100,
			BatchTimeout:        5 * time.Minute,
			PendingWorkCapacity: 1000,
			BatchMaker: func() muster.Batch {
				return &batch{
					logger:   logger,
					bucket:   bucket,
					s3Client: s3Client,
				}
			},
		},
	}
	e := result.musterClient.Start()
	if e != nil {
		panic(e)
	}
	return result
}

func (p *Persistence) LogDefineIntentRequest(logEntry alexa.Interaction) {
	p.musterClient.Work <- logEntry
}

func (p *Persistence) ShutDown() {
	p.musterClient.Stop()
}

type batch struct {
	s3Client *s3.S3
	bucket   string
	logger   *zap.SugaredLogger

	entries []alexa.Interaction
}

func (bm *batch) Add(item interface{}) {
	bm.entries = append(bm.entries, item.(alexa.Interaction))
}

func (bm *batch) Fire(notifier muster.Notifier) {
	defer notifier.Done()

	b, e := json.Marshal(bm.entries)
	if e != nil {
		bm.logger.Errorw("Error while trying to marshal log entries", "error", e)
	}

	_, e = bm.s3Client.PutObject(&s3.PutObjectInput{
		Bucket: aws.String(bm.bucket),
		Key:    aws.String(time.Now().Format(time.RFC3339Nano) + "_" + fmt.Sprintf("%v", rand.Intn(10000)) + ".json"),
		Body:   bytes.NewReader(b),
	})
	if e != nil {
		bm.logger.Errorw("Error while trying to upload DefineIntent request data", "bucket", bm.bucket, "error", e)
	}

}
