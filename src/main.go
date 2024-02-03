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

var currentInv []invItem
type invItem struct {
    Ingredient Ingredient
    Amount int
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

// not done
func getIngredients(query string, key string) []Ingredient{
    var results []Ingredient

    rows, err := db.Query(query, key)
    if err != nil {
        log.Fatal(err)
    }
    defer rows.Close()

    return results
}

func getIngredient(query string, key string) Ingredient{
    var result Ingredient

    row := db.QueryRow(query, key)
    err := row.Scan(
        &result.Id,
        &result.Name,
        &result.Description,
    )
    if err != nil {
        log.Fatal(err)
    }

    return result
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
	r.HandleFunc("/", indexHandler).Methods("GET")
	r.HandleFunc("/", addInventory).Methods("POST")
	http.Handle("/", r)

	log.Println("App running on localhost:8000")
	log.Fatal(http.ListenAndServe(":8000", nil))
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
    data := map[string]interface{} {
        "Title": "Index",
        "Ingredients": getAllIngredients(),
        "Inventory": currentInv,
    }

    templ.ExecuteTemplate(w, "base", data)
}

func addInventory(w http.ResponseWriter, r *http.Request) {
    key := r.FormValue("ingredient")

    query := "select name from ingredients where name == ?"
    ing := getIngredient(query, key)
    item := invItem{Ingredient: ing, Amount: 1} 

    currentInv = append(currentInv, item)

    templ.ExecuteTemplate(w, "invItem", item)
}
