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

		Context("locale: en-US", func() {
			It("returns a StatusOK and an American welcome message", func() {
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
					"locale": "en-US"
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
						"text": "This is Wikipedia. To read an article say e.g. \"What is a cheese cake?\"."
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

			Context("locale: en-US", func() {
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
								"word": "tree",
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
							"locale": "en-US",
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
								"text": "Table of contents. section 1: Definition.\nsection 2: Overview.\nsection 3: Distribution.\nsection 4: Parts and function.\nsection 5: Evolutionary history.\nsection 6: Tree ecology.\nsection 7: Uses.\nsection 8: Care.\nsection 9: Mythology.\nsection 10: Superlative trees.\nsection 11: See also.\nsection 12: Notes.\nsection 13: References.\nsection 14: Further reading.\n Which section do you want to go to?"
							},
							"shouldEndSession": false
						},
						"sessionAttributes": {
							"position": 0,
							"word": "tree",
							"last_question": "jump_where"
						}
					}`))
				})
			})
		})

		Context("DefineIntent", func() {
			Context("locale: de-DE", func() {
				It("reads the intro", func() {
					response, e := client.Post("http://127.0.0.1:4443/", "", strings.NewReader(`{
						"version": "1.0",
						"session": {
							"new": false,
							"sessionId": "xxx",
							"application": {
								"applicationId": "xxx"
							},
							"attributes": {
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
								"apiAccessToken": "zzz"
							}
						},
						"request": {
							"type": "IntentRequest",
							"requestId": "zzz",
							"timestamp": "`+time.Now().UTC().Format("2006-01-02T15:04:05Z")+`",
							"locale": "de-DE",
							"intent": {
								"name": "DefineIntent",
								"confirmationStatus": "NONE",
								"slots": {
									"word": {
										"name": "word",
										"value": "a cheese cake",
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
						"sessionAttributes": {
						  "last_question": "should_continue",
						  "position": 0,
						  "word": "a cheese cake"
						},
						"response": {
						  "outputSpeech": {
							"type": "PlainText",
							"text": "Der Dreikönigskuchen oder Königskuchen (englisch King Cake oder King’s Cake, französisch Galette des Rois, spanisch Roscón de Reyes, portugiesisch Bolo Rei) ist ein traditionelles Festtagsgebäck, das zum 6. Januar, dem Tag der Erscheinung des Herrn (Epiphanias), dem Festtag der heiligen drei Könige gebacken wird. Mit seiner Hilfe wird der Bohnenkönig gelost. Zur weiteren Navigation kannst Du jederzeit zum Inhaltsverzeichnis springen indem Du \"Inhaltsverzeichnis\" oder \"nächster Abschnitt\" sagst. Soll ich zunächst einfach weiterlesen?"
						  },
						  "shouldEndSession": false
						}
					  }`))
				})
			})

			Context("locale: en-US", func() {
				It("reads the intro", func() {
					response, e := client.Post("http://127.0.0.1:4443/", "", strings.NewReader(`{
						"version": "1.0",
						"session": {
							"new": false,
							"sessionId": "xxx",
							"application": {
								"applicationId": "xxx"
							},
							"attributes": {
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
								"apiAccessToken": "zzz"
							}
						},
						"request": {
							"type": "IntentRequest",
							"requestId": "zzz",
							"timestamp": "`+time.Now().UTC().Format("2006-01-02T15:04:05Z")+`",
							"locale": "en-US",
							"intent": {
								"name": "DefineIntent",
								"confirmationStatus": "NONE",
								"slots": {
									"word": {
										"name": "word",
										"value": "a cheese cake",
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
						"sessionAttributes": {
						  "last_question": "should_continue",
						  "position": 0,
						  "word": "a cheese cake"
						},
						"response": {
						  "outputSpeech": {
							"type": "PlainText",
							"text": "Cheesecake is a sweet dessert consisting of one or more layers. The main, and thickest layer, consists of a mixture of soft, fresh cheese (typically cream cheese or ricotta), eggs, vanilla and sugar; if there is a bottom layer it often consists of a crust or base made from crushed cookies (or digestive biscuits), graham crackers, pastry, or sponge cake. It may be baked or unbaked (usually refrigerated). Cheesecake is usually sweetened with sugar and may be flavored or topped with fruit, whipped cream, nuts, cookies, fruit sauce, or chocolate syrup. Cheesecake can be prepared in many flavors, such as strawberry, pumpkin, key lime, lemon, chocolate, Oreo, chestnut, or toffee. For further navigation you can jump to the table of contents any time or jump to a specific section. For now, shall I simply continue reading?"
						  },
						  "shouldEndSession": false
						}
					  }`))
				})
			})
		})

		Context("GoToSectionIntent", func() {
			Context("locale: de-DE", func() {
				It("returns a StatusOK and a message with content of the section", func() {
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
								"last_question": "jump_where"
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
										"value": "superlative",
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
								"text": "Abschnitt acht. Superlative. Der höchste Baum der Welt ist der „Hyperion“, ein Küstenmammutbaum (Sequoia sempervirens) im Redwood-Nationalpark in Kalifornien mit 115,5 Meter Wuchshöhe.\nDer höchste Baum Deutschlands, vielleicht sogar des Kontinents, ist die „Waldtraut vom Mühlwald“, eine 63,33 Meter (Stand: 18. August 2008) hohe Douglasie (Pseudotsuga menziesii) im  Arboretum Freiburg-Günterstal, einem Teil des Freiburger Stadtwalds.\nDer voluminöseste Baum der Welt ist angeblich der General Sherman Tree, ein Riesenmammutbaum im Sequoia National Park, Kalifornien, USA: Volumen etwa 1489 Kubikmeter, Gewicht etwa 1385 Tonnen (US), Alter rund 2500 Jahre.\nDer dickste Baum ist der „Baum von Tule“, eine Mexikanische Sumpfzypresse (Taxodium mucronatum) in Santa María del Tule im mexikanischen Staat Oaxaca. Sein Durchmesser an der dicksten Stelle beträgt 14,05 Meter.\nDie ältesten Bäume bezogen auf einen einzelnen Baumstamm sind – gemäß verbürgter Jahresringzählung – über 4800 Jahre alte Langlebige Kiefern (Pinus longaeva, früher als Varietät der Grannen-Kiefer angesehen) in den White Mountains in Kalifornien.\nDer älteste Baum bezogen auf den lebenden Organismus ist die Amerikanische Zitterpappelkolonie „Pando“ in Utah, USA, deren Alter auf mindestens 80.000 Jahre geschätzt wird. Aus den Wurzeln sprießen immer wieder neue, genetisch identische Baumstämme (vegetative Vermehrung), die etwa 100–150 Jahre alt werden.  Bei einem Individuum der Art „Huon Pine“ in Tasmanien, das mindestens 10.500 Jahre (vielleicht sogar 50.000 Jahre) alt ist, ist der älteste Baumstamm etwa 2000 Jahre alt. Die ältesten Bäume Europas stehen in der Provinz Dalarna in Schweden. 2008 wurden dort etwa 20 gemeine Fichten auf über 8000 Jahre datiert, die älteste auf 9550 Jahre. Die einzelnen Baumstämme sterben dabei nach etwa 600 Jahren ab und werden aus der Wurzel neu gebildet. \nDie winterhärtesten Bäume sind die Dahurische Lärche (Larix gmelinii) und die Ostasiatische Zwerg-Kiefer (Pinus pumila): Sie widerstehen Temperaturen bis zu −70 °C.\nDie Dahurische Lärche ist auch die Baumart, die am weitesten im Norden überleben kann: 72° 30' N, 102° 27' O.\nDie Bäume in der größten Höhe finden sich auf 4600 Meter Seehöhe am Osthimalaya in Sichuan, dort gedeiht die Schuppenrindige Tanne (Abies squamata).\nDas Holz geringster Dichte ist jenes des Balsabaumes.\nBäume, die bis dahin kahle Flächen besiedeln können, sogenannte Pionierpflanzen, sind zum Beispiel bestimmte Birken-, Weiden- und Pappelarten.\nIn der Bonsai­kunst versucht man, das Abbild eines uralten und erhabenen Baumes in klein in der Schale nachzuahmen.\nDie älteste Baumart der Erde und vermutlich das älteste lebende Fossil in der Pflanzenwelt ist der Ginkgo-Baum (Ginkgo biloba). Soll ich noch weiterlesen?"
							},
							"shouldEndSession": false
						},
						"sessionAttributes": {
							"position": 22,
							"word": "baum",
							"last_question": "should_continue"
						}
					}`))
				})
			})

			Context("Title cannot be found", func() {
				Context("locale: de-DE", func() {
					It("returns a StatusOK and a message that title cannot be found", func() {
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
