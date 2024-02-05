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
	r := mux.NewRouter()
	r.HandleFunc("/", indexHandler)
	r.HandleFunc("/search", searchHandler)
    r.HandleFunc("/ingredient/{name}", ingredientDetailHandler)
    r.HandleFunc("/potion/{name}", potionDetailHandler)
    r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	http.Handle("/", r)


	log.Println("App running on localhost:8000")
	log.Fatal(http.ListenAndServe(":8000", nil))
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	tmpl := template.Must(template.ParseFiles("./templates/index.html"))
	tmpl.ExecuteTemplate(w, "index", nil)
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
    defer db.Close()
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


type Potion struct {
    Type string
    Name string
    Description string
    Recipe []IngredientMap
}

type IngredientMap struct {
    Name string
    Amount string
}

type Ingredient struct {
    Type string
	Name        string
	Description string
    IngredientFor []Potion
}

func ingredientDetailHandler(w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open("sqlite3", "./db/herbarium.db")
    check_err(err)
	defer db.Close()

    vars := mux.Vars(r)

    // get the ingredient
	query_string := "select rowid, name, description from ingredients where name = ?"
	query_key := vars["name"]

    ingRow := db.QueryRow(query_string, query_key)

    var ingredientId string
    var name string
    var description string
    err = ingRow.Scan(&ingredientId, &name, &description)
    check_err(err)


    query_string = `
        SELECT potion.name, potion.description
        from ingredients
        INNER JOIN ingredientsMap
        on ingredients.rowid = ingredientsMap.ingredient_id
        INNER JOIN potion
        on ingredientsMap.potion_id = potion.rowid
        WHERE ingredients.rowid = ?
        `
    query_key = ingredientId
    mapRows, err := db.Query(query_string, query_key)
    check_err(err)

    var possible_potions []Potion
	for mapRows.Next() {
		var name string
        var description string
		err = mapRows.Scan(&name, &description)
        check_err(err)

        potion := Potion{Type: "potion", Name: name, Description: description}
        possible_potions = append(possible_potions, potion)
	}

    ingredient := Ingredient{Type: "ingredient", Name: name, Description: description, IngredientFor: possible_potions}

    tmpl := template.Must(template.ParseFiles("./templates/index.html"))
    if r.Header.Get("HX-Request") == "true" {
        tmpl.ExecuteTemplate(w, "ingredientDetail" ,ingredient)
    } else {
        tmpl.ExecuteTemplate(w, "index", ingredient)
    }
}


func potionDetailHandler(w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open("sqlite3", "./db/herbarium.db")
    check_err(err)
	defer db.Close()

    vars := mux.Vars(r)
	query_string := "select rowid, name, description from potion where name = ?"
	query_key := vars["name"]

    potRow := db.QueryRow(query_string, query_key)
    var potionId string
    var name string
    var description string
    err = potRow.Scan(&potionId, &name, &description)
    check_err(err)

    query_string = `
        SELECT ingredients.name, ingredientsMap.amount
        FROM ingredientsMap
        INNER JOIN ingredients
        ON ingredients.rowid = ingredientsMap.ingredient_id
        INNER JOIN potion
        ON potion.rowid = ingredientsMap.potion_id
        WHERE potion.rowid = ?;
        `
    query_key = potionId
    mapRows, err := db.Query(query_string, query_key)
    check_err(err)

    var ingredientsUsed []IngredientMap
	for mapRows.Next() {
		var name string
        var amount string
		err = mapRows.Scan(&name, &amount)
        check_err(err)

        ingredient := IngredientMap{Name: name, Amount: amount}
        ingredientsUsed = append(ingredientsUsed, ingredient)
	}

    potion := Potion{Type: "potion", Name: name, Description: description, Recipe: ingredientsUsed}

    tmpl := template.Must(template.ParseFiles("./templates/index.html"))
    if r.Header.Get("HX-Request") == "true" {
        tmpl.ExecuteTemplate(w, "potionDetail", potion)
    } else {
        tmpl.ExecuteTemplate(w, "index", potion)
    }
}
