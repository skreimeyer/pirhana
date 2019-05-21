// forms.go
package main

import (
	"fmt"
	"html/template"
	"net/http"
)

func main() {
	tmpl := template.Must(template.ParseFiles("main.html", "page1.html", "page2.html"))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		tmpl.ExecuteTemplate(w, "main.html", nil)
	})
	http.HandleFunc("/page1", func(w http.ResponseWriter, r *http.Request) {
		tmpl.ExecuteTemplate(w, "page1.html", nil)
	})
	http.HandleFunc("/page2", func(w http.ResponseWriter, r *http.Request) {
		tmpl.ExecuteTemplate(w, "page2.html", nil)
	})
	http.HandleFunc("/first", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("First form triggered")
		fmt.Println("Data:",
			r.FormValue("email"),
			r.FormValue("name"),
			r.FormValue("address"),
		)
		http.Redirect(w, r, "page1", 300)
	})
	http.HandleFunc("/second", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("[ERROR] Second form triggered")
		fmt.Println("Data:",
			r.FormValue("favcolor"),
			r.FormValue("name"),
			r.FormValue("address"),
		)
		http.Redirect(w, r, "page1", 300)
	})

	http.ListenAndServe(":8000", nil)
}
