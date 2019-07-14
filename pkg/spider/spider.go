// Package spider is the web crawler at the core of the project. It deals with
// identifying potential mailing-list or registration forms that take an email
// field.
package spider

import (
	"encoding/csv"
	"fmt"
	"net/http"
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

// Form is a container for HTML form inputs. Input tag names are saved in the
// Fields slice.
type Form struct {
	URL    *url.URL
	Action string
	Fields []string
}

// Crawl traverses a pre-defined list of malicious websites and attempts to
// identify URLs containing Form input for email contact inFormation such as
// mailing lists or registrations. This is expected to be a time-consuming
// process, so results are saved and POST-related functionality is its own
// function. Failure is expected to occur frequently, due to rate limiting or
// dead links, so errors are passed over.
//
// 	spider.Crawl([]string{"google.com","yahoo.com"})
//
func Crawl(domains []string, w *csv.Writer) error {
	log.Info("Executing crawl")
	mx := &sync.Mutex{}

	c := colly.NewCollector(
		colly.AllowedDomains(domains...),
		colly.Async(true),
		colly.MaxDepth(1),
	)

	c.OnRequest(func(r *colly.Request) {
		log.Info("Request:", r.URL.String())
	})

	c.OnResponse(func(r *colly.Response) {
		log.Info(fmt.Sprintf("Response URL:%s Code:%d", r.Request.URL.String(), r.StatusCode))
	})

	c.OnHTML("form", func(e *colly.HTMLElement) {
		// log.Info("form found")
		fields := e.ChildAttrs("input", "name")
		for _, field := range fields {
			if field == "email" {
				f := Form{
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
		if len(link) > 0 && string(link[0]) == "#" {
			return
		}
		log.Info(fmt.Sprintf("Following: %s", e.Request.AbsoluteURL(link)))
		c.Visit(e.Request.AbsoluteURL(link))
	})
	for _, site := range domains {
		c.Visit("http://" + site)
		// c.Wait()
	}
	c.Wait()
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

// SignUp submits POST data to contact registration Forms. Cookies or other
// contextual variables are not attempted. Secondary input validation, like
// captchas, are not attempted.
func SignUp(f Form, c Contact) error {
	// check f.Action for full domain name. If not, prepend domain from URL
	var act string
	if !strings.HasPrefix(f.Action, "http") {
		act = "http://" + f.URL.Hostname()
		if f.URL.Port() != "" {
			act += ":" + f.URL.Port() + "/" + f.Action
		} else {
			act += "/" + f.Action
		}
	} else {
		act = f.Action
	}
	data := url.Values{}
	for _, fd := range f.Fields {
		data.Add(fd, matcher(fd, c))
	}
	_, err := http.PostForm(act, data) // TODO optional error checking for verbose
	if err != nil {
		log.Error(err)
		return err
	}
	return nil
}

// matcher is a helper function for SignUp. It matches contact fields based on
// probable input field names from our forms.
func matcher(field string, c Contact) string {
	f := strings.ToLower(field)

	switch true {
	case strings.Contains(f, "first"):
		return c.First
	case strings.Contains(f, "last"):
		return c.Last
	case strings.Contains(f, "name"):
		return fmt.Sprintf("%s %s", c.First, c.Last)
	case strings.Contains(f, "add"):
		return c.Street
	case strings.Contains(f, "city"):
		return c.City
	case strings.Contains(f, "state"):
		return c.State
	case strings.Contains(f, "zip"):
		return c.Zip
	case strings.Contains(f, "email"):
		return c.Email
	case strings.Contains(f, "home"):
		return c.HomePhone
	case strings.Contains(f, "mobile"):
		return c.MobilePhone
	case strings.Contains(f, "cell"):
		return c.MobilePhone
	case strings.Contains(f, "phone"):
		return c.MobilePhone
	default:
		return ""
	}
}

// Unpack assigns elements in a string slice to a new Contact
func Unpack(arg []string) Contact {
	if len(arg) == 9 {
		return Contact{
			First:       arg[0],
			Last:        arg[1],
			Street:      arg[2],
			City:        arg[3],
			State:       arg[4],
			Email:       arg[5],
			Zip:         arg[6],
			HomePhone:   arg[7],
			MobilePhone: arg[8],
		}
	}
	return Contact{}
}
