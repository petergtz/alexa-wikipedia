package mediawiki

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"regexp"
	"strings"

	"github.com/pkg/errors"

	wiki "github.com/petergtz/alexa-wikipedia"
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

func (mw *MediaWiki) GetPage(word string) (wiki.Page, error) {
	extract := ExtractQuery{}
	e := makeJsonRequest("https://de.wikipedia.org/w/api.php?format=json&action=query&prop=extracts&titles="+url.QueryEscape(word)+"&redirects=true&formatversion=2&explaintext=true", &extract)
	if e != nil {
		return wiki.Page{}, e
	}
	if extract.Query.Pages[0].Missing {
		return wiki.Page{}, errors.New("Page not found on Wikipedia")
	}
	return WikiPageFrom(extract.Query.Pages[0]), nil
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

func WikiPageFrom(mediawikipage Page) wiki.Page {
	page := wiki.Page{
		Title: mediawikipage.Title,
		// It's not obvious, but suffixing a \n to the text helps with regexes below
		Body: mediawikipage.Extract + "\n",
	}
	parse((*wiki.Section)(&page), 2)
	return page
}

func parse(section *wiki.Section, level int) {
	var sectionTitleRegex = regexp.MustCompile(fmt.Sprintf("\n={%v} (.*) ={%v}\n", level, level))
	sectionTitles := sectionTitleRegex.FindAllStringSubmatch(section.Body, -1)
	sections := sectionTitleRegex.Split(section.Body, -1)
	section.Body = strings.Trim(sections[0], "\n")
	section.Subsections = make([]wiki.Section, len(sections)-1)
	for i, s := range sections[1:] {
		if level == 2 {
			section.Subsections[i].Number = wiki.Convert(i + 1)
		} else {
			section.Subsections[i].Number = section.Number + "." + wiki.Convert(i+1)
		}
		section.Subsections[i].Title = strings.TrimSuffix(strings.TrimPrefix(strings.Trim(sectionTitles[i][1], "\n"), "== "), " ==")
		section.Subsections[i].Body = s
		parse(&section.Subsections[i], level+1)
	}
}
