package main

import (
	"context"

	"github.com/aws/aws-lambda-go/lambdacontext"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/petergtz/alexa-wikipedia/cmd/skill/factory"
	"github.com/petergtz/go-alexa"
)

func main() {
	logger := factory.CreateLoggerWith("debug")
	defer logger.Sync()

	skill := factory.CreateSkill(logger)
	lambda.Start(func(ctx context.Context, requestEnv alexa.RequestEnvelope) (alexa.ResponseEnvelope, error) {
		lc, _ := lambdacontext.FromContext(ctx)
		logger.Infow("Request",
			"aws-request-id", lc.AwsRequestID,
			"alexa-request-id", requestEnv.Request.RequestID,
			"type", requestEnv.Request.Type,
			"intent", requestEnv.Request.Intent,
			"session-attributes", requestEnv.Session.Attributes,
			"locale", requestEnv.Request.Locale)

		return *skill.ProcessRequest(&requestEnv), nil
	})
}
