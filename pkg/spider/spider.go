// Package spider is the web crawler at the core of the project. It deals with
// identifying potential mailing-list or registration forms that take an email
// field.
package spider

import (
	"encoding/csv"
	"fmt"
	"net/url"
	"strings"
	"sync"

	"github.com/gocolly/colly"
	log "github.com/sirupsen/logrus"
)

// Contact is a collection of identifying information for an individual. Fields
// are self-explanatory. All fields, including 'numeric' fields are represented
// as strings for simplicity.
type Contact struct {
	First       string
	Last        string
	Street      string
	City        string
	State       string
	Email       string
	Zip         string
	HomePhone   string
	MobilePhone string
}

type form struct {
	URL    *url.URL
	Action string
	Fields []string
}

// Crawl traverses a pre-defined list of malicious websites and attempts to
// identify URLs containing form input for email contact information such as
// mailing lists or registrations. This is expected to be a time-consuming
// process, so results are saved and POST-related functionality is its own
// function. Failure is expected to occur frequently, due to rate limiting or
// dead links, so errors are passed over.
//
// 	spider.Crawl([]string{"google.com","yahoo.com"})
//
func Crawl(domains []string, w *csv.Writer) error {
	mx := &sync.Mutex{}
	var hosts []string
	for _, d := range domains {
		h, err := url.Parse(d)
		if err != nil {
			continue
		}
		hosts = append(hosts, h.Hostname())
	}
	c := colly.NewCollector(
		colly.AllowedDomains(hosts...),
		// FIXME: No depth limit
		colly.MaxDepth(1),
	)

	c.OnHTML("form", func(e *colly.HTMLElement) {
		fields := e.ChildAttrs("input", "name")
		// is there an email field at all?
		for _, field := range fields {
			if field == "email" {
				f := form{
					URL:    e.Request.URL,
					Action: e.Attr("action"),
					Fields: fields,
				}
				mx.Lock()
				defer mx.Unlock()
				w.Write([]string{f.URL.String(), f.Action, strings.Join(f.Fields, "|")})
			}
		}
	})
	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		link := e.Attr("href")
		if string(link[0]) == "#" {
			return
		}
		log.Info(fmt.Sprintf("Visiting:\t%s", e.Request.AbsoluteURL(link)))
		c.Visit(e.Request.AbsoluteURL(link))
	})
	for _, site := range domains {
		c.Visit(site)
	}
	return nil
}

// Leak posts a contact to locations on the internet likely to be scraped by
// others. Failure is expected to occur frequently, so errors are logged and
// then passed over. This is the low-hanging fruit of pirhana.
//
// Posts to:
// 1. Craigslist.org
// 2. Pastebin.org
func Leak(contacts []Contact) error {
	return nil
}

// SignUp submits POST data to contact registration forms. Cookies or other
// contextual variables are not attempted. Secondary input validation, like
// captchas, are not attempted.
func SignUp(f form) error {
	return nil
}
