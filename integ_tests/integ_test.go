package main_test

import (
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"testing"
	"time"

	"github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo"
	"github.com/onsi/gomega"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"
)

func TestEndToEnd(t *testing.T) {
	gomega.RegisterFailHandler(ginkgo.Fail)
	ginkgo.RunSpecs(t, "EndToEnd")
}

var _ = Describe("Skill", func() {

	var (
		session *gexec.Session
		client  *http.Client
	)

	BeforeSuite(func() {
		pathToWebserver, err := gexec.Build("github.com/petergtz/alexa-wikipedia")
		Ω(err).ShouldNot(HaveOccurred())

		os.Setenv("PORT", "4443")
		os.Setenv("SKILL_ADDR", "127.0.0.1")
		os.Setenv("SKILL_SKIP_REQUEST_VALIDATION", "true")
		os.Setenv("APPLICATION_ID", "xxx")

		session, err = gexec.Start(exec.Command(pathToWebserver), GinkgoWriter, GinkgoWriter)
		Ω(err).ShouldNot(HaveOccurred())
		time.Sleep(200 * time.Millisecond)
		Expect(session.ExitCode()).To(Equal(-1), "Webserver error message: %s", string(session.Err.Contents()))

		client = &http.Client{}
	})

	AfterSuite(func() {
		if session != nil {
			session.Kill()
		}
		gexec.CleanupBuildArtifacts()
	})

	Describe("LaunchRequest", func() {
		Context("locale: de-DE", func() {
			It("returns a StatusOK and a German welcome message", func() {
				response, e := client.Post("http://127.0.0.1:4443/", "", strings.NewReader(`{
					"version": "1.0",
					"session": {
					"new": true,
					"sessionId": "xxx",
					"application": {
						"applicationId": "xxx"
					},
					"user": {
						"userId": "xxx"
					}
					},
					"context": {
					"AudioPlayer": {
						"playerActivity": "IDLE"
					},
					"Display": {
						"token": ""
					},
					"System": {
						"application": {
						"applicationId": "xxx"
						},
						"user": {
						"userId": "xxx"
						},
						"device": {
						"deviceId": "xxx",
						"supportedInterfaces": {
							"AudioPlayer": {},
							"Display": {
							"templateVersion": "1.0",
							"markupVersion": "1.0"
							}
						}
						},
						"apiEndpoint": "https://api.eu.amazonalexa.com",
						"apiAccessToken": "xxx"
					}
					},
					"request": {
					"type": "LaunchRequest",
					"requestId": "xxx",
					"timestamp": "`+time.Now().UTC().Format("2006-01-02T15:04:05Z")+`",
					"locale": "de-DE"
					}
				}`))
				Expect(e).NotTo(HaveOccurred())

				Expect(response.StatusCode).To(Equal(http.StatusOK))

				Expect(ioutil.ReadAll(response.Body)).To(MatchJSON(`{
					"version": "1.0",
					"sessionAttributes": {
					"last_question": "none"
					},
					"response": {
					"outputSpeech": {
						"type": "PlainText",
						"text": "Du befindest Dich jetzt bei Wikipedia. Um einen Artikel vorgelesen zu bekommen, sage z.B. \"Suche nach Käsekuchen.\" oder \"Was ist Käsekuchen?\". Du kannst jederzeit zum Inhaltsverzeichnis springen, indem Du \"Inhaltsverzeichnis\" sagst. Oder sage \"Springe zu Abschnitt 3.2\", um direkt zu diesem Abschnitt zu springen."
					},
					"shouldEndSession": false
					}
				}`))
			})
		})
	})

	Describe("IntentRequest", func() {
		Context("TocIntent", func() {
			Context("locale: de-DE", func() {
				It("reads the table of contents", func() {
					response, e := client.Post("http://127.0.0.1:4443/", "", strings.NewReader(`{
						"version": "1.0",
						"session": {
							"new": false,
							"sessionId": "xxx",
							"application": {
								"applicationId": "xxx"
							},
							"attributes": {
								"position": 0,
								"word": "baum",
								"last_question": "none"
							},
							"user": {
								"userId": "xxx"
							}
						},
						"context": {
							"System": {
								"application": {
									"applicationId": "xxx"
								},
								"user": {
									"userId": "xxx"
								},
								"device": {
									"deviceId": "xxx",
									"supportedInterfaces": {}
								},
								"apiEndpoint": "https://api.eu.amazonalexa.com",
								"apiAccessToken": "xxx"
							}
						},
						"request": {
							"type": "IntentRequest",
							"requestId": "xxx",
							"timestamp": "`+time.Now().UTC().Format("2006-01-02T15:04:05Z")+`",
							"locale": "de-DE",
							"intent": {
								"name": "TocIntent",
								"confirmationStatus": "NONE"
							}
						}
					}`))

					Expect(e).NotTo(HaveOccurred())
					Expect(response.StatusCode).To(Equal(http.StatusOK))
					Expect(ioutil.ReadAll(response.Body)).To(MatchJSON(`{
						"version": "1.0",
						"response": {
							"outputSpeech": {
								"type": "PlainText",
								"text": "Inhaltsverzeichnis. Abschnitt 1: Etymologie.\nAbschnitt 2: Definition und taxonomische Verbreitung.\nAbschnitt 3: Die besonderen Merkmale der Bäume.\nAbschnitt 4: Entwicklung baumförmiger Pflanzen in der Erdgeschichte.\nAbschnitt 5: Physiologie.\nAbschnitt 6: Ökologie.\nAbschnitt 7: Bäume und Menschen.\nAbschnitt 8: Superlative.\nAbschnitt 9: Filmografie.\nAbschnitt 10: Literatur.\nAbschnitt 11: Weblinks.\nAbschnitt 12: Einzelnachweise.\n Zu welchem Abschnitt möchtest Du springen?"
							},
							"shouldEndSession": false
						},
						"sessionAttributes": {
							"position": 0,
							"word": "baum",
							"last_question": "jump_where"
						}
					}`))
				})
			})
		})

		Context("GoToSectionIntent", func() {
			Context("Title cannot be found", func() {
				Context("locale: de-DE", func() {
					It("returns a StatusOK and a German welcome message", func() {
						response, e := client.Post("http://127.0.0.1:4443/", "", strings.NewReader(`{
							"version": "1.0",
							"session": {
								"new": false,
								"sessionId": "xxx",
								"application": {
									"applicationId": "xxx"
								},
								"attributes": {
									"position": 0,
									"word": "baum",
									"last_question": "should_continue"
								},
								"user": {
									"userId": "xxx"
								}
							},
							"context": {
								"System": {
									"application": {
										"applicationId": "xxx"
									},
									"user": {
										"userId": "xxx"
									},
									"device": {
										"deviceId": "xxx",
										"supportedInterfaces": {}
									},
									"apiEndpoint": "https://api.eu.amazonalexa.com",
									"apiAccessToken": "xxx"
								}
							},
							"request": {
								"type": "IntentRequest",
								"requestId": "xxx",
								"timestamp": "`+time.Now().UTC().Format("2006-01-02T15:04:05Z")+`",
								"locale": "de-DE",
								"intent": {
									"name": "GoToSectionIntent",
									"confirmationStatus": "NONE",
									"slots": {
										"section_title_or_number": {
											"name": "section_title_or_number",
											"value": "bla",
											"confirmationStatus": "NONE"
										}
									}
								}
							}
						}`))

						Expect(e).NotTo(HaveOccurred())
						Expect(response.StatusCode).To(Equal(http.StatusOK))
						Expect(ioutil.ReadAll(response.Body)).To(MatchJSON(`{
							"version": "1.0",
							"response": {
								"outputSpeech": {
									"type": "PlainText",
									"text": "Ich konnte den angegebenen Abschnitt \"bla\" nicht finden."
								},
								"shouldEndSession": false
							},
							"sessionAttributes": {
								"position": 0,
								"word": "baum",
								"last_question": "none"
							}
						}`))
					})
				})
			})
		})
	})

	Context("Invalid body", func() {
		It("returns a StatusBadRequest", func() {
			response, e := client.Post("http://127.0.0.1:4443/", "", strings.NewReader("Hello"))
			Expect(e).NotTo(HaveOccurred())

			Expect(response.StatusCode).To(Equal(http.StatusBadRequest))
		})
	})

})
