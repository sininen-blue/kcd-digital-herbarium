package main

import (
	"html/template"
	"log"
	"net/http"
)

func main() {
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("statis"))))

    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        tmpl := template.Must(template.ParseFiles("./templates/index.html"))
        tmpl.Execute(w, nil)
    })

    http.HandleFunc("/search", func(w http.ResponseWriter, r *http.Request) {
        tmpl := template.Must(
            template.ParseFiles("./templates/fragments/results.html"),
        )
        data := map[string][]Stock {
            "Results": SearchTicker(r.URL.Query().Get("key")),
    }

        tmpl.Execute(w, data)
    })

    log.Println("App running on localhost:8000")
    log.Fatal(http.ListenAndServe(":8000", nil))
}
