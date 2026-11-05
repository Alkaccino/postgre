package main

import (
	"context"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"path"
	"strconv"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
)

type Asset struct {
	Id              int
	Asset_name      string
	Category        string
	Model           *string
	Manufacturer    *string
	Assigned_user   *string
	Asset_location  *string
	Purchased_at    *string
	Purchased_price *float32
	Quantity        *int
	Asset_status    *string
	Created_at      string
	Updated_at      string
	Created_by      string
	Updated_by      string
}

var conn *pgx.Conn
var ctx = context.Background()

func main() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)

	InitDatabase()
	defer conn.Close(ctx)

	http.Handle("/static/", http.FileServer(http.Dir(".")))
	http.HandleFunc("/", LandingHandler)
	http.HandleFunc("/create/", CreateHandler)
	http.HandleFunc("/update/", UpdateHandler)
	http.ListenAndServe(":8080", nil)
}

func InitDatabase() {
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

	create_sql, err := os.ReadFile("./sql/create_tables.sql")
	if err != nil {
		log.Fatal(err)
	}

	_, err = conn.Exec(ctx, string(create_sql))
	if err != nil {
		log.Fatal(err)
	}
}

func LandingHandler(w http.ResponseWriter, r *http.Request) {
	var assets []Asset

	select_sql, err := os.ReadFile("./sql/select_assets.sql")
	if err != nil {
		log.Fatal(err)
	}

	selected_assets, err := conn.Query(ctx, string(select_sql))
	if err != nil {
		log.Fatal(err)
	}

	for selected_assets.Next() {
		var asset Asset
		if err := selected_assets.Scan(&asset.Id, &asset.Asset_name, &asset.Category, &asset.Model, &asset.Manufacturer, &asset.Assigned_user, &asset.Asset_location, &asset.Purchased_price, &asset.Purchased_at, &asset.Quantity, &asset.Asset_status, &asset.Created_by, &asset.Updated_by, nil, nil); err != nil {
			log.Fatal(err)
		}

		assets = append(assets, asset)
	}

	defer selected_assets.Close()

	if r.Method == "POST" {
		delete_sql, err := os.ReadFile("./sql/delete_asset.sql")
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

	RenderTemplate(w, "index.html", assets)
}

func CreateHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {

		insert_query, err := os.ReadFile("./sql/insert_asset.sql")
		if err != nil {
			log.Fatal(err)
		}

		quantity, err := strconv.Atoi(r.FormValue("quantity"))
		if err != nil {
			http.Error(w, "Wrong input: Quantity has to be a number!", http.StatusBadRequest)
			return
		}

		purchased_price_string := r.FormValue("purchased_price")

		if strings.Contains(purchased_price_string, ",") {
			purchased_price_string = strings.Replace(purchased_price_string, ",", ".", -1)
		}

		purchased_price, err := strconv.ParseFloat(purchased_price_string, 32)
		if err != nil {
			http.Error(w, "Wrong input: Purchased price has to be a Float!", http.StatusBadRequest)
			return
		}

		insert_args := pgx.NamedArgs{
			"asset_name":      r.FormValue("asset_name"),
			"category":        r.FormValue("category"),
			"model":           r.FormValue("model"),
			"manufacturer":    r.FormValue("manufacturer"),
			"assigned_user":   r.FormValue("assigned_user"),
			"asset_location":  r.FormValue("asset_location"),
			"purchased_price": purchased_price,
			"purchased_at":    r.FormValue("purchased_at"),
			"quantity":        quantity,
			"asset_status":    r.FormValue("asset_status"),
			"created_by":      r.FormValue("created_by"),
			"updated_by":      r.FormValue("updated_by"),
		}

		_, err = conn.Exec(ctx, string(insert_query), insert_args)
		if err != nil {
			log.Fatal(err)
		}

		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	RenderTemplate(w, "create_asset.html", nil)
}

func UpdateHandler(w http.ResponseWriter, r *http.Request) {
	var asset Asset

	select_sql, err := os.ReadFile("./sql/select_asset.sql")
	if err != nil {
		log.Fatal(err)
	}

	id_value, err := strconv.Atoi(path.Base(r.URL.Path))
	if err != nil {
		log.Fatal(err)
	}

	selected_asset := conn.QueryRow(ctx, string(select_sql), pgx.NamedArgs{"id": id_value})
	if err != nil {
		log.Fatal(err)
	}

	err = selected_asset.Scan(&asset.Id, &asset.Asset_name, &asset.Category, &asset.Model, &asset.Manufacturer, &asset.Assigned_user, &asset.Asset_location, &asset.Purchased_price, &asset.Purchased_at, &asset.Quantity, &asset.Asset_status, &asset.Created_by, &asset.Updated_by, &asset.Created_at, &asset.Updated_at)

	if r.Method == "POST" {

		update_query, err := os.ReadFile("./sql/update_asset.sql")
		if err != nil {
			log.Fatal(err)
		}

		quantity, err := strconv.Atoi(r.FormValue("quantity"))
		if err != nil {
			http.Error(w, "Wrong input: Quantity has to be a number!", http.StatusBadRequest)
			return
		}

		purchased_price_string := r.FormValue("purchased_price")

		if strings.Contains(purchased_price_string, ",") {
			purchased_price_string = strings.Replace(purchased_price_string, ",", ".", -1)
		}

		purchased_price, err := strconv.ParseFloat(purchased_price_string, 32)
		if err != nil {
			http.Error(w, "Wrong input: Purchased price has to be a Float!", http.StatusBadRequest)
			return
		}

		id_value, err := strconv.Atoi(path.Base(r.URL.Path))
		if err != nil {
			log.Fatal(err)
		}

		update_args := pgx.NamedArgs{
			"id":              id_value,
			"asset_name":      r.FormValue("asset_name"),
			"category":        r.FormValue("category"),
			"model":           r.FormValue("model"),
			"manufacturer":    r.FormValue("manufacturer"),
			"assigned_user":   r.FormValue("assigned_user"),
			"asset_location":  r.FormValue("asset_location"),
			"purchased_price": purchased_price,
			"purchased_at":    r.FormValue("purchased_at"),
			"quantity":        quantity,
			"asset_status":    r.FormValue("asset_status"),
			"updated_at":      time.Now(),
			"created_by":      r.FormValue("created_by"),
			"updated_by":      r.FormValue("updated_by"),
		}

		_, err = conn.Exec(ctx, string(update_query), update_args)
		if err != nil {
			log.Fatal(err)
		}

		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	RenderTemplate(w, "update_asset.html", asset)
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
