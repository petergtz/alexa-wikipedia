package main

import (
	"bytes"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"go.uber.org/zap"

	"github.com/BurntSushi/toml"
	"github.com/petergtz/alexa-wikipedia/locale"

	"golang.org/x/text/language"

	"github.com/nicksnyder/go-i18n/v2/i18n"
	. "github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/petergtz/alexa-wikipedia/mediawiki"
	"github.com/petergtz/alexa-wikipedia/wiki"
	"github.com/petergtz/go-alexa"
)

var (
	logger *zap.SugaredLogger
)

func main() {
	l := createLoggerWith("debug")
	defer l.Sync()
	logger = l.Sugar()

	i18nBundle := &i18n.Bundle{DefaultLanguage: language.English}
	i18nBundle.RegisterUnmarshalFunc("toml", toml.Unmarshal)

	i18nBundle.MustParseMessageFileBytes(locale.DeDe, "active.de.toml")
	i18nBundle.MustParseMessageFileBytes(locale.EnUs, "active.en.toml")

	var e error
	skipRequestValidation := false
	if os.Getenv("SKILL_SKIP_REQUEST_VALIDATION") != "" {
		skipRequestValidation, e = strconv.ParseBool(os.Getenv("SKILL_SKIP_REQUEST_VALIDATION"))
		if e != nil {
			logger.Fatalw("Invalid env var SKILL_SKIP_REQUEST_VALIDATION", "value", os.Getenv("SKILL_SKIP_REQUEST_VALIDATION"))
		}
		if skipRequestValidation {
			logger.Info("Skipping request validation. THIS SHOULD ONLY BE USED IN TESTING")
		}
	}

	if os.Getenv("APPLICATION_ID") == "" {
		logger.Fatal("env var APPLICATION_ID not provided.")
	}

	handler := &alexa.Handler{
		Skill: &WikipediaSkill{
			i18nBundle: i18nBundle,
			wiki:       &mediawiki.MediaWiki{},
		},
		Log: logger,
		ExpectedApplicationID: os.Getenv("APPLICATION_ID"),
		SkipRequestValidation: skipRequestValidation,
	}
	serveMux := http.NewServeMux()
	serveMux.HandleFunc("/", handler.Handle)

	port := os.Getenv("PORT")
	if port == "" {
		logger.Fatal("No env variable PORT specified")
	}
	addr := os.Getenv("SKILL_ADDR")
	if addr == "" {
		addr = "0.0.0.0"
		logger.Infow("No SKILL_ADDR provided. Using default.", "addr", addr)
	} else {
		logger.Infow("SKILL_ADDR provided.", "addr", addr)
	}

	httpServer := &http.Server{
		Handler:      serveMux,
		Addr:         addr + ":" + port,
		WriteTimeout: 60 * time.Minute,
		ReadTimeout:  60 * time.Minute,
		ErrorLog:     NewStdLog(l),
	}

	if os.Getenv("SKILL_USE_TLS") == "true" {
		logger.Infof("Certificate path: %v", os.Getenv("CERT"))
		logger.Infof("Private key path: %v", os.Getenv("KEY"))
		e = httpServer.ListenAndServeTLS(os.Getenv("CERT"), os.Getenv("KEY"))
	} else {
		e = httpServer.ListenAndServe()
	}
	logger.Fatal(e)
}

type WikipediaSkill struct {
	wiki       wiki.Wiki
	i18nBundle *i18n.Bundle
}

const helpText = "Um einen Artikel vorgelesen zu bekommen, " +
	"sage z.B. \"Suche nach Käsekuchen.\" oder \"Was ist Käsekuchen?\". " +
	"Du kannst jederzeit zum Inhaltsverzeichnis springen, indem Du \"Inhaltsverzeichnis\" sagst. " +
	"Oder sage \"Springe zu Abschnitt 3.2\", um direkt zu diesem Abschnitt zu springen."

const quickHelpText = "Suche zunächst nach einem Begriff. " +
	"Sage z.B. \"Suche nach Käsekuchen.\" oder \"Was ist Käsekuchen?\"."

