package github

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"runtime/debug"
	"sort"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sns"
	"github.com/davecgh/go-spew/spew"

	"go.uber.org/zap"

	"github.com/google/go-github/github"
	"github.com/pkg/errors"
	"golang.org/x/oauth2"
)

type ErrorReporter struct {
	ghClient    *github.Client
	logger      *zap.SugaredLogger
	ctx         context.Context
	owner       string
	repo        string
	logsURL     string
	snsClient   *sns.SNS
	snsTopicArn string
}

func NewErrorReporter(owner, repo, token string, logger *zap.SugaredLogger, logsURL string, snsClient *sns.SNS, snsTopicArn string) *ErrorReporter {
	ctx := context.TODO()
	spew.Config.ContinueOnMethod = true
	return &ErrorReporter{
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

func (r *ErrorReporter) ReportPanic(e interface{}, context interface{}) {
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
		attributes["issue-url"] = issue.GetHTMLURL()

		r.logger.Errorw("Internal Server Error", slicify(attributes)...)
	}

	_, snsErr := r.snsClient.Publish(&sns.PublishInput{
		TopicArn: aws.String(r.snsTopicArn),
		Subject:  aws.String(fmt.Sprintf(r.repo+": Internal Server Error (ErrID: %v)", errorID)),
		Message: aws.String(fmt.Sprintf(`ERROR DETAILS:
%s

CONTEXT:
%v

CLOUDWATCH QUERY:
%v`, stringify(attributes), marshalContext(context), fmt.Sprintf(r.logsURL, errorID))),
	})

	if snsErr != nil {
		attributes["sns-error"] = snsErr

		r.logger.Errorw("Error while trying to publish Internal Server Error via SNS", slicify(attributes)...)
	}
}

func (r *ErrorReporter) ReportError(e error) {
	r.ReportPanic(e, nil)
}

func errorStringFrom(e interface{}) string {
	if _, hasStackTrace := e.(interface{ StackTrace() errors.StackTrace }); hasStackTrace {
		return fmt.Sprintf("STRING: %v\nSTACKTRACE:\n%+v\nINTROSPECTION:\n%v", e, e, spew.Sdump(e))
	}
	return fmt.Sprintf("STRING: %v\nSTACKTRACE:\n%s\nINTROSPECTION:\n%v", e, debug.Stack(), spew.Sdump(e))
}

func stringify(m map[string]interface{}) string {
	var slice []string
	for k, v := range m {
		slice = append(slice, fmt.Sprintf("%v: %v", k, v))
	}
	sort.Strings(slice)
	return strings.Join(slice, "\n")
}

func slicify(m map[string]interface{}) []interface{} {
	result := make([]interface{}, 0)
	for k, v := range m {
		result = append(result, k, v)
	}
	return result
}

func marshalContext(context interface{}) string {
	if context == nil {
		return "Not available."
	}
	buf, e := json.MarshalIndent(context, "", "  ")
	if e != nil {
		return "Error while marshalling context. Error: " + e.Error()
	}
	return string(buf)
}
