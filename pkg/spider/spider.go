// Package spider is the web crawler at the core of the project. It deals with
// identifying potential mailing-list or registration forms that take an email
// field.
package spider

import (
	"net/url"

	"github.com/gocolly/colly"
)

type contact struct {
	First       string
	Last        string
	Street      string
	City        string
	State       string
	Email       string
	Zip         string // represent as string so we don't have to convert
	HomePhone   string
	MobilePhone string
}

type form struct {
	URL    *url.URL
	Action string
	Method string
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
func Crawl(d []string) error {
	c := colly.NewCollector(
		colly.AllowedDomains(d...),
	)
	c.OnHTML("form", func(e *colly.HTMLElement) {
		fields := e.ChildAttrs("input", "name")
		// is there an email field at all?
		for _, field := range fields {
			if field == "email" {
				f := form{
					URL:    e.Request.URL,
					Action: e.Attr("action"),
					Method: e.Attr("method"),
					Fields: fields,
				}
				// Save this to a file or sqlite.
				// Will need locking
			}
		}
	})
	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		link := e.Attr("href")
		c.Visit(e.Request.AbsoluteURL(link))
	})
	for _, site := range d {
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
// 3. Ghostbin.org
func Leak(c contact) error {
	return nil
}

// SignUp submits POST data to contact registration forms. Cookies or other
// contextual variables are not attempted. Secondary input validation, like
// captchas, are not attempted.
func SignUp(f form) error {
	return nil
}
