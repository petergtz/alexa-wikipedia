package mediawiki

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/pkg/errors"

	"github.com/jaytaylor/html2text"
	wiki "github.com/petergtz/alexa-wikipedia"
)

type MediaWiki struct {
}

func (mw *MediaWiki) GetPage(word string) (wiki.Page, error) {
	r, e := http.Get("https://de.wikipedia.org/w/api.php?format=json&action=query&prop=extracts&titles=" + word + "&redirects=true&formatversion=2")
	// r, e := http.Get("https://de.wikipedia.org/w/api.php?action=query&titles=" + intent.Slots["word"].Value + "&prop=revisions&rvprop=content&format=json&formatversion=2")
	// r, e := http.Get("https://de.wikipedia.org/w/api.php?action=parse&page=" + intent.Slots["word"].Value + "&contentmodel=wikitext&section=0&prop=text|sections&format=json")
	if e != nil {
		return wiki.Page{}, errors.Wrap(e, "Could not request Wikipedia page")
	}
	content, e := ioutil.ReadAll(r.Body)
	if e != nil {
		return wiki.Page{}, errors.Wrap(e, "Could not read body of wikipedia page")
	}
	// log.Infof("%s", content)
	t := struct {
		Query struct {
			Pages []struct {
				Extract string
			}
		}
	}{}
	// t := struct {
	// 	Query struct {
	// 		Pages []struct {
	// 			Revisions []struct {
	// 				Content string
	// 			}
	// 		}
	// 	}
	// }{}
	// t := struct {
	// 	Parse struct {
	// 		Text struct {
	// 			Body string `json:"*"`
	// 		}
	// 	}
	// }{}
	e = json.Unmarshal(content, &t)
	if e != nil {
		return wiki.Page{}, errors.Wrap(e, "Could not unmarshal body of wikipedia page")
	}
	// log.Info(t.Parse.Text.Body)
	// log.Info(t.Query.Pages[0].Revisions[0].Content)
	// article, e := gowiki.ParseArticle("Bla", t.Query.Pages[0].Revisions[0].Content, &gowiki.DummyPageGetter{})
	// log.Debugf("%#v", article.GetText())
	// gowiki.ParseArticle(title string, text string, g gowiki.PageGetter)

	// article := strings.Replace(html2text.HTML2Text(t.Query.Pages[0].Extract), "\r\n", "\n", -1)
	// html2text.FromString(input string, options ...html2text.Options)

	text, e := html2text.FromString(t.Query.Pages[0].Extract, html2text.Options{OmitLinks: true})
	// text, e := html2text.FromString(t.Parse.Text.Body, html2text.Options{OmitLinks: true})
	if e != nil {
		panic(e)
	}
	article := strings.Replace(text, "\r\n", "\n", -1)
	return wiki.Page{
		Sections: []wiki.Section{wiki.Section{Body: article}},
	}, nil
}
