package main

import (
	"html/template"
	"log"
	"net/http"
)

func main() {
	renderPage()
	servePages()
}

func renderPage() {
	tmpl := template.Must(template.ParseFiles("index.html"))
	pagePath := "/"
	http.HandleFunc(pagePath, func(w http.ResponseWriter, r *http.Request) {
		if err := tmpl.Execute(w, nil); err != nil {
			log.Fatalln("Failed template execution %v", err)
		}
	})
}

func servePages() {
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalln("Receive execute listen and serve failed %v", err)
	}
}
