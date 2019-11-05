package github

import (
	"context"
	"regexp"
	"strings"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

type ErrorReporter interface {
	ReportError(error)
}

type GithubPersistence struct {
	ghClient       *github.Client
	ctx            context.Context
	owner          string
	repo           string
	issueCommentID int64
	errorReporter  ErrorReporter
}

func NewGithubPersistence(owner, repo string, issueCommentID int64, token string, errorReporter ErrorReporter) *GithubPersistence {
	ctx := context.TODO()
	return &GithubPersistence{
		ghClient:       github.NewClient(oauth2.NewClient(ctx, oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token}))),
		ctx:            ctx,
		errorReporter:  errorReporter,
		repo:           repo,
		owner:          owner,
		issueCommentID: issueCommentID,
	}
}

func (pp *GithubPersistence) Persist(findings []string) {
	comment, _, e := pp.ghClient.Issues.GetComment(pp.ctx, pp.owner, pp.repo, pp.issueCommentID)
	if e != nil {
		pp.errorReporter.ReportError(e)
		return
	}
	updated := AddMissing(comment.GetBody(), findings)
	if updated == comment.GetBody() {
		return
	}
	_, _, e = pp.ghClient.Issues.EditComment(pp.ctx, pp.owner, pp.repo, pp.issueCommentID, &github.IssueComment{
		ID:   github.Int64(pp.issueCommentID),
		Body: github.String(updated),
	})
	if e != nil {
		pp.errorReporter.ReportError(e)
		return
	}
}

var linebreakPattern = regexp.MustCompile("\r?\n")

func AddMissing(body string, findings []string) string {
	lines := linebreakPattern.Split(strings.Trim(strings.ReplaceAll(body, "```", ""), "\n\r"), -1)
	set := make(map[string]bool)
	for _, finding := range findings {
		// TODO this should also replace \r. But we'll neglect this corner case for now
		set[strings.ReplaceAll(finding, "\n", `\n`)] = true
	}
	for _, line := range lines {
		if set[line] {
			delete(set, line)
		}
	}
	for line := range set {
		// fmt.Printf("Appending: %#v\n", line)
		lines = append(lines, strings.ReplaceAll(line, "\n", `\n`))
	}
	// if len(set) > 0 {
	// 	fmt.Printf("%#v\n\n", body)
	// 	fmt.Printf("%#v\n\n", lines)
	// }
	return "```\n" + strings.Join(lines, "\n") + "\n```"
}
