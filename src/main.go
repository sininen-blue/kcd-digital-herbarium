package main

import (
	"database/sql"
	"html/template"
	"log"
	"net/http"
	"slices"

	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"
)

var templ *template.Template
var db *sql.DB

type Potion struct {
	Id          string
	Name        string
	Description string
	Effects     string
	Url         string
}

type IngredientMap struct {
	IngredientName string
	Amount         int
	Got            bool
}

func getAllPotions() []Potion {
	var possibles []Potion
	var results []Potion

	inventory := getAllInventory()

	for _, item := range inventory {
		query := `
            select potions.* from potions
            inner join ingredientsMap on ingredientsMap.potion_id = potions.rowid
            inner join ingredients on ingredients.rowid = ingredientsMap.ingredient_id
            where ingredients.name = ?;
        `
		rows, err := db.Query(query, item.IngredientName)
		if err != nil {
			log.Fatal(err)
		}
		defer rows.Close()

		for rows.Next() {
			var i Potion
			err = rows.Scan(
				&i.Id,
				&i.Name,
				&i.Description,
				&i.Effects,
			)
			if err != nil {
				log.Fatal(err)
			}
			possibles = append(possibles, i)
		}
	}

	var ingredientsNeeded []IngredientMap
	for _, potion := range possibles {
		query := `
            select ingredients.name, ingredientsMap.amount from potions
            inner join ingredientsMap on ingredientsMap.potion_id = potions.rowid
            inner join ingredients on ingredients.rowid = ingredientsMap.ingredient_id
            where potions.name = ?;
        `
		rows, err := db.Query(query, potion.Name)
		if err != nil {
			log.Fatal(err)
		}
		defer rows.Close()

		for rows.Next() {
			var i IngredientMap
			err = rows.Scan(
				&i.IngredientName,
				&i.Amount,
			)
			if err != nil {
				log.Fatal(err)
			}
			ingredientsNeeded = append(ingredientsNeeded, i)
		}

		for i, ingN := range ingredientsNeeded {
			for _, item := range inventory {
				if item.IngredientName == ingN.IngredientName && item.Amount >= ingN.Amount {
					ingredientsNeeded[i].Got = true
				}
			}
		}

		check := true
		for _, ingN := range ingredientsNeeded {
			// log.Println(potion.Name, ingN)
			if ingN.Got == false {
				check = false
			}
		}

		if check == true {
            if slices.Contains(results, potion) == false {
                results = append(results, potion)
            }
		} else {
			ingredientsNeeded = ingredientsNeeded[:0]
		}
	}

	return results
}

type Ingredient struct {
	Id          string
	Name        string
	Description string
	url         string
}

type invItem struct {
	Id             string
	IngredientName string
	Amount         int
}

func (i invItem) Save() {
	tx, err := db.Begin()
	if err != nil {
		log.Fatal(err)
	}

	stmt, err := tx.Prepare("insert into inventory(ingredientName, amount) values(?, ?)")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(i.IngredientName, i.Amount)
	if err != nil {
		log.Fatal(err)
	}

	err = tx.Commit()
	if err != nil {
		log.Fatal(err)
	}
}

func (i invItem) Update() {
	tx, err := db.Begin()
	if err != nil {
		log.Fatal(err)
	}

	stmt, err := tx.Prepare("update inventory set ingredientName = ?, amount = ? where rowid = ?")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(i.IngredientName, i.Amount, i.Id)
	if err != nil {
		log.Fatal(err)
	}

	err = tx.Commit()
	if err != nil {
		log.Fatal(err)
	}
}

func (i invItem) Delete() {
	tx, err := db.Begin()
	if err != nil {
		log.Fatal(err)
	}

	stmt, err := tx.Prepare("delete from inventory where rowid = ?")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(i.Id)
	if err != nil {
		log.Fatal(err)
	}

	err = tx.Commit()
	if err != nil {
		log.Fatal(err)
	}
}

func (i invItem) Inc() {
	i.Amount += 1
	i.Update()
}

func (i invItem) Dec() {
	i.Amount -= 1

	if i.Amount <= 0 {
		i.Delete()
	} else {
		i.Update()
	}
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

func getIngredient(query string, key string) Ingredient {
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

func getAllInventory() []invItem {
	var results []invItem

	query := "select * from inventory"
	rows, err := db.Query(query)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		var i invItem
		err = rows.Scan(
			&i.Id,
			&i.IngredientName,
			&i.Amount,
		)
		if err != nil {
			log.Fatal(err)
		}
		results = append(results, i)
	}

	return results
}

func getInventory(query string, key string) invItem {
	var result invItem

	row := db.QueryRow(query, key)
	err := row.Scan(
		&result.Id,
		&result.IngredientName,
		&result.Amount,
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
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	r.HandleFunc("/", indexHandler).Methods("GET")
	r.HandleFunc("/", addInventory).Methods("POST")
	r.HandleFunc("/{ingredient}/add", itemAdd).Methods("POST")
	r.HandleFunc("/{ingredient}/subtract", itemSubtract).Methods("POST")
	http.Handle("/", r)

	log.Println("App running on localhost:8000")
	log.Fatal(http.ListenAndServe(":8000", nil))
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	data := map[string]interface{}{
		"Title":       "Index",
		"Ingredients": getAllIngredients(),
		"Inventory":   getAllInventory(),
		"Potions":     getAllPotions(),
	}

	templ.ExecuteTemplate(w, "base", data)
}

func addInventory(w http.ResponseWriter, r *http.Request) {
	// still need true validation
	query := "select * from ingredients where name like ?"
	key := "%" + r.FormValue("ingredient") + "%"
	ing := getIngredient(query, key)

	exists := false
	for _, invItem := range getAllInventory() {
		if invItem.IngredientName == ing.Name {
			exists = true
			invItem.Inc()
		}
	}

	if exists == false {
		item := invItem{IngredientName: ing.Name, Amount: 1}
		item.Save()
		templ.ExecuteTemplate(w, "invItem", item)
	} else {
		w.Header().Add("hx-trigger", "changedInv")
	}
}

func itemAdd(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	for _, item := range getAllInventory() {
		if item.IngredientName == vars["ingredient"] {
			item.Inc()
		}
	}

	w.Header().Add("hx-trigger", "changedInv")
}

func itemSubtract(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	for _, item := range getAllInventory() {
		if item.IngredientName == vars["ingredient"] {
			item.Dec()
		}
	}

	w.Header().Add("hx-trigger", "changedInv")
}
