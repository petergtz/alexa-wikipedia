package github

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"runtime/debug"

	"github.com/petergtz/go-alexa"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sns"

	"go.uber.org/zap"

	"github.com/google/go-github/github"
	"github.com/pkg/errors"
	"golang.org/x/oauth2"
)

type GithubErrorReporter struct {
	ghClient    *github.Client
	logger      *zap.SugaredLogger
	ctx         context.Context
	owner       string
	repo        string
	logsURL     string
	snsClient   *sns.SNS
	snsTopicArn string
}

func NewGithubErrorReporter(owner, repo, token string, logger *zap.SugaredLogger, logsURL string, snsClient *sns.SNS, snsTopicArn string) *GithubErrorReporter {
	ctx := context.TODO()
	return &GithubErrorReporter{
		ghClient:    github.NewClient(oauth2.NewClient(ctx, oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token}))),
		ctx:         ctx,
		logger:      logger,
		repo:        repo,
		owner:       owner,
		logsURL:     logsURL,
		snsClient:   snsClient,
		snsTopicArn: snsTopicArn,
	}
}

func (r *GithubErrorReporter) ReportPanic(e interface{}, requestEnv *alexa.RequestEnvelope) {
	errorID := rand.Int63()
	errorString := errorStringFrom(e)

	issue, _, ghErr := r.ghClient.Issues.Create(r.ctx, r.owner, r.repo, &github.IssueRequest{
		Title: github.String(fmt.Sprintf("Internal Server Error (ErrID: %v)", errorID)),
		Body:  github.String(fmt.Sprintf("An error occurred and it can be found using %v", fmt.Sprintf(r.logsURL, errorID))),
	})

	attributes := map[string]interface{}{
		"error-id": errorID,
		"error":    errorString,
	}
	if ghErr != nil {
		attributes["github-error"] = ghErr

		r.logger.Errorw("Error while trying to report Internal Server Error", slicify(attributes)...)
	} else {
		attributes["github-issue-id"] = *issue.Number
		attributes["github-issue-url"] = issue.GetHTMLURL()

		r.logger.Errorw("Internal Server Error", slicify(attributes)...)
	}

	_, snsErr := r.snsClient.Publish(&sns.PublishInput{
		TopicArn: aws.String(r.snsTopicArn),
		Subject:  aws.String(fmt.Sprintf(r.repo+": Internal Server Error (ErrID: %v)", errorID)),
		Message: aws.String(fmt.Sprintf(`ERROR DETAILS:
%s

ALEXA REQUEST:
%v

CLOUDWATCH QUERY:
%v`, stringify(attributes), alexaRequestString(requestEnv), fmt.Sprintf(r.logsURL, errorID))),
	})

	if snsErr != nil {
		attributes["sns-error"] = snsErr

		r.logger.Errorw("Error while trying to publish Internal Server Error via SNS", slicify(attributes)...)
	}
}

func (r *GithubErrorReporter) ReportError(e error) {
	r.ReportPanic(e, nil)
}

func errorStringFrom(e interface{}) string {
	if _, hasStackTrace := e.(interface{ StackTrace() errors.StackTrace }); hasStackTrace {
		return fmt.Sprintf("%+v", e)
	}
	return fmt.Sprintf("%v\n%s", e, debug.Stack())
}

func stringify(m map[string]interface{}) string {
	result := ""
	for k, v := range m {
		result += fmt.Sprintf("%v: %v\n", k, v)
	}
	return result
}

func slicify(m map[string]interface{}) []interface{} {
	result := make([]interface{}, 0)
	for k, v := range m {
		result = append(result, k, v)
	}
	return result
}

func alexaRequestString(requestEnv *alexa.RequestEnvelope) string {
	if requestEnv == nil {
		return "Not available."
	}
	r := *requestEnv
	r.Session.User.AccessToken = "<REDACTED>"
	buf, e := json.MarshalIndent(r, "", "  ")
	if e != nil {
		return "Error while marshalling request. Error: " + e.Error()
	}
	return string(buf)
}
