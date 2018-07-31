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

	"github.com/pkg/errors"

	"github.com/petergtz/alexa-wikipedia/locale"
	"github.com/petergtz/alexa-wikipedia/wiki"
)

type MediaWiki struct {
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
	var extract ExtractQuery
	e := makeJsonRequest("https://"+localizer.WikiEndpoint()+"/w/api.php?format=json&action=query&prop=extracts&titles="+url.QueryEscape(word)+"&redirects=true&formatversion=2&explaintext=true", &extract)
	if e != nil {
		return wiki.Page{}, e
	}
	if extract.Query.Pages[0].Missing {
		return wiki.Page{}, errors.New("Page not found on Wikipedia")
	}
	return WikiPageFrom(extract.Query.Pages[0], localizer), nil
}

func (mw *MediaWiki) SearchPage(word string, localizer *locale.Localizer) (wiki.Page, error) {
	var search SearchQuery
	e := makeJsonRequest("https://"+localizer.WikiEndpoint()+"/w/api.php?format=json&action=query&list=search&srsearch="+url.QueryEscape(word)+"&srprop=&utf8=", &search)
	if e != nil {
		return wiki.Page{}, e
	}
	if len(search.Query.Search) == 0 {
		return wiki.Page{}, errors.New("Page not found on Wikipedia")
	}

	var extract ExtractQuery
	e = makeJsonRequest("https://"+localizer.WikiEndpoint()+"/w/api.php?format=json&action=query&prop=extracts&pageids="+strconv.Itoa(search.Query.Search[0].Pageid)+"&redirects=true&formatversion=2&explaintext=true", &extract)
	if e != nil {
		return wiki.Page{}, e
	}
	if extract.Query.Pages[0].Missing {
		return wiki.Page{}, errors.New("Page not found on Wikipedia")
	}
	return WikiPageFrom(extract.Query.Pages[0], localizer), nil

}

func makeJsonRequest(url string, data interface{}) error {
	r, e := http.Get(url)
	if e != nil {
		return errors.Wrapf(e, "Could not request url: \"%v\"", url)
	}
	content, e := ioutil.ReadAll(r.Body)
	if e != nil {
		return errors.Wrap(e, "Could not read body of page")
	}
	e = json.Unmarshal(content, data)
	if e != nil {
		return errors.Wrapf(e, "Could not unmarshal body of page. body was: \"%v\"", content)
	}
	return nil
}

func WikiPageFrom(mediawikipage Page, localizer *locale.Localizer) wiki.Page {
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