func (h *WikipediaSkill) ProcessRequest(requestEnv *alexa.RequestEnvelope) *alexa.ResponseEnvelope {
	logger.Infow("Request", "Type", requestEnv.Request.Type, "Intent", requestEnv.Request.Intent,
		"SessionAttributes", requestEnv.Session.Attributes, "locale", requestEnv.Request.Locale)

	l := locale.NewLocalizer(h.i18nBundle, requestEnv.Request.Locale, logger)

	switch requestEnv.Request.Type {

	case "LaunchRequest":
		return &alexa.ResponseEnvelope{Version: "1.0",
			Response: &alexa.Response{
				OutputSpeech: plainText(
					l.MustLocalize(&LocalizeConfig{DefaultMessage: &Message{
						ID: "YouAreAtWikipediaNow",
					}})),
			},
			SessionAttributes: map[string]interface{}{
				"last_question": "none",
			},
		}

	case "IntentRequest":
		intent := requestEnv.Request.Intent
		switch intent.Name {
		case "DefineIntent":
			var definition string
			page, e := h.wiki.GetPage(intent.Slots["word"].Value, l)
			switch {
			case isNotFoundError(e):
				page, e = h.wiki.SearchPage(intent.Slots["word"].Value, l)
				switch {
				case isNotFoundError(e):
					return &alexa.ResponseEnvelope{Version: "1.0",
						Response: &alexa.Response{
							OutputSpeech: plainText(
								l.MustLocalize(&LocalizeConfig{DefaultMessage: &Message{
									ID:    "CouldNotFindExpression",
									Other: "Diesen Begriff konnte ich bei Wikipedia leider nicht finden. Versuche es doch mit einem anderen Begriff.",
								}}),
							),
						},
					}
				case e != nil:
					logger.Errorw("Could not get Wikipedia page", "error", e)
					return internalError()
				default:
					definition = page.Body
				}
			case e != nil:
				logger.Errorw("Could not get Wikipedia page", "error", e)
				return internalError()
			default:
				definition = page.Body
			}
			return &alexa.ResponseEnvelope{Version: "1.0",
				Response: &alexa.Response{
					OutputSpeech: plainText(definition + " " + l.MustLocalize(&LocalizeConfig{DefaultMessage: &Message{
						ID: "FurtherNavigationHints",
						Other: "Zur weiteren Navigation kannst Du jederzeit zum Inhaltsverzeichnis springen" +
							" indem Du \"Inhaltsverzeichnis\" oder \"nächster Abschnitt\" sagst. " +
							"Soll ich zunächst einfach weiterlesen?",
					}})),
				},
				SessionAttributes: map[string]interface{}{
					"word":          intent.Slots["word"].Value,
					"position":      0,
					"last_question": "should_continue",
				},
			}
		case "AMAZON.YesIntent", "AMAZON.ResumeIntent":
			if lastQuestionIn(requestEnv.Session) != "should_continue" {
				return &alexa.ResponseEnvelope{Version: "1.0",
					Response: &alexa.Response{OutputSpeech: plainText(l.MustLocalize(&LocalizeConfig{DefaultMessage: &Message{
						ID:    "What",
						Other: "Wie meinen?",
					}}))},
					SessionAttributes: requestEnv.Session.Attributes,
				}
			}
			page, resp := h.pageFromSession(requestEnv.Session, l)
			if resp != nil {
				return resp
			}
			newPosition := int(requestEnv.Session.Attributes["position"].(float64)) + 1
			return &alexa.ResponseEnvelope{Version: "1.0",
				Response: &alexa.Response{
					OutputSpeech: plainText(page.TextForPosition(newPosition) + " " +
						l.MustLocalize(&LocalizeConfig{DefaultMessage: &Message{
							ID:    "ShouldIContinue",
							Other: "Soll ich noch weiterlesen?",
						}})),
				},
				SessionAttributes: map[string]interface{}{
					"word":          requestEnv.Session.Attributes["word"],
					"position":      newPosition,
					"last_question": "should_continue",
				},
			}
		case "AMAZON.RepeatIntent":
			page, resp := h.pageFromSession(requestEnv.Session, l)
			if resp != nil {
				return resp
			}
			newPosition := int(requestEnv.Session.Attributes["position"].(float64))
			return &alexa.ResponseEnvelope{Version: "1.0",
				Response: &alexa.Response{
					OutputSpeech: plainText(page.TextForPosition(newPosition) +
						" " + l.MustLocalize(&LocalizeConfig{DefaultMessage: &Message{
						ID:    "ShouldIContinue",
						Other: "Soll ich noch weiterlesen?",
					}})),
				},
				SessionAttributes: map[string]interface{}{
					"word":          requestEnv.Session.Attributes["word"],
					"position":      newPosition,
					"last_question": "should_continue",
				},
			}
		case "AMAZON.NextIntent":
			page, resp := h.pageFromSession(requestEnv.Session, l)
			if resp != nil {
				return resp
			}
			newPosition := int(requestEnv.Session.Attributes["position"].(float64)) + 1
			return &alexa.ResponseEnvelope{Version: "1.0",
				Response: &alexa.Response{
					OutputSpeech: plainText(page.TextForPosition(newPosition) + " " +
						l.MustLocalize(&LocalizeConfig{DefaultMessage: &Message{
							ID:    "ShouldIContinue",
							Other: "Soll ich noch weiterlesen?",
						}})),
				},
				SessionAttributes: map[string]interface{}{
					"word":          requestEnv.Session.Attributes["word"],
					"position":      newPosition,
					"last_question": "should_continue",
				},
			}
		case "AMAZON.NoIntent":
			if lastQuestionIn(requestEnv.Session) != "should_continue" {
				return &alexa.ResponseEnvelope{Version: "1.0",
					Response: &alexa.Response{OutputSpeech: plainText(l.MustLocalize(&LocalizeConfig{DefaultMessage: &Message{
						ID:    "What",
						Other: "Wie meinen?",
					}}))},
					SessionAttributes: requestEnv.Session.Attributes,
				}
			}
			delete(requestEnv.Session.Attributes, "last_question")
			return &alexa.ResponseEnvelope{Version: "1.0",
				Response: &alexa.Response{
					OutputSpeech: plainText(l.MustLocalize(&LocalizeConfig{DefaultMessage: &Message{
						ID:    "No?Okay",
						Other: "Nein? Okay.",
					}})),
					ShouldSessionEnd: true,
				},
				SessionAttributes: requestEnv.Session.Attributes,
			}
		case "AMAZON.PauseIntent":
			delete(requestEnv.Session.Attributes, "last_question")
			return &alexa.ResponseEnvelope{Version: "1.0",
				Response:          &alexa.Response{OutputSpeech: plainText(" ")},
				SessionAttributes: requestEnv.Session.Attributes,
			}
		case "TocIntent":
			page, resp := h.pageFromSession(requestEnv.Session, l)
			if resp != nil {
				return resp
			}
			requestEnv.Session.Attributes["last_question"] = "jump_where"
			return &alexa.ResponseEnvelope{Version: "1.0",
				Response: &alexa.Response{
					OutputSpeech: plainText(page.Toc(l) + " " + l.MustLocalize(&LocalizeConfig{DefaultMessage: &Message{
						ID:    "WhichSectionToJump",
						Other: "Zu welchem Abschnitt möchtest Du springen?",
					}})),
				},
				SessionAttributes: requestEnv.Session.Attributes,
			}
		case "GoToSectionIntent":
			page, resp := h.pageFromSession(requestEnv.Session, l)
			if resp != nil {
				return resp
			}
			sectionTitleOrNumber := intent.Slots["section_title_or_number"].Value
			s, position := page.TextAndPositionFromSectionNumber(sectionTitleOrNumber, l)
			if s == "" {
				s, position = page.TextAndPositionFromSectionName(sectionTitleOrNumber, l)
			}
			lastQuestion := ""
			if s != "" {
				s += " " + l.MustLocalize(&LocalizeConfig{DefaultMessage: &Message{
					ID:    "ShouldIContinue",
					Other: "Soll ich noch weiterlesen?",
				}})
				lastQuestion = "should_continue"
			} else {
				s = l.MustLocalize(&LocalizeConfig{
					DefaultMessage: &Message{
						ID:    "CouldNotFindSection",
						Other: "Ich konnte den angegebenen Abschnitt \"{{.SectionTitleOrNumber}}\" nicht finden.",
					},
					TemplateData: map[string]string{"SectionTitleOrNumber": sectionTitleOrNumber},
				})
				position = int(requestEnv.Session.Attributes["position"].(float64))
				lastQuestion = "none"
			}
			return &alexa.ResponseEnvelope{Version: "1.0",
				Response: &alexa.Response{OutputSpeech: plainText(s)},
				SessionAttributes: map[string]interface{}{
					"word":          requestEnv.Session.Attributes["word"],
					"position":      position,
					"last_question": lastQuestion,
				},
			}
		case "AMAZON.HelpIntent":
			return &alexa.ResponseEnvelope{Version: "1.0",
				Response: &alexa.Response{
					OutputSpeech: plainText(l.MustLocalize(&LocalizeConfig{DefaultMessage: &Message{
						ID:    "HelpText",
						Other: helpText,
					}})),
				},
			}
		case "AMAZON.CancelIntent", "AMAZON.StopIntent":
			return &alexa.ResponseEnvelope{Version: "1.0",
				Response: &alexa.Response{ShouldSessionEnd: true},
			}
		default:
			return internalError()
		}

	case "SessionEndedRequest":
		return &alexa.ResponseEnvelope{Version: "1.0"}

	default:
		return &alexa.ResponseEnvelope{Version: "1.0"}
	}
}

