// forms.go
package main

import (
	"fmt"
	"html/template"
	"net/http"
)

// ContactDetails Yep
type ContactDetails struct {
	Email   string
	Subject string
	Message string
}

func main() {
	tmpl := template.Must(template.ParseFiles("main.html"))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			tmpl.Execute(w, nil)
			return
		}

		details := ContactDetails{
			Email:   r.FormValue("email"),
			Subject: r.FormValue("subject"),
			Message: r.FormValue("message"),
		}

		fmt.Printf("Received form with fields: Email:%s | Subject:%s | Message:%s\n", details.Email, details.Subject, details.Message)

		tmpl.Execute(w, struct{ Success bool }{true})
	})
	http.HandleFunc("/callout", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Received callout")
		fmt.Println("Data:",
			r.FormValue("email"),
			r.FormValue("firstname"),
			r.FormValue("last-name"))
	})

	http.ListenAndServe(":8080", nil)
}
