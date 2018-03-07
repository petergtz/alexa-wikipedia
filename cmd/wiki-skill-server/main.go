package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"

	"go.uber.org/zap"

	"github.com/petergtz/alexa-wikipedia/mediawiki"
)

var (
	log                   *zap.SugaredLogger
	expectedApplicationID = os.Getenv("APPLICATION_ID")
)

func main() {
	l, e := zap.NewDevelopment()
	if e != nil {
		panic(e)
	}
	defer l.Sync()
	log = l.Sugar()

	http.HandleFunc("/", handler)
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

func handler(w http.ResponseWriter, req *http.Request) {
	requestBody, e := ioutil.ReadAll(req.Body)
	if e != nil {
		panic(e)
	}
	var alexaRequest RequestEnvelope
	e = json.Unmarshal(requestBody, &alexaRequest)
	if e != nil {
		panic(e)
	}
	if alexaRequest.Session == nil {
		panic("Empty session")
	}
	output, e := json.Marshal(processRequest(&alexaRequest))
	if e != nil {
		panic(e)
	}
	w.Header().Set("Content-Type", "text/plain")
	w.Write(output)
}

func processRequest(requestEnv *RequestEnvelope) *ResponseEnvelope {
	if requestEnv.Session.Application.ApplicationID != expectedApplicationID {
		log.Fatalf("ApplicationID does not match: %v", requestEnv.Session.Application.ApplicationID)
		return internalError()
	}

	wiki := &mediawiki.MediaWiki{}

	switch requestEnv.Request.Type {

	case "LaunchRequest":
		return &ResponseEnvelope{Version: "1.0",
			Response: &Response{
				OutputSpeech: plainText("Du befindest Dich jetzt bei Wikipedia. Um einen Artikel vorgelesen zu bekommen, " +
					"sage z.B. \"Suche nach Käsekuchen.\" oder \"Was ist Käsekuchen?\"."),
			},
		}

	case "IntentRequest":
		intent := requestEnv.Request.Intent
		switch intent.Name {
		case "define":
			log.Infof("Slot word: %v", intent.Slots["word"].Value)

			page, e := wiki.GetPage(intent.Slots["word"].Value)
			if e != nil {
				log.Errorw("Could not get Wikipedia page", "error", e)
				return internalError()
			}
			return &ResponseEnvelope{Version: "1.0",
				Response: &Response{
					OutputSpeech:     plainText(page.Body + " Soll ich einfach weiterlesen oder soll ich zunächst das Inhaltsverzeichnis vorlesen?"),
					ShouldSessionEnd: false,
				},
				SessionAttributes: map[string]interface{}{
					"word":     intent.Slots["word"].Value,
					"position": 0,
				},
			}
		case "AMAZON.ResumeIntent":
			page, e := wiki.GetPage(requestEnv.Session.Attributes["word"].(string))
			if e != nil {
				log.Errorw("Could not get Wikipedia page", "error", e)
				return internalError()
			}
			position := int(requestEnv.Session.Attributes["position"].(float64))
			s := page.TextForPosition(position)
			return &ResponseEnvelope{Version: "1.0",
				Response: &Response{
					OutputSpeech:     plainText(s),
					ShouldSessionEnd: false,
				},
				SessionAttributes: map[string]interface{}{
					"word":     requestEnv.Session.Attributes["word"],
					"position": position + 1,
				},
			}
		case "toc":
			page, e := wiki.GetPage(requestEnv.Session.Attributes["word"].(string))
			if e != nil {
				log.Errorw("Could not get Wikipedia page", "error", e)
				return internalError()
			}
			s := ""
			for i, section := range page.Subsections {
				s += fmt.Sprintf("Abschnitt %v: %v.\n", i+1, section.Title)
			}
			return &ResponseEnvelope{Version: "1.0",
				Response: &Response{
					OutputSpeech:     plainText(s),
					ShouldSessionEnd: false,
				},
				SessionAttributes: requestEnv.Session.Attributes,
			}
		case "goto_section":
			page, e := wiki.GetPage(requestEnv.Session.Attributes["word"].(string))
			if e != nil {
				log.Errorw("Could not get Wikipedia page", "error", e)
				return internalError()
			}
			sectionTitleOrNumber := intent.Slots["section_title_or_number"].Value
			s := ""
			sectionNumber, e := strconv.Atoi(sectionTitleOrNumber)
			if e == nil {
				s = page.Subsections[sectionNumber].Body
				if s == "" && len(page.Subsections[sectionNumber].Subsections) > 0 {
					s = page.Subsections[sectionNumber].Subsections[0].Body
				}
			} else {
				for _, section := range page.Subsections {
					if strings.ToLower(section.Title) == strings.ToLower(sectionTitleOrNumber) {
						s = section.Body
						break
					}
				}
			}
			if s == "" {
				s = "Ich konnte den angegebenen Abschnitt \"" + sectionTitleOrNumber + "\" nicht finden."
			}
			return &ResponseEnvelope{Version: "1.0",
				Response: &Response{
					OutputSpeech:     plainText(s),
					ShouldSessionEnd: false,
				},
				SessionAttributes: map[string]interface{}{
					"word":     requestEnv.Session.Attributes["word"],
					"position": sectionNumber,
				},
			}
		case "AMAZON.HelpIntent":
			return &ResponseEnvelope{Version: "1.0",
				Response: &Response{
					OutputSpeech: plainText("Hier muss noch ein Hilfetext her"),
				},
			}
		case "AMAZON.CancelIntent":
			return &ResponseEnvelope{Version: "1.0",
				Response: &Response{
					OutputSpeech:     plainText("Hier muss noch ein Text her."),
					ShouldSessionEnd: true,
				},
			}
		case "AMAZON.StopIntent":
			return &ResponseEnvelope{Version: "1.0",
				Response: &Response{
					OutputSpeech:     plainText("Hier muss noch ein Text her."),
					ShouldSessionEnd: true,
				},
			}
		case "AMAZON.YesIntent":
			log.Debugf("%#v", requestEnv.Session.Attributes)
			return &ResponseEnvelope{Version: "1.0",
				Response: &Response{
					OutputSpeech: plainText(fmt.Sprintf("Ich würde jetzt weiterlesen bei Wort %v und Position  %v",

						requestEnv.Session.Attributes["word"], requestEnv.Session.Attributes["position"])),
				},
			}
		default:
			return internalError()
		}

	case "SessionEndedRequest":
		return &ResponseEnvelope{Version: "1.0"}

	default:
		return &ResponseEnvelope{Version: "1.0"}
	}

}

type RequestEnvelope struct {
	Version string   `json:"version"`
	Session *Session `json:"session"`
	Request *Request `json:"request"`
	// TODO Add Request Context
}

// Session containes the session data from the Alexa request.
type Session struct {
	New        bool                   `json:"new"`
	SessionID  string                 `json:"sessionId"`
	Attributes map[string]interface{} `json:"attributes"`
	User       struct {
		UserID      string `json:"userId"`
		AccessToken string `json:"accessToken"`
	} `json:"user"`
	Application struct {
		ApplicationID string `json:"applicationId"`
	} `json:"application"`
}

// Request contines the data in the request within the main request.
type Request struct {
	Locale      string `json:"locale"`
	Timestamp   string `json:"timestamp"`
	Type        string `json:"type"`
	RequestID   string `json:"requestId"`
	DialogState string `json:"dialogState"`
	Intent      Intent `json:"intent"`
	Name        string `json:"name"`
}

// Intent contains the data about the Alexa Intent requested.
type Intent struct {
	Name               string                `json:"name"`
	ConfirmationStatus string                `json:"confirmationStatus,omitempty"`
	Slots              map[string]IntentSlot `json:"slots"`
}

// IntentSlot contains the data for one Slot
type IntentSlot struct {
	Name               string `json:"name"`
	ConfirmationStatus string `json:"confirmationStatus,omitempty"`
	Value              string `json:"value"`
}

// ResponseEnvelope contains the Response and additional attributes.
type ResponseEnvelope struct {
	Version           string                 `json:"version"`
	SessionAttributes map[string]interface{} `json:"sessionAttributes,omitempty"`
	Response          *Response              `json:"response"`
}

// Response contains the body of the response.
type Response struct {
	OutputSpeech     *OutputSpeech `json:"outputSpeech,omitempty"`
	Card             *Card         `json:"card,omitempty"`
	Reprompt         *Reprompt     `json:"reprompt,omitempty"`
	Directives       []interface{} `json:"directives,omitempty"`
	ShouldSessionEnd bool          `json:"shouldEndSession"`
}

// OutputSpeech contains the data the defines what Alexa should say to the user.
type OutputSpeech struct {
	Type string `json:"type"`
	Text string `json:"text,omitempty"`
	SSML string `json:"ssml,omitempty"`
}

// Card contains the data displayed to the user by the Alexa app.
type Card struct {
	Type    string `json:"type"`
	Title   string `json:"title,omitempty"`
	Content string `json:"content,omitempty"`
	Text    string `json:"text,omitempty"`
	Image   *Image `json:"image,omitempty"`
}

// Image provides URL(s) to the image to display in resposne to the request.
type Image struct {
	SmallImageURL string `json:"smallImageUrl,omitempty"`
	LargeImageURL string `json:"largeImageUrl,omitempty"`
}

// Reprompt contains data about whether Alexa should prompt the user for more data.
type Reprompt struct {
	OutputSpeech *OutputSpeech `json:"outputSpeech,omitempty"`
}

// AudioPlayerDirective contains device level instructions on how to handle the response.
type AudioPlayerDirective struct {
	Type         string     `json:"type"`
	PlayBehavior string     `json:"playBehavior,omitempty"`
	AudioItem    *AudioItem `json:"audioItem,omitempty"`
}

// AudioItem contains an audio Stream definition for playback.
type AudioItem struct {
	Stream Stream `json:"stream,omitempty"`
}

// Stream contains instructions on playing an audio stream.
type Stream struct {
	Token                string `json:"token"`
	URL                  string `json:"url"`
	OffsetInMilliseconds int    `json:"offsetInMilliseconds"`
}

// DialogDirective contains directives for use in Dialog prompts.
type DialogDirective struct {
	Type          string  `json:"type"`
	SlotToElicit  string  `json:"slotToElicit,omitempty"`
	SlotToConfirm string  `json:"slotToConfirm,omitempty"`
	UpdatedIntent *Intent `json:"updatedIntent,omitempty"`
}

func plainText(text string) *OutputSpeech {
	return &OutputSpeech{Type: "PlainText", Text: text}
}

func internalError() *ResponseEnvelope {
	return &ResponseEnvelope{Version: "1.0",
		Response: &Response{
			OutputSpeech:     plainText("Es ist ein interner Fehler aufgetreten bei der Benutzung von Wikipedia."),
			ShouldSessionEnd: false,
		},
	}
}
