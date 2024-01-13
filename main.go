package main

import (
	"database/sql"
	"html/template"
	"log"
	"net/http"

	_ "github.com/mattn/go-sqlite3"
)


type Ingredient struct {
    Name string
    Description string
}

func main() {
    db, err := sql.Open("sqlite3", "./db/herbarium.db")
    if err != nil {
        log.Fatal(err)
    }
    defer db.Close()

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("statis"))))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		tmpl := template.Must(template.ParseFiles("./templates/index.html"))
		tmpl.Execute(w, nil)
	})

	http.HandleFunc("/search", func(w http.ResponseWriter, r *http.Request) {
        tmpl := template.Must(template.ParseFiles("./templates/fragments/results.html"))

        query_string := "select name, description from ingredients where name like ?"
        query_key := "%"+r.URL.Query().Get("key")+"%"

        rows, err := db.Query(query_string, query_key)
        if err != nil {
            log.Println("first")
            log.Fatal(err)
        }
        defer rows.Close()

        var results []Ingredient
        for rows.Next() {
            var name string
            var description string
            err = rows.Scan(&name, &description)
            if err != nil {
                log.Fatal(err)
            }
            ingredient := Ingredient{Name: name, Description: description}
            results = append(results, ingredient)
        }

        data := map[string][]Ingredient {
            "results": results,
        }
        tmpl.Execute(w, data)

	})

	log.Println("App running on localhost:8000")
	log.Fatal(http.ListenAndServe(":8000", nil))
}