func createLoggerWith(logLevel string) *zap.Logger {
	loggerConfig := zap.NewProductionConfig()
	loggerConfig.Level = zapLogLevelFrom(logLevel)
	loggerConfig.DisableStacktrace = true
	logger, e := loggerConfig.Build()
	if e != nil {
		log.Panic(e)
	}
	return logger
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

func isNotFoundError(e error) bool {
	return e != nil && e.Error() == "Page not found on Wikipedia"
}

func lastQuestionIn(session *alexa.Session) string {
	if session.Attributes["last_question"] == nil {
		return ""
	}
	return session.Attributes["last_question"].(string)
}

func (h *WikipediaSkill) pageFromSession(session *alexa.Session, localizer *locale.Localizer) (wiki.Page, *alexa.ResponseEnvelope) {
	if !wordIn(session) {
		return wiki.Page{}, quickHelp(session.Attributes)
	}

	page, e := h.wiki.GetPage(session.Attributes["word"].(string), localizer)
	switch {
	case isNotFoundError(e):
		page, e = h.wiki.SearchPage(session.Attributes["word"].(string), localizer)
		switch {
		case isNotFoundError(e):
			return wiki.Page{}, &alexa.ResponseEnvelope{Version: "1.0",
				Response: &alexa.Response{
					OutputSpeech: plainText("Diesen Begriff konnte ich bei Wikipedia leider nicht finden. Versuche es doch mit einem anderen Begriff."),
				},
			}
		case e != nil:
			logger.Errorw("Could not get Wikipedia page", "error", e)
			return wiki.Page{}, internalError()
		}
	case e != nil:
		logger.Errorw("Could not get Wikipedia page", "error", e)
		return wiki.Page{}, internalError()
	}
	return page, nil
}

func quickHelp(sessionAttributes map[string]interface{}) *alexa.ResponseEnvelope {
	return &alexa.ResponseEnvelope{Version: "1.0",
		Response:          &alexa.Response{OutputSpeech: plainText(quickHelpText)},
		SessionAttributes: sessionAttributes,
	}
}

func wordIn(session *alexa.Session) bool {
	return session.Attributes["word"] != nil
}

func plainText(text string) *alexa.OutputSpeech {
	return &alexa.OutputSpeech{Type: "PlainText", Text: text}
}

func internalError() *alexa.ResponseEnvelope {
	return &alexa.ResponseEnvelope{Version: "1.0",
		Response: &alexa.Response{
			OutputSpeech:     plainText("Es ist ein interner Fehler aufgetreten bei der Benutzung von Wikipedia."),
			ShouldSessionEnd: false,
		},
	}
}

const (
	_stdLogDefaultDepth = 2
	_loggerWriterDepth  = 1
)

// Copied from go.uber.org/zap/global.go and changed to use Error instead of Info:
func NewStdLog(l *zap.Logger) *log.Logger {
	return log.New(&loggerWriter{l.WithOptions(
		zap.AddCallerSkip(_stdLogDefaultDepth + _loggerWriterDepth),
	)}, "" /* prefix */, 0 /* flags */)
}

type loggerWriter struct{ logger *zap.Logger }

func (l *loggerWriter) Write(p []byte) (int, error) {
	p = bytes.TrimSpace(p)
	l.logger.Error(string(p))
	return len(p), nil
}
