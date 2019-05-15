// Package spider is the web crawler at the core of the project. It deals with
// identifying potential mailing-list or registration forms that take an email
// field.
package spider

import "github.com/gocolly/colly"

// "github.com/gocolly/colly"
// "net/http"
// "fmt"

func crawl(d []string) {
	c := colly.NewCollector(
		colly.AllowedDomains(d...),
	)
	for _, site := range d {
		c.Visit(site)
	}
}
