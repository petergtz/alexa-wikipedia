package skill

import (
	"strings"
	"time"

	"github.com/petergtz/alexa-wikipedia/bodychoppers/dumb"
	"github.com/petergtz/alexa-wikipedia/bodychoppers/paragraph"
	"github.com/petergtz/alexa-wikipedia/locale"
	"go.uber.org/zap"

	"github.com/nicksnyder/go-i18n/v2/i18n"
	. "github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/petergtz/alexa-wikipedia/wiki"
	"github.com/petergtz/go-alexa"
)

type WikipediaSkill struct {
	wiki               wiki.Wiki
	i18nBundle         *i18n.Bundle
	interactionLogger  alexa.InteractionLogger
	interactionHistory alexa.InteractionHistory
	bodyChopper        wiki.BodyChopper
	logger             *zap.SugaredLogger
}

func NewWikipediaSkill(
	wiki wiki.Wiki,
	i18nBundle *i18n.Bundle,
	interactionLogger alexa.InteractionLogger,
	interactionHistory alexa.InteractionHistory,
	logger *zap.SugaredLogger,
) *WikipediaSkill {
	return &WikipediaSkill{
		i18nBundle:         i18nBundle,
		wiki:               wiki,
		interactionLogger:  interactionLogger,
		interactionHistory: interactionHistory,
		bodyChopper: &paragraph.BodyChopper{
			MaxBodyPartLen: 6000,
			Fallback: dumb.BodyChopper{
				MaxBodyPartLen: 6000,
			},
		},
		logger: logger,
	}
}

