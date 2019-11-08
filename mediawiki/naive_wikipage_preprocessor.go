package mediawiki

import (
	"regexp"

	"github.com/petergtz/go-alexa"

	"github.com/pkg/math"
)

type Persistence interface{ Persist([]string) error }

type ErrorReporter interface {
	ReportPanic(interface{}, *alexa.RequestEnvelope)
}

type HighlightMissingSpacesNaivelyWikiPagePreProcessor struct {
	Persistence   Persistence
	pageChan      chan *Page
	errorReporter ErrorReporter
}

func NewHighlightMissingSpacesNaivelyWikiPagePreProcessor(persistence Persistence, errorReporter ErrorReporter) *HighlightMissingSpacesNaivelyWikiPagePreProcessor {
	pp := &HighlightMissingSpacesNaivelyWikiPagePreProcessor{
		Persistence:   persistence,
		pageChan:      make(chan *Page, 100),
		errorReporter: errorReporter,
	}
	go func() {
		for page := range pp.pageChan {
			pp.doProcess(page)
		}
	}()
	return pp
}

var pattern = regexp.MustCompile(`[a-z]\.[A-Z]`)

func (pp *HighlightMissingSpacesNaivelyWikiPagePreProcessor) Process(page *Page) *Page {
	pp.pageChan <- page
	return page
}

func (pp *HighlightMissingSpacesNaivelyWikiPagePreProcessor) doProcess(page *Page) {
	defer func() {
		if e := recover(); e != nil {
			pp.errorReporter.ReportPanic(e, nil)
		}
	}()

	var findings []string
	for _, index := range pattern.FindAllStringIndex(page.Extract, -1) {
		finding := page.Extract[math.MaxInt(index[0]-10, 0):math.MinInt(index[1]+10, len(page.Extract))]
		findings = append(findings, finding)
	}
	if len(findings) > 0 {
		e := pp.Persistence.Persist(findings)
		PanicOnError(e)
	}
}

func PanicOnError(e error) {
	if e != nil {
		panic(e)
	}
}
