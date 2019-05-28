package spider

import (
	"encoding/csv"
	"io"
	"net/url"
	"os"
	"strings"

	log "github.com/sirupsen/logrus"
)

// MassSign takes csv files created using Crawl
func MassSign(workers int, target, contact *os.File) {
	tReader := csv.NewReader(target)
	cReader := csv.NewReader(contact)
	// skip first line
	_, err := tReader.Read()
	if err != nil {
		log.Fatal(err)
	}
	_, err = cReader.Read()
	if err != nil {
		log.Fatal(err)
	}
	var contacts []Contact
	// load contacts
	for {
		row, err := cReader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Error(err)
			continue
		}
		c := Unpack(row)
		contacts = append(contacts, c)
	}
	// loop over targets
	forms := make(chan Form)
	done := make(chan bool)
	for i := 0; i < workers; i++ {
		go worker(contacts, forms, done)
	}
	for {
		target, err := tReader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Error(err)
			break // reader problem, so break loop
		}
		u, err := url.Parse(target[0])
		if err != nil {
			log.Error(err)
			continue
		}
		f := Form{
			URL:    u,
			Action: target[1],
			Fields: strings.Split(target[2], "|"),
		}
		forms <- f
	}
	close(forms)
	<-done
}

func worker(cons []Contact, forms <-chan Form, done chan<- bool) {
	for f := range forms {
		for _, c := range cons {
			err := SignUp(f, c)
			if err != nil {
				log.Error(err)
			}
		}
		done <- true
	}
}