func (h *WikipediaSkill) ProcessRequest(requestEnv *alexa.RequestEnvelope) *alexa.ResponseEnvelope {
	logger := h.logger.With("alexa-request-id", requestEnv.Request.RequestID)
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
			definition, e := h.findDefinition(intent.Slots["word"].Value, l)
			if e != nil {
				logger.Errorw("Could not get Wikipedia page", "error", e)
				return internalError(l)
			}
			if definition == nil {
				h.interactionLogger.Log(alexa.InteractionFrom(requestEnv).WithAttributes(map[string]interface{}{
					"Intent":      intent.Name,
					"SearchQuery": intent.Slots["word"].Value,
					"ActualTitle": "NOT_FOUND",
				}))

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
			}
			if titleWasAlreadyRecentlyFound(definition.Title, h.interactionHistory.GetInteractionsByUser(requestEnv.Session.User.UserID)) {
				return &alexa.ResponseEnvelope{Version: "1.0",
					Response: &alexa.Response{
						OutputSpeech: plainText(
							l.MustLocalize(&LocalizeConfig{
								DefaultMessage: &Message{
									ID: "SpellingHint",
									Other: "Ich habe den Artikel, \"{{.Title}}\", gerade erst gelesen. " +
										"Falls ich nicht Deinen gewünschten Artikel gefunden habe, unterbrich mich und sage: " +
										"\"Alexa, Suche buchstabieren\", um Deine Suchanfrage zu buchstabieren. Hier ist der Artikel:",
								},
								TemplateData: map[string]string{"Title": definition.Title},
							}) + "\n\n" +
								strings.TrimRight(h.bodyChopper.FetchBodyPart(definition.Body, 0), ". ") + ". " +
								l.MustLocalize(&LocalizeConfig{DefaultMessage: &Message{
									ID: "FurtherNavigationHints",
									Other: "Zur weiteren Navigation kannst Du jederzeit zum Inhaltsverzeichnis springen" +
										" indem Du \"Inhaltsverzeichnis\" oder \"nächster Abschnitt\" sagst. " +
										"Soll ich zunächst einfach weiterlesen?",
								}})),
					},
					SessionAttributes: map[string]interface{}{
						"word":                         intent.Slots["word"].Value,
						"position":                     0,
						"position_within_section_body": 0,
						"last_question":                "should_continue",
					},
				}
			}
			h.interactionLogger.Log(alexa.InteractionFrom(requestEnv).WithAttributes(map[string]interface{}{
				"Intent":      intent.Name,
				"SearchQuery": intent.Slots["word"].Value,
				"ActualTitle": definition.Title,
			}))
			return &alexa.ResponseEnvelope{Version: "1.0",
				Response: &alexa.Response{
					OutputSpeech: plainText(strings.TrimRight(h.bodyChopper.FetchBodyPart(definition.Body, 0), ". ") + ". " + l.MustLocalize(&LocalizeConfig{DefaultMessage: &Message{
						ID: "FurtherNavigationHints",
						Other: "Zur weiteren Navigation kannst Du jederzeit zum Inhaltsverzeichnis springen" +
							" indem Du \"Inhaltsverzeichnis\" oder \"nächster Abschnitt\" sagst. " +
							"Soll ich zunächst einfach weiterlesen?",
					}})),
				},
				SessionAttributes: map[string]interface{}{
					"word":                         intent.Slots["word"].Value,
					"position":                     0,
					"position_within_section_body": 0,
					"last_question":                "should_continue",
				},
			}
		case "SpellIntent":
			assembledSearchQuery := l.AssembleTermFromSpelling(intent.Slots["spelled_term"].Value)
			definition, e := h.findDefinition(assembledSearchQuery, l)
			if e != nil {
				logger.Errorw("Could not get Wikipedia page", "error", e)
				return internalError(l)
			}
			if definition == nil {
				h.interactionLogger.Log(alexa.InteractionFrom(requestEnv).WithAttributes(map[string]interface{}{
					"Intent":               intent.Name,
					"SpelledSearchQuery":   intent.Slots["spelled_term"].Value,
					"AssembledSearchQuery": assembledSearchQuery,
					"ActualTitle":          "NOT_FOUND",
				}))

				return &alexa.ResponseEnvelope{Version: "1.0",
					Response: &alexa.Response{
						OutputSpeech: plainText(
							l.MustLocalize(&LocalizeConfig{
								DefaultMessage: &Message{
									ID:    "CouldNotFindSpelledTerm",
									Other: "Den buchstabierten Begriff {{.SpelledTerm}} konnte ich bei Wikipedia leider nicht finden. Versuche es doch mit einem anderen Begriff.",
								},
								TemplateData: map[string]string{"SpelledTerm": intent.Slots["spelled_term"].Value},
							}),
						),
					},
				}
			}
			h.interactionLogger.Log(alexa.InteractionFrom(requestEnv).WithAttributes(map[string]interface{}{
				"Intent":               intent.Name,
				"SpelledSearchQuery":   intent.Slots["spelled_term"].Value,
				"AssembledSearchQuery": assembledSearchQuery,
				"ActualTitle":          definition.Title,
			}))
			return &alexa.ResponseEnvelope{Version: "1.0",
				Response: &alexa.Response{
					OutputSpeech: plainText(strings.TrimRight(h.bodyChopper.FetchBodyPart(definition.Body, 0), ". ") + ". " + l.MustLocalize(&LocalizeConfig{DefaultMessage: &Message{
						ID: "FurtherNavigationHints",
						Other: "Zur weiteren Navigation kannst Du jederzeit zum Inhaltsverzeichnis springen" +
							" indem Du \"Inhaltsverzeichnis\" oder \"nächster Abschnitt\" sagst. " +
							"Soll ich zunächst einfach weiterlesen?",
					}})),
				},
				SessionAttributes: map[string]interface{}{
					"word":                         assembledSearchQuery,
					"position":                     0,
					"position_within_section_body": 0,
					"last_question":                "should_continue",
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
			page, resp := h.pageFromSession(requestEnv.Session, l, logger)
			if resp != nil {
				return resp
			}

			newPosition, newPositionWithinSectionBody := h.bodyChopper.MoveToNextBodyPart(
				page.TextForPosition(int(requestEnv.Session.Attributes["position"].(float64))),
				int(requestEnv.Session.Attributes["position"].(float64)),
				int(requestEnv.Session.Attributes["position_within_section_body"].(float64)))

			bodyPart := h.bodyChopper.FetchBodyPart(page.TextForPosition(newPosition), newPositionWithinSectionBody)

			if bodyPart == "" {
				return &alexa.ResponseEnvelope{Version: "1.0",
					Response: &alexa.Response{
						OutputSpeech: plainText(l.MustLocalize(&LocalizeConfig{DefaultMessage: &Message{
							ID:    "EndOfArticle",
							Other: "Oh! Wir sind bereits am Ende angelangt. Wenn Du noch einen weiteren Artikel vorgelesen kriegen möchtest, sage z.B. \"Suche nach Elefant\"",
						}})),
					},
					SessionAttributes: map[string]interface{}{
						"word":                         requestEnv.Session.Attributes["word"],
						"position":                     newPosition,
						"position_within_section_body": newPositionWithinSectionBody,
					},
				}
			}

			return &alexa.ResponseEnvelope{Version: "1.0",
				Response: &alexa.Response{
					OutputSpeech: plainText(bodyPart + " " +
						l.MustLocalize(&LocalizeConfig{DefaultMessage: &Message{
							ID:    "ShouldIContinue",
							Other: "Soll ich noch weiterlesen?",
						}})),
				},
				SessionAttributes: map[string]interface{}{
					"word":                         requestEnv.Session.Attributes["word"],
					"position":                     newPosition,
					"position_within_section_body": newPositionWithinSectionBody,
					"last_question":                "should_continue",
				},
			}
		case "AMAZON.RepeatIntent":
			page, resp := h.pageFromSession(requestEnv.Session, l, logger)
			if resp != nil {
				return resp
			}
			return &alexa.ResponseEnvelope{Version: "1.0",
				Response: &alexa.Response{
					OutputSpeech: plainText(h.bodyChopper.FetchBodyPart(
						page.TextForPosition(int(requestEnv.Session.Attributes["position"].(float64))),
						int(requestEnv.Session.Attributes["position_within_section_body"].(float64))) +
						" " + l.MustLocalize(&LocalizeConfig{DefaultMessage: &Message{
						ID:    "ShouldIContinue",
						Other: "Soll ich noch weiterlesen?",
					}})),
				},
				SessionAttributes: map[string]interface{}{
					"word":                         requestEnv.Session.Attributes["word"],
					"position":                     requestEnv.Session.Attributes["position"],
					"position_within_section_body": requestEnv.Session.Attributes["position_within_section_body"],
					"last_question":                "should_continue",
				},
			}
		case "AMAZON.NextIntent":
			page, resp := h.pageFromSession(requestEnv.Session, l, logger)
			if resp != nil {
				return resp
			}
			newPosition, newPositionWithinSectionBody := h.bodyChopper.MoveToNextBodyPart(
				page.TextForPosition(int(requestEnv.Session.Attributes["position"].(float64))),
				int(requestEnv.Session.Attributes["position"].(float64)),
				int(requestEnv.Session.Attributes["position_within_section_body"].(float64)))

			bodyPart := h.bodyChopper.FetchBodyPart(page.TextForPosition(newPosition), newPositionWithinSectionBody)

			if bodyPart == "" {
				return &alexa.ResponseEnvelope{Version: "1.0",
					Response: &alexa.Response{
						OutputSpeech: plainText(l.MustLocalize(&LocalizeConfig{DefaultMessage: &Message{
							ID:    "EndOfArticle",
							Other: "Oh! wir sind bereits am Ende angelangt. Wenn Du noch einen weiteren Artikel vorgelesen kriegen möchtest, sage z.B. \"Suche nach Elefant\".",
						}})),
					},
					SessionAttributes: map[string]interface{}{
						"word":                         requestEnv.Session.Attributes["word"],
						"position":                     newPosition,
						"position_within_section_body": newPositionWithinSectionBody,
					},
				}
			}

			return &alexa.ResponseEnvelope{Version: "1.0",
				Response: &alexa.Response{
					OutputSpeech: plainText(bodyPart + " " +
						l.MustLocalize(&LocalizeConfig{DefaultMessage: &Message{
							ID:    "ShouldIContinue",
							Other: "Soll ich noch weiterlesen?",
						}})),
				},
				SessionAttributes: map[string]interface{}{
					"word":                         requestEnv.Session.Attributes["word"],
					"position":                     newPosition,
					"position_within_section_body": newPositionWithinSectionBody,
					"last_question":                "should_continue",
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
			page, resp := h.pageFromSession(requestEnv.Session, l, logger)
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
			page, resp := h.pageFromSession(requestEnv.Session, l, logger)
			if resp != nil {
				return resp
			}
			sectionTitleOrNumber := intent.Slots["section_title_or_number"].Value
			s, position := page.TextAndPositionFromSectionNumber(sectionTitleOrNumber, l)
			if s == "" {
				s, position = page.TextAndPositionFromSectionName(sectionTitleOrNumber, l)
			}
			s = h.bodyChopper.FetchBodyPart(s, 0)
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
					"word":                         requestEnv.Session.Attributes["word"],
					"position":                     position,
					"position_within_section_body": 0,
					"last_question":                lastQuestion,
				},
			}
		case "AMAZON.HelpIntent":
			return &alexa.ResponseEnvelope{Version: "1.0",
				Response: &alexa.Response{
					OutputSpeech: plainText(l.MustLocalize(&LocalizeConfig{DefaultMessage: &Message{
						ID: "HelpText",
						Other: "Um einen Artikel vorgelesen zu bekommen, " +
							"sage z.B. \"Suche nach Käsekuchen.\" oder \"Was ist Käsekuchen?\". " +
							"Du kannst jederzeit zum Inhaltsverzeichnis springen, indem Du \"Inhaltsverzeichnis\" sagst. " +
							"Oder sage \"Springe zu Abschnitt 3.2\", um direkt zu diesem Abschnitt zu springen.",
					}})),
				},
			}
		case "AMAZON.FallbackIntent":
			return &alexa.ResponseEnvelope{Version: "1.0",
				Response: &alexa.Response{
					OutputSpeech: plainText(l.MustLocalize(&LocalizeConfig{DefaultMessage: &Message{
						ID:    "FallbackText",
						Other: "Meine Enzyklopädie kann hiermit nicht weiterhelfen. Aber Du kannst z.B. sagen \"Suche nach Käsekuchen\".",
					}})),
				},
			}
		case "AMAZON.CancelIntent", "AMAZON.StopIntent":
			return &alexa.ResponseEnvelope{Version: "1.0",
				Response: &alexa.Response{ShouldSessionEnd: true},
			}
		default:
			return internalError(l)
		}

	case "SessionEndedRequest":
		return &alexa.ResponseEnvelope{Version: "1.0"}

	default:
		return &alexa.ResponseEnvelope{Version: "1.0"}
	}
}

