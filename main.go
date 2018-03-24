package main

import (
	"net/http"
	"os"

	"go.uber.org/zap"

	"github.com/petergtz/alexa-wikipedia/mediawiki"
	"github.com/petergtz/alexa-wikipedia/wiki"
	"github.com/petergtz/go-alexa"
)

var (
	log *zap.SugaredLogger
)

func main() {
	l, e := zap.NewDevelopment()
	if e != nil {
		panic(e)
	}
	defer l.Sync()
	log = l.Sugar()

	handler := &alexa.Handler{
		Skill: &WikipediaSkill{
			wiki: &mediawiki.MediaWiki{},
		},
		Log: log,
		ExpectedApplicationID: os.Getenv("APPLICATION_ID"),
	}
	http.HandleFunc("/", handler.Handle)
	port := os.Getenv("PORT")
	if port == "" { // the port variable lets us distinguish between a local server an done in CF
		log.Debugf("Certificate path: %v", os.Getenv("cert"))
		log.Debugf("Private key path: %v", os.Getenv("key"))
		e = http.ListenAndServeTLS("0.0.0.0:4443", os.Getenv("cert"), os.Getenv("key"), nil)
	} else {
		e = http.ListenAndServe("0.0.0.0:"+port, nil)
	}
	log.Fatal(e)
}

type WikipediaSkill struct {
	wiki wiki.Wiki
}

const helpText = "Um einen Artikel vorgelesen zu bekommen, " +
	"sage z.B. \"Suche nach Käsekuchen.\" oder \"Was ist Käsekuchen?\". " +
	"Du kannst jederzeit zum Inhaltsverzeichnis springen, indem Du \"Inhaltsverzeichnis\" sagst. " +
	"Oder sage \"Springe zu Abschnitt 3.2\", um direkt zu diesem Abschnitt zu springen."

const quickHelpText = "Suche zunächst nach einem Begriff. " +
	"Sage z.B. \"Suche nach Käsekuchen.\" oder \"Was ist Käsekuchen?\"."

