package main

import (
	"github.com/aws/aws-lambda-go/lambda"
	. "github.com/petergtz/alexa-wikipedia/skill"
	"github.com/petergtz/go-alexa"
)

func main() {
	skill := CreateSkill()
	lambda.Start(func(requestEnv alexa.RequestEnvelope) (alexa.ResponseEnvelope, error) {
		return *skill.ProcessRequest(&requestEnv), nil
	})
}
