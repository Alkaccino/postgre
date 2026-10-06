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
	Id            int
	Item_name     string
	Item_category string
	Item_amount   int
}

var conn *pgx.Conn
var ctx = context.Background()

func main() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)

	InitDataBase()
	defer conn.Close(ctx)

	http.Handle("/static/", http.FileServer(http.Dir(".")))
	http.HandleFunc("/", LandingHandler)
	http.HandleFunc("/create/", CreateHandler)
	http.ListenAndServe(":8080", nil)
}

func InitDataBase() {
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

	conn, err = pgx.Connect(ctx, conn_str)
	if err != nil {
		log.Fatal(err)
	}

	create_sql, err := os.ReadFile("./sql/create_table.sql")
	if err != nil {
		log.Fatal(err)
	}

	_, err = conn.Exec(ctx, string(create_sql))
	if err != nil {
		log.Fatal(err)
	}
}

func LandingHandler(w http.ResponseWriter, r *http.Request) {
	var items []Item

	select_sql, err := os.ReadFile("./sql/select_items.sql")
	if err != nil {
		log.Fatal(err)
	}

	selected_items, err := conn.Query(ctx, string(select_sql))
	if err != nil {
		log.Fatal(err)
	}

	for selected_items.Next() {
		var item Item
		if err := selected_items.Scan(&item.Id, &item.Item_name, &item.Item_category, &item.Item_amount, nil); err != nil {
			log.Fatal(err)
		}

		items = append(items, item)
	}

	defer selected_items.Close()

	if r.Method == "POST" {
		delete_sql, err := os.ReadFile("./sql/delete_item.sql")
		if err != nil {
			log.Fatal(err)
		}

		id_value, err := strconv.Atoi(r.FormValue("delete"))
		if err != nil {
			log.Fatal(err)
		}

		delete_args := pgx.NamedArgs{
			"id": id_value,
		}

		_, err = conn.Exec(ctx, string(delete_sql), delete_args)
		if err != nil {
			log.Fatal(err)
		}

		http.Redirect(w, r, "/", http.StatusSeeOther)

		return
	}

	RenderTemplate(w, "index.html", items)
}

func CreateHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {

		// POST && Editmode (START) -------------------------------------

		id_value, err := strconv.Atoi(r.FormValue("edit"))
		if err != nil {
			log.Fatal(err)
		}

		select_sql, err := os.ReadFile("./sql/select_item.sql")
		if err != nil {
			log.Fatal(err)
		}

		select_args := pgx.NamedArgs{
			"id": id_value,
		}

		selected_item, err := conn.Query(ctx, string(select_sql), select_args)
		if err != nil {
			log.Fatal(err)
		}

		var item Item
		if err := selected_item.Scan(&item.Id, &item.Item_name, &item.Item_category, &item.Item_amount, nil); err != nil {
			log.Fatal(err)
		}

		fmt.Print(item)

		// POST && Editmode (END) ----------------------------------------

		if id_value == 0 {
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

			_, err = conn.Exec(ctx, string(insert_query), insert_args)
			if err != nil {
				log.Fatal(err)
			}

			r.ParseForm()
			http.Redirect(w, r, "/", http.StatusSeeOther)
		}
	}

	RenderTemplate(w, "create_item.html", nil)
}

func RenderTemplate(w http.ResponseWriter, templ_path string, templ_data any) {
	templ, err := template.ParseFiles("./templates/" + templ_path)
	if err != nil {
		log.Fatal(err)
	}

	err = templ.Execute(w, templ_data)
	if err != nil {
		log.Fatal(err)
	}
}