const recentlyFoundThreshold = 20 * time.Second

func titleWasAlreadyRecentlyFound(currentTitle string, userInteractions []*alexa.Interaction) bool {
	for i := len(userInteractions) - 1; i >= 0; i-- {
		if userInteractions[i].RequestType == "IntentRequest" &&
			userInteractions[i].Attributes["Intent"] == "DefineIntent" &&
			userInteractions[i].Attributes["ActualTitle"] == currentTitle &&
			time.Now().Sub(userInteractions[i].Timestamp) < recentlyFoundThreshold {

			return true
		}
	}
	return false
}

func (h *WikipediaSkill) findDefinition(word string, l *locale.Localizer) (*wiki.Page, error) {
	page, e := h.wiki.GetPage(word, l)
	switch {
	case isNotFoundError(e):
		page, e = h.wiki.SearchPage(word, l)
		switch {
		case isNotFoundError(e):
			return nil, nil
		case e != nil:
			return nil, e
		default:
			return &page, nil
		}
	case e != nil:
		return nil, e
	default:
		return &page, nil
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

func (h *WikipediaSkill) pageFromSession(session *alexa.Session, l *locale.Localizer, logger *zap.SugaredLogger) (wiki.Page, *alexa.ResponseEnvelope) {
	if !wordIn(session) {
		return wiki.Page{}, quickHelp(session.Attributes, l)
	}

	page, e := h.wiki.GetPage(session.Attributes["word"].(string), l)
	switch {
	case isNotFoundError(e):
		page, e = h.wiki.SearchPage(session.Attributes["word"].(string), l)
		switch {
		case isNotFoundError(e):
			return wiki.Page{}, &alexa.ResponseEnvelope{Version: "1.0",
				Response: &alexa.Response{
					OutputSpeech: plainText("Diesen Begriff konnte ich bei Wikipedia leider nicht finden. Versuche es doch mit einem anderen Begriff."),
				},
			}
		case e != nil:
			logger.Errorw("Could not get Wikipedia page", "error", e)
			return wiki.Page{}, internalError(l)
		}
	case e != nil:
		logger.Errorw("Could not get Wikipedia page", "error", e)
		return wiki.Page{}, internalError(l)
	}
	return page, nil
}

func quickHelp(sessionAttributes map[string]interface{}, l *locale.Localizer) *alexa.ResponseEnvelope {
	return &alexa.ResponseEnvelope{Version: "1.0",
		Response: &alexa.Response{OutputSpeech: plainText(l.MustLocalize(&LocalizeConfig{
			DefaultMessage: &Message{
				ID: "QuickHelpText",
				Other: "Suche zunächst nach einem Begriff. " +
					"Sage z.B. \"Suche nach Käsekuchen.\" oder \"Was ist Käsekuchen?\".",
			},
		}))},
		SessionAttributes: sessionAttributes,
	}
}

func wordIn(session *alexa.Session) bool {
	return session.Attributes["word"] != nil
}

func plainText(text string) *alexa.OutputSpeech {
	return &alexa.OutputSpeech{Type: "PlainText", Text: text}
}

func internalError(l *locale.Localizer) *alexa.ResponseEnvelope {
	return &alexa.ResponseEnvelope{Version: "1.0",
		Response: &alexa.Response{
			OutputSpeech: plainText(l.MustLocalize(&LocalizeConfig{DefaultMessage: &Message{
				ID:    "InternalError",
				Other: "Es ist ein interner Fehler aufgetreten bei der Benutzung von Wikipedia.",
			}})),
			ShouldSessionEnd: false,
		},
	}
}
