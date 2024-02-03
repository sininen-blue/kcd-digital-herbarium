package main

import (
	"html/template"
	"log"
	"net/http"
    "database/sql"

	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"
)

var templ *template.Template
var db *sql.DB

type Potion struct {
	Name        string
	Description string
	Effects     string
	url         string
}

type Ingredient struct {
    Id string
    Name string
    Description string
	url         string
}

func getAllIngredients() []Ingredient {
    var results []Ingredient

    query := "select * from ingredients"
    rows, err := db.Query(query)
    if err != nil {
        log.Fatal(err)
    }
    defer rows.Close()

    for rows.Next() {
        var i Ingredient
        err = rows.Scan(
            &i.Id,
            &i.Name,
            &i.Description,
        )
        if err != nil {
            log.Fatal(err)
        }
        results = append(results, i)
    }

    return results
}

func getIngredients(query string, key string) []Ingredient{
    var results []Ingredient

    rows, err := db.Query(query, key)
    if err != nil {
        log.Fatal(err)
    }
    defer rows.Close()

    return results
}

func getPotions(query string) []Potion {
	var results []Potion

	return results
}

func init() {
	var err error

	funcMap := template.FuncMap{
		"inc": func(i int) int { return i + 1 },
        "dec": func(i int) int { return i - 1 },
	}

    templ = template.New("template")
    templ = templ.Funcs(funcMap)
	templ, err = templ.ParseGlob("./templates/*.html")
	if err != nil {
		log.Println("template error")
		log.Fatal(err)
	}

    db, err = sql.Open("sqlite3", "./db/herbarium.db")
    if err != nil {
        log.Fatal(err)
    }
}

func main() {
    defer db.Close()

	r := mux.NewRouter()
	r.HandleFunc("/", indexHandler)
	http.Handle("/", r)

	log.Println("App running on localhost:8000")
	log.Fatal(http.ListenAndServe(":8000", nil))
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
    data := map[string]interface{} {
        "Title": "Index",
        "Ingredients": getAllIngredients(),
    }

    templ.ExecuteTemplate(w, "base", data)
}
