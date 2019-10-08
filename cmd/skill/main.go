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

	invocationCount := 0

	skill := factory.CreateSkill(logger)
	lambda.Start(func(ctx context.Context, requestEnv alexa.RequestEnvelope) (alexa.ResponseEnvelope, error) {
		invocationCount++
		lc, _ := lambdacontext.FromContext(ctx)

		if requestEnv.Request == nil {

			logger.Infow("Keep-alive CloudWatch Request",
				"aws-request-id", lc.AwsRequestID,
				"function-invocation-count", invocationCount)

			return alexa.ResponseEnvelope{}, nil
		}
		logger.Infow("Alexa Request",
			"aws-request-id", lc.AwsRequestID,
			"alexa-request-id", requestEnv.Request.RequestID,
			"function-invocation-count", invocationCount,
			"type", requestEnv.Request.Type,
			"intent", requestEnv.Request.Intent,
			"session-attributes", requestEnv.Session.Attributes,
			"locale", requestEnv.Request.Locale,
			"user-id", requestEnv.Session.User.UserID,
			"session-id", requestEnv.Session.SessionID)

		return *skill.ProcessRequest(&requestEnv), nil
	})
}
