package dynamodb

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/petergtz/alexa-wikipedia/persistence"
	"go.uber.org/zap"
)

type Persistence struct {
	dynamo *dynamodb.DynamoDB
	logger *zap.SugaredLogger
}

func NewPersistence(accessKeyID, secretAccessKey string, logger *zap.SugaredLogger) *Persistence {

	dynamoClient := dynamodb.New(session.Must(session.NewSession(&aws.Config{
		Region:      aws.String("eu-central-1"),
		Credentials: credentials.NewStaticCredentials(accessKeyID, secretAccessKey, ""),
	})))

	return &Persistence{
		dynamo: dynamoClient,
		logger: logger,
	}
}

func (p *Persistence) LogDefineIntentRequest(logEntry persistence.LogEntry) {
	input, e := dynamodbattribute.MarshalMap(logEntry)
	if e != nil {
		p.logger.Errorw("Could not marshal entry", "error", e)
		return
	}
	_, e = p.dynamo.PutItem(&dynamodb.PutItemInput{
		Item:      input,
		TableName: aws.String("AlexaWikipediaRequests"),
	})
	if e != nil {
		p.logger.Errorw("Could not log requests", "error", e)
		return
	}
}
