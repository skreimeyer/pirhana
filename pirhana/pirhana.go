package main

//go:generate go run ../scripts/addstatic.go

import (
	"bufio"
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/skreimeyer/pirhana/pkg/spider"
)

func main() {
	// Flags and help text
	crawl := flag.Bool("crawl", false,
		`Crawls malicious sites for forms. Saves to csv. The spider makes
		requests to sites associated with spam and ransomware, so this should be
		done within	a virtual machine to reduce risk of infection of the host
		machine.

		There are over one thousand sites, so this could take a long time`,
	)
	signup := flag.Bool("signup", false,
		`Fills out all forms in targets.csv with information from contacts.csv.
		Both of these files must be present to function.`,
	)
	leak := flag.Bool("leak", false,
		`Posts information for all contacts to multiple sites which are common
		targets for scrapers. (craigslist, pastebin)`,
	)
	entry := flag.Bool("entry", false,
		`Data-entry mode for contacts. Typical spreadsheet software is probably
		a more efficient way of editing a contact list. Saves to contacts.csv`,
	)
	verbose := flag.Bool("v", false, "verbose output")

	flag.Parse()

	// Setup logging
	logfile := "log.log"
	f, err := os.OpenFile(logfile, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
	Formatter := new(log.TextFormatter)
	Formatter.TimestampFormat = "2006-01-02 15:04:05"
	Formatter.FullTimestamp = true
	log.SetFormatter(Formatter)
	if err != nil {
		fmt.Println(err)
	}
	if *verbose {
		mw := io.MultiWriter(os.Stdout, f)
		log.SetOutput(mw)
	} else {
		log.SetOutput(f)
	}

	// Select execution options
	if *crawl {
		startCrawler()
	}
	if *signup {
		t, err := os.Open("targets.csv")
		if err != nil {
			log.Fatal("Failed to open targets.csv")
		}
		c, err := os.Open("contacts.csv")
		if err != nil {
			log.Fatal("Failed to open contacts.csv", err)
		}
		spider.MassSign(8, t, c)
		fmt.Println("Signup complete") // maybe do something more interesting.
	}
	if *leak {
		cmd := "cat contacts.csv | curl -F 'clbin=<-' https://clbin.com"
		out, err := exec.Command("bash", "-c", cmd).Output()
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(string(out))
	}
	if *entry {
		enter()
	}

}

func startCrawler() {
	fmt.Println("Starting crawler...")
	file, err := os.Create("targets.csv")
	if err != nil {
		log.Fatal("Cannot create savefile. Does we have write permission?")
		return
	}
	defer file.Close()
	w := csv.NewWriter(file)
	defer w.Flush()
	w.Write([]string{
		"url",
		"action",
		"fields",
	})
	sitelist := strings.Split(ransomware, "\n")
	sitelist = append(sitelist, strings.Split(suspicious, "\n")...)
	// //TEST ONLY
	// sitelist := []string{
	// 	"localhost:8080",
	// 	"localhost:8000",
	// }
	// //TEST ONLY
	log.Info("START CRAWLER")
	spider.Crawl(sitelist, w)
	fmt.Println("END CRAWL")
}

// enter handles data entry. Not very necessary, but possibly convenient.
func enter() {
	filename := "contacts.csv"
	newFile := false // flag for creating a new csv
	file, err := os.OpenFile(filename, os.O_WRONLY|os.O_APPEND, 0644)
	// TODO: test for and instantiate the file elsewhere.
	if os.IsNotExist(err) {
		file, err = os.Create(filename)
		if err != nil {
			log.Fatal(err)
		}
		newFile = true
	}
	defer file.Close()
	w := csv.NewWriter(file)
	if newFile {
		w.Write([]string{"First", "Last", "Addr", "City", "State", "Email",
			"Zip", "Phone", "Mobile"})
	}
	defer w.Flush()
	// Get user input
	fmt.Println("Create a new contact entry")
	s := bufio.NewScanner(os.Stdin)
	var data []string
	fields := []string{
		"First Name",
		"Last Name",
		"Address (number and street)",
		"City",
		"State",
		"Email",
		"Zip",
		"Home phone number",
		"Mobile phone number",
	}
	for _, f := range fields {
		fmt.Print(f + ":\t")
		s.Scan()
		datum := s.Text()
		data = append(data, datum)
	}
	err = w.Write(data)
	if err != nil {
		fmt.Println(err)
		log.Fatal(err)
	}

}
