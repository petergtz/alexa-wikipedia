package mediawiki

import (
	"regexp"

	"github.com/pkg/math"
)

type Persistence interface {
	Persist([]string)
}

type HighlightMissingSpacesNaivelyWikiPagePreProcessor struct {
	Persistence Persistence
	pageChan    chan *Page
}

func NewHighlightMissingSpacesNaivelyWikiPagePreProcessor(Persistence Persistence) *HighlightMissingSpacesNaivelyWikiPagePreProcessor {
	pp := &HighlightMissingSpacesNaivelyWikiPagePreProcessor{
		Persistence: Persistence,
		pageChan:    make(chan *Page, 100),
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

func (pp *HighlightMissingSpacesNaivelyWikiPagePreProcessor) doProcess(page *Page) *Page {
	var findings []string
	for _, index := range pattern.FindAllStringIndex(page.Extract, -1) {
		finding := page.Extract[math.MaxInt(index[0]-10, 0):math.MinInt(index[1]+10, len(page.Extract))]
		findings = append(findings, finding)
	}
	if len(findings) > 0 {
		pp.Persistence.Persist(findings)
	}

	return page
}
