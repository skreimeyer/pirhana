package main

import (
	"bufio"
	"os"
	"strings"
)

func main() {
	f, err := os.Open("../static/suspicious")
	if err != nil {
		panic(err)
	}
	o, err := os.Create("../static/suspiciousFormatted")
	if err != nil {
		panic(err)
	}
	s := bufio.NewScanner(f)
	for s.Scan() {
		o.WriteString(strings.Split(strings.TrimSpace(s.Text()), "\t")[0] + "\n")
	}

}
