package main

import (
	"database/sql"
	"html/template"
	"log"
	"net/http"

    _ "github.com/lib/pq"
)

func main() {
    connStr := "postgres://postgres:secret@localhost:5432/gopgtest?sslmode=disable"

    db, err := sql.Open("postgres", connStr)
    defer db.Close()

    if err != nil {
        log.Fatal(err)
    }
    if err = db.Ping(); err != nil {
        log.Fatal(err)
    }


	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("statis"))))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		tmpl := template.Must(template.ParseFiles("./templates/index.html"))
		tmpl.Execute(w, nil)
	})

	http.HandleFunc("/search", func(w http.ResponseWriter, r *http.Request) {

	})

	log.Println("App running on localhost:8000")
	log.Fatal(http.ListenAndServe(":8000", nil))
}
