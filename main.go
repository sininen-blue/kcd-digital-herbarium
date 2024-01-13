package main

import (
	"database/sql"
	"html/template"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("statis"))))

	r := mux.NewRouter()
	r.HandleFunc("/", indexHandler)
	r.HandleFunc("/search", searchHandler)
    r.HandleFunc("/page/ingredient/{id:[0-9]+}", ingredientDetailHandler)

	http.Handle("/", r)


	log.Println("App running on localhost:8000")
	log.Fatal(http.ListenAndServe(":8000", nil))
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	tmpl := template.Must(template.ParseFiles("./templates/index.html"))
	tmpl.Execute(w, nil)
}

type Ingredient struct {
	Name        string
	Description string
}

func searchHandler(w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open("sqlite3", "./db/herbarium.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	tmpl := template.Must(template.ParseFiles("./templates/fragments/results.html"))

	query_string := `
    select name, description from ingredients where name like ?
    union
    select name, description from potion where name like ?
    `
	query_key := "%" + r.URL.Query().Get("key") + "%"

	rows, err := db.Query(query_string, query_key, query_key)
	if err != nil {
		log.Println("Database query error")
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

	data := map[string][]Ingredient{
		"results": results,
	}
	tmpl.Execute(w, data)
}

func ingredientDetailHandler(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
	query_string := "select name, description from ingredients where rowid = ?"
	query_key := vars["id"]

}
