package mediawiki

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/pkg/errors"
	"go.uber.org/zap"

	"github.com/petergtz/alexa-wikipedia/locale"
	"github.com/petergtz/alexa-wikipedia/wiki"
)

type WikiPagePreProcessor interface{ Process(*Page) *Page }

type MediaWiki struct {
	Logger               *zap.SugaredLogger
	WikiPagePreProcessor WikiPagePreProcessor
}

type Page struct {
	Title   string
	Extract string
	Missing bool
}

type ExtractQuery struct {
	Query struct {
		Pages []Page
	}
}

type RevisionsQuery struct {
	Query struct {
		Pages []struct {
			Revisions []struct {
				Content string
			}
		}
	}
}

type ParseQuery struct {
	Parse struct {
		Text struct {
			Body string `json:"*"`
		}
	}
}

type SearchQuery struct {
	Query struct {
		Search []struct {
			Title  string
			Pageid int
		}
	}
}

func (mw *MediaWiki) GetPage(word string, localizer *locale.Localizer) (wiki.Page, error) {
	return mw.getPage("titles="+url.QueryEscape(strings.Title(word)), localizer)
}

func (mw *MediaWiki) SearchPage(word string, localizer *locale.Localizer) (wiki.Page, error) {
	var search SearchQuery
	e := makeJsonRequest("https://"+localizer.WikiEndpoint()+"/w/api.php?format=json&action=query&list=search&srsearch="+url.QueryEscape(word)+"&srprop=&utf8=&srlimit=1", &search, mw.Logger)
	if e != nil {
		return wiki.Page{}, e
	}
	if len(search.Query.Search) == 0 {
		return wiki.Page{}, errors.New("Page not found on Wikipedia")
	}
	return mw.getPage("pageids="+strconv.Itoa(search.Query.Search[0].Pageid), localizer)
}

func (mw *MediaWiki) getPage(query string, localizer *locale.Localizer) (wiki.Page, error) {
	var extract ExtractQuery
	e := makeJsonRequest("https://"+localizer.WikiEndpoint()+"/w/api.php?format=json&action=query&prop=extracts&"+query+"&redirects=true&formatversion=2&explaintext=true&exlimit=1", &extract, mw.Logger)
	if e != nil {
		return wiki.Page{}, e
	}
	if extract.Query.Pages[0].Missing {
		return wiki.Page{}, errors.New("Page not found on Wikipedia")
	}
	return WikiPageFrom(mw.WikiPagePreProcessor.Process(&extract.Query.Pages[0]), localizer), nil
}

func makeJsonRequest(url string, data interface{}, logger *zap.SugaredLogger) error {
	logger = logger.With("url", url)
	logger.Debug("Before http Get")
	startTime := time.Now()
	r, e := http.Get(url)
	logger.Debugw("After http Get", "duration", time.Since(startTime).String())
	if e != nil {
		return errors.Wrapf(e, "Could not request url: \"%v\"", url)
	}
	logger.Debug("Before read body")
	startTime = time.Now()
	content, e := ioutil.ReadAll(r.Body)
	logger.Debugw("After read body", "duration", time.Since(startTime).String(), "body-size", len(content))
	if e != nil {
		return errors.Wrap(e, "Could not read body of page")
	}
	logger.Debug("Before json Unmarhsal")
	startTime = time.Now()
	e = json.Unmarshal(content, data)
	logger.Debugw("After json Unmarshal", "duration", time.Since(startTime).String())
	if e != nil {
		return errors.Wrapf(e, "Could not unmarshal body of page. body was: \"%v\"", content)
	}
	return nil
}

func WikiPageFrom(mediawikipage *Page, localizer *locale.Localizer) wiki.Page {
	page := wiki.Page{
		Title: mediawikipage.Title,
		// It's not obvious, but suffixing a \n to the text helps with regexes below
		Body: mediawikipage.Extract + "\n",
	}
	parse((*wiki.Section)(&page), 2, localizer)
	return page
}

func parse(section *wiki.Section, level int, localizer *locale.Localizer) {
	var sectionTitleRegex = regexp.MustCompile(fmt.Sprintf("\n={%v} (.*) ={%v}\n", level, level))
	sectionTitles := sectionTitleRegex.FindAllStringSubmatch(section.Body, -1)
	sections := sectionTitleRegex.Split(section.Body, -1)
	section.Body = strings.Trim(sections[0], "\n")
	section.Subsections = make([]wiki.Section, len(sections)-1)
	for i, s := range sections[1:] {
		if level == 2 {
			section.Subsections[i].Number = localizer.Spell(i + 1)
		} else {
			section.Subsections[i].Number = section.Number + "." + localizer.Spell(i+1)
		}
		section.Subsections[i].Title = strings.TrimSuffix(strings.TrimPrefix(strings.Trim(sectionTitles[i][1], "\n"), "== "), " ==")
		section.Subsections[i].Body = s
		parse(&section.Subsections[i], level+1, localizer)
	}
}
