package factory

import (
	"net/http"
	"os"

	"github.com/aws/aws-sdk-go/service/sns"

	"github.com/petergtz/alexa-wikipedia/github"

	"github.com/BurntSushi/toml"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/petergtz/alexa-wikipedia/locale"
	"github.com/petergtz/alexa-wikipedia/mediawiki"
	"github.com/petergtz/alexa-wikipedia/skill"
	"golang.org/x/text/language"

	"go.uber.org/zap"

	awsdyndb "github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/petergtz/go-alexa"
	"github.com/petergtz/go-alexa/decorator"
	"github.com/petergtz/go-alexa/dynamodb"
)

func CreateSkill(logger *zap.SugaredLogger) *decorator.InteractionLoggingSkill {
	// this is a pure priming call to make subsequent calls faster
	http.Get("https://en.wikipedia.org/w/api.php?format=json&action=query&prop=extracts&titles=Keepalive&redirects=true&formatversion=2&explaintext=true&exlimit=1")

	tableName := "AlexaWikipediaRequests"
	if os.Getenv("TABLE_NAME_OVERRIDE") != "" {
		tableName = os.Getenv("TABLE_NAME_OVERRIDE")
		logger.Infow("Using DynamoDB table override", "table", tableName)
	}

	githubToken := os.Getenv("GITHUB_TOKEN")
	if githubToken == "" {
		logger.Fatal("GITHUB_TOKEN not set. Please set it to a valid token from Github.")
	}

	interactionLogger := dynamodb.NewInteractionLogger(
		awsdyndb.New(session.Must(session.NewSession(&aws.Config{Region: aws.String("eu-central-1")}))),
		logger,
		tableName)

	return decorator.ForSkillWithInteractionLogging(
		skill.NewWikipediaSkill(
			&mediawiki.MediaWiki{
				Logger: logger,
				WikiPagePreProcessor: mediawiki.NewHighlightMissingSpacesNaivelyWikiPagePreProcessor(
					github.NewGithubPersistence(
						"petergtz",
						"alexa-wikipedia",
						549126277,
						githubToken,
					),
					github.NewGithubErrorReporter(
						"petergtz",
						"alexa-wikipedia",
						githubToken,
						logger,
						"``fields @timestamp, @message | filter `error-id` = %v``",
						sns.New(session.Must(session.NewSession(&aws.Config{Region: aws.String("eu-west-1")}))),
						"arn:aws:sns:eu-west-1:512841817041:AlexaWikipediaErrors")),
			},
			createI18nBundle(),
			interactionLogger,
			interactionLogger,
			logger,
		),
		interactionLogger,
		func(requestEnv *alexa.RequestEnvelope) bool {
			return !(requestEnv.Request.Type == "IntentRequest" && requestEnv.Request.Intent.Name == "DefineIntent")
		},
	)
}

func createI18nBundle() *i18n.Bundle {
	i18nBundle := i18n.NewBundle(language.English)
	i18nBundle.RegisterUnmarshalFunc("toml", toml.Unmarshal)

	i18nBundle.MustParseMessageFileBytes(locale.DeDe, "active.de.toml")
	i18nBundle.MustParseMessageFileBytes(locale.EnUs, "active.en.toml")
	return i18nBundle
}
