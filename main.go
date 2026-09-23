package main

import (
	"html/template"
	"log"
	"net/http"
)

type Item struct {
	ID       int
	Itemname string
	Category string
	Amount   int
}

func main() {
	http.HandleFunc("GET /", handler)
	http.ListenAndServe(":8080", nil)
}

func handler(w http.ResponseWriter, r *http.Request) {
	p := Item{1, "Testitem", "Testcategory", 2}

	templ, err := template.ParseFiles("./templates/index.html")
	if err != nil {
		log.Fatal(err)
	}
	err = templ.Execute(w, p)
}
