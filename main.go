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
    r.HandleFunc("/ingredient/{name}", ingredientDetailHandler)

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

type SearchResult struct {
	Name        string
	Description string
    Category string
}

func check_err(err error) {
    if err != nil {
        log.Fatal(err)
    }
}

func searchHandler(w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open("sqlite3", "./db/herbarium.db")
    check_err(err)

	tmpl := template.Must(template.ParseFiles("./templates/fragments/results.html"))

	query_string := "select name, description from ingredients where name like ?"
	query_key := "%" + r.URL.Query().Get("key") + "%"

	ingRows, err := db.Query(query_string , query_key)
    check_err(err)
	defer ingRows.Close()

	var results []SearchResult

	for ingRows.Next() {
		var name string
		var description string
		err = ingRows.Scan(&name, &description)
        check_err(err)

        ingredient := SearchResult{Name: name, Description: description, Category: "ingredient"}
		results = append(results, ingredient)
	}

	query_string = "select name, description from potion where name like ?"
	potRows, err := db.Query(query_string , query_key)
    check_err(err)
	defer ingRows.Close()

	for potRows.Next() {
		var name string
		var description string
		err = potRows.Scan(&name, &description)
        check_err(err)

        potion := SearchResult{Name: name, Description: description, Category: "potion"}
		results = append(results, potion)
	}

	data := map[string][]SearchResult{
		"results": results,
	}
	tmpl.Execute(w, data)
}

func ingredientDetailHandler(w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open("sqlite3", "./db/herbarium.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

    vars := mux.Vars(r)
	query_string := "select name, description from ingredients where name = ?"
	query_key := vars["name"]


    row := db.QueryRow(query_string, query_key)

    var name string
    var description string
    err = row.Scan(&name, &description)
    if err != nil {
        log.Fatal(err)
    }
    ingredient := Ingredient{Name: name, Description: description}


	tmpl := template.Must(template.ParseFiles("./templates/fragments/ingredient-detail.html"))
    tmpl.Execute(w, ingredient)
}
