package factory

import (
	"log"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/petergtz/alexa-wikipedia/mediawiki"
	"github.com/petergtz/alexa-wikipedia/skill"

	"go.uber.org/zap"

	"github.com/BurntSushi/toml"
	"github.com/petergtz/alexa-wikipedia/locale"

	"golang.org/x/text/language"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	awsdyndb "github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/petergtz/go-alexa"
	"github.com/petergtz/go-alexa/decorator"
	"github.com/petergtz/go-alexa/dynamodb"
)

func CreateSkill(logger *zap.SugaredLogger) *decorator.InteractionLoggingSkill {
	// this is a pure priming call to make subsequent calls faster
	http.Get("https://en.wikipedia.org/w/api.php?format=json&action=query&prop=extracts&titles=Keepalive&redirects=true&formatversion=2&explaintext=true&exlimit=1")

	interactionLogger := CreateInteractionLogger(logger)
	return decorator.ForSkillWithInteractionLogging(
		skill.NewWikipediaSkill(
			&mediawiki.MediaWiki{
				Logger: logger,
			},
			CreateI18nBundle(),
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

func CreateInteractionLogger(logger *zap.SugaredLogger) *dynamodb.RequestLogger {
	tableName := "AlexaWikipediaRequests"
	if os.Getenv("TABLE_NAME_OVERRIDE") != "" {
		tableName = os.Getenv("TABLE_NAME_OVERRIDE")
		logger.Infow("Using DynamoDB table override", "table", tableName)
	}
	il := dynamodb.NewInteractionLogger(
		awsdyndb.New(session.Must(session.NewSession(&aws.Config{Region: aws.String("eu-central-1")}))),
		logger,
		tableName)
	go il.GetInteractionsByUser("thisisjustaprimingcall")
	return il
}

func CreateI18nBundle() *i18n.Bundle {
	i18nBundle := i18n.NewBundle(language.English)
	i18nBundle.RegisterUnmarshalFunc("toml", toml.Unmarshal)

	i18nBundle.MustParseMessageFileBytes(locale.DeDe, "active.de.toml")
	i18nBundle.MustParseMessageFileBytes(locale.EnUs, "active.en.toml")
	return i18nBundle
}

func CreateLoggerWith(logLevel string) *zap.SugaredLogger {
	loggerConfig := zap.NewProductionConfig()
	loggerConfig.Level = zapLogLevelFrom(logLevel)
	loggerConfig.DisableStacktrace = true
	logger, e := loggerConfig.Build()
	if e != nil {
		log.Panic(e)
	}
	rand.Seed(time.Now().UnixNano())
	return logger.Sugar().With("function-instance-id", rand.Int63())
}

func zapLogLevelFrom(configLogLevel string) zap.AtomicLevel {
	switch strings.ToLower(configLogLevel) {
	case "", "debug":
		return zap.NewAtomicLevelAt(zap.DebugLevel)
	case "info":
		return zap.NewAtomicLevelAt(zap.InfoLevel)
	case "warn":
		return zap.NewAtomicLevelAt(zap.WarnLevel)
	case "error":
		return zap.NewAtomicLevelAt(zap.ErrorLevel)
	case "fatal":
		return zap.NewAtomicLevelAt(zap.FatalLevel)
	default:
		log.Fatal("Invalid log level in config", "log-level", configLogLevel)
		return zap.NewAtomicLevelAt(-1)
	}
}
