package main

import (
	"log"
	"net/http"

    "github.com/a-h/templ"
	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"
)

type Potion struct {
    Name string
    Description string
}

func main() {
	r := mux.NewRouter()
	// r.PathPrefix("/").Handler(templ.Handler(home()))
    r.HandleFunc("/potions/", potionHandler)



	http.Handle("/", r)


	log.Println("App running on localhost:8000")
	log.Fatal(http.ListenAndServe(":8000", nil))
}


func potionHandler(w http.ResponseWriter, r *http.Request) {
    var potionSlice []Potion
    pot := Potion{Name: "fuck", Description: "it makes you fuck"}
    potionSlice = append(potionSlice, pot)

    templ.Handler(potions(potionSlice)).ServeHTTP(w, r)
}
