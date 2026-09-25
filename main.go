package main

import (
	"context"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/jackc/pgx/v5"
)

type Item struct {
	ID       int
	Itemname string
	Category string
	Amount   int
}

func main() {
	env, err := os.ReadFile(".env")
	if err != nil {
		log.Fatal(err)
	}

	envData := strings.Split(string(env), "\n")

	_, db_user, _ := strings.Cut(envData[0], "=")
	_, db_password, _ := strings.Cut(envData[1], "=")
	_, db_host, _ := strings.Cut(envData[2], "=")
	_, db_port, _ := strings.Cut(envData[3], "=")
	_, db_name, _ := strings.Cut(envData[4], "=")

	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s", db_user, db_password, db_host, db_port, db_name)

	conn, err := pgx.Connect(context.Background(), connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close(context.Background())

	create_sql, err := os.ReadFile("./sql/create_table.sql")
	if err != nil {
		log.Fatal(err)
	}

	_, err = conn.Exec(context.Background(), string(create_sql))
	if err != nil {
		log.Fatal(err)
	}

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
	if err != nil {
		log.Fatal(err)
	}
}