func (h *WikipediaSkill) ProcessRequest(requestEnv *alexa.RequestEnvelope) *alexa.ResponseEnvelope {
	wiki := &mediawiki.MediaWiki{}

	log.Infow("Request", "Type", requestEnv.Request.Type, "Intent", requestEnv.Request.Intent,
		"SessionAttributes", requestEnv.Session.Attributes)
	switch requestEnv.Request.Type {

	case "LaunchRequest":
		return &alexa.ResponseEnvelope{Version: "1.0",
			Response: &alexa.Response{
				OutputSpeech: plainText("Du befindest Dich jetzt bei Wikipedia. " + helpText),
			},
			SessionAttributes: map[string]interface{}{
				"last_question": "none",
			},
		}

	case "IntentRequest":
		intent := requestEnv.Request.Intent
		switch intent.Name {
		case "DefineIntent":
			page, e := wiki.GetPage(intent.Slots["word"].Value)
			if e != nil {
				if e.Error() == "Page not found on Wikipedia" {
					return &alexa.ResponseEnvelope{Version: "1.0",
						Response: &alexa.Response{
							OutputSpeech: plainText("Diesen Begriff konnte ich bei Wikipedia leider nicht finden. Versuche es doch mit einem anderen Begriff."),
						},
					}
				}
				log.Errorw("Could not get Wikipedia page", "error", e)
				return internalError()
			}
			return &alexa.ResponseEnvelope{Version: "1.0",
				Response: &alexa.Response{
					OutputSpeech: plainText(page.Body +
						" Zur weiteren Navigation kannst Du jederzeit zum Inhaltsverzeichnis springen" +
						" indem Du \"Inhaltsverzeichnis\" oder \"nächster Abschnitt\" sagst. " +
						"Soll ich zunächst einfach weiterlesen?"),
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
					Response:          &alexa.Response{OutputSpeech: plainText("Wie meinen?")},
					SessionAttributes: requestEnv.Session.Attributes,
				}
			}
			page, resp := h.pageFromSession(requestEnv.Session)
			if resp != nil {
				return resp
			}
			newPosition := int(requestEnv.Session.Attributes["position"].(float64)) + 1
			return &alexa.ResponseEnvelope{Version: "1.0",
				Response: &alexa.Response{
					OutputSpeech: plainText(page.TextForPosition(newPosition) +
						" Soll ich noch weiterlesen?"),
				},
				SessionAttributes: map[string]interface{}{
					"word":          requestEnv.Session.Attributes["word"],
					"position":      newPosition,
					"last_question": "should_continue",
				},
			}
		case "AMAZON.RepeatIntent":
			page, resp := h.pageFromSession(requestEnv.Session)
			if resp != nil {
				return resp
			}
			newPosition := int(requestEnv.Session.Attributes["position"].(float64))
			return &alexa.ResponseEnvelope{Version: "1.0",
				Response: &alexa.Response{
					OutputSpeech: plainText(page.TextForPosition(newPosition) +
						" Soll ich noch weiterlesen?"),
				},
				SessionAttributes: map[string]interface{}{
					"word":          requestEnv.Session.Attributes["word"],
					"position":      newPosition,
					"last_question": "should_continue",
				},
			}
		case "AMAZON.NextIntent":
			page, resp := h.pageFromSession(requestEnv.Session)
			if resp != nil {
				return resp
			}
			newPosition := int(requestEnv.Session.Attributes["position"].(float64)) + 1
			return &alexa.ResponseEnvelope{Version: "1.0",
				Response: &alexa.Response{
					OutputSpeech: plainText(page.TextForPosition(newPosition) +
						" Soll ich noch weiterlesen?"),
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
					Response:          &alexa.Response{OutputSpeech: plainText("Wie meinen?")},
					SessionAttributes: requestEnv.Session.Attributes,
				}
			}
			delete(requestEnv.Session.Attributes, "last_question")
			return &alexa.ResponseEnvelope{Version: "1.0",
				Response: &alexa.Response{
					OutputSpeech:     plainText("Nein? Okay."),
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
			page, resp := h.pageFromSession(requestEnv.Session)
			if resp != nil {
				return resp
			}
			requestEnv.Session.Attributes["last_question"] = "jump_where"
			return &alexa.ResponseEnvelope{Version: "1.0",
				Response: &alexa.Response{
					OutputSpeech: plainText(page.Toc() + " Zu welchem Abschnitt möchtest Du springen?"),
				},
				SessionAttributes: requestEnv.Session.Attributes,
			}
		case "GoToSectionIntent":
			page, resp := h.pageFromSession(requestEnv.Session)
			if resp != nil {
				return resp
			}
			sectionTitleOrNumber := intent.Slots["section_title_or_number"].Value
			s, position := page.TextAndPositionFromSectionNumber(sectionTitleOrNumber)
			if s == "" {
				s, position = page.TextAndPositionFromSectionName(sectionTitleOrNumber)
			}
			lastQuestion := ""
			if s != "" {
				s += "Soll ich noch weiterlesen?"
				lastQuestion = "should_continue"
			} else {
				s = "Ich konnte den angegebenen Abschnitt \"" + sectionTitleOrNumber + "\" nicht finden."
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
					OutputSpeech: plainText(helpText),
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

func lastQuestionIn(session *alexa.Session) string {
	if session.Attributes["last_question"] == nil {
		return ""
	}
	return session.Attributes["last_question"].(string)
}

func (h *WikipediaSkill) pageFromSession(session *alexa.Session) (wiki.Page, *alexa.ResponseEnvelope) {
	if !wordIn(session) {
		return wiki.Page{}, quickHelp(session.Attributes)
	}

	page, e := h.wiki.GetPage(session.Attributes["word"].(string))
	if e != nil {
		log.Errorw("Could not get Wikipedia page", "error", e)
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
