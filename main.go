package main

import (
	"context"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/jackc/pgx/v5"
)

type Item struct {
	ID            int
	Item_Name     string
	Item_Category string
	Item_Amount   int
}

var conn *pgx.Conn

func main() {
	env, err := os.ReadFile(".env")
	if err != nil {
		log.Fatal(err)
	}

	env_data := strings.Split(string(env), "\n")

	_, db_user, _ := strings.Cut(env_data[0], "=")
	_, db_password, _ := strings.Cut(env_data[1], "=")
	_, db_host, _ := strings.Cut(env_data[2], "=")
	_, db_port, _ := strings.Cut(env_data[3], "=")
	_, db_name, _ := strings.Cut(env_data[4], "=")

	conn_str := fmt.Sprintf("postgres://%s:%s@%s:%s/%s", db_user, db_password, db_host, db_port, db_name)

	conn, err = pgx.Connect(context.Background(), conn_str)
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

	http.HandleFunc("/", landingHandler)
	http.HandleFunc("/create/", createHandler)
	http.ListenAndServe(":8080", nil)
}

func landingHandler(w http.ResponseWriter, r *http.Request) {
	var Items []Item

	select_sql, err := os.ReadFile("./sql/select_items.sql")
	if err != nil {
		log.Fatal(err)
	}

	selected_items, err := conn.Query(context.Background(), string(select_sql))
	if err != nil {
		log.Fatal(err)
	}
	defer selected_items.Close()

	for selected_items.Next() {
		var itm Item
		if err := selected_items.Scan(&itm.ID, &itm.Item_Name, &itm.Item_Category, &itm.Item_Amount, nil); err != nil {
			log.Fatal(err)
		}

		Items = append(Items, itm)
	}

	renderTemplate(w, "index.html", Items)
}

func createHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {

		insert_query, err := os.ReadFile("./sql/insert_items.sql")
		if err != nil {
			log.Fatal(err)
		}

		item_amount, err := strconv.Atoi(r.FormValue("item_amount"))
		if err != nil {
			log.Fatal(err)
		}

		insert_args := pgx.NamedArgs{
			"item_name":     r.FormValue("item_name"),
			"item_category": r.FormValue("item_category"),
			"item_amount":   item_amount,
		}

		_, err = conn.Exec(context.Background(), string(insert_query), insert_args)

		if err != nil {
			log.Fatal(err)
		}

		r.ParseForm()
	}

	renderTemplate(w, "create_item.html", nil)
}

func renderTemplate(w http.ResponseWriter, templ string, templ_data any) {
	tmpl, err := template.ParseFiles("./templates/" + templ)
	if err != nil {
		log.Fatal(err)
	}

	err = tmpl.Execute(w, templ_data)
	if err != nil {
		log.Fatal(err)
	}
}
