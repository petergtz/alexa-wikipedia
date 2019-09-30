package mediawiki_test

import (
	"fmt"
	"io/ioutil"

	"go.uber.org/zap"

	"github.com/nicksnyder/go-i18n/v2/i18n"

	"github.com/petergtz/alexa-wikipedia/locale"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/petergtz/alexa-wikipedia/mediawiki"
	"github.com/petergtz/alexa-wikipedia/wiki"
	"golang.org/x/text/language"
)

var _ = Describe("Mediawiki", func() {
	var (
		localizer *locale.Localizer
		logger    *zap.Logger
	)
	BeforeEach(func() {
		var e error
		logger, e = zap.NewDevelopment()
		Expect(e).NotTo(HaveOccurred())
		localizer = locale.NewLocalizer(i18n.NewBundle(language.English), "de-DE", logger.Sugar())
	})

	It("returns the page even when it's not an exact match", func() {
		page, e := (&mediawiki.MediaWiki{logger.Sugar()}).SearchPage("Der Baum", localizer)
		Expect(e).NotTo(HaveOccurred())
		Expect(page.Title).To(Equal("Baum"))
	})

	It("returns the page when it finds it", func() {
		page, e := (&mediawiki.MediaWiki{logger.Sugar()}).GetPage("Baum", localizer)
		Expect(e).NotTo(HaveOccurred())
		Expect(page.Title).To(Equal("Baum"))
	})

	It("returns an error when it cannot find the page", func() {
		_, e := (&mediawiki.MediaWiki{logger.Sugar()}).GetPage("NotExistingWikiPage", localizer)
		Expect(e).To(HaveOccurred())
		Expect(e.Error()).To(Equal("Page not found on Wikipedia"))
	})

	It("properly escapes wearch words", func() {
		page, e := (&mediawiki.MediaWiki{logger.Sugar()}).GetPage("Albert Einstein", localizer)
		Expect(e).NotTo(HaveOccurred())
		Expect(page.Title).To(Equal("Albert Einstein"))
	})

	Describe("WikiPageFrom", func() {
		It("works", func() {
			text, e := ioutil.ReadFile("testdata/extract-baum.wiki.txt")
			Expect(e).NotTo(HaveOccurred())
			page := mediawiki.WikiPageFrom(mediawiki.Page{Extract: string(text), Title: "Baum"}, localizer)

			Expect(page.Title).To(Equal("Baum"))
			Expect(page.Body).To(Equal("Als Baum wird im allgemeinen Sprachgebrauch eine verholzte Pflanze verstanden, die aus einer Wurzel, einem daraus emporsteigenden, hochgewachsenen Stamm und einer belaubten Krone besteht."))
			Expect(page.Subsections).To(HaveLen(12))
			Expect(page.Subsections[0].Title).To(Equal("Etymologie"))
			Expect(page.Subsections[1].Title).To(Equal("Definition und taxonomische Verbreitung"))
			Expect(page.Subsections[2].Title).To(Equal("Die besonderen Merkmale der Bäume"))
			Expect(page.Subsections[10].Title).To(Equal("Weblinks"))
			Expect(page.Subsections[11].Title).To(Equal("Einzelnachweise"))
			Expect(fmt.Sprintf("%#v", page)).To(Equal(fmt.Sprintf("%#v", wiki.Page{
				Title: "Baum",
				Body:  "Als Baum wird im allgemeinen Sprachgebrauch eine verholzte Pflanze verstanden, die aus einer Wurzel, einem daraus emporsteigenden, hochgewachsenen Stamm und einer belaubten Krone besteht.",
				Subsections: []wiki.Section{
					wiki.Section{
						Number:      "eins",
						Title:       "Etymologie",
						Body:        "Die Herkunft des westgerm. Wortes mhd., ahd. boum ist ungeklärt. Zum engl. tree siehe Teer#Etymologie. Baum als Begriff ist Teil der Swadesh-Liste.",
						Subsections: []wiki.Section{},
					},
					wiki.Section{
						Number:      "zwei",
						Title:       "Definition und taxonomische Verbreitung",
						Body:        "Die Botanik definiert Bäume als ausdauernde und verholzende Samenpflanzen, die eine dominierende Sprossachse aufweisen, die durch sekundäres Dickenwachstum an Umfang zunimmt. Diese Merkmale unterscheiden einen Baum von Sträuchern, Farnen, Palmen und anderen verholzenden Pflanzen. Im Gegensatz zu ihren entwicklungsgeschichtlichen Vorläufern verfügen die meisten Bäume zudem über wesentlich differenziertere Blattorgane, die mehrfach verzweigten Seitentrieben (Lang- und Kurztrieben) entspringen. Stamm, Äste und Zweige verlängern sich jedes Jahr durch Austreiben von End- und Seitenknospen, verholzen dabei und nehmen kontinuierlich an Umfang zu. Im Gegensatz zum Strauch ist es besonderes Merkmal der Bäume, dass die Endknospen über die Seitenknospen dominieren (Apikaldominanz) und sich dadurch ein vorherrschender Haupttrieb herausbildet (Akrotonie).\nBaumförmige Lebensformen kommen in verschiedenen Pflanzengruppen vor: „Echte“ Bäume sind die Laubbäume unter den Bedecktsamern und die baumförmigen Nacktsamer, zu denen Nadelholzgewächse wie die Koniferen gehören, aber auch Ginkgo biloba (als einziger noch existierender Vertreter der Ginkgogewächse) sowie zahlreiche Vertreter der fiederblättrigen Nacktsamer (Cycadophytina). Eigentümlichster Baum ist wohl die in Namibia vorkommende Welwitschia mirabilis, deren Stamm im Boden verbleibt. Daneben können auch die Palmen und die Baumfarne eine baumähnliche Form ausbilden. Diese Gruppen besitzen aber kein echtes Holz (sekundäres Xylem) und gelten daher nicht als Bäume. Eine Sonderstellung nimmt der Drachenbaum (Dracaena) ein. Dieser gehört zwar zu den Einkeimblättrigen, hat aber ein atypisches sekundäres Dickenwachstum.\nBaumähnliche Formen finden sich hauptsächlich in rund 50 höheren Pflanzenfamilien. Dagegen fehlt die Baumform bei Algen, Moosen, Liliengewächsen, Iridaceae, Hydrocharitaceae, Orchideen, Chenopodiaceae, Primelgewächsen und meist auch bei den Convolvulaceae, Glockenblumengewächsen, Cucurbitaceae, Doldengewächsen, Saxifragaceae, Papaveraceae, Ranunculaceae oder Caryophyllaceae.\nBäume kommen heute innerhalb der Nacktsamer (Gymnospermae) einerseits in Form der Ginkgoopsida mit der Art Ginkgo, andererseits der nadelblättrigen Nacktsamer (Coniferopsida, „Nadelbäume“) vor. Dominiert werden die Arten vor allem von der Ordnung Pinales mit den Familien Pinaceae (Fichten, Kiefern, Tannen, Douglasien, Lärchen, Goldlärche), Cupressaceae (Zypressen, Scheinzypressen, Sumpfzypressen, Lebensbäume, Wacholder, Mammutbäume), Podocarpaceae (Steineiben, Harzeiben), Araucariaceae (Araukarien, Kauri-Bäume), Taxaceae (Eiben) und Cephalotaxaceae (Kopfeiben).\nViele Baumarten kommen aber auch innerhalb der Bedecktsamer (Angiospermen) vor. Die verschiedenen Unterklassen haben hier unterschiedliche Laubbaumtypen hervorgebracht. Zu den bedeutendsten gehören die Buchengewächse (Fagaceae), zu denen neben den Buchen (Fagus spp.) auch die Eichen (Quercus spp.) und die Kastanien (Castanea) gezählt werden. Ebenfalls bedeutend sind die Birkengewächse (Betulaceae) mit den Birken und Erlen sowie die Nussbäume (Juglandaceae), die Ulmen (Ulmaceae) und die Maulbeergewächse (Moraceae). Zu den Rosiden zählen die Linden aus der Familie der Malvengewächse, die Obstgehölze aus der Familie der Rosengewächse (Rosaceae) sowie die Leguminosen (Fabales) mit sehr zahlreichen, vor allem tropischen Arten. Neben der Gattung Dalbergia (Palisanderbäume) gehört auch die Gattung Robinia in diese Gruppe. Wirtschaftlich bedeutsam sind die Zedrachgewächse (Meliaceae) mit den Gattungen Entandrophragma (Mahagonibäume) und Cedrela sowie die Familie der Dipterocarpaceae mit der Gattung Shorea (Meranti, Bangkirai).",
						Subsections: []wiki.Section{},
					},
					wiki.Section{
						Number: "drei",
						Title:  "Die besonderen Merkmale der Bäume",
						Body:   "",
						Subsections: []wiki.Section{
							wiki.Section{
								Number: "drei.eins",
								Title:  "Morphologie baumförmiger Lebensformen",
								Body:   "Baumartige Lebensformen zeigen eine große Variationsbreite in ihrem Aufbau (Morphologie). Assoziiert wird mit dem Begriff Baum der Aufbau aus Baumkrone, Baumstamm und Baumwurzeln. Bei den baumartigen Farnen und den meisten Palmen finden sich einfache Stämme, die keine Äste ausbilden, sondern schopfartig angeordnete, häufig gefiederte Blätter. Vor allem zeigen sie kein sekundäres Dickenwachstum und sind damit keine echten Bäume.",
								Subsections: []wiki.Section{
									wiki.Section{
										Number:      "drei.eins.eins",
										Title:       "Wachstum",
										Body:        "Bei den echten Bäumen wächst aus dem Spross der Keimpflanze durch Längen- und sekundäres Dickenwachstum der künftige Baumstamm heran: Es bildet sich der Spross an der Spitze durch die sich ständig erneuernde Gipfelknospe aufrecht weiter und wird zum geraden, bis zur höchsten Kronenspitze durchgehenden Baumstamm (Monopodium). In der Spitzenknospe gebildete Wuchsstoffe (Auxine) unterdrücken die Aktivität der Seitenknospen. Bei vielen Baumarten lässt diese Dominanz des Haupttriebs mit dem Alter nach und es bildet sich eine typische, verzweigte Laubbaumkrone.\n\nBei anderen Gehölzen wie der Buche oder der Hainbuche übernimmt eine subterminale Seitenknospe die Führung (Sympodium). Bei Bäumen entsteht so eine aufrechte „Scheinachse“ (Monochasium). Im späteren Verlauf lässt die Dominanz der führenden Knospe nach und aus weiteren Seitenknospen entwickeln sich stärkere Äste, die schließlich eine Krone bilden. Dies geschieht meist früher als bei Bäumen mit monopodialem Wuchs.\nSträucher hingegen sind durch das völlige Fehlen der apikalen Dominanz gekennzeichnet. Zahlreiche bodenbürtige Seitentriebe bilden hier eine weit verzweigte Wuchsform.\nBei Gehölzen bildet sich an den Wuchsachsen während der Vegetationsperiode je ein Triebabschnitt (Jahrestrieb), dessen Beginn lange an den schmalen ringförmigen Blattnarben der ehemaligen Knospenschuppen erkennbar ist. Ein weiterer Austrieb nach der Vegetationsperiode wird als Johannistrieb (Prolepsis) bezeichnet. Tropische Arten neigen zu mehrfachem Austrieb.",
										Subsections: []wiki.Section{},
									},
									wiki.Section{
										Number:      "drei.eins.zwei",
										Title:       "Alter",
										Body:        "Aus der Zahl der Jahrestriebe und dem Grad der Verzweigung lässt sich das Alter eines Astes ermitteln. Diese Altersbestimmung wird jedoch bei zahlreichen Arten (zum Beispiel Fichte oder Tanne) und regelmäßig bei älteren Bäumen durch die Ausbildung von sogenannten Proventivtrieben erschwert, die aus „schlafenden“ Knospen austreiben. Die regelmäßige Bildung von Proventivtrieben wird als Reiteration (sprich: Re-Iteration) bezeichnet. Diese Wiederholungstriebe dienen der Erneuerung der Krone und verschaffen Bäumen die Möglichkeit, alternde Äste zu ersetzen sowie auf Stress (Schneebruch, Insektenkalamitäten) zu reagieren.\n\nBäume können ein Alter von mehreren 100 Jahren, an bestimmten Standorten sogar von mehreren 1000 Jahren erreichen. Als ältester Baum der Welt gilt (Stand: 2008) die 9550 Jahre alte Fichte Old Tjikko im Nationalpark Fulufjället im mittelschwedischen Bezirk Dalarna. Unter dieser Fichte wurden drei weitere „Generationen“ (375, 5660 und 9000 Jahre alt) mit identischem Erbmaterial gefunden. Die Zahl der über 8000 Jahre alten Fichten wird auf etwa 20 Stück geschätzt. Damit ist die Fichte rund doppelt so alt wie die nordamerikanischen Kiefern, die mit 4000 bis 5000 Jahren bislang als die ältesten lebenden Bäume galten. Die nachweislich ältesten Bäume Mitteleuropas werden auf etwa 600 bis 700 Jahre datiert.\nWächst der Baum unter im Jahresrhythmus schwankenden klimatischen Bedingungen, wird während der Vegetationsperiode ein Jahresring angelegt. Mit Hilfe dieser Ringe lassen sich das Alter eines Baumes und dessen Wuchsbedingungen in den einzelnen Jahren ablesen. Die Dendrochronologie nutzt dies, um altes Holz zu datieren und das Klima einer Region bis zu mehreren 1000 Jahren zu rekonstruieren.",
										Subsections: []wiki.Section{},
									},
									wiki.Section{
										Number:      "drei.eins.drei",
										Title:       "Baumschädigungen",
										Body:        "Seine Entwicklung bringt für den Baum zahlreiche Probleme und Schädigungen mit sich. Hierunter fallen vor allem:\n\n Pilzbefall,\n Insektenschaden,\n Windbruch (Baumteile brechen ab),\n Windwurf (der Baum wird mit den Wurzeln aus dem Boden gehebelt),\n Schneebruch (Baumteile unter schweren Schneelasten brechen ab),\n Blitzschaden (Stammteile werden abgesprengt),\n Frost (Trockenschaden durch Transpiration bei gefrorenem Boden, Stammrisse).Bei Jungbäumen kommt es insbesondere zu:\n\n übermäßigem Wildverbiss,\n Schälung der Rinde,\n Wühlmausschaden an der Wurzel.Einige wichtige Krankheiten, von denen Bäume befallen werden können, sind Brand, Krebs, Rost, Mehltau, Rotfäule, Weißfäule, Braunfäule und Harzfluss. Zu Missbildungen an Bäumen zählen die Maserkröpfe, die Hexenbesen oder Wetterbüsche sowie die Gallen.",
										Subsections: []wiki.Section{},
									},
								},
							},
							wiki.Section{
								Number:      "drei.zwei",
								Title:       "Aufbau des Baumstammes",
								Body:        "Ein Querschnitt durch einen Baumstamm, die verholzende Hauptachse (Caulom) – in der Dendrologie Schaft genannt, zeigt verschiedene Zonen. Ganz innen befinden sich das aus Primärgewebe bestehende Mark und das tote Kernholz. Bestimmte Baumarten (z. B. Buche, Esche) bilden fakultativ einen Falschkern aus, der sich in den Eigenschaften vom echten Kernholz unterscheidet. Weiter außen befindet sich das Splintholz, das der Leitung und Speicherung dient und sich bei sogenannten Kernholzbäumen farblich meist deutlich vom Kernholz abhebt. Bei der Eiche, der Eibe und der Robinie ist dies sehr gut sichtbar. Die Fichte hat einen farblosen Kern (Reifholz).\nDie äußerste Schicht bildet die Baumrinde. Sie besteht aus der Bastschicht, die in Wasser gelöste Nährstoffe transportiert, und der Borke, die den Stamm vor Umwelteinflüssen (UV-Einstrahlung, Hitze, mechanische und biotische Schäden) schützt.\nZwischen der Bastschicht und dem Holz befindet sich bei Gymnospermen und Dikotyledonen das Kambium. Diese Wachstumsschicht bildet durch sekundäres Dickenwachstum nach innen Holz (Xylem) und nach außen Bast (Phloem). Das Holz zeichnet sich durch die Einlagerung von Lignin in die Zellwand aus. Dadurch werden die Zellen versteift und bilden ein festes Dauergewebe. Das sekundäre Dickenwachstum, die Lignifizierung der hölzernen Zellwand und die Vermehrung durch Samen verschafften den Bäumen in den meisten Biomen der Erde einen Vorteil gegenüber anderen Pflanzen und haben dort zur Entwicklung großflächiger Waldbestände geführt. Ausnahmen bilden die Wüsten, die arktischen Tundren und die zentralkontinentalen Steppen.\nHinsichtlich des inneren Baus des Baumstamms weichen die zu den Einkeimblättrigen gehörenden Palmen von den echten Bäumen erheblich ab. Bei ersteren stehen die Gefäßbündel im Grundgewebe zerstreut, weshalb es keinen Kambium\u00adring, keinen Holzzylinder und somit kein fortdauerndes sekundäres Dickenwachstum des Stammes gibt. Bei den zu den Dikotyledonen oder Gymnospermen gehörenden Bäumen besitzt der Stamm schon in der frühesten Jugend als dünner Stängel einen unter der Rinde gelegenen Kreis von Leitbündeln, der den Rindenbereich vom innen liegenden Mark scheidet. Dieser Leitbündelring stellt in seiner inneren, dem Mark anliegenden Hälfte das Holz und im äußeren, an die Rinde angrenzenden Teil den Bast dar; zwischen beiden zieht sich der Kambiumring hindurch. Dieser wird aus zarten, saftreichen, sich ständig teilenden Zellen gebildet und vergrößert durch seinen laufenden Zellvermehrungsprozess die beiderseits ihm anliegenden Gewebe. So wird alljährlich an der Außenseite des Holzringes eine neue Zone Holzgewebe angelegt, wodurch die Jahresringe des auf diese Weise erstarkenden Holzkörpers entstehen, die als konzentrische Linien am Stammquerschnitt wahrnehmbar sind. Andererseits erhält aber auch der weiter außen liegende Bast an seiner Innenseite einen jährlichen, wenn auch weit geringeren Zuwachs. Auf diese Weise kommt die dauernde Verdickung des Stammes und aller Äste sowie auch der Wurzeln zustande.",
								Subsections: []wiki.Section{},
							},
							wiki.Section{
								Number:      "drei.drei",
								Title:       "Wurzel",
								Body:        "Auch in der Wurzelbildung unterscheiden sich die Bäume untereinander. Neben der genetischen Festlegung steuern die Erfordernisse der Verankerung des Baumes im Boden ebenso wie die Notwendigkeit der Versorgung der Pflanze mit Wasser und Nährstoffen die Intensität und Art des Wurzelwachstums. Man spricht entsprechend der Form des Wurzelstocks von Pfahlwurzel, Flachwurzel oder Herzwurzel. Bei der Pfahlwurzel wächst die Hauptwurzel senkrecht in den Boden hinab, was besonders für die Eiche charakteristisch ist. Flachgründige Böden und hoch anstehendes Grundgestein oder Grundwasser begünstigen z. B. die Bildung von Flachwurzeln. Trockene Böden begünstigen eine Bildung von Pfahlwurzeln. Die überwiegende Masse des Wurzelstocks machen bei den Bäumen nicht die verholzten Wurzelteile, sondern die mit einer Mykorrhiza vergesellschafteten Feinwurzeln aus. Die Gesamtwurzelmasse reicht oft an die Masse der oberirdischen Pflanzenteile heran. Bei einkeimblättrigen baumähnlichen Lebensformen endet der Stamm nahe unter der Bodenfläche und es entwickelt sich ein sprossbürtiges Wurzelsystem (Homorhizie).\nAn alten Bäumen finden sich meist junge Adventivwurzeln, die alte, ineffektive Wurzeln ersetzen. Bei einigen Baumarten bilden oberflächennahe Wurzeln eine sogenannte Wurzelbrut, eine Form der vegetativen Vermehrung. Wurzelkappungen infolge von Baumaßnahmen können das Absterben von Wurzelteilen bewirken und führen zum Eindringen von holzzerstörenden Pilzen in den Baum. Dies ist die häufigste Ursache von irreparablen Baumschäden im städtischen Bereich.",
								Subsections: []wiki.Section{},
							},
							wiki.Section{
								Number:      "drei.vier",
								Title:       "Blätter",
								Body:        "Bäume tragen Laubblätter oder Nadelblätter, die entweder mehrjährig am Baum verbleiben (immergrüne Arten) oder am Ende einer Vegetationsperiode abgeworfen werden (laubabwerfende Arten). Dazwischen liegen noch die halbimmergrünen Arten, die am Ende einer Vegetationsperiode nur einen Teil ihrer Blätter verlieren, bei Neuaustrieb dann aber die vorjährigen ersetzen. Die Nadelgehölze sind mit Ausnahme der Gattungen Lärchen (Larix) und Goldlärchen (Pseudolarix) immergrüne Arten. In den borealen und hochmontanen Biomen der Nordhalbkugel haben sich die immergrünen Nadelgehölze durchgesetzt, da sie zu Beginn der Vegetationsperiode bei ausreichender Temperatur sofort mit der Assimilation beginnen können, ohne zunächst Assimilationsorgane bilden zu müssen wie die laubabwerfenden Baumarten.\nDie Gestalt der Blätter (Laub) ist ein wichtiges Bestimmungsmerkmal. Anordnung, Form, Größe, Farbe, Nervatur und Zähnung sowie haptische Eigenschaften können zur Differenzierung herangezogen werden. Nicht minder brauchbar zur Unterscheidung im winterlichen Zustand sind die (Blatt-)Knospen des Baumes. Eine eindeutige taxonomische Identifizierung der Arten ist allerdings nur anhand der Blüten oder Früchte möglich. Manche Bäume sind mit Dornen ausgestattet. Dies sind entweder kurze Zweige, die mit dorniger Spitze enden (Weißdorne, Wildformen von Obstbäumen) oder es sind stachelartig ausgebildete Nebenblätter wie etwa bei der Gewöhnlichen Robinie.\nEin europäischer Laubbaum trägt durchschnittlich 30.000 Blätter, die zusammen eine enorme Transpirationskapazität haben. An warmen Sommertagen kann der Baum mehrere hundert Liter Wasser verdunsten. Beispiel einer 80-jährigen, alleinstehenden Rotbuche: In diesem Lebensalter ist der Baum 25 Meter hoch, und seine Baumkrone mit einem Durchmesser von 15 Meter bedeckt eine Standfläche von 160 m². In ihren 2700 m³ Rauminhalt finden sich 800.000 Blätter mit einer gesamten Blattoberfläche von 1600 m², deren Zellwände zusammen eine Fläche von 160.000 m² ergibt. Pro Stunde verbraucht diese Buche 2,352 kg Kohlenstoffdioxid, 0,96 kg Wasser und 25.435 Kilojoule Energie (das ist die in Form von Traubenzucker gespeicherte Energie, die eingestrahlte Sonnenenergie ist etwa siebenmal größer); im gleichen Zeitraum stellt sie 1,6 kg Traubenzucker her und deckt mit 1,712 kg Sauerstoff den Verbrauch von zehn Menschen. Die 15 m³ Holz des Baumes wiegen trocken 12.000 kg, allein 6000 kg davon sind Kohlenstoff.",
								Subsections: []wiki.Section{},
							},
							wiki.Section{
								Number:      "drei.fünf",
								Title:       "Blüten",
								Body:        "Die Blüten der Bäume aus gemäßigten Breiten sind manchmal verhältnismäßig unscheinbar; bei einigen Taxa sind einzelne Blütenblattkreise reduziert. Einige Baumarten gemäßigter Breiten haben eingeschlechtige Blüten. Dabei sitzen die Blüten beider Geschlechter entweder auf demselben Baum (einhäusig getrenntgeschlechtig, zum Beispiel Eiche, Buche, Hainbuche, Birke, Erle und Nussbaum) oder auf verschiedenen (zweihäusig getrenntgeschlechtig), so dass man männliche und weibliche Bäume zu unterscheiden hat (unter anderen bei Weiden und Pappeln). Andere Bäume wie Obstbäume, Rosskastanie und viele Bäume der wärmeren Klimate haben Zwitterblüten, die sowohl Staub- als auch Fruchtblätter ausbilden.",
								Subsections: []wiki.Section{},
							},
							wiki.Section{
								Number:      "drei.sechs",
								Title:       "Frucht- und Samenbildung",
								Body:        "Die Frucht- und Samenbildung zeigt weniger Eigentümlichkeiten. Bei den meisten Bäumen fällt die Reife in den Sommer oder Herbst desselben Jahres; nur bei den Kiefernarten erlangen die Samen und die sie enthaltenden Zapfen erst im zweiten Herbst nach der Blüte vollständige Ausbildung. Die Früchte sind meistens nussartig mit einem einzigen ausgebildeten Samen, oder sie bestehen aus mehreren einsamigen, nussartigen Teilen, wie bei den Ahornen. Saftige Steinfrüchte, ebenfalls mit einem oder wenigen Samen, finden sich bei den Obstbäumen, Kapseln mit zahlreichen Samen bei den Weiden und Pappeln.",
								Subsections: []wiki.Section{},
							},
						},
					},
					wiki.Section{
						Number:      "vier",
						Title:       "Entwicklung baumförmiger Pflanzen in der Erdgeschichte",
						Body:        "Die Voraussetzungen für die Entstehung und Verbreitung der Bäume waren:\n\n die Entwicklung des Kormus (Differenzierung zwischen Blatt, Spross und Wurzel) als Organisationsform der höheren Pflanzen,\n die Entwicklung des Samens als Fortpflanzungsmethode,\n die Entwicklung des Lignins für die Bildung von Dauergewebe,\n die Entwicklung des sekundären Dickenwachstums für die Bildung mehrjähriger Organismen.Die Vorläufer der Bäume kennt man aus dem Karbon. Sie gehörten zu den Schachtelhalmgewächsen, den Bärlappgewächsen und den Farnen. Sie besaßen verholzte Stämme, die auch ein sekundäres Dickenwachstum aufwiesen. Fossile Gattungen sind beispielsweise Lepidodendron und Sigillaria. Die verdichteten Sedimente dieser Wälder bilden die Steinkohle.\nDie weitere Evolution der Pflanzen brachte im Perm die Samenpflanzen hervor. Die Nacktsamer breiteten sich als erste Bäume rasch aus, erreichten wohl in der Trias (vor etwa 200 Millionen Jahren) ihre größte Artenvielfalt, bis sie im Tertiär (vor etwa 60 Millionen Jahren) von den Angiospermen in ihrer Bedeutung abgelöst wurden. Von den bekannten 220.000 Blütenpflanzen sind etwa 30.000 Holzarten, so dass etwa jede achte Blütenpflanze ein Baum oder Strauch ist. Die meisten Baumarten zählen zu den Bedecktsamern (Angiospermen). Die Gymnospermen (Nacktsamer) umfassen nur ungefähr 800 Arten, bedecken aber immerhin ein Drittel der Waldfläche der Erde.\nDie globale Verteilung der Baumarten wurde vor allem durch die klimatischen Verhältnisse und durch die Kontinentalverschiebung geprägt. Während zum Beispiel die Buchengewächse (Fagaceae) eine typische Familie der Nordhemisphäre sind, ist beispielsweise die Familie Podocarpaceae vorwiegend in der Südhemisphäre verbreitet. Die heutige natürliche Artenverteilung wurde stark von den quartären Eiszeiten beeinflusst. Das gleichzeitige Vordringen der skandinavischen und alpinen Gletschermassen Europas hat zu einer Verdrängung zahlreicher Spezies geführt und die im Vergleich zu Nordamerika auffällige Artenarmut in Zentraleuropa verursacht. So stehen etwa der einzigen in den montanen Regionen Mitteleuropas heimischen Fichtenart, der Gemeinen Fichte (Picea abies), zahlreiche Fichtenarten auf dem nordamerikanischen Kontinent gegenüber.",
						Subsections: []wiki.Section{},
					},
					wiki.Section{
						Number: "fünf",
						Title:  "Physiologie",
						Body:   "",
						Subsections: []wiki.Section{
							wiki.Section{
								Number:      "fünf.eins",
								Title:       "Wuchs",
								Body:        "Wie bei allen Pflanzen unterliegen auch bei Bäumen der Stoffwechsel und das Wachstum sowohl endogenen (genetisch festgelegten) als auch äußeren Einflüssen. Zu letzteren zählen vor allem die Standortverhältnisse, das Klima und die Konkurrenz mit anderen Organismen beziehungsweise deren schädigende Wirkung. Während der Vegetationsperiode sorgen die Spitzenmeristeme und das Kambium für stetigen Längen- und Dickenzuwachs. Beginn und Ende der Vegetationsperiode sind je nach Baumart durch die Witterung und die Wasserverfügbarkeit beziehungsweise durch die Tageslänge bestimmt. Das Wachstum wird dabei durch Phytohormone gesteuert und die Akkumulation von Biomasse gezielt optimiert. Bäume sind so in der Lage, sich an ändernde Wuchsbedingungen anzupassen und gerichtet Festigungs-, Leit-, Speicher- oder Assimilationsgewebe anzulegen.\nDie Produktion neuen Gewebes mit dem sekundären Dickenwachstum und die Anlage neuer Jahrestriebe bewirkt, dass sich ein Baum ständig von innen nach außen erneuert. Der amerikanische Baumbiologe Alex Shigo hat daraus das Konzept der Kompartimentierung entwickelt, das den Baum als ein Ensemble zusammenwirkender Kompartimente sieht. Auf Verletzungen reagiert der Baum, anders als Tiere und Menschen, durch Abschottungsreaktionen und Aufgabe der eingekapselten Kompartimente (CODIT-Modell). Durch adaptives Wachstum optimiert er zudem seine Gestalt.\nComputermodellierungen des Karlsruher Physikers und Biomechanikers Claus Mattheck konnten zeigen, dass Bäume durch adaptives Wachstum eine mechanisch optimale Gestalt anstreben und zum Beispiel Kerbspannungen in Verzweigungen vermeiden, so dass die Gefahr von Brüchen minimiert wird. Diese Erkenntnisse haben zu Optimierungen unter anderem im Maschinenbau geführt.",
								Subsections: []wiki.Section{},
							},
							wiki.Section{
								Number:      "fünf.zwei",
								Title:       "Wasserleitung",
								Body:        "Der Wassertransport wird in den Nadelgehölzen durch die Tracheiden, in den Laubbäumen durch die effektiveren Gefäße (Poren) bewerkstelligt. Letztere sind bei den Laubbäumen entweder zerstreut (zum Beispiel bei Buche, Ahorn, Pappel) oder ringförmig (zum Beispiel bei Eiche, Ulme, Esche) im Jahresring angeordnet. Beispielsweise kann eine Eichenpore mit 400 µm Durchmesser 160.000-mal mehr Wasser als eine Nadelholztracheide mit 20 µm Durchmesser im gleichen Zeitraum transportieren.\nNach überwiegend vertretener Lehre funktioniert der Wassertransport der Bäume durch Saugspannungen in den Leitgeweben infolge Verdunstung an den Stomata der Blätter (Kohäsionstheorie). Dabei müssen Baumhöhen bis über 100 Meter überwunden werden können, was nach dieser Theorie nur mit enormen Drücken möglich ist. Kritiker dieser Lehre behaupten, dass schon bei wesentlich geringeren Höhen die Saugspannung zum Abriss des Wasserfadens in den Kapillaren führen müsste. Als gesichert gilt allerdings, dass im Frühjahr Zucker in den Speicherzellen mobilisiert werden und durch den aufgebauten osmotischen Druck Wasser aus den Wurzeln nachfließt. Dabei werden im Bodenwasser gelöste Nährsalze (vor allem K, Ca, Mg, Fe) vom Baum aufgenommen. Erst nach Ausdifferenzierung der Blätter werden die in der Krone erzeugten Assimilate über den Bast stammabwärts transportiert und stehen für das Dickenwachstum zur Verfügung. Eine Ausnahme bilden die ringporigen Laubbäume, bei denen die ersten Frühholzporen aus den im Vorjahr gebildeten Reservestoffen gebildet werden.\nDie süßen „Baumsäfte“ wurden von Menschen durch Einschneiden der Rinde abgezapft und durch Einkochen zu Sirupen weiterverarbeitet, beispielsweise Ahornsirup oder der Saft der Manna-Esche. Palmzucker oder Palmsirup allerdings ist ein Extrakt aus dem Blütensaft der Nipa- und Zuckerpalme (Unterfamilie Arecoideae), Agavensirup stammt aus dem „Saft“ der zu den Stauden gehörenden Agaven, Birkenzucker wurde ursprünglich in Finnland direkt aus der Birkenrinde gewonnen.\nDie Hydrologie beziehungsweise Bodenökologie unterscheidet zwischen dem Niederschlag, welcher im Bereich der Baumkrone auf den Boden trifft (Kronendurchlass) und dem Anteil, welche am Stamm herabfließt (Stammabfluss). Ein Teil des Niederschlags verdunstet direkt vom Baum (Interzeption) und erreicht den Boden nicht.",
								Subsections: []wiki.Section{},
							},
						},
					},
					wiki.Section{
						Number: "sechs",
						Title:  "Ökologie",
						Body:   "",
						Subsections: []wiki.Section{
							wiki.Section{
								Number:      "sechs.eins",
								Title:       "Wald",
								Body:        "Dort wo Bäume ausreichend Licht, Wärme und Wasser vorfinden, bilden sie Wälder. Im Jahr 2000 waren laut FAO 30 Prozent der Festlandmasse der Erde bewaldet. Pro Hektar binden Waldbäume zwischen 60 und 2000 Tonnen organisches Material und sind damit die größten Biomassespeicher der Kontinente. Die Gesamtmenge der 2005 weltweit in den Wäldern akkumulierten Holzmasse betrug 422 Gigatonnen. Da etwa die Hälfte der Holzsubstanz aus Kohlenstoff besteht, sind Wälder nach den Ozeanen die größten Kohlenstoffsenken der Biosphäre und damit für die CO2-Bilanz der Erdatmosphäre bedeutsam.\nDie mit der Bestandsbildung von Bäumen einhergehende Konkurrenz um Ressourcen führt zu einer Anpassung des Habitus gegenüber den freistehenden Exemplaren (Solitäre). Natürlicher Astabwurf innerhalb der Schattenkrone sowie Verlagerung der Assimilation in die Lichtkrone sind Optimierungsreaktionen der Bäume, die zu einem hohen, schlanken Wuchs mit kleinen Kronen und oft zu hallenartigen Beständen führen (zum Beispiel Buchen-Altbestände).\nDie heutige Ausbreitung und Artenzusammensetzung der Wälder steht stark unter dem Einfluss der wirtschaftlichen Tätigkeit des Menschen. Der Übergang von der Jäger- und Sammlerkultur zum Ackerbau ging in den dicht besiedelten Regionen mit der Zurückdrängung der Wälder einher. Nützlich waren Bäume den Menschen zunächst vorwiegend als Brennholz (Niederwald\u00adwirtschaft). Im Laufe der Entwicklung wurde die Gewinnung von Nutzholz aus Hochwäldern immer wichtiger. Diese Entwicklung hält an. Laut FAO wurden noch Ende der 1990er-Jahre weltweit 46 Prozent des weltweiten Holzeinschlags (3,2 Milliarden m³) als Brennholz genutzt, in den Tropen waren es sogar 86 Prozent. Die extensive Waldvernichtung in Zentraleuropa während des Mittelalters hat in der Neuzeit zur Einführung des Prinzips der nachhaltigen Waldbewirtschaftung geführt, nach dem nur so viel Holz entnommen werden darf, wie nachwächst.",
								Subsections: []wiki.Section{},
							},
							wiki.Section{
								Number:      "sechs.zwei",
								Title:       "Verbreitungszentren, Diversität",
								Body:        "In den Primärwäldern der feuchten Tropen findet sich die größte Artenvielfalt aller Waldtypen. Wichtige tropische Familien sind die Wolfsmilchgewächse (Euphorbiaceae), Seifenbaumgewächse (Sapindaceae), Bombacaceae, Byttneriaceae, Mahagonigewächse (Meliaceae), Hülsenfrüchtler (Fabaceae), Caesalpiniaceae, Verbenaceae, Sterculiaceae, Dipterocarpaceae und Sapotaceae.\nIn der subtropischen Zone findet man Bäume unter den immergrünen Myrtengewächsen (Myrtaceae) und Lorbeergewächsen (Lauraceae) sowie Silberbaumgewächsen (Proteaceae), denen sich in der wärmeren gemäßigten Zone andere immergrüne Bäume anschließen, so die immergrünen Eichen, Granatbäume, Orangen- und Zitronenbäume sowie Ölbäume.\nDagegen sind in der gemäßigten Zone die laubwechselnden Bäume vorherrschend. Hier sind Wälder von Eichen, Buchen und Hainbuchen charakteristisch. Zu den in Mitteleuropa heimischen Laubbäumen zählen die Ahorne, Birken, Buchen, Eichen, Erlen, Eschen, Linden, Mehlbeeren, Pappeln, Ulmen und Weiden. Typische Nadelbäume sind die Fichten, Kiefern, Lärchen, Tannen und Eiben. In Mitteleuropa häufig vorkommende Baumarten, die in diesem Gebiet ursprünglich nicht beheimatet sind, sind die Gewöhnliche Robinie, der Walnussbaum und viele Obstbäume. Sie alle sind Neophyten. Eine detaillierte Aufstellung bietet die Liste von Bäumen und Sträuchern in Mitteleuropa.\nUnd obgleich auch hier bereits Nadelhölzer in zusammenhängenden Waldungen auftreten, werden die Nadelwälder erst in der subarktischen (borealen) Zone vorherrschend, wo die Laubbäume nach und nach verdrängt werden. Artenvielfalt wie auch Wuchshöhe der Bäume nehmen mit zunehmender Annäherung an den Polarkreis ab. Eichen, Linden, Eschen, Ahorne und Buchen finden sich in Schweden nur noch diesseits des 64. Grades nördlicher Breite. Jenseits dieser Breite besteht die Baumvegetation hauptsächlich aus Fichten und Tannen, die in zusammenhängenden Wäldern nordöstlich noch über den 60. Grad hinausreichen, sowie aus Birken, die in zusammenhängenden Beständen sich fast bis zum 71. Grad nördlicher Breite erstrecken, und zum Teil aus Erlen und Weiden.\nAuch die Höhe über dem Meeresspiegel hat auf die Ausbreitung und Höhe der Bäume (in Abhängigkeit von der geographischen Breite) einen bedeutenden Einfluss. In den Anden finden sich noch bis in 5000 m Höhe Polylepis-Bäume. Unter 30 Grad nördlicher Breite, wo die Schneegrenze bei 4048–4080 m liegt, kommen auf dem Himalaja, nördlich von Indien, noch in 3766 m Höhe Baumgruppen vor, die aus Eichen und Fichten bestehen. Ebenso sind in Mexiko, unter 25–28 Grad nördlicher Breite, die Gebirge bis 3766 m mit Fichten und bis 2825 m hoch mit mexikanischen Eichen bedeckt. In den Alpen des mittleren Europas endet der Holzwuchs bei einer Höhe von 1570 m, im Riesengebirge bei 1193 m und auf dem Brocken bei 1005 m. Eichen und Tannen stehen auf den Pyrenäen noch bis zu einer Höhe von 1883 m; dagegen wächst die Fichte auf dem Sulitelma in Lappland, bei 68 Grad nördlicher Breite, kaum in einer Höhe von 188 m, die Birke kaum in einer von 376 m.",
								Subsections: []wiki.Section{},
							},
						},
					},
					wiki.Section{
						Number: "sieben",
						Title:  "Bäume und Menschen",
						Body:   "Die wissenschaftliche Lehre von den Bäumen (Gehölzen) ist die Dendrologie. Anpflanzungen von Bäumen in systematischer oder pflanzengeographischer Anordnung, die Arboreten, dienen ihr zu Beobachtungs- und Versuchszwecken. Gehölze können vegetativ, das heißt durch Pflanzenteile, oder generativ durch Aussaat vermehrt werden. In Baumschulen findet eine gezielte Auslese, Anzucht und Vermehrung von Bäumen und Sträuchern statt. Neben der forstlichen Nutzung finden Bäume reichliche Verwendung im Garten- und Landschaftsbau. Mit der Baumpflege hat sich ein eigener Berufsstand zum Erhalt und zur fachgerechten Behandlung von Bäumen in urbanen Regionen entwickelt.\n\nDas schrieb der Historiker Alexander Demandt und hat dem Baum mit Über allen Wipfeln – Der Baum in der Kulturgeschichte ein umfangreiches Werk gewidmet. Für ihn beginnt die Kulturgeschichte mit dem Feuer, das der Blitz in die Bäume schlug, und mit dem Werkzeug, für das Holz zu allen Zeiten unentbehrlich war.",
						Subsections: []wiki.Section{
							wiki.Section{
								Number:      "sieben.eins",
								Title:       "Nutzung",
								Body:        "Neben der wichtigen Funktion der Bäume bei der Gestaltung von Kulturlandschaften begleitet vor allem die Holznutzung die Entwicklung der Menschheit. Abgesehen von der vor allem in Entwicklungsländern immer noch weit verbreiteten Brennholznutzung ist Holz ein vielseitiger Bau- und Werkstoff, dessen produzierte Menge die Produktionsmengen von Stahl, Aluminium und Beton weit übersteigt. Damit ist Holz nach wie vor der wichtigste Bau- und Werkstoff weltweit; Bäume sind dementsprechend eine bedeutende Rohstoffquelle.\nNeben der Holznutzung dienen Bäume vor allem der Gewinnung von Blüten, Früchten, Samen oder einzelnen chemischen Bestandteilen (Terpentin, Zucker, Kautschuk, Balsame, Alkaloide und so weiter). In der Forstwirtschaft der industrialisierten Länder spielen diese Nutzungen eine untergeordnete Rolle. Lediglich der Obstbau als Teilbereich der Landwirtschaft ist in vielen Regionen ein wichtiger Wirtschaftsfaktor. Der Anbau erfolgt in Form von Plantagen. Hochwertige Obstsorten werden meist durch Okulation oder Pfropfen veredelt. Dies erfolgt durch den Einsatz ausgewählter Obstsorten, wobei die bekannten und gewollten Eigenschaften der Früchte einer Obstsorte auf einen jungen Baum übertragen werden. Zurückgegangen ist dagegen die Nutzung von Streuobstwiesen, die früher in vielen Gebieten Mitteleuropas landschaftsprägend waren.",
								Subsections: []wiki.Section{},
							},
							wiki.Section{
								Number:      "sieben.zwei",
								Title:       "Gesellschaftliches",
								Body:        "Dieser Bedeutung entsprechend ist ein vielfältiges Brauchtum mit dem Baum verknüpft. Das reicht vom Baum, der zur Geburt eines Kindes zu pflanzen ist, über den Maibaum, der in manchen Regionen in der Nacht zum ersten Mai der Liebsten verehrt wird, über Kirmesbaum und Weihnachtsbaum, unter denen man feiert, und über den Richtbaum auf dem Dachstuhl eines neu errichteten Hauses bis zum Baum, der auf dem Grab gepflanzt wird. Nationen und Völkern werden bestimmte, für sie charakteristische Bäume zugeordnet. Eiche und Linde gelten als typisch „deutsche“ Bäume. Die Birke symbolisiert Russland, und der Baobab gilt als der typische Baum der afrikanischen Savanne. Unter der Gerichtslinde wurde Recht gesprochen (siehe auch → Thing) und unter der Tanzlinde gefeiert.\nSeit 1989 wird jedes Jahr im Oktober für das darauffolgende Jahr der Baum des Jahres bestimmt, zunächst vom „Verein Baum des Jahres e. V.“, seit 2008 von der „Dr. Silvius Wodarz Stiftung“ und durch deren Fachbeirat, das „Kuratorium Baum des Jahres“ (KBJ). Im Jahr 2000 wählte die Stiftung den Ginkgo-Baum (Ginkgo biloba) zum Baum des Jahrtausends als Mahnmal für Umweltschutz und Frieden.",
								Subsections: []wiki.Section{},
							},
							wiki.Section{
								Number:      "sieben.drei",
								Title:       "Mythologie und Religion",
								Body:        "Zahlreiche Mythen erzählen von einem Lebens- oder Weltenbaum, der die Weltachse im Zentrum des Kosmos darstellt. Bei den nordischen Völkern war es zum Beispiel die Weltesche Yggdrasil, unter deren Krone die Asen ihr Gericht abhielten. So spielt der Baum in den Mythen der Völker als Lebensbaum wie die Sykomore bei den Ägyptern oder in der jüdischen Mythologie eine Rolle. Kelten, Slawen, Germanen und Balten haben einst in Götterhainen Bäume verehrt, und das Fällen solcher Götzenbäume ist der Stoff zahlreicher Legenden, die von der Missionierung Nord- und Mitteleuropas berichten.\nIn vielen alten Kulturen und Religionen wurden Bäume oder Haine als Sitz der Götter oder anderer übernatürlicher Wesen verehrt. Solche Vorstellungen haben sich als abgesunkenes religiöses Gut bis in die heutige Zeit erhalten. Als Baum der Unsterblichkeit gilt der Pfirsichbaum in China. Der Bodhibaum, unter dem Buddha Erleuchtung fand, ist im Buddhismus ein Symbol des Erwachens.\nAuch in der Bibel werden Bäume immer wieder erwähnt. Tanach wie auch das Neue Testament nennen unterschiedliche Baumarten, wie zum Beispiel den Olivenbaum oder den Feigenbaum, mit dessen relativ großen Blättern das erste Menschenpaar Adam und Eva laut 1. Mose/Genesis 3:7 nach ihrem Sündenfall ihre Blöße bedeckte. Im 1. Buch Mose, der Genesis, wird in Kapitel 1 in den Versen 11 und 12 berichtet, dass Gott die Bäume und insbesondere die fruchttragenden Bäume in seiner Schöpfung der Welt hervorbrachte. Zwei Bäume spielen in der Schöpfungsgeschichte eine entscheidende Rolle: Der Baum des Lebens und der Baum der Erkenntnis von Gut und Böse. So hat der Baum auch in der christlichen Ikonographie eine besondere Bedeutung. Dem Baum als Symbol des Sündenfalls, um dessen Stamm sich eine Schlange windet, steht häufig das hölzerne Kreuz als Symbol der Erlösung gegenüber. Ein dürrer und ein grünender Baum symbolisieren in den Dogmenallegorien der Reformationszeit den Alten und den Neuen Bund. In der Pflanzensymbolik haben verschiedene Baumarten wie auch ihre Blätter, Zweige und Früchte eine besondere Bedeutung. So weist die Akazie auf die Unsterblichkeit der menschlichen Seele hin, der Ölbaum auf den Frieden und ist ein altes marianisches Symbol für die Verkündigung an Maria. Der Zapfen der Pinie weist auf die Leben spendende Gnade und Kraft Gottes hin, die Stechpalme, aus deren Zweigen nach der Legende die Dornenkrone gefertigt war, auf die Passion Christi.",
								Subsections: []wiki.Section{},
							}, wiki.Section{
								Number:      "sieben.vier",
								Title:       "In der Geschichte",
								Body:        "Der Arbre de Diane (Dianes Baum) ist eine Platane in Les Clayes-sous-Bois, Frankreich, die 1556 von Diana von Poitiers, der Mätresse Heinrichs II., gepflanzt worden sein soll.\nGedenkbäume sind Bäume, die zum Gedenken an ein Ereignis oder zum Gedenken an eine Person gepflanzt wurden.",
								Subsections: []wiki.Section{},
							},
						},
					},
					wiki.Section{
						Number:      "acht",
						Title:       "Superlative",
						Body:        " Der höchste Baum der Welt ist der „Hyperion“, ein Küstenmammutbaum (Sequoia sempervirens) im Redwood-Nationalpark in Kalifornien mit 115,5 Meter Wuchshöhe.\n Der höchste Baum Deutschlands, vielleicht sogar des Kontinents, ist die „Waldtraut vom Mühlwald“, eine 63,33 Meter (Stand: 18. August 2008) hohe Douglasie (Pseudotsuga menziesii) im  Arboretum Freiburg-Günterstal, einem Teil des Freiburger Stadtwalds.\n Der voluminöseste Baum der Welt ist angeblich der General Sherman Tree, ein Riesenmammutbaum im Sequoia National Park, Kalifornien, USA: Volumen etwa 1489 Kubikmeter, Gewicht etwa 1385 Tonnen (US), Alter rund 2500 Jahre.\n Der dickste Baum ist der „Baum von Tule“, eine Mexikanische Sumpfzypresse (Taxodium mucronatum) in Santa María del Tule im mexikanischen Staat Oaxaca. Sein Durchmesser an der dicksten Stelle beträgt 14,05 Meter.\n Die ältesten Bäume bezogen auf einen einzelnen Baumstamm sind – gemäß verbürgter Jahresringzählung – über 4800 Jahre alte Langlebige Kiefern (Pinus longaeva, früher als Varietät der Grannen-Kiefer angesehen) in den White Mountains in Kalifornien.\n Der älteste Baum bezogen auf den lebenden Organismus ist die Amerikanische Zitterpappelkolonie „Pando“ in Utah, USA, deren Alter auf mindestens 80.000 Jahre geschätzt wird. Aus den Wurzeln sprießen immer wieder neue, genetisch identische Baumstämme (vegetative Vermehrung), die etwa 100-150 Jahre alt werden.  Bei einem Individuum der Art „Huon Pine“ in Tasmanien, das mindestens 10.500 Jahre (vielleicht sogar 50.000 Jahre) alt ist, ist der älteste Baumstamm etwa 2000 Jahre alt. Die ältesten Bäume Europas stehen in der Provinz Dalarna in Schweden. 2008 wurden dort etwa 20 gemeine Fichten auf über 8000 Jahre datiert, die älteste auf 9550 Jahre. Die einzelnen Baumstämme sterben dabei nach etwa 600 Jahren ab und werden aus der Wurzel neu gebildet.\n Die winterhärtesten Bäume sind die Dahurische Lärche (Larix gmelinii) und die Ostasiatische Zwerg-Kiefer (Pinus pumila): Sie widerstehen Temperaturen bis zu −70 °C.\n Die Dahurische Lärche ist auch die Baumart, die am weitesten im Norden überleben kann: 72° 30' N, 102° 27' O.\n Die Bäume in der größten Höhe finden sich auf 4600 Meter Seehöhe am Osthimalaya in Sichuan, dort gedeiht die Schuppenrindige Tanne (Abies squamata).\n Das Holz geringster Dichte ist jenes des Balsabaumes.\n Bäume, die bis dahin kahle Flächen besiedeln können, sogenannte Pionierpflanzen, sind zum Beispiel bestimmte Birken-, Weiden- und Pappelarten.\n In der Bonsai\u00adkunst versucht man, das Abbild eines uralten und erhabenen Baumes in klein in der Schale nachzuahmen.\n Die älteste Baumart der Erde und vermutlich das älteste lebende Fossil in der Pflanzenwelt ist der Ginkgo-Baum (Ginkgo biloba).",
						Subsections: []wiki.Section{},
					},
					wiki.Section{
						Number:      "neun",
						Title:       "Filmografie",
						Body:        " Deutschlands älteste Bäume. Dokumentation, 45 Minuten. Ein Film von Jan Haft. Produktion: Bayerischer Rundfunk, Sendung am 23. April 2007.\n Planet Erde: Waldwelten. Dokumentation, 45 Minuten. Ein Film von Alastair Fothergill. Produktion: BBC, 2006, deutsche Erstausstrahlung: ARD, am 26. März 2007.",
						Subsections: []wiki.Section{},
					},
					wiki.Section{
						Number: "zehn",
						Title:  "Literatur",
						Body:   "",
						Subsections: []wiki.Section{
							wiki.Section{
								Number:      "zehn.eins",
								Title:       "Einführungen/Übersichten",
								Body:        " Horst Bartels: Gehölzkunde. Einführung in die Dendrologie. 1. Aufl., Ulmer, Stuttgart 1993, ISBN 978-3-8252-1720-4 (Hervorragende Einführung, bestehend aus einem systematischen Teil und einem Wörterbuch der Dendrologie).\n Helmut Josef Braun: Bau und Leben der Bäume. 4. Aufl., Rombach, Freiburg 1998, ISBN 978-3-7930-9184-4 (Allgemeinverständliche und reichhaltig illustrierte Einführung in Baumanatomie und -physiologie).\n Alex L. Shigo: Die neue Baumbiologie. Fachbegriffe von A bis Z. Haymarket Media Verlag Bernhard Thalacker, Braunschweig 1990, ISBN 978-3-8781-5022-0 (Darstellung des Kompartimentkonzepts und der Wundreaktionen von Bäumen, zahlreiche Abbildungen).\n Claus Mattheck: Design in der Natur – Der Baum als Lehrmeister. 4. Neuaufl., Rombach, Freiburg im Breisgau / Berlin 1997, ISBN 978-3-7930-9470-8 (Einführung in die Baummechanik).\n Peter Schütt, Hans Joachim Schuck, Bernd Stimm: Lexikon der Forstbotanik. Morphologie, Pathologie, Ökologie und Systematik wichtiger Baum- und Straucharten. 1. Aufl., ecomed, Landsberg/Lech 1992, ISBN 978-3-609-65800-1.\n Dietrich Böhlmann: Warum Bäume nicht in den Himmel wachsen. Eine Einführung in das Leben unserer Gehölze. Quelle & Meyer, Wiebelsheim 2009, ISBN 978-3-494-01420-3.",
								Subsections: []wiki.Section{},
							},
							wiki.Section{
								Number:      "zehn.zwei",
								Title:       "Bestimmungsbücher",
								Body:        " Andreas Roloff, Andreas Bärtels: Flora der Gehölze, Bestimmung, Eigenschaften und Verwendung. 2. Aufl., Ulmer, Stuttgart 2006, ISBN 3-8001-4832-3 (Die aktuelle und zugleich umfassendste Gehölzflora, mit einem Winterbestimmungsschlüssel von Bernd Schulz).\n Ulrich Hecker: BLV Handbuch Bäume und Sträucher. BLV, München 1995, ISBN 3-405-14738-7 (Bestimmungsbuch und Nachschlagewerk in einem).\n Alan Mitchell, John Wilkinson, Peter Schütt: Pareys Buch der Bäume. Nadel- und Laubbäume in Europa nördlich des Mittelmeeres. (The Trees of Britain and Northern Europe). Paul Parey, Hamburg / Berlin 1987, ISBN 3-490-19518-3.",
								Subsections: []wiki.Section{},
							},
							wiki.Section{
								Number:      "zehn.drei",
								Title:       "Kulturgeschichte",
								Body:        " Alexander Demandt: Über allen Wipfeln. Der Baum in der Kulturgeschichte. Böhlau, Köln 2002, ISBN 3-412-13501-1.\n Doris Laudert: Mythos Baum. Was Bäume uns Menschen bedeuten. Geschichte, Brauchtum, 30 Baumporträts. BLV, München 2001, ISBN 3-405-15350-6.\n Graeme Matthews, David Bellamy: Bäume. Eine Weltreise in faszinierenden Fotos. (Trees of the World.) BLV, München 1993, ISBN 3-405-14479-5.\n Gerd und Marlene Haerkötter: Macht und Magie der Bäume. Sagen – Geschichte – Beschreibungen. Eichborn, Frankfurt am Main 1989, ISBN 3-8218-1226-5.\n Fred Hageneder: Die Weisheit der Bäume. Mythos, Geschichte, Heilkraft. Franckh-Kosmos, Stuttgart 2006, ISBN 3-440-10728-0.\n Klaus Offenberg: Das Jahrtausendtreffen: Ein Baummärchen. Agenda Verlag, 2011, ISBN 3-89688-437-9.",
								Subsections: []wiki.Section{},
							},
						},
					},
					wiki.Section{
						Number:      "elf",
						Title:       "Weblinks",
						Body:        " Baumkunde\n Schaubild zum Aufbau eines Baumstamms\n Baum des JahresInformationen über verschiedene Baumarten:\n\n Baumliste\n Bundesamt für Wald Österreich\n Bäume – für Kinder und Jugendliche\n 680 Tree Fact Sheets, University of Florida (englisch)\n GlobalTreeSearch, Botanic Gardens Conservation International (BGCI) (englisch)Informationen über seltene mitteleuropäische Baumarten:\n\n Projekt Förderung seltener Baumarten (Schweiz)",
						Subsections: []wiki.Section{},
					},
					wiki.Section{
						Number:      "zwölf",
						Title:       "Einzelnachweise",
						Body:        "",
						Subsections: []wiki.Section{},
					},
				},
			})))
		})
	})
})
