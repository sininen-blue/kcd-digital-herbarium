package main

import (
	"html/template"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"
)

var templ *template.Template

type Page struct {
    Title string
}

type Potion struct {
	Name        string
	Description string
	Effects     string
	url         string
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

}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/", indexHandler)
	http.Handle("/", r)

	log.Println("App running on localhost:8000")
	log.Fatal(http.ListenAndServe(":8000", nil))
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
    page := Page{Title: "test"}

    data := map[string]interface{} {
        "Page": page,
    }

    templ.ExecuteTemplate(w, "base", data)
}
